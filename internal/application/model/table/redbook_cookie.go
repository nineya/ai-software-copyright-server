package table

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"time"
)

type RedbookCookie struct {
	Id         int64             `json:"id" xorm:"<- PK AUTOINCR"` //主键
	AdminId    int64             `json:"adminId,omitempty" xorm:"notnull unique(admin_id_xhs_user_id) comment('管理员id')" label:"管理员id"`
	Nickname   string            `json:"nickname" xorm:"VARCHAR(127) notnull comment('昵称')" binding:"required,lte=127" label:"昵称"` //昵称
	XhsUserId  string            `json:"xhsUserId" xorm:"VARCHAR(50) notnull unique(admin_id_xhs_user_id) comment('小红书用户ID')" binding:"required,lte=50" label:"小红书用户ID"`
	Cookie     string            `json:"cookie" xorm:"LONGTEXT notnull comment('cookie内容')" binding:"required" label:"cookie内容"`
	Status     enum.CookieStatus `json:"status" xorm:"SMALLINT notnull comment('任务状态')" label:"任务状态"`
	CreateTime *time.Time        `json:"createTime" xorm:"DATETIME created"`
	UpdateTime *time.Time        `json:"updateTime" xorm:"DATETIME updated"`
}

func (RedbookCookie) TableName() string {
	return "redbook_cookie"
}

func (t *RedbookCookie) SetAdminId(adminId int64) {
	t.AdminId = adminId
}
