package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type UserLog struct {
	Id         int64            `json:"id" xorm:"<- PK AUTOINCR"`
	UserId     int64            `json:"userId" xorm:"notnull comment('用户id')" label:"用户id"`
	Content    string           `json:"content" xorm:"VARCHAR(1023) notnull comment('日志内容')" binding:"lte=1023" label:"日志内容"`
	IpAddress  string           `json:"ipAddress" xorm:"VARCHAR(127) notnull comment('请求ip')" binding:"lte=127" label:"请求ip"`
	Type       enum.UserLogType `json:"type" xorm:"SMALLINT notnull comment('日志类型')" label:"日志类型"`
	CreateTime *time.Time       `json:"createTime" xorm:"DATETIME created"`
}

func (UserLog) TableName() string {
	return "user_log"
}

func (t *UserLog) SetUserId(userId int64) {
	t.UserId = userId
}
