package request

type PageableParam struct {
	Page int `json:"page"  form:"page" binding:"gte=0" label:"分页序号"` // 从0开始
	Size int `json:"size"  form:"size" binding:"required,gte=1" label:"分页大小"`
}
