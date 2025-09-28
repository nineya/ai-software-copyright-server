package table

import "time"

type Qrcode struct {
	Id         int64      `json:"id" xorm:"<- PK AUTOINCR"`                                     //主键
	UserId     int64      `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"` // 用户id允许为空
	Alias      string     `json:"alias" xorm:"VARCHAR(24) notnull unique(alias) comment('别名')" binding:"lte=24" label:"别名"`
	Title      string     `json:"title" xorm:"VARCHAR(50) notnull comment('标题')" binding:"required,lte=50" label:"标题"`
	TargetUrls []string   `json:"targetUrls" xorm:"TEXT notnull comment('目标地址列表')" label:"目标地址列表"`
	Visits     int        `json:"visits" xorm:"INT notnull comment('浏览量')" binding:"gte=0" label:"浏览量"`
	CreateTime *time.Time `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time `json:"updateTime" xorm:"DATETIME updated"`
}

func (Qrcode) TableName() string {
	return "qrcode"
}

func (t *Qrcode) SetUserId(userId int64) {
	t.UserId = userId
}
