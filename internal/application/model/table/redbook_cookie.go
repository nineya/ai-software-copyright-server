package table

import "time"

type RedbookCookie struct {
	Id         int64      `json:"id" xorm:"<- PK AUTOINCR"` //主键
	AdminId    int64      `json:"adminId,omitempty" xorm:"-> notnull comment('管理员id')"`
	Cookie     string     `json:"cookie" xorm:"LONGTEXT notnull comment('cookie内容')" binding:"required"`
	CreateTime *time.Time `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time `json:"updateTime" xorm:"DATETIME updated"`
}

func (RedbookCookie) TableName() string {
	return "RedbookCookie"
}
