package request

type QueryParam struct {
	Keyword string `json:"keyword" form:"keyword" label:"关键词"`
}

type QueryPageParam struct {
	PageableParam
	QueryParam
}
