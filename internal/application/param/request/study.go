package request

type StudyResourceQueryPageParam struct {
	QueryPageParam
	Type string `json:"type" form:"type"`
}
