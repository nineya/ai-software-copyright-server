package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type AdminLog struct {
	Id         int64             `json:"id" xorm:"<- PK AUTOINCR"`
	AdminId    int64             `json:"adminId" xorm:"notnull comment('管理员id')" label:"管理员id"`
	Content    string            `json:"content" xorm:"VARCHAR(1023) notnull comment('日志内容')" binding:"lte=1023" label:"日志内容"`
	IpAddress  string            `json:"ipAddress" xorm:"VARCHAR(127) notnull comment('请求ip')" binding:"lte=127" label:"请求ip"`
	Type       enum.AdminLogType `json:"type" xorm:"SMALLINT notnull comment('日志类型')" label:"日志类型"`
	CreateTime *time.Time        `json:"createTime" xorm:"DATETIME created"`
}

func (AdminLog) TableName() string {
	return "admin_log"
}

func (t *AdminLog) SetAdminId(adminId int64) {
	t.AdminId = adminId
}
