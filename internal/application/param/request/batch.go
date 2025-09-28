package request

type BatchParam struct {
	Ids []int64 `json:"ids" label:"ID列表"`
}
