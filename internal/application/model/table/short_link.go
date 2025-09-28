package table

import "time"

type ShortLink struct {
	Id         int64      `json:"id" xorm:"<- PK AUTOINCR"`                                     //主键
	UserId     int64      `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"` // 用户id允许为空
	Alias      string     `json:"alias" xorm:"VARCHAR(24) notnull unique(alias) comment('别名')" binding:"lte=24" label:"别名"`
	TargetUrl  string     `json:"targetUrl" xorm:"VARCHAR(512) notnull comment('目标地址')" binding:"required,lte=512" label:"目标地址"`
	Visits     int        `json:"visits" xorm:"INT notnull comment('浏览量')" binding:"gte=0" label:"浏览量"`
	CreateTime *time.Time `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time `json:"updateTime" xorm:"DATETIME updated"`
}

func (ShortLink) TableName() string {
	return "short_link"
}

func (t *ShortLink) SetUserId(userId int64) {
	t.UserId = userId
}

type ShortLinkStatistic struct {
	Date     string `json:"date,omitempty"`
	Origin   string `json:"origin,omitempty"`
	Pv       int    `json:"pv"`
	Uv       int    `json:"uv"`
	Mobile   int    `json:"mobile"`
	MobileUv int    `json:"mobileUv"`
	PcUv     int    `json:"pcUv"`
}

func (ShortLinkStatistic) TableName() string {
	return "statistic"
}

type ShortLinkActive struct {
	Total       int `json:"total"`
	TodayActive int `json:"todayActive"`
	TotalActive int `json:"totalActive"`
}

func (ShortLinkActive) TableName() string {
	return "short_link"
}
