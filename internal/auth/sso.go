// Copyright 2023 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"net/http"
	"strings"
	"time"

	log "unknwon.dev/clog/v2"

	"gogs.io/gogs/internal/conf"
	"gopkg.in/macaron.v1"
)

// GetJWTTokenFromRequest 从请求中提取JWT令牌，优先从Authorization头获取，其次从Cookie获取
func GetJWTTokenFromRequest(req *http.Request) string {
	// 尝试从Authorization头获取Bearer令牌
	bearerToken := ""
	authHeader := req.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		bearerToken = strings.TrimPrefix(authHeader, "Bearer ")
	}

	// 如果没有Bearer令牌，尝试从Cookie获取
	if bearerToken == "" {
		cookie, err := req.Cookie(conf.Auth.SSOCookieName)
		if err == nil && cookie != nil {
			bearerToken = cookie.Value
		}
	}

	return bearerToken
}

// VerifyJWTToken 验证来自主站的JWT令牌，并返回用户名
func VerifyJWTToken(token string) (bool, string) {
	// 如果未启用SSO，直接返回
	if !conf.Auth.EnableSSOWithMainSite {
		return false, ""
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 创建请求
	req, err := http.NewRequest("GET", conf.Auth.MainSiteVerifyURL, nil)
	if err != nil {
		log.Error("Failed to create HTTP request for JWT verification: %v", err)
		return false, ""
	}

	// 添加令牌到请求头
	req.Header.Add("Authorization", "Bearer "+token)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Failed to send HTTP request for JWT verification: %v", err)
		return false, ""
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		log.Trace("JWT verification failed with status: %d", resp.StatusCode)
		return false, ""
	}

	// 从响应头获取用户名
	username := resp.Header.Get("X-Username")
	if username == "" {
		log.Error("JWT verification succeeded but no username returned")
		return false, ""
	}

	log.Trace("JWT verification succeeded for user: %s", username)
	return true, username
}

// SetSSOCookie 设置SSO认证Cookie
func SetSSOCookie(c *macaron.Context, token string, maxAge int) {
	if conf.Auth.EnableSSOWithMainSite {
		c.SetCookie(conf.Auth.SSOCookieName, token, maxAge, "/", conf.Auth.SSOCookieDomain, conf.Security.CookieSecure, true)
	}
}

// ClearSSOCookie 清除SSO认证Cookie
func ClearSSOCookie(c *macaron.Context) {
	if conf.Auth.EnableSSOWithMainSite {
		c.SetCookie(conf.Auth.SSOCookieName, "", -1, "/", conf.Auth.SSOCookieDomain, conf.Security.CookieSecure, true)
	}
}
