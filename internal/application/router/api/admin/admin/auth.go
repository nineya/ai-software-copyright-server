package admin

import (
	"github.com/gin-gonic/gin"
	"tool-server/internal/application/param/response"
	"tool-server/internal/application/router/api"
	adminSev "tool-server/internal/application/service/admin"
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
	adminSev.GetAuthService().Logout(claims)
	m.Log(c, "ADMIN_LOGOUT", "注销登录成功")
	response.Ok(c)
}
