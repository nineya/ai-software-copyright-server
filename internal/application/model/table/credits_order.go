package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

// 积分订单
type CreditsOrder struct {
	Id          int64            `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	UserId      int64            `json:"userId" xorm:"notnull comment('用户id')" label:"用户id"`
	TradeNo     string           `json:"tradeNo" xorm:"VARCHAR(32) notnull unique(trade_no) comment('订单编号')" label:"订单编号"`
	ClientType  enum.ClientType  `json:"clientType" xorm:"SMALLINT notnull comment('客户端类型')" label:"客户端类型"`
	WxOpenid    string           `json:"wxOpenid" xorm:"VARCHAR(127) notnull comment('微信OpenId')" binding:"lte=127" label:"微信OpenId"` //微信用户id
	Description string           `json:"description" xorm:"VARCHAR(127) notnull comment('订单描述')" label:"订单描述"`
	Credits     int              `json:"credits" xorm:"INT notnull comment('积分数量')" label:"积分数量"`
	OrderAmount string           `json:"orderAmount" xorm:"DECIMAL(11,2) notnull comment('订单总金额')" label:"订单总金额"`
	Status      enum.OrderStatus `json:"status" xorm:"SMALLINT notnull comment('订单状态')" label:"订单状态"`
	CreateTime  *time.Time       `json:"createTime" xorm:"DATETIME created"`
	UpdateTime  *time.Time       `json:"updateTime" xorm:"DATETIME updated"`
}

func (CreditsOrder) TableName() string {
	return "credits_order"
}

func (t *CreditsOrder) SetUserId(userId int64) {
	t.UserId = userId
}
