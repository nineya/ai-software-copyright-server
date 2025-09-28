package response

type WeixinSessionResponse struct {
	SessionKey string `json:"session_key"`
	ExpiresIn  int    `json:"expires_in"`
	Openid     string `json:"openid"`
	Unionid    string `json:"unionid"`
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
}

type WeixinCgiBinAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type WeixinMsgSecCheckResponse struct {
	Errcode int               `json:"errcode"`
	Errmsg  string            `json:"errmsg"`
	Result  MsgSecCheckResult `json:"result"`
}

type MsgSecCheckResult struct {
	Suggest string `json:"suggest"`
	Label   int    `json:"label"`
}
