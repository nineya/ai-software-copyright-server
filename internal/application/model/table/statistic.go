package table

import (
	"time"
)

type Statistic struct {
	Id         int64      `json:"id" xorm:"<- PK AUTOINCR"`
	Url        string     `json:"url" xorm:"VARCHAR(511) index notnull comment('请求地址')" binding:"required,lte=511" label:"请求地址"`
	IpAddress  string     `json:"ipAddress" xorm:"VARCHAR(127) notnull comment('请求ip')" binding:"lte=127" label:"请求ip"`
	Referrer   string     `json:"referrer" xorm:"VARCHAR(512) notnull comment('请求来源')" binding:"lte=512" label:"请求来源"`
	Origin     string     `json:"origin" xorm:"VARCHAR(127) notnull comment('请求来源站')" binding:"lte=127" label:"请求来源站"`
	HttpStatus int        `json:"httpStatus" xorm:"INT notnull comment('请求响应码')" label:"请求响应码"`
	UserAgent  string     `json:"userAgent,omitempty" xorm:"VARCHAR(1024) default('') comment('用户标识')" binding:"lte=1024" label:"用户标识"`
	CreateTime *time.Time `json:"createTime" xorm:"DATETIME index created"`
}

func (Statistic) TableName() string {
	return "statistic"
}
