package yz_go

import (
	"errors"
	"fmt"
	"net/url"
)

func (c *Client)Call(api, version string, params interface{}) ([]byte, error) {
	if api == "" {
		return nil, errors.New("api is required")
	}
	if version == "" {
		return nil, errors.New("version is required")
	}
	path := fmt.Sprintf("%s/%s", api, version)
	req := url.Values{}
	if c.Config.Type == Self {
		req.Add("access_token", c.SelfAccessToken)
	} else {
		req.Add("access_token", c.ToolAccessToken)
	}
	return c.httpRequest(path,req, params)
}
