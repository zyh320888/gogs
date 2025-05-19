// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package repo

import (
	"net/http"
	"path"

	api "github.com/gogs/go-gogs-client"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"gogs.io/gogs/internal/conf"
	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/form"
	"gogs.io/gogs/internal/route/api/v1/convert"
)

// 临时添加，直到go-gogs-client的更新被正确集成
type ForkRepoOption struct {
	RepoName     string `json:"repo_name" binding:"Required;AlphaDashDot;MaxSize(100)"`
	Description  string `json:"description" binding:"MaxSize(255)"`
	Organization string `json:"organization"` // 可选，如果要fork到组织
}

func Search(c *context.APIContext) {
	opts := &database.SearchRepoOptions{
		Keyword:  path.Base(c.Query("q")),
		OwnerID:  c.QueryInt64("uid"),
		PageSize: convert.ToCorrectPageSize(c.QueryInt("limit")),
		Page:     c.QueryInt("page"),
	}

	// Check visibility.
	if c.IsLogged && opts.OwnerID > 0 {
		if c.User.ID == opts.OwnerID {
			opts.Private = true
		} else {
			u, err := database.Handle.Users().GetByID(c.Req.Context(), opts.OwnerID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, map[string]any{
					"ok":    false,
					"error": err.Error(),
				})
				return
			}
			if u.IsOrganization() && u.IsOwnedBy(c.User.ID) {
				opts.Private = true
			}
			// FIXME: how about collaborators?
		}
	}

	repos, count, err := database.SearchRepositoryByName(opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}

	if err = database.RepositoryList(repos).LoadAttributes(); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"ok":    false,
			"error": err.Error(),
		})
		return
	}

	results := make([]*api.Repository, len(repos))
	for i := range repos {
		results[i] = repos[i].APIFormatLegacy(nil)
	}

	c.SetLinkHeader(int(count), opts.PageSize)
	c.JSONSuccess(map[string]any{
		"ok":   true,
		"data": results,
	})
}

// listUserRepositories 列出指定用户的所有仓库
func listUserRepositories(c *context.APIContext, username string) {
	// 根据用户名获取用户信息
	user, err := database.Handle.Users().GetByUsername(c.Req.Context(), username)
	if err != nil {
		c.NotFoundOrError(err, "get user by name")
		return
	}

	// 如果请求的是其他用户的仓库列表，或者不是组织成员，则只列出公开仓库
	var ownRepos []*database.Repository
	if user.IsOrganization() {
		// 获取组织的仓库
		ownRepos, _, err = user.GetUserRepositories(c.User.ID, 1, user.NumRepos)
	} else {
		// 获取个人用户的仓库
		ownRepos, err = database.GetUserRepositories(&database.UserRepoOptions{
			UserID:   user.ID,
			Private:  c.User.ID == user.ID || c.User.IsAdmin, // 如果是自己或管理员，可以查看私有仓库
			Page:     1,
			PageSize: user.NumRepos,
		})
	}
	if err != nil {
		c.Error(err, "get user repositories")
		return
	}

	// 加载仓库的额外属性
	if err = database.RepositoryList(ownRepos).LoadAttributes(); err != nil {
		c.Error(err, "load attributes")
		return
	}

	// 如果是查询其他用户的仓库且不是管理员，则提前返回(只返回公开仓库)
	if c.User.ID != user.ID && !c.User.IsAdmin {
		repos := make([]*api.Repository, len(ownRepos))
		for i := range ownRepos {
			repos[i] = ownRepos[i].APIFormatLegacy(&api.Permission{Admin: true, Push: true, Pull: true})
		}
		c.JSONSuccess(&repos)
		return
	}

	// 获取用户作为协作者的仓库及其访问权限
	accessibleRepos, err := database.Handle.Repositories().GetByCollaboratorIDWithAccessMode(c.Req.Context(), user.ID)
	if err != nil {
		c.Error(err, "get repositories accesses by collaborator")
		return
	}

	// 合并个人仓库和协作仓库
	numOwnRepos := len(ownRepos)
	repos := make([]*api.Repository, 0, numOwnRepos+len(accessibleRepos))
	
	// 添加个人仓库，拥有完全权限
	for _, r := range ownRepos {
		repos = append(repos, r.APIFormatLegacy(&api.Permission{Admin: true, Push: true, Pull: true}))
	}

	// 将仓库对象转换为 RepositoryList 以便加载属性  
	repoList := make(database.RepositoryList, 0, len(accessibleRepos))  
	for repo := range accessibleRepos {  
		repoList = append(repoList, repo)  
	}  
	
	// 加载仓库的关联属性（包括 Owner 信息）  
	if err = repoList.LoadAttributes(); err != nil {  
		c.Error(err, "load repository attributes")  
		return  
	}  

	// 添加协作仓库，根据访问权限设置相应权限
	for repo, access := range accessibleRepos {
		repos = append(repos,
			repo.APIFormatLegacy(&api.Permission{
				Admin: access >= database.AccessModeAdmin,
				Push:  access >= database.AccessModeWrite,
				Pull:  true,
			}),
		)
	}

	// 返回合并后的仓库列表
	c.JSONSuccess(&repos)
}

