package redbook

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	rbSev "ai-software-copyright-server/internal/application/service/redbook"
	"github.com/gin-gonic/gin"
)

type ProhibitedApiRouter struct {
	api.BaseApi
}

func (m *ProhibitedApiRouter) InitProhibitedApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("redbook/prohibited")
	m.Router = router
	router.POST("save", m.Save)
}

// @summary Save redbook prohibited
// @description Save redbook prohibited
// @tags redbook
// @accept json
// @param param body table.RedbookProhibited true "Redbook prohibited information"
// @success 200 {object} response.Response{}
// @security admin
// @router /redbook/prohibited/save [post]
func (m *ProhibitedApiRouter) Save(c *gin.Context) {
	var param table.RedbookProhibited
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = rbSev.GetProhibitedService().Save(param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.Ok(c)
}
