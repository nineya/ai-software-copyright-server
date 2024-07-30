package admin

import (
	"github.com/gin-gonic/gin"
	"tool-server/internal/application/param/response"
	"tool-server/internal/application/router/api"
	adminSev "tool-server/internal/application/service/admin"
)

type AdminApiRouter struct {
	api.BaseApi
}

func (m *AdminApiRouter) InitAdminApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("admin")
	m.Router = router
	router.GET("profiles", m.GetProfiles)
}

// @summary Get admin profiles
// @description Get admin profiles
// @tags admin
// @Produce json
// @success 200 {object} response.Response{data=table.Admin}
// @Security admin
// @router /admin/profiles [get]
func (m *AdminApiRouter) GetProfiles(c *gin.Context) {
	user, err := adminSev.GetAdminService().GetById(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	user.Password = ""
	response.OkWithData(user, c)
}