func ListMyRepos(c *context.APIContext) {
	listUserRepositories(c, c.User.Name)
}

func ListUserRepositories(c *context.APIContext) {
	listUserRepositories(c, c.Params(":username"))
}

func ListOrgRepositories(c *context.APIContext) {
	listUserRepositories(c, c.Params(":org"))
}

func CreateUserRepo(c *context.APIContext, owner *database.User, opt api.CreateRepoOption) {
	repo, err := database.CreateRepository(c.User, owner, database.CreateRepoOptionsLegacy{
		Name:        opt.Name,
		Description: opt.Description,
		Gitignores:  opt.Gitignores,
		License:     opt.License,
		Readme:      opt.Readme,
		IsPrivate:   opt.Private,
		AutoInit:    opt.AutoInit,
	})
	if err != nil {
		if database.IsErrRepoAlreadyExist(err) ||
			database.IsErrNameNotAllowed(err) {
			c.ErrorStatus(http.StatusUnprocessableEntity, err)
		} else {
			if repo != nil {
				if err = database.DeleteRepository(c.User.ID, repo.ID); err != nil {
					log.Error("Failed to delete repository: %v", err)
				}
			}
			c.Error(err, "create repository")
		}
		return
	}

	c.JSON(201, repo.APIFormatLegacy(&api.Permission{Admin: true, Push: true, Pull: true}))
}

func Create(c *context.APIContext, opt api.CreateRepoOption) {
	// Shouldn't reach this condition, but just in case.
	if c.User.IsOrganization() {
		c.ErrorStatus(http.StatusUnprocessableEntity, errors.New("Not allowed to create repository for organization."))
		return
	}
	CreateUserRepo(c, c.User, opt)
}

func CreateOrgRepo(c *context.APIContext, opt api.CreateRepoOption) {
	org, err := database.GetOrgByName(c.Params(":org"))
	if err != nil {
		c.NotFoundOrError(err, "get organization by name")
		return
	}

	if !org.IsOwnedBy(c.User.ID) {
		c.ErrorStatus(http.StatusForbidden, errors.New("Given user is not owner of organization."))
		return
	}
	CreateUserRepo(c, org, opt)
}

func Migrate(c *context.APIContext, f form.MigrateRepo) {
	ctxUser := c.User
	// Not equal means context user is an organization,
	// or is another user/organization if current user is admin.
	if f.Uid != ctxUser.ID {
		org, err := database.Handle.Users().GetByID(c.Req.Context(), f.Uid)
		if err != nil {
			if database.IsErrUserNotExist(err) {
				c.ErrorStatus(http.StatusUnprocessableEntity, err)
			} else {
				c.Error(err, "get user by ID")
			}
			return
		} else if !org.IsOrganization() && !c.User.IsAdmin {
			c.ErrorStatus(http.StatusForbidden, errors.New("Given user is not an organization."))
			return
		}
		ctxUser = org
	}

	if c.HasError() {
		c.ErrorStatus(http.StatusUnprocessableEntity, errors.New(c.GetErrMsg()))
		return
	}

	if ctxUser.IsOrganization() && !c.User.IsAdmin {
		// Check ownership of organization.
		if !ctxUser.IsOwnedBy(c.User.ID) {
			c.ErrorStatus(http.StatusForbidden, errors.New("Given user is not owner of organization."))
			return
		}
	}

	remoteAddr, err := f.ParseRemoteAddr(c.User)
	if err != nil {
		if database.IsErrInvalidCloneAddr(err) {
			addrErr := err.(database.ErrInvalidCloneAddr)
			switch {
			case addrErr.IsURLError:
				c.ErrorStatus(http.StatusUnprocessableEntity, err)
			case addrErr.IsPermissionDenied:
				c.ErrorStatus(http.StatusUnprocessableEntity, errors.New("You are not allowed to import local repositories."))
			case addrErr.IsInvalidPath:
				c.ErrorStatus(http.StatusUnprocessableEntity, errors.New("Invalid local path, it does not exist or not a directory."))
			case addrErr.IsBlockedLocalAddress:
				c.ErrorStatus(http.StatusUnprocessableEntity, errors.New("Clone address resolved to a local network address that is implicitly blocked."))
			default:
				c.Error(err, "unexpected error")
			}
		} else {
			c.Error(err, "parse remote address")
		}
		return
	}

	repo, err := database.MigrateRepository(c.User, ctxUser, database.MigrateRepoOptions{
		Name:        f.RepoName,
		Description: f.Description,
		IsPrivate:   f.Private || conf.Repository.ForcePrivate,
		IsMirror:    f.Mirror,
		RemoteAddr:  remoteAddr,
	})
	if err != nil {
		if repo != nil {
			if errDelete := database.DeleteRepository(ctxUser.ID, repo.ID); errDelete != nil {
				log.Error("DeleteRepository: %v", errDelete)
			}
		}

		if database.IsErrReachLimitOfRepo(err) {
			c.ErrorStatus(http.StatusUnprocessableEntity, err)
		} else {
			c.Error(errors.New(database.HandleMirrorCredentials(err.Error(), true)), "migrate repository")
		}
		return
	}

	log.Trace("Repository migrated: %s/%s", ctxUser.Name, f.RepoName)
	c.JSON(201, repo.APIFormatLegacy(&api.Permission{Admin: true, Push: true, Pull: true}))
}

