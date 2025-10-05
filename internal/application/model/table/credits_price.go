package table

import (
	"time"
)

type CreditsPrice struct {
	Id         int64      `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	Credits    int        `json:"credits" xorm:"INT notnull comment('积分数量')" label:"积分数量"`
	Price      string     `json:"price" xorm:"DECIMAL(11,2) notnull comment('售价')" label:"售价"`
	CreateTime *time.Time `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time `json:"updateTime" xorm:"DATETIME updated"`
}

func (CreditsPrice) TableName() string {
	return "credits_price"
}
