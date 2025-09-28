package request

import "time"

type CdkeyCreateParam struct {
	CreditsNum        int        `json:"creditsNum" form:"creditsNum" binding:"gte=0" label:"币数量"`
	HelperStandardDay int        `json:"helperStandardDay" form:"helperStandardDay" binding:"gte=0" label:"网盘助手标准版赠送天数"`
	HelperWechatDay   int        `json:"helperWechatDay" form:"helperWechatDay" binding:"gte=0" label:"网盘助手微信版赠送天数"`
	Count             int        `json:"count" form:"count" binding:"required,gte=0" label:"总兑换次数"`
	ExpireTime        *time.Time `json:"expireTime" form:"expireTime" label:"失效时间"`
	CdkeyNum          int        `json:"cdkeyNum" form:"cdkeyNum" binding:"required,gte=0" label:"Cdkey数量"`
}

type CdkeyUseParam struct {
	Cdkey string `json:"cdkey" form:"cdkey" binding:"lte=52" label:"Cdkey"`
}
