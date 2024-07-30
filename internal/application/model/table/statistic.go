package table

import "time"

type Statistic struct {
	Id         int64      `json:"id" xorm:"<- PK AUTOINCR comment('站点id')"`
	Url        string     `json:"url" xorm:"VARCHAR(511) notnull comment('请求地址')" binding:"required,lte=511"`
	IpAddress  string     `json:"ipAddress" xorm:"VARCHAR(127) notnull comment('请求ip')" binding:"lte=127"`
	HttpStatus int        `json:"httpStatus" xorm:"INT notnull comment('请求响应码')"`
	UserAgent  string     `json:"userAgent,omitempty" xorm:"VARCHAR(255) default('') comment('用户标识')" binding:"lte=255"`
	CreateTime *time.Time `json:"createTime" xorm:"DATETIME created"`
}

func (Statistic) TableName() string {
	return "statistic"
}
