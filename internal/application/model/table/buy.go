package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type Buy struct {
	Id         int64        `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	UserId     int64        `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"`
	Type       enum.BuyType `json:"type" xorm:"SMALLINT notnull comment('类型')" label:"类型"`
	PayCredits int          `json:"payCredits" xorm:"INT notnull comment('支付积分数量')" label:"支付积分数量"`
	Remark     string       `json:"remark,omitempty" xorm:"VARCHAR(100) comment('备注')" binding:"lte=100" label:"备注"`
	CreateTime *time.Time   `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time   `json:"updateTime" xorm:"DATETIME updated"`
}

func (Buy) TableName() string {
	return "buy"
}

func (t *Buy) SetUserId(userId int64) {
	t.UserId = userId
}
