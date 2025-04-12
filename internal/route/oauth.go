// Copyright 2025 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package route

import (
	"net/http"

	"gogs.io/gogs/internal/auth/oauth"
	"gogs.io/gogs/internal/context"
	"gopkg.in/macaron.v1"
)

func OAuthCallback(c *context.Context) {
	// 1. 验证state参数
	// 2. 获取授权码
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, map[string]string{
			"message": "missing authorization code",
		})
		return
	}

	// 3. 获取OAuth Provider
	provider, ok := c.Data["OAuthProvider"].(*oauth.Provider)
	if !ok {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "OAuth provider not configured",
		})
		return
	}

	// 4. 认证用户
	account, err := provider.Authenticate("", "") // 参数由授权码流程处理
	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{
			"message": err.Error(),
		})
		return
	}

	// 5. 登录用户
	c.Session.Set("uid", account.Login)
	c.Redirect("/")
}

func RegisterOAuthRoutes(m *macaron.Macaron) {
	m.Get("/auth/oauth/callback", OAuthCallback)
}