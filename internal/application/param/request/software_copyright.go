package request

type SCTriggerParam struct {
	Id     int64   `json:"id" form:"id" binding:"required" label:"软著申请ID"`
	ApiKey *string `json:"apiKey" form:"apiKey" binding:"lte=50"  label:"ApiKey"`
}
