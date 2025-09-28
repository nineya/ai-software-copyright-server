package request

type QrcodeBuildParam struct {
	Content string `json:"content" form:"content" binding:"required,lte=300" label:"二维码内容"`
}

type QrcodeDeleteImageParam struct {
	TargetUrl string `json:"targetUrl" form:"targetUrl" binding:"required,lte=300" label:"目标路径"`
}
