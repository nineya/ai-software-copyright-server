package request

import "ai-software-copyright-server/internal/application/model/enum"

type UserLoginParam struct {
	Phone     string `json:"phone" form:"phone" binding:"required,lte=15" label:"手机"`
	Password  string `json:"password" form:"password" binding:"required,gte=6,lte=100" label:"密码"`
	Captcha   string `json:"captcha" form:"captcha" binding:"required" label:"验证码"`       // 验证码
	CaptchaId string `json:"captchaId" form:"captchaId" binding:"required" label:"验证码ID"` // 验证码ID
}

type UserAddNyCreditsParam struct {
	InviteCodes []string               `json:"inviteCodes" form:"inviteCodes" label:"邀请码列表"`
	Type        enum.CreditsChangeType `json:"type" form:"type" binding:"required" label:"金额变动类型"`
	AddCredits  int                    `json:"addCredits" form:"addCredits" binding:"required" label:"加币数量"`
	Remark      string                 `json:"remark,omitempty" form:"remark" binding:"lte=100" label:"备注"`
}

type UserRewardNyCreditsParam struct {
	RewardCredits int    `json:"rewardCredits" form:"rewardCredits" binding:"required,gt=0,lte=300" label:"激励金额"`
	Remark        string `json:"remark,omitempty" form:"remark" binding:"lte=100" label:"备注"`
}

type UserRewardGoodsParam struct {
	Name string `json:"name,omitempty" form:"name" binding:"lte=50" label:"物品名称"`
}

type UserInviterCreditsParam struct {
	Inviter       string          `json:"inviter" form:"inviter" binding:"lte=10" label:"邀请人"`
	Type          enum.InviteType `json:"type" form:"type" label:"激励类型"`
	RewardCredits int             `json:"rewardCredits" form:"rewardCredits" binding:"required,gt=0" label:"激励金额"`
	Remark        string          `json:"remark,omitempty" form:"remark" binding:"lte=100" label:"备注"`
}

type UserUpdateInfoParam struct {
	Phone    string `json:"phone" form:"phone" label:"手机"` // 不能加限制，可能为空
	Email    string `json:"email" form:"email" binding:"lte=127" label:"邮箱"`
	Password string `json:"password" form:"password" label:"密码"` // 不能加限制，可能为空
}
