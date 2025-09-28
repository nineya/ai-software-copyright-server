package table

import (
	"time"
)

type Cdkey struct {
	Id                int64      `json:"id,omitempty" xorm:"<- PK AUTOINCR"` //主键
	AdminId           int64      `json:"adminId,omitempty" xorm:"notnull comment('管理员id')" label:"管理员id"`
	Cdkey             string     `json:"cdkey,omitempty" xorm:"VARCHAR(24) unique(cdkey) comment('cdkey')" binding:"lte=24" label:"Cdkey"`
	CreditsNum        int        `json:"creditsNum" xorm:"INT notnull comment('币数量')" label:"币数量"`
	HelperStandardDay int        `json:"helperStandardDay" xorm:"INT notnull comment('网盘助手标准版赠送天数')" label:"网盘助手标准版赠送天数"`
	HelperWechatDay   int        `json:"helperWechatDay" xorm:"INT notnull comment('网盘助手微信版赠送天数')" label:"网盘助手微信版赠送天数"`
	TotalCount        int        `json:"totalCount" xorm:"INT notnull comment('总兑换次数')" label:"总兑换次数"`
	SurplusCount      int        `json:"surplusCount" xorm:"INT notnull comment('剩余兑换次数')" label:"剩余兑换次数"`
	Remark            string     `json:"remark,omitempty" xorm:"VARCHAR(100) comment('备注')" binding:"lte=100" label:"备注"`
	ExpireTime        *time.Time `json:"expireTime" xorm:"DATETIME comment('失效时间')" label:"失效时间"`
	CreateTime        *time.Time `json:"createTime" xorm:"DATETIME created"`
	UpdateTime        *time.Time `json:"updateTime" xorm:"DATETIME updated"`
}

func (Cdkey) TableName() string {
	return "cdkey"
}

func (t *Cdkey) SetAdminId(adminId int64) {
	t.AdminId = adminId
}
