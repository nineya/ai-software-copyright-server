package request

type ShortLinkCreateCloudDiskParam struct {
	Content string `json:"content" form:"content" binding:"required,lte=20000" label:"文案内容"`
}

type ShortLinkRedirectParam struct {
	SourceUrl string `json:"sourceUrl" form:"sourceUrl" binding:"required,lte=100" label:"源链接"`   // 源URL
	TargetUrl string `json:"targetUrl" form:"targetUrl" binding:"required,lte=500" label:"新目标链接"` // 目标URL
}
