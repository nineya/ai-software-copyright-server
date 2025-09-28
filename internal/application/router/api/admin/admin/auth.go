package admin

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

// @summary Administrator logout
// @description Administrator logout
// @tags auth
// @accept json
// @success 200 {object} response.Response
// @security admin
// @router /auth/logout [post]
func (m *AuthApiRouter) Logout(c *gin.Context) {
	claims := m.GetClaims(c)
	utils.RemoveToken(claims)
	m.AdminLog(c, "ADMIN_LOGOUT", "注销登录成功")
	response.Ok(c)
}
