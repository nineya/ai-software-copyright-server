package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type CreditsChange struct {
	Id             int64                  `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	UserId         int64                  `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"`
	Type           enum.CreditsChangeType `json:"type" xorm:"SMALLINT notnull comment('类型')" label:"类型"`
	OriginCredits  int                    `json:"originCredits" xorm:"INT notnull comment('原积分数量')" label:"原积分数量"`
	ChangeCredits  int                    `json:"changeCredits" xorm:"INT notnull comment('变动积分数量')" label:"变动积分数量"`
	BalanceCredits int                    `json:"balanceCredits" xorm:"INT notnull comment('积分余额数量')" label:"积分余额数量"`
	Remark         string                 `json:"remark,omitempty" xorm:"VARCHAR(100) comment('备注')" binding:"lte=100" label:"备注"`
	CreateTime     *time.Time             `json:"createTime" xorm:"DATETIME created"`
	UpdateTime     *time.Time             `json:"updateTime" xorm:"DATETIME updated"`
}

func (CreditsChange) TableName() string {
	return "credits_change"
}

func (t *CreditsChange) SetUserId(userId int64) {
	t.UserId = userId
}
