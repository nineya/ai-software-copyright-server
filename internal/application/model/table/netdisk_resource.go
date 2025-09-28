package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type NetdiskResource struct {
	Id       int64  `json:"id" xorm:"<- PK AUTOINCR"`                                     //主键
	UserId   int64  `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"` // 用户id允许为空
	UserName string `json:"userName" xorm:"VARCHAR(25) notnull comment('资源所属用户名称')" binding:"lte=25" label:"资源所属用户名称"`
	// CREATE FULLTEXT INDEX FT_netdisk_resource_name ON netdisk_resource(name) WITH PARSER ngram
	Name             string             `json:"name" xorm:"VARCHAR(255) notnull comment('资源名称')" binding:"lte=255" label:"资源名称"`
	TargetUrl        string             `json:"targetUrl" xorm:"VARCHAR(512) notnull comment('目标地址')" binding:"required,lte=512" label:"目标地址"`
	ShareTargetUrl   string             `json:"shareTargetUrl" xorm:"VARCHAR(512) notnull comment('外部地址')" binding:"lte=512" label:"外部地址"` //外部地址
	SharePwd         string             `json:"sharePwd" xorm:"VARCHAR(15) notnull comment('外部分享密码')" binding:"lte=15" label:"外部分享密码"`     // 外部分享的密码
	ShortLink        string             `json:"shortLink" xorm:"VARCHAR(15) notnull comment('短链地址')" binding:"lte=15" label:"短链地址"`
	ShortLinkPwd     string             `json:"shortLinkPwd" xorm:"VARCHAR(15) notnull comment('短链访问密码')" binding:"lte=15" label:"短链访问密码"`
	ShortLinkPwdTips string             `json:"shortLinkPwdTips" xorm:"VARCHAR(128) notnull comment('短链密码提示')" binding:"lte=128" label:"短链密码提示"`
	Type             enum.NetdiskType   `json:"type" xorm:"SMALLINT notnull comment('网盘类型')" label:"网盘类型"`
	Origin           enum.NetdiskOrigin `json:"origin" xorm:"SMALLINT notnull comment('资源来源')" label:"资源来源"`
	Visits           int                `json:"visits" xorm:"INT notnull comment('浏览量')" binding:"gte=0" label:"浏览量"`
	Status           enum.NetdiskStatus `json:"status" xorm:"SMALLINT notnull comment('资源状态')" label:"资源状态"`
	CheckTime        *time.Time         `json:"checkTime" xorm:"DATETIME comment('检查资源的时间')" label:"检查资源的时间"`
	CreateTime       *time.Time         `json:"createTime" xorm:"DATETIME created"`
	UpdateTime       *time.Time         `json:"updateTime" xorm:"DATETIME updated"`
}

func (NetdiskResource) TableName() string {
	return "netdisk_resource"
}

func (t *NetdiskResource) SetUserId(userId int64) {
	t.UserId = userId
}

type NetdiskShortLink struct {
	Id               int64              `json:"id"`               //主键
	UserId           int64              `json:"userId,omitempty"` // 用户id允许为空
	Alias            string             `json:"alias"`
	TargetUrl        string             `json:"targetUrl"`
	Visits           int                `json:"visits"`
	CreateTime       *time.Time         `json:"createTime"`
	UpdateTime       *time.Time         `json:"updateTime"`
	UserName         string             `json:"userName"`
	Name             string             `json:"name"`
	ShortLinkUrl     string             `json:"shortLinkUrl"`
	NetdiskTargetUrl string             `json:"netdiskTargetUrl"`
	Type             enum.NetdiskType   `json:"type"`
	Status           enum.NetdiskStatus `json:"status"`
}

func (NetdiskShortLink) TableName() string {
	return "netdisk_resource"
}

func (t *NetdiskShortLink) SetUserId(userId int64) {
	t.UserId = userId
}
