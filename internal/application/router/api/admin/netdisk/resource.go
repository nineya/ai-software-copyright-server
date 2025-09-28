package netdisk

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	netdSev "ai-software-copyright-server/internal/application/service/netdisk"
	"github.com/gin-gonic/gin"
)

type ResourceApiRouter struct {
	api.BaseApi
}

func (m *ResourceApiRouter) InitResourceApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("netdisk/resource")
	m.Router = router
	router.POST("changeAccount", m.ChangeAccount)
}

// @summary 网盘资源更换账号
// @description 网盘资源更换账号(有bug，获取资源的逻辑不对，没有判断夸克网盘，也没有判断资源状态)
// @tags netdisk
// @accept json
// @param param body request.NetdiskResourceChangeAccountParam true "网盘资源账号信息"
// @success 200 {object} response.Response{data=[]table.NetdiskResource}
// @security user
// @router /netdisk/resource/changeAccount [post]
func (m *ResourceApiRouter) ChangeAccount(c *gin.Context) {
	var param request.NetdiskResourceChangeAccountParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = netdSev.GetResourceService().ChangeAccount(param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.Ok(c)
}
