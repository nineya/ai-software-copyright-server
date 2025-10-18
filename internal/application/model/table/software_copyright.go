package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type SoftwareCopyright struct {
	Id             int64                        `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	UserId         int64                        `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"`
	Name           string                       `json:"name" xorm:"VARCHAR(50) notnull comment('软件全称')" binding:"required,lte=50" label:"软件全称"`
	ShortName      string                       `json:"shortName" xorm:"VARCHAR(50) notnull comment('软件简称')" binding:"lte=50" label:"软件简称"`
	Version        string                       `json:"version" xorm:"VARCHAR(10) notnull comment('版本号')" binding:"lte=10" label:"版本号"`
	Category       string                       `json:"category" xorm:"VARCHAR(20) notnull comment('软件分类')" binding:"lte=20" label:"软件分类"`
	CodeLang       string                       `json:"codeLang" xorm:"VARCHAR(50) notnull comment('开发语言')" binding:"required,lte=50" label:"开发语言"`
	Description    string                       `json:"description" xorm:"VARCHAR(500) notnull comment('软件功能描述')" binding:"required,lte=500" label:"软件功能描述"`
	Owner          string                       `json:"owner" xorm:"VARCHAR(50) notnull comment('著作权人')" binding:"required,lte=50" label:"著作权人"`
	Progress       int                          `json:"progress" xorm:"notnull comment('生成进度')" label:"生成进度"`
	Status         enum.SoftwareCopyrightStatus `json:"status" xorm:"SMALLINT notnull comment('状态')" label:"状态"`
	ApiKey         string                       `json:"-" xorm:"VARCHAR(128) notnull comment('应用ApiKey')" binding:"lte=128" label:"应用ApiKey"`
	ConversationId string                       `json:"-" xorm:"VARCHAR(50) notnull comment('会话id')" binding:"lte=50" label:"会话id"`
	CreateTime     *time.Time                   `json:"createTime" xorm:"DATETIME created"`
	UpdateTime     *time.Time                   `json:"updateTime" xorm:"DATETIME updated"`
}

func (SoftwareCopyright) TableName() string {
	return "software_copyright"
}

func (t *SoftwareCopyright) SetUserId(userId int64) {
	t.UserId = userId
}

type SoftwareCopyrightStatistic struct {
	TotalCount    int    `json:"totalCount"`    // 总数量
	GenerateCount int    `json:"generateCount"` // 生成中数量
	CompleteCount string `json:"completeCount"` // 已完成数量
}

func (SoftwareCopyrightStatistic) TableName() string {
	return "software_copyright"
}
