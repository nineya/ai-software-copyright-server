package public

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"tool-server/internal/application/param/request"
	"tool-server/internal/application/param/response"
	"tool-server/internal/application/router/api"
	adminSev "tool-server/internal/application/service/admin"
	"tool-server/internal/global"
	"tool-server/internal/utils"
)

type AuthApiRouter struct {
	api.BaseApi
}

var store = utils.NewCaptchaStore()

func (m *AuthApiRouter) InitAuthApiRouter(router *gin.RouterGroup) {
	router.POST("login", m.Login)
	router.POST("refresh_token", m.RefreshToken)
	router.POST("captcha", m.Captcha)
}

// @summary Administrator login
// @description Administrator login
// @tags public,auth
// @accept json
// @param param body request.AdminLoginParam true "Login user name and password information"
// @success 200 {object} response.Response{data=response.TokenResponse}
// @router /public/login [post]
func (m *AuthApiRouter) Login(c *gin.Context) {
	var param request.AdminLoginParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	//if !store.Verify(param.CaptchaId, param.Captcha, true) {
	//	response.FailWithMessage("验证码错误", c)
	//	return
	//}
	token, err := adminSev.GetAuthService().Login(param)
	if err != nil {
		m.Log(c, "FAILED_LOGIN", fmt.Sprintf("试图登录 %s 账号失败：%s", param.Username, err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.Log(c, "ADMIN_LOGIN", fmt.Sprintf("账号 %s 登录成功", param.Username))
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
	token, err := adminSev.GetAuthService().RefreshToken(param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
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
