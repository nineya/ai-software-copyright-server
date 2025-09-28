package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type NetdiskResourceSearch struct {
	Id            int64           `json:"id" xorm:"<- PK AUTOINCR"`                                     //主键
	UserId        int64           `json:"userId,omitempty" xorm:"notnull comment('用户id')" label:"用户id"` // 用户id允许为空
	Keyword       string          `json:"keyword" xorm:"VARCHAR(255) notnull comment('搜索关键字')" binding:"required,lte=255" label:"搜索关键字"`
	ResourceCount int             `json:"resourceCount" xorm:"notnull comment('资源数量')" binding:"gte=0" label:"资源数量"`
	Origin        enum.ClientType `json:"origin" xorm:"SMALLINT notnull comment('搜索客户端来源')" label:"搜索客户端来源"`
	CreateTime    *time.Time      `json:"createTime" xorm:"DATETIME created"`
}

func (NetdiskResourceSearch) TableName() string {
	return "netdisk_resource_search"
}

func (t *NetdiskResourceSearch) SetUserId(userId int64) {
	t.UserId = userId
}
