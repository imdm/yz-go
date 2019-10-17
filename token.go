package yz_go

type GenTokenBaseParams struct {
	AuthorizeType string `json:"authorize_type"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
}

// GenSelfTokenParams 获取自用型AccessToken
type GenSelfTokenParams struct {
	GenTokenBaseParams
	GrantID string `json:"grant_id"`
}

// GenToolTokenParams  获取工具型AccessToken
type GenToolTokenParams struct {
	GenTokenBaseParams
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

// RefreshTokenParams 工具型应用刷新token
type RefreshTokenParams struct {
	GenTokenBaseParams
	RefreshToken string `json:"refresh_token"`
}

// GenToolTokenResponse 工具型AccessToken响应参数结构体
type GenToolTokenResponse struct {
	YZBaseResponse
	Data ToolToken
}

type ToolToken struct {
	AccessToken  string `json:"access_token"`
	Expires      int64  `json:"expires"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

func (t *GenToolTokenResponse) ExpiresIn() int64 {
	return t.Data.Expires
}

// GenSelfTokenResponse 自用型AccessToken响应参数结构体
type GenSelfTokenResponse struct {
	YZBaseResponse
	Data SelfToken
}

type SelfToken struct {
	AccessToken string `json:"access_token"`
	Expires     int64  `json:"expires"`
	Scope       string `json:"scope"`
	Created     int64
}

func (t *GenSelfTokenResponse) ExpiresIn() int64 {
	return t.Data.Expires
}
