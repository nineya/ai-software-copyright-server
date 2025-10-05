package public

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type AuthApiRouter struct {
	api.BaseApi
}

var store = utils.NewCaptchaStore()

func (m *AuthApiRouter) InitAuthApiRouter(router *gin.RouterGroup) {
	router.POST("login", m.Login)
	router.POST("refresh_token", m.RefreshToken)
	router.POST("captcha", m.Captcha)
	router.POST("register", m.Register)
}

// @summary User login
// @description User login
// @tags public,auth
// @accept json
// @param param body request.AdminLoginParam true "Login user name and password information"
// @success 200 {object} response.Response{data=response.TokenResponse}
// @router /public/login [post]
func (m *AuthApiRouter) Login(c *gin.Context) {
	var param request.UserLoginParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	//if !store.Verify(param.CaptchaId, param.Captcha, true) {
	//	response.FailWithMessage("验证码错误", c)
	//	return
	//}
	token, err := userSev.GetAuthService().Login(param)
	if err != nil {
		m.UserLog(c, "FAILED_LOGIN", fmt.Sprintf("试图登录 %s 账号失败，原因：%s", param.Phone, err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "USER_LOGIN", fmt.Sprintf("账号 %s 登录成功", param.Phone))
	response.OkWithData(token, c)
}

// @summary Refreshing an Administrator token
// @description Refreshing an Administrator token
// @tags public,auth
// @accept json
// @param param body request.RefreshTokenParam true "Refresh token"
// @success 200 {object} response.Response{data=response.TokenResponse}
// @router /public/refresh_token [post]
func (m *AuthApiRouter) RefreshToken(c *gin.Context) {
	var param request.RefreshTokenParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	token, err := utils.RefreshToken(param)
	if err != nil {
		m.UserLog(c, "FAILED_LOGIN", fmt.Sprintf("刷新登录Token失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "USER_LOGIN", fmt.Sprintf("刷新登录Token成功，AccessToken：%s", token.AccessToken))
	response.OkWithData(token, c)
}

// @summary Generate user login captcha
// @description Generate user login captcha
// @tags public,auth
// @accept json
// @success 200 {object} response.Response{data=response.CaptchaResponse}
// @router /public/captcha [post]
func (m *AuthApiRouter) Captcha(c *gin.Context) {
	driver := base64Captcha.NewDriverDigit(140, 400, 6, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	if id, b64s, _, err := cp.Generate(); err != nil {
		global.LOG.Error("验证码获取失败", zap.Any("err", err))
		response.FailWithMessage("验证码获取失败", c)
	} else {
		response.OkWithDetailed(response.CaptchaResponse{
			CaptchaId: id,
			Path:      b64s,
		}, "验证码获取成功", c)
	}
}

// @summary User register
// @description User register
// @tags public,user
// @accept json
// @param param body request.UserInfoParam true "User information"
// @success 200 {object} response.Response{data=response.UserLoginResponse}
// @router /public/register [post]
func (m *AuthApiRouter) Register(c *gin.Context) {
	var param request.UserInfoParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = userSev.GetAuthService().Register(param)
	if err != nil {
		m.UserLog(c, "FAILED_LOGIN", fmt.Sprintf("用户注册失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "USER_LOGIN", fmt.Sprintf("用户 %s 注册", param.Nickname))
	response.Ok(c)
}
