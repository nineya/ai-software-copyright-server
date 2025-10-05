package table

import (
	"time"
)

type CdkeyRecord struct {
	Id         int64      `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	UserId     int64      `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"`
	Cdkey      string     `json:"cdkey,omitempty" xorm:"VARCHAR(24) comment('cdkey')" binding:"lte=24" label:"Cdkey"`
	Credits    int        `json:"credits" xorm:"INT notnull comment('积分数量')" label:"积分数量"`
	CreateTime *time.Time `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time `json:"updateTime" xorm:"DATETIME updated"`
}

func (CdkeyRecord) TableName() string {
	return "cdkey_record"
}

func (t *CdkeyRecord) SetUserId(userId int64) {
	t.UserId = userId
}
