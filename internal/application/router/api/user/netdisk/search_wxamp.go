package netdisk

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	netdSev "ai-software-copyright-server/internal/application/service/netdisk"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SearchWxampApiRouter struct {
	api.BaseApi
}

func (m *SearchWxampApiRouter) InitSearchWxampApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("netdisk/search/wxamp")
	m.Router = router
	router.POST("configure/save", m.SaveConfigure)
	router.GET("configure", m.GetConfigureByUserId)
}

// @summary 保存网盘搜索微信小程序配置
// @description 保存网盘搜索微信小程序配置
// @tags netdisk
// @accept json
// @param param body table.NetdiskSearchWxampConfigure true "网盘资源配置信息"
// @success 200 {object} response.Response
// @security user
// @router /netdisk/search/wxamp/configure/save [post]
func (m *SearchWxampApiRouter) SaveConfigure(c *gin.Context) {
	var param table.NetdiskSearchWxampConfigure
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = netdSev.GetSearchWxampConfigureService().SaveConfigure(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "NETDISK_CONFIGURE_SAVE", fmt.Sprintf("保存网盘资源配置失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_CONFIGURE_SAVE", fmt.Sprintf("保存网盘资源配置，用户 Id 为 %d", m.GetUserId(c)))
	response.Ok(c)
}

// @summary 取得网盘搜索微信小程序配置
// @description 取得网盘搜索微信小程序配置
// @tags netdisk
// @success 200 {object} response.Response{data=table.NetdiskSearchWxampConfigure}
// @security user
// @router /netdisk/search/wxamp/configure [get]
func (m *SearchWxampApiRouter) GetConfigureByUserId(c *gin.Context) {
	mod, err := netdSev.GetSearchWxampConfigureService().GetByUserId(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
