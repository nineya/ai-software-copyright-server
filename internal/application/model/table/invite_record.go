package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

// 邀请奖励记录
type InviteRecord struct {
	Id            int64           `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	UserId        int64           `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"`
	InviteeId     int64           `json:"inviteeId" xorm:"notnull comment('受邀人id')" label:"受邀人id"`
	Type          enum.InviteType `json:"type" xorm:"SMALLINT notnull comment('邀请类型')" label:"邀请类型"`
	RewardCredits int             `json:"rewardCredits" xorm:"INT notnull comment('奖赏积分')" label:"奖赏积分"`
	Remark        string          `json:"remark,omitempty" xorm:"VARCHAR(50) comment('备注')" binding:"lte=50" label:"备注"`
	CreateTime    *time.Time      `json:"createTime" xorm:"DATETIME created"`
	UpdateTime    *time.Time      `json:"updateTime" xorm:"DATETIME updated"`
}

func (InviteRecord) TableName() string {
	return "invite_record"
}

func (t *InviteRecord) SetUserId(userId int64) {
	t.UserId = userId
}

type InviteInfo struct {
	InviteNum     int `json:"inviteNum"`
	InviteCredits int `json:"inviteCredits"`
	ActiveCredits int `json:"activeCredits"`
}
