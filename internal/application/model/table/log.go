package table

import (
	"time"
	"tool-server/internal/application/model/enum"
)

type Log struct {
	Id         int64        `json:"id" xorm:"<- PK AUTOINCR comment('元数据id')"`
	AdminId    int64        `json:"adminId" xorm:"notnull comment('用户id')"`
	Content    string       `json:"content" xorm:"VARCHAR(1023) notnull comment('日志内容')" binding:"lte=1023"`
	IpAddress  string       `json:"ipAddress" xorm:"VARCHAR(127) notnull comment('请求ip')" binding:"lte=127"`
	Type       enum.LogType `json:"type" xorm:"SMALLINT notnull comment('日志类型')"`
	CreateTime *time.Time   `json:"createTime" xorm:"DATETIME created"`
}

func (Log) TableName() string {
	return "log"
}

func (t *Log) SetAdminId(adminId int64) {
	t.AdminId = adminId
}
