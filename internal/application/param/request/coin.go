package request

type CreditsCreateOrderParam struct {
	CreditsPriceId int64 `json:"creditsPriceId" form:"creditsPriceId" binding:"required" label:"购买积分ID"`
}
