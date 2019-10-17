package yz_go

import (
	"net/http"
	"sync"
	"time"
)

type Type uint32

const (
	Self Type = iota
	Tool
)

type Client struct {
	Config               *YZConfig
	HTTPClient           *http.Client
	SelfAccessToken      string
	ToolAccessToken      string
	ToolRefreshToken     string
	SelfAccessTokenCache Cache
	ToolAccessTokenCache Cache
	RefreshTokenCache    Cache
	Locker               *sync.Mutex
}

type YZConfig struct {
	Type
	ClientID     string
	ClientSecret string
	YZID         string
	Code         string
	RedirectURI  string
}

func NewClient(config *YZConfig) *Client {
	c := &Client{
		Config: config,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		SelfAccessTokenCache: newFileCache("." + "self_access_token"),
		ToolAccessTokenCache: newFileCache("." + "tool_access_token"),
		RefreshTokenCache:    newFileCache("." + "refresh_token"),
		Locker:               new(sync.Mutex),
	}
	return c
}

func (c *Client) RefreshToken() error {
	if c.Config.Type == Self {
		return c.refreshSelfAccessToken()
	}
	return c.refreshToolAccessToken()
}

func (c *Client) refreshSelfAccessToken() error {
	c.Locker.Lock()
	defer c.Locker.Unlock()
	res := &GenSelfTokenResponse{}
	err := c.SelfAccessTokenCache.Get(res)
	if err == nil {
		c.SelfAccessToken = res.Data.AccessToken
		return nil
	}
	res, err = c.genSelfToken()
	if err != nil {
		return err
	}
	c.SelfAccessToken = res.Data.AccessToken
	res.Data.Expires = res.Data.Expires / 1000
	err = c.SelfAccessTokenCache.Set(res)
	return err
}

func (c *Client) refreshToolAccessToken() error {
	c.Locker.Lock()
	defer c.Locker.Unlock()
	res := &GenToolTokenResponse{}
	err := c.ToolAccessTokenCache.Get(res)
	if err == nil {
		c.ToolAccessToken = res.Data.AccessToken
		return nil
	}
	refresh := &GenToolTokenResponse{}
	err = c.RefreshTokenCache.Get(refresh)
	if err == nil {
		c.ToolRefreshToken = refresh.Data.RefreshToken
		res, err = c.refreshToken()
		if err != nil {
			res, err = c.genToolToken()
			if err != nil {
				return err
			}
		}
	} else {
		res, err = c.genToolToken()
		if err != nil {
			return err
		}
	}
	c.ToolAccessToken = res.Data.AccessToken
	res.Data.Expires = res.Data.Expires / 1000
	err = c.ToolAccessTokenCache.Set(res)
	if err != nil {
		return err
	}
	res.Data.Expires = res.Data.Expires + 3600*24*21
	return c.RefreshTokenCache.Set(res)
}

func (c *Client) genSelfToken() (*GenSelfTokenResponse, error) {
	res := &GenSelfTokenResponse{}
	req := &GenSelfTokenParams{
		GenTokenBaseParams: GenTokenBaseParams{
			AuthorizeType: "silent",
			ClientID:      c.Config.ClientID,
			ClientSecret:  c.Config.ClientSecret,
		},
		GrantID: c.Config.YZID,
	}
	err := c.httpJSON(TokenPath, nil, req, res)
	return res, err
}

func (c *Client) genToolToken() (*GenToolTokenResponse, error) {
	res := &GenToolTokenResponse{}
	req := &GenToolTokenParams{
		GenTokenBaseParams: GenTokenBaseParams{
			AuthorizeType: "authorization_code",
			ClientID:      c.Config.ClientID,
			ClientSecret:  c.Config.ClientSecret,
		},
		Code:        c.Config.Code,
		RedirectURI: c.Config.RedirectURI,
	}
	err := c.httpJSON(TokenPath, nil, req, res)
	return res, err
}

func (c *Client) refreshToken() (*GenToolTokenResponse, error) {
	res := &GenToolTokenResponse{}
	req := &RefreshTokenParams{
		GenTokenBaseParams: GenTokenBaseParams{
			AuthorizeType: "refresh_token",
			ClientID:      c.Config.ClientID,
			ClientSecret:  c.Config.ClientSecret,
		},
		RefreshToken: c.ToolRefreshToken,
	}
	err := c.httpJSON(TokenPath, nil, req, res)
	return res, err
}
