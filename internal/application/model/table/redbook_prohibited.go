package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type RedbookProhibited struct {
	Id         int64               `json:"id" xorm:"<- PK AUTOINCR"` //主键
	UserId     int64               `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"`
	Type       enum.ProhibitedType `json:"type" xorm:"SMALLINT notnull comment('敏感词类型')" label:"敏感词类型"`
	Words      []string            `json:"words" xorm:"TEXT notnull comment('敏感词列表')" label:"敏感词列表"`
	CreateTime *time.Time          `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time          `json:"updateTime" xorm:"DATETIME updated"`
}

func (RedbookProhibited) TableName() string {
	return "redbook_prohibited"
}

func (t *RedbookProhibited) SetUserId(userId int64) {
	t.UserId = userId
}
