package request

type AiWriteMessageParam struct {
	System string `json:"system" form:"system" binding:"lte=2000" label:"系统提示词"`
	User   string `json:"user" form:"user" binding:"required,lte=2000"  label:"用户提示词"`
}
