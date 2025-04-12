// Copyright 2025 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package oauth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"gogs.io/gogs/internal/auth"
)

type Config struct {
	Enabled       bool
	Name          string
	Icon          string
	ClientID      string
	ClientSecret  string
	AuthURL       string
	TokenURL      string
	APIURL        string
	Scopes        string
	UserNameField string
	UserEmailField string
	SkipVerify    bool
}

type Provider struct {
	config *Config
}

func NewProvider(cfg *Config) auth.Provider {
	// 验证必要配置
	if cfg.ClientID == "" || cfg.ClientSecret == "" ||
	   cfg.AuthURL == "" || cfg.TokenURL == "" || cfg.APIURL == "" {
		return nil
	}
	
	return &Provider{
		config: cfg,
	}
}

func (p *Provider) Authenticate(login, password string) (*auth.ExternalAccount, error) {
	// 1. 获取授权码
	authCode, err := p.getAuthCode()
	if err != nil {
		return nil, err
	}

	// 2. 交换访问令牌
	token, err := p.exchangeToken(authCode)
	if err != nil {
		return nil, err
	}

	// 3. 获取用户信息
	userInfo, err := p.getUserInfo(token)
	if err != nil {
		return nil, err
	}

	return &auth.ExternalAccount{
		Login:    userInfo[p.config.UserNameField],
		Name:     userInfo[p.config.UserNameField],
		Email:    userInfo[p.config.UserEmailField],
	}, nil
}

func (p *Provider) getAuthCode() (string, error) {
	// 重定向用户到授权页面获取授权码
	authURL, err := url.Parse(p.config.AuthURL)
	if err != nil {
		return "", err
	}
	
	params := url.Values{}
	params.Add("client_id", p.config.ClientID)
	params.Add("redirect_uri", "http://gogs-server/auth/oauth/callback")
	params.Add("response_type", "code")
	params.Add("scope", p.config.Scopes)
	params.Add("state", "random-state-string")
	
	authURL.RawQuery = params.Encode()
	return authURL.String(), nil
}

func (p *Provider) exchangeToken(code string) (string, error) {
	reqBody := url.Values{}
	reqBody.Set("grant_type", "authorization_code")
	reqBody.Set("code", code)
	reqBody.Set("client_id", p.config.ClientID)
	reqBody.Set("client_secret", p.config.ClientSecret)
	reqBody.Set("redirect_uri", "http://gogs-server/auth/oauth/callback")

	req, err := http.NewRequest("POST", p.config.TokenURL, bytes.NewBufferString(reqBody.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", auth.ErrBadCredentials{}
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func (p *Provider) getUserInfo(token string) (map[string]string, error) {
	req, err := http.NewRequest("GET", p.config.APIURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, auth.ErrBadCredentials{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	userInfo := make(map[string]string)
	for k, v := range result {
		if strVal, ok := v.(string); ok {
			userInfo[k] = strVal
		}
	}

	return userInfo, nil
}

func (p *Provider) Config() any {
	return p.config
}

func (*Provider) HasTLS() bool {
	return true
}

func (*Provider) UseTLS() bool {
	return true
}

func (p *Provider) SkipTLSVerify() bool {
	return p.config.SkipVerify
}