package request

type PageableParam struct {
	Page int `json:"page"  form:"page" binding:"gte=0"` // 从0开始
	Size int `json:"size"  form:"size" binding:"required,gte=1"`
}
