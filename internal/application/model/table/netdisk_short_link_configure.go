package table

import (
	"time"
)

type NetdiskShortLinkConfigure struct {
	Id               int64      `json:"id" xorm:"<- PK AUTOINCR"`                                             //主键
	UserId           int64      `json:"userId,omitempty" xorm:"unique(user_id) comment('用户id')" label:"用户id"` // 用户id允许为空
	MoreResource     string     `json:"moreResource" xorm:"VARCHAR(255) notnull comment('短链更多资源')" binding:"lte=255" label:"短链更多资源"`
	FailureText      string     `json:"failureText" xorm:"VARCHAR(255) notnull comment('短链失效说明')" binding:"lte=255" label:"短链失效说明"`
	Tips             string     `json:"tips" xorm:"VARCHAR(255) notnull comment('短链提示')" binding:"lte=255" label:"短链提示"`
	CustomExpireTime *time.Time `json:"customExpireTime" xorm:"<- DATETIME comment('定制版过期时间')" label:"定制版过期时间"`
	CustomHost       string     `json:"customHost" xorm:"VARCHAR(64) notnull comment('短链自定义主机名')" binding:"lte=255" label:"短链自定义主机名"`
	CustomFavicon    string     `json:"customFavicon" xorm:"VARCHAR(255) notnull comment('短链自定义小图标')" binding:"lte=255" label:"短链自定义小图标"`
	CustomHead       string     `json:"customHead" xorm:"VARCHAR(2048) notnull comment('自定义Head')" binding:"lte=2048" label:"自定义Head"`
	InlineCss        string     `json:"inlineCss" xorm:"VARCHAR(2048) notnull comment('内嵌CSS')" binding:"lte=2048" label:"内嵌CSS"`
	InlineHtmlBody   string     `json:"inlineHtmlBody" xorm:"VARCHAR(4096) notnull comment('内嵌HTML')" binding:"lte=4096" label:"内嵌HTML"`
	CreateTime       *time.Time `json:"createTime" xorm:"DATETIME created"`
	UpdateTime       *time.Time `json:"updateTime" xorm:"DATETIME updated"`
}

func (NetdiskShortLinkConfigure) TableName() string {
	return "netdisk_short_link_configure"
}

func (t *NetdiskShortLinkConfigure) SetUserId(userId int64) {
	t.UserId = userId
}
