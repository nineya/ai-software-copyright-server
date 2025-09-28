package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type ClientInfo struct {
	Id         int64           `json:"id" xorm:"<- PK AUTOINCR"`                                                  //主键
	UserId     int64           `json:"userId,omitempty" xorm:"unique(user_id_type) comment('用户id')" label:"用户id"` // 用户id允许为空
	Type       enum.ClientType `json:"type" xorm:"SMALLINT notnull unique(user_id_type) comment('客户端类型')" label:"客户端类型"`
	WxOpenid   string          `json:"wxOpenid" xorm:"VARCHAR(127) notnull unique(wx_openid) comment('微信OpenId')" binding:"lte=127" label:"微信OpenId"` //微信用户id
	QrCodeUrl  string          `json:"qrCodeUrl" xorm:"VARCHAR(512) notnull comment('邀请码地址')" binding:"required,lte=512" label:"邀请码地址"`
	CreateTime *time.Time      `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time      `json:"updateTime" xorm:"DATETIME updated"`
}

func (*ClientInfo) TableName() string {
	return "client_info"
}

func (t *ClientInfo) SetUserId(userId int64) {
	t.UserId = userId
}