// FIXME: inject in the handler chain
func parseOwnerAndRepo(c *context.APIContext) (*database.User, *database.Repository) {
	owner, err := database.Handle.Users().GetByUsername(c.Req.Context(), c.Params(":username"))
	if err != nil {
		if database.IsErrUserNotExist(err) {
			c.ErrorStatus(http.StatusUnprocessableEntity, err)
		} else {
			c.Error(err, "get user by name")
		}
		return nil, nil
	}

	repo, err := database.GetRepositoryByName(owner.ID, c.Params(":reponame"))
	if err != nil {
		c.NotFoundOrError(err, "get repository by name")
		return nil, nil
	}

	return owner, repo
}

func Get(c *context.APIContext) {
	_, repo := parseOwnerAndRepo(c)
	if c.Written() {
		return
	}

	c.JSONSuccess(repo.APIFormatLegacy(&api.Permission{
		Admin: c.Repo.IsAdmin(),
		Push:  c.Repo.IsWriter(),
		Pull:  true,
	}))
}

func Delete(c *context.APIContext) {
	owner, repo := parseOwnerAndRepo(c)
	if c.Written() {
		return
	}

	if owner.IsOrganization() && !owner.IsOwnedBy(c.User.ID) {
		c.ErrorStatus(http.StatusForbidden, errors.New("Given user is not owner of organization."))
		return
	}

	if err := database.DeleteRepository(owner.ID, repo.ID); err != nil {
		c.Error(err, "delete repository")
		return
	}

	log.Trace("Repository deleted: %s/%s", owner.Name, repo.Name)
	c.NoContent()
}

func CreateFork(c *context.APIContext, form ForkRepoOption) {
	// 获取要fork的源仓库
	_, repo := parseOwnerAndRepo(c)
	if c.Written() {
		return
	}

	// 确认仓库可以被fork
	if !repo.CanBeForked() {
		c.ErrorStatus(http.StatusBadRequest, errors.New("Repository cannot be forked"))
		return
	}

	// 确定目标用户/组织
	var targetOwner *database.User
	var err error

	if len(form.Organization) > 0 {
		// Fork到指定组织
		targetOwner, err = database.Handle.Users().GetByUsername(c.Req.Context(), form.Organization)
		if err != nil {
			c.NotFoundOrError(err, "get organization by name")
			return
		}

		// 检查用户是否有权限向该组织贡献
		if !targetOwner.IsOrganization() || !targetOwner.IsOwnedBy(c.User.ID) {
			c.ErrorStatus(http.StatusForbidden, errors.New("Given organization is either not an organization or you're not an owner"))
			return
		}
	} else {
		// Fork到当前用户
		targetOwner = c.User
	}

	// 检查是否已经fork过
	existRepo, has, err := database.HasForkedRepo(targetOwner.ID, repo.ID)
	if err != nil {
		c.Error(err, "check if already forked")
		return
	} else if has {
		c.JSONSuccess(existRepo.APIFormatLegacy(&api.Permission{Admin: true, Push: true, Pull: true}))
		return
	}

	// 创建fork
	repoName := form.RepoName
	if len(repoName) == 0 {
		repoName = repo.Name
	}

	forkedRepo, err := database.ForkRepository(c.User, targetOwner, repo, repoName, form.Description)
	if err != nil {
		if database.IsErrRepoAlreadyExist(err) {
			c.ErrorStatus(http.StatusUnprocessableEntity, err)
		} else if database.IsErrNameNotAllowed(err) {
			c.ErrorStatus(http.StatusUnprocessableEntity, err)
		} else {
			c.Error(err, "fork repository")
		}
		return
	}

	log.Trace("Repository forked from '%s' -> '%s'", repo.FullName(), forkedRepo.FullName())
	c.JSON(201, forkedRepo.APIFormatLegacy(&api.Permission{Admin: true, Push: true, Pull: true}))
}

