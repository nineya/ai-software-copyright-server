package table

import (
	"time"
)

type CdkeyRecord struct {
	Id                int64      `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	UserId            int64      `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"`
	Cdkey             string     `json:"cdkey,omitempty" xorm:"VARCHAR(24) comment('cdkey')" binding:"lte=24" label:"Cdkey"`
	CreditsNum        int        `json:"creditsNum" xorm:"INT notnull comment('币数量')" label:"币数量"`
	HelperStandardDay int        `json:"helperStandardDay" xorm:"INT notnull comment('网盘助手标准版赠送天数')" label:"网盘助手标准版赠送天数"`
	HelperWechatDay   int        `json:"helperWechatDay" xorm:"INT notnull comment('网盘助手微信版赠送天数')" label:"网盘助手微信版赠送天数"`
	CreateTime        *time.Time `json:"createTime" xorm:"DATETIME created"`
	UpdateTime        *time.Time `json:"updateTime" xorm:"DATETIME updated"`
}

func (CdkeyRecord) TableName() string {
	return "cdkey_record"
}

func (t *CdkeyRecord) SetUserId(userId int64) {
	t.UserId = userId
}
