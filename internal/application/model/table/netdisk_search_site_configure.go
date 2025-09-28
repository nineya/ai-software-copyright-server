package table

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type NetdiskSearchSiteConfigure struct {
	Id             int64              `json:"id" xorm:"<- PK AUTOINCR"`                                             //主键
	UserId         int64              `json:"userId,omitempty" xorm:"unique(user_id) comment('用户id')" label:"用户id"` // 用户id允许为空
	ExpireTime     *time.Time         `json:"expireTime" xorm:"<- DATETIME comment('过期时间')" label:"过期时间"`
	Title          string             `json:"title" xorm:"VARCHAR(64) notnull comment('网站标题')" binding:"lte=64" label:"网站标题"`
	Subtitle       string             `json:"subtitle" xorm:"VARCHAR(128) notnull comment('网站副标题')" binding:"lte=128" label:"网站副标题"`
	Favicon        string             `json:"favicon" xorm:"VARCHAR(255) notnull comment('网站小图标')" binding:"lte=255" label:"网站小图标"`
	SeoKeywords    string             `json:"seoKeywords" xorm:"VARCHAR(255) notnull comment('SEO关键词')" binding:"lte=255" label:"SEO关键词"`
	SeoDescription string             `json:"seoDescription" xorm:"VARCHAR(512) notnull comment('SEO描述')" binding:"lte=512" label:"SEO描述"`
	UseShare       bool               `json:"useShare" xorm:"notnull comment('是否使用共享资料')" label:"是否使用共享资料"`
	CollectTypes   []enum.NetdiskType `json:"collectTypes" xorm:"VARCHAR(128) notnull comment('采集网盘类型列表')" form:"collectTypes" label:"采集网盘类型列表"`
	UseShortLink   bool               `json:"useShortLink" xorm:"notnull comment('是否使用短链')" label:"是否使用短链"`
	BrowserTips    bool               `json:"browserTips" xorm:"notnull comment('是否提示浏览器')" label:"是否提示浏览器"`
	Notice         string             `json:"notice" xorm:"VARCHAR(512) notnull comment('网站公告')" binding:"lte=512" label:"网站公告"`
	WechatQrcode   string             `json:"wechatQrcode" xorm:"VARCHAR(255) notnull comment('微信二维码')" binding:"lte=255" label:"微信二维码"`
	Menus          []common.Link      `json:"menus" xorm:"TEXT json notnull comment('顶部菜单')" label:"顶部菜单"`
	Friends        []common.Link      `json:"friends" xorm:"TEXT json notnull comment('友链列表')" label:"友链列表"`
	Beian          string             `json:"beian" xorm:"VARCHAR(64) notnull comment('备案号')" binding:"lte=64" label:"备案号"`
	CustomHead     string             `json:"customHead" xorm:"VARCHAR(2048) notnull comment('自定义Head')" binding:"lte=2048" label:"自定义Head"`
	InlineCss      string             `json:"inlineCss" xorm:"VARCHAR(2048) notnull comment('内嵌CSS')" binding:"lte=2048" label:"内嵌CSS"`
	InlineHtmlBody string             `json:"inlineHtmlBody" xorm:"VARCHAR(4096) notnull comment('内嵌HTML')" binding:"lte=4096" label:"内嵌HTML"`
	CreateTime     *time.Time         `json:"createTime" xorm:"DATETIME created"`
	UpdateTime     *time.Time         `json:"updateTime" xorm:"DATETIME updated"`
}

func (NetdiskSearchSiteConfigure) TableName() string {
	return "netdisk_search_site_configure"
}

func (t *NetdiskSearchSiteConfigure) SetUserId(userId int64) {
	t.UserId = userId
}
