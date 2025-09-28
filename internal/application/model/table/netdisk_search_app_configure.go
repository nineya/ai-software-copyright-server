package table

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type NetdiskSearchAppConfigure struct {
	Id                 int64                 `json:"id" xorm:"<- PK AUTOINCR"`                                             //主键
	UserId             int64                 `json:"userId,omitempty" xorm:"unique(user_id) comment('用户id')" label:"用户id"` // 用户id允许为空
	ExpireTime         *time.Time            `json:"expireTime" xorm:"<- DATETIME comment('过期时间')" label:"过期时间"`
	UseShare           bool                  `json:"useShare" xorm:"notnull comment('是否使用共享资料')" label:"是否使用共享资料"`
	CollectTypes       []enum.NetdiskType    `json:"collectTypes" xorm:"VARCHAR(128) notnull comment('采集网盘类型列表')" form:"collectTypes" label:"采集网盘类型列表"`
	ExtendTitle        string                `json:"extendTitle" xorm:"VARCHAR(50) notnull comment('推广引流标题')" binding:"lte=50" label:"推广引流标题"`
	ExtendImageUrl     string                `json:"extendImageUrl" xorm:"VARCHAR(255) notnull comment('推广引流图片url')" binding:"lte=255" label:"推广引流图片url"`
	BannerAdId         string                `json:"bannerAdId" xorm:"VARCHAR(50) notnull comment('小程序横幅广告位ID')" binding:"lte=50" label:"横幅广告位ID"`
	RewardAdId         string                `json:"rewardAdId" xorm:"VARCHAR(50) notnull comment('小程序激励广告位ID')" binding:"lte=50" label:"激励广告位ID"`
	ItemAdId           string                `json:"itemAdId" xorm:"VARCHAR(50) notnull comment('小程序列表项广告位ID')" binding:"lte=50" label:"列表项广告位ID"`
	RewardScore        int                   `json:"rewardScore" xorm:"INT notnull comment('奖赏积分')" label:"奖赏积分"`
	SearchScore        int                   `json:"searchScore" xorm:"INT notnull comment('搜索所需积分')" label:"搜索所需积分"`
	CopyScore          int                   `json:"copyScore" xorm:"INT notnull comment('复制所需积分')" label:"复制所需积分"`
	WelfareConfig      []common.WelfareLabel `json:"welfareConfig" xorm:"TEXT notnull comment('小程序羊毛配置')" label:"羊毛配置"`
	WelfareMoreWxAppid string                `json:"welfareMoreWxAppid" xorm:"VARCHAR(50) notnull comment('更多羊毛跳转小程序Appid')" binding:"lte=50" label:"更多羊毛跳转小程序Appid"`
	WelfareMoreWxPath  string                `json:"welfareMoreWxPath" xorm:"VARCHAR(256) notnull comment('更多羊毛跳转小程序Path')" binding:"lte=256" label:"更多羊毛跳转小程序Path"`
	CreateTime         *time.Time            `json:"createTime" xorm:"DATETIME created"`
	UpdateTime         *time.Time            `json:"updateTime" xorm:"DATETIME updated"`
}

func (NetdiskSearchAppConfigure) TableName() string {
	return "netdisk_search_app_configure"
}

func (t *NetdiskSearchAppConfigure) SetUserId(userId int64) {
	t.UserId = userId
}
