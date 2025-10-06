package table

import "time"

type User struct {
	Id         int64      `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	Avatar     string     `json:"avatar,omitempty" xorm:"VARCHAR(511) notnull comment('头像地址')" binding:"required,lte=511" label:"头像地址"`
	Nickname   string     `json:"nickname,omitempty" xorm:"VARCHAR(127) notnull comment('昵称')" binding:"required,lte=127" label:"昵称"`
	Phone      *string    `json:"phone,omitempty" xorm:"VARCHAR(15) unique(phone) comment('用户手机号')" label:"用户手机号"`
	Email      *string    `json:"email,omitempty" xorm:"VARCHAR(127) unique(email) comment('邮箱')" label:"邮箱"`
	Password   string     `json:"-" xorm:"VARCHAR(255) notnull comment('密码')" binding:"lte=100" label:"密码"`
	WxUnionid  string     `json:"wxUnionid" xorm:"VARCHAR(127) notnull unique(wx_unionid) comment('微信Unionid')" binding:"required,lte=127" label:"微信Unionid"` //微信用户id
	WxOpenid   string     `json:"wxOpenid" xorm:"VARCHAR(127) notnull unique(wx_openid) comment('微信Openid')" binding:"lte=127" label:"微信Openid"`              //微信用户id
	InviteCode string     `json:"inviteCode,omitempty" xorm:"VARCHAR(10) notnull unique(invite_code) comment('邀请码')" binding:"lte=10" label:"邀请码"`
	Inviter    string     `json:"inviter,omitempty" xorm:"VARCHAR(10) comment('邀请人')" binding:"lte=10" label:"邀请人"`
	Credits    int        `json:"credits" xorm:"INT notnull comment('积分')" binding:"gte=0" label:"积分"`
	ActiveTime time.Time  `json:"activeTime" xorm:"DATETIME notnull comment('活跃时间')" label:"活跃时间"`
	CreateTime *time.Time `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time `json:"updateTime" xorm:"DATETIME updated"`
}

func (User) TableName() string {
	return "user"
}
