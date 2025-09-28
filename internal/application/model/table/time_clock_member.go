package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type TimeClockMember struct {
	Id         int64                      `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	UserId     int64                      `json:"userId,omitempty" xorm:"unique(user_id_clock_id) notnull comment('用户id')" label:"用户id"`
	ClockId    int64                      `json:"clockId" xorm:"unique(user_id_clock_id) notnull comment('打卡id')" binding:"required" label:"打卡id"` //主键
	Status     enum.TimeClockMemberStatus `json:"status" xorm:"SMALLINT notnull comment('成员状态')" label:"成员状态"`
	JoinTime   *time.Time                 `json:"joinTime" xorm:"DATETIME comment('加入时间')" label:"加入时间"`
	Remark     string                     `json:"remark,omitempty" xorm:"VARCHAR(20) comment('备注')" binding:"lte=20" label:"备注"`
	CreateTime *time.Time                 `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time                 `json:"updateTime" xorm:"DATETIME updated"`
}

func (TimeClockMember) TableName() string {
	return "time_clock_member"
}

func (t *TimeClockMember) SetUserId(userId int64) {
	t.UserId = userId
}
