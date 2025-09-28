package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type StudyResource struct {
	Id         int64          `json:"id" xorm:"<- PK AUTOINCR"` //主键
	Name       string         `json:"name" xorm:"VARCHAR(255) notnull comment('资源名称')" binding:"required,lte=255" label:"资源名称"`
	TargetUrl  string         `json:"targetUrl" xorm:"VARCHAR(512) notnull comment('目标地址')" binding:"required,lte=512" label:"目标地址"`
	Type       enum.StudyType `json:"type" xorm:"SMALLINT notnull comment('学习类目')" label:"学习类目"`
	CreateTime *time.Time     `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time     `json:"updateTime" xorm:"DATETIME updated"`
}

func (StudyResource) TableName() string {
	return "study_resource"
}
