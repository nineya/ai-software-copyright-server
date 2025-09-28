package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type Notice struct {
	Id         int64           `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	ClientType enum.ClientType `json:"clientType" xorm:"SMALLINT notnull comment('客户端类型')" label:"客户端类型"`
	Content    string          `json:"content" xorm:"VARCHAR(255) comment('通知内容')" binding:"lte=255" label:"通知内容"`
	CreateTime *time.Time      `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time      `json:"updateTime" xorm:"DATETIME updated"`
}

func (Notice) TableName() string {
	return "notice"
}
