// Copyright 2023 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package database

import (
	"context"

	gouuid "github.com/satori/go.uuid"
	log "unknwon.dev/clog/v2"

	"gogs.io/gogs/internal/auth"
	"gogs.io/gogs/internal/conf"
)

// AuthenticateByJWT 通过JWT令牌认证用户
// 返回认证的用户和认证是否成功
func AuthenticateByJWT(store *DB, ctx context.Context, token string, autoRegister bool) (*User, bool) {
	if !conf.Auth.EnableSSOWithMainSite || token == "" {
		return nil, false
	}

	// 验证令牌
	isValid, username := auth.VerifyJWTToken(token)
	if !isValid || username == "" {
		return nil, false
	}

	// 查找用户
	user, err := store.Users().GetByUsername(ctx, username)
	if err != nil {
		if !IsErrUserNotExist(err) {
			log.Error("Failed to get user by name: %v", err)
			return nil, false
		}

		// 用户不存在，如果启用了自动注册则创建用户
		if autoRegister {
			user, err = store.Users().Create(
				ctx,
				username,
				gouuid.NewV4().String()+"@localhost",
				CreateUserOptions{
					Activated: true,
				},
			)
			if err != nil {
				log.Error("Failed to create user %q: %v", username, err)
				return nil, false
			}
		} else {
			// 用户不存在且未启用自动注册
			log.Trace("User %q not found and auto-registration is disabled", username)
			return nil, false
		}
	}

	log.Trace("User %q authenticated via JWT token", username)
	return user, true
}
