package request

import "time"

type CdkeyCreateParam struct {
	Credits    int        `json:"credits" form:"credits" binding:"gte=0" label:"积分数量"`
	Count      int        `json:"count" form:"count" binding:"required,gte=0" label:"总兑换次数"`
	ExpireTime *time.Time `json:"expireTime" form:"expireTime" label:"失效时间"`
	CdkeyNum   int        `json:"cdkeyNum" form:"cdkeyNum" binding:"required,gte=0" label:"Cdkey数量"`
}

type CdkeyUseParam struct {
	Cdkey string `json:"cdkey" form:"cdkey" binding:"lte=52" label:"Cdkey"`
}
