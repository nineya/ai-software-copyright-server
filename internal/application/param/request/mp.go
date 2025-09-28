package request

type MpImageTextMessageParam struct {
	Message string `json:"message" form:"message" binding:"required,lte=2000" label:"消息"`
}

type MpVerifyParam struct {
	Signature string `json:"signature" form:"signature" binding:"required,lte=200" label:"签名"`
	Timestamp string `json:"timestamp" form:"timestamp" binding:"required" label:"时间戳"`
	Nonce     string `json:"nonce" form:"nonce" binding:"required" label:"随机数"`
	Echostr   string `json:"echostr" form:"echostr" binding:"required" label:"随机字符串"`
}
