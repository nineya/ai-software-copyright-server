package request

type TimeClockQueryPageParam struct {
	QueryPageParam
	ClockId int64 `json:"clockId" form:"clockId"`
}
