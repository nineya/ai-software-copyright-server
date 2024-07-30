package table

import "time"

type Admin struct {
	Id         int64      `json:"id" xorm:"<- PK AUTOINCR"`                                                                      //主键
	Email      string     `json:"email" xorm:"VARCHAR(127) notnull unique(email) comment('邮箱')" binding:"email,lte=127"`         //邮箱
	Nickname   string     `json:"nickname" xorm:"VARCHAR(127) notnull comment('昵称')" binding:"required,lte=127"`                 //昵称
	Username   string     `json:"username" xorm:"VARCHAR(50) notnull unique(username) comment('用户名')" binding:"required,lte=50"` //用户名
	Password   string     `json:"password,omitempty" xorm:"VARCHAR(255) notnull comment('密码')" binding:"required,gte=8,lte=100"`
	CreateTime *time.Time `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time `json:"updateTime" xorm:"DATETIME updated"`
}

func (Admin) TableName() string {
	return "admin"
}
