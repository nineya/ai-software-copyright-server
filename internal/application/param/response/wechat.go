package response

type WechatOauth2AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Unionid      string `json:"unionid"` //针对一个微信开放平台账号下的应用，同一用户的 unionid 是唯一的
	Scope        string `json:"scope"`
}

type WechatUserInfoResponse struct {
	Openid     string `json:"openid"`
	Unionid    string `json:"unionid"` //针对一个微信开放平台账号下的应用，同一用户的 unionid 是唯一的
	Nickname   string `json:"nickname"`
	HeadImgUrl string `json:"headimgurl"` // 用户头像
}

type WechatSessionResponse struct {
	SessionKey string `json:"session_key"`
	ExpiresIn  int    `json:"expires_in"`
	Openid     string `json:"openid"`
	Unionid    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}

type WechatCgiBinAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type WechatMsgSecCheckResponse struct {
	Errcode int               `json:"errcode"`
	Errmsg  string            `json:"errmsg"`
	Result  MsgSecCheckResult `json:"result"`
}

type MsgSecCheckResult struct {
	Suggest string `json:"suggest"`
	Label   int    `json:"label"`
}