func ListForks(c *context.APIContext) {
	forks, err := c.Repo.Repository.GetForks()
	if err != nil {
		c.Error(err, "get forks")
		return
	}

	apiForks := make([]*api.Repository, len(forks))
	for i := range forks {
		if err := forks[i].GetOwner(); err != nil {
			c.Error(err, "get owner")
			return
		}

		accessMode := database.Handle.Permissions().AccessMode(
			c.Req.Context(),
			c.User.ID,
			forks[i].ID,
			database.AccessModeOptions{
				OwnerID: forks[i].OwnerID,
				Private: forks[i].IsPrivate,
			},
		)

		apiForks[i] = forks[i].APIFormatLegacy(
			&api.Permission{
				Admin: accessMode >= database.AccessModeAdmin,
				Push:  accessMode >= database.AccessModeWrite,
				Pull:  true,
			},
		)
	}

	c.JSONSuccess(&apiForks)
}

func IssueTracker(c *context.APIContext, form api.EditIssueTrackerOption) {
	_, repo := parseOwnerAndRepo(c)
	if c.Written() {
		return
	}

	if form.EnableIssues != nil {
		repo.EnableIssues = *form.EnableIssues
	}
	if form.EnableExternalTracker != nil {
		repo.EnableExternalTracker = *form.EnableExternalTracker
	}
	if form.ExternalTrackerURL != nil {
		repo.ExternalTrackerURL = *form.ExternalTrackerURL
	}
	if form.TrackerURLFormat != nil {
		repo.ExternalTrackerFormat = *form.TrackerURLFormat
	}
	if form.TrackerIssueStyle != nil {
		repo.ExternalTrackerStyle = *form.TrackerIssueStyle
	}

	if err := database.UpdateRepository(repo, false); err != nil {
		c.Error(err, "update repository")
		return
	}

	c.NoContent()
}

func Wiki(c *context.APIContext, form api.EditWikiOption) {
	_, repo := parseOwnerAndRepo(c)
	if c.Written() {
		return
	}

	if form.AllowPublicWiki != nil {
		repo.AllowPublicWiki = *form.AllowPublicWiki
	}
	if form.EnableExternalWiki != nil {
		repo.EnableExternalWiki = *form.EnableExternalWiki
	}
	if form.EnableWiki != nil {
		repo.EnableWiki = *form.EnableWiki
	}
	if form.ExternalWikiURL != nil {
		repo.ExternalWikiURL = *form.ExternalWikiURL
	}
	if err := database.UpdateRepository(repo, false); err != nil {
		c.Error(err, "update repository")
		return
	}

	c.NoContent()
}

func MirrorSync(c *context.APIContext) {
	_, repo := parseOwnerAndRepo(c)
	if c.Written() {
		return
	} else if !repo.IsMirror {
		c.NotFound()
		return
	}

	go database.MirrorQueue.Add(repo.ID)
	c.Status(http.StatusAccepted)
}

func Releases(c *context.APIContext) {
	_, repo := parseOwnerAndRepo(c)
	releases, err := database.GetReleasesByRepoID(repo.ID)
	if err != nil {
		c.Error(err, "get releases by repository ID")
		return
	}
	apiReleases := make([]*api.Release, 0, len(releases))
	for _, r := range releases {
		publisher, err := database.Handle.Users().GetByID(c.Req.Context(), r.PublisherID)
		if err != nil {
			c.Error(err, "get release publisher")
			return
		}
		r.Publisher = publisher
	}
	for _, r := range releases {
		apiReleases = append(apiReleases, r.APIFormat())
	}

	c.JSONSuccess(&apiReleases)
}
