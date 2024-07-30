package request

// Admin login structure
type AdminLoginParam struct {
	Username string `json:"username" form:"username" binding:"required,lte=50"`
	Password string `json:"password" form:"password" binding:"required,gte=8,lte=100"`
	//Captcha   string `json:"captcha" form:"captcha" binding:"required"`     // 验证码
	//CaptchaId string `json:"captchaId" form:"captchaId" binding:"required"` // 验证码ID
}
