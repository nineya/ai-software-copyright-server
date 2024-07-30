package request

type QueryParam struct {
	Keyword string `json:"keyword" form:"keyword"`
}

type QueryPageParam struct {
	PageableParam
	QueryParam
}
