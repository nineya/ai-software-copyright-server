package table

import (
	"time"
)

type TimeClock struct {
	Id          int64      `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	UserId      int64      `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"`
	Name        string     `json:"name" xorm:"VARCHAR(50) notnull comment('打卡名称')" binding:"required,lte=50" label:"打卡名称"`
	Description string     `json:"description" xorm:"VARCHAR(500) comment('打卡描述')" binding:"lte=500" label:"打卡描述"`
	StartTime   time.Time  `json:"startTime" xorm:"DATETIME notnull comment('打卡开始时间')" label:"打卡开始时间"`
	EndTime     time.Time  `json:"endTime" xorm:"DATETIME notnull comment('打卡结束时间')" label:"打卡结束时间"`
	CreateTime  *time.Time `json:"createTime" xorm:"DATETIME created"`
	UpdateTime  *time.Time `json:"updateTime" xorm:"DATETIME updated"`
}

func (TimeClock) TableName() string {
	return "time_clock"
}

func (t *TimeClock) SetUserId(userId int64) {
	t.UserId = userId
}
