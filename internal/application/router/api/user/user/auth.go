package user

import (
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	"ai-software-copyright-server/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthApiRouter struct {
	api.BaseApi
}

func (m *AuthApiRouter) InitAuthApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("auth")
	router.POST("logout", m.Logout)
}

// @summary User logout
// @description User logout
// @tags auth
// @accept json
// @success 200 {object} response.Response
// @security user
// @router /auth/logout [post]
func (m *AuthApiRouter) Logout(c *gin.Context) {
	claims := m.GetClaims(c)
	utils.RemoveToken(claims)
	m.UserLog(c, "USER_LOGOUT", "注销登录成功")
	response.Ok(c)
}
