// Copyright 2015 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package admin

import (
	"net/http"

	api "github.com/gogs/go-gogs-client"
	"github.com/pkg/errors"
	log "unknwon.dev/clog/v2"

	"gogs.io/gogs/internal/context"
	"gogs.io/gogs/internal/database"
	"gogs.io/gogs/internal/route/api/v1/repo"
	"gogs.io/gogs/internal/route/api/v1/user"
)

// 临时添加，直到go-gogs-client的更新被正确集成
type ForkRepoOption struct {
	RepoName     string `json:"repo_name" binding:"Required;AlphaDashDot;MaxSize(100)"`
	Description  string `json:"description" binding:"MaxSize(255)"`
	Organization string `json:"organization"` // 可选，如果要fork到组织
}

func CreateRepo(c *context.APIContext, form api.CreateRepoOption) {
	owner := user.GetUserByParams(c)
	if c.Written() {
		return
	}

	repo.CreateUserRepo(c, owner, form)
}

// ForkRepo creates a fork from source repository to target user
func ForkRepo(c *context.APIContext, form ForkRepoOption) {
	// 获取目标用户（接收fork的用户）
	targetOwner := user.GetUserByParams(c)
	if c.Written() {
		return
	}

	// 获取源仓库所有者和仓库
	repoUser, repoName := c.Query("repo_owner"), c.Query("repo_name")
	if len(repoUser) == 0 || len(repoName) == 0 {
		c.ErrorStatus(http.StatusBadRequest, errors.New("repo_owner and repo_name parameters are required"))
		return
	}

	repoOwner, err := database.Handle.Users().GetByUsername(c.Req.Context(), repoUser)
	if err != nil {
		c.NotFoundOrError(err, "get source repo owner")
		return
	}

	sourceRepo, err := database.GetRepositoryByName(repoOwner.ID, repoName)
	if err != nil {
		c.NotFoundOrError(err, "get source repository")
		return
	}

	// 确认仓库可以被fork
	if !sourceRepo.CanBeForked() {
		c.ErrorStatus(http.StatusBadRequest, errors.New("Repository cannot be forked"))
		return
	}

	// 检查是否已经fork过
	existRepo, has, err := database.HasForkedRepo(targetOwner.ID, sourceRepo.ID)
	if err != nil {
		c.Error(err, "check if already forked")
		return
	} else if has {
		c.JSONSuccess(existRepo.APIFormatLegacy(&api.Permission{Admin: true, Push: true, Pull: true}))
		return
	}

	// 创建fork
	forkName := form.RepoName
	if len(forkName) == 0 {
		forkName = sourceRepo.Name
	}

	forkedRepo, err := database.ForkRepository(c.User, targetOwner, sourceRepo, forkName, form.Description)
	if err != nil {
		if database.IsErrRepoAlreadyExist(err) || database.IsErrNameNotAllowed(err) {
			c.ErrorStatus(http.StatusUnprocessableEntity, err)
		} else {
			c.Error(err, "fork repository")
		}
		return
	}

	log.Trace("[Admin] Repository forked from '%s' -> '%s' by admin '%s'", sourceRepo.FullName(), forkedRepo.FullName(), c.User.Name)
	c.JSON(201, forkedRepo.APIFormatLegacy(&api.Permission{Admin: true, Push: true, Pull: true}))
}
