package request

type RedbookProhibitedDetectionParam struct {
	Content string `json:"content" form:"content" binding:"required,lte=2000" label:"内容"`
}

type RedbookWriteMessageParam struct {
	Message string `json:"message" form:"message" binding:"required,lte=2000" label:"消息"`
}

type RedbookParam struct {
	Html string `json:"html" form:"html" label:"网页内容"`
	Url  string `json:"url" form:"url" binding:"lte=1000" label:"目标路径"`
}
