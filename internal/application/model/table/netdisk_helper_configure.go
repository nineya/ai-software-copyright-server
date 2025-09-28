package table

import (
	"ai-software-copyright-server/internal/application/model/common"
	"time"
)

type NetdiskHelperConfigure struct {
	Id               int64                            `json:"id" xorm:"<- PK AUTOINCR"`                                             //主键
	UserId           int64                            `json:"userId,omitempty" xorm:"unique(user_id) comment('用户id')" label:"用户id"` // 用户id允许为空
	ExpireTime       *time.Time                       `json:"expireTime" xorm:"<- DATETIME comment('基础版过期时间')" label:"过期时间"`
	WechatExpireTime *time.Time                       `json:"wechatExpireTime" xorm:"<- DATETIME comment('微信工具人过期时间')" label:"微信工具人过期时间"`
	Quark            common.NetdiskHelperConfigQuark  `json:"quark" xorm:"LONGTEXT json notnull comment('夸克网盘配置')" label:"夸克网盘配置"`
	Baidu            common.NetdiskHelperConfigBaidu  `json:"baidu" xorm:"LONGTEXT json notnull comment('百度网盘配置')" label:"百度网盘配置"`
	Ai               common.NetdiskHelperConfigAi     `json:"ai" xorm:"TEXT json notnull comment('AI配置')" label:"AI配置"`
	Mail             common.NetdiskHelperConfigMail   `json:"mail" xorm:"TEXT json notnull comment('邮箱配置')" label:"邮箱配置"`
	Wechat           common.NetdiskHelperConfigWechat `json:"wechat" xorm:"TEXT json notnull comment('微信工具人配置')" label:"微信工具人配置"`
	CreateTime       *time.Time                       `json:"createTime" xorm:"DATETIME created"`
	UpdateTime       *time.Time                       `json:"updateTime" xorm:"DATETIME updated"`
}

func (NetdiskHelperConfigure) TableName() string {
	return "netdisk_helper_configure"
}

func (t *NetdiskHelperConfigure) SetUserId(userId int64) {
	t.UserId = userId
}
