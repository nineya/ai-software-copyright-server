package table

import (
	"time"
)

type TimeClockRecord struct {
	Id         int64      `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	UserId     int64      `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"`
	ClockId    int64      `json:"clockId" xorm:"notnull comment('打卡id')" binding:"required" label:"打卡id"` //主键
	CreateTime *time.Time `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time `json:"updateTime" xorm:"DATETIME updated"`
}

func (TimeClockRecord) TableName() string {
	return "time_clock_record"
}

func (t *TimeClockRecord) SetUserId(userId int64) {
	t.UserId = userId
}
