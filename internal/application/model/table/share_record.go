package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

// 用户分享记录
type ShareRecord struct {
	Id            int64            `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	UserId        int64            `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"`
	ShareUrl      string           `json:"shareUrl" xorm:"VARCHAR(500) notnull comment('分享链接')" label:"分享链接"`
	RewardCredits int              `json:"rewardCredits" xorm:"INT notnull comment('奖赏积分')" label:"奖赏积分"`
	Status        enum.ShareStatus `json:"status" xorm:"SMALLINT notnull comment('分享状态')" label:"分享状态"`
	Remark        string           `json:"remark,omitempty" xorm:"VARCHAR(50) comment('备注')" binding:"lte=50" label:"备注"`
	CreateTime    *time.Time       `json:"createTime" xorm:"DATETIME created"`
	UpdateTime    *time.Time       `json:"updateTime" xorm:"DATETIME updated"`
}

func (ShareRecord) TableName() string {
	return "share_record"
}

func (t *ShareRecord) SetUserId(userId int64) {
	t.UserId = userId
}

type ShareStatistic struct {
	ShareCredits int `json:"shareCredits"` //分享积分
	AwaitCount   int `json:"awaitCount"`   // 等待人数
	PassCount    int `json:"passCount"`    //通过人数
}

func (ShareStatistic) TableName() string {
	return "share_record"
}
