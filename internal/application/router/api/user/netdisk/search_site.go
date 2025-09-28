package netdisk

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	netdSev "ai-software-copyright-server/internal/application/service/netdisk"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SearchSiteApiRouter struct {
	api.BaseApi
}

func (m *SearchSiteApiRouter) InitSearchSiteApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("netdisk/search/site")
	m.Router = router
	router.POST("configure/save", m.SaveConfigure)
	router.GET("configure", m.GetConfigureByUserId)
}

// @summary 保存网盘搜索站点配置
// @description 保存网盘搜索站点配置
// @tags netdisk
// @accept json
// @param param body table.NetdiskSearchSiteConfigure true "网盘资源配置信息"
// @success 200 {object} response.Response
// @security user
// @router /netdisk/search/site/configure/save [post]
func (m *SearchSiteApiRouter) SaveConfigure(c *gin.Context) {
	var param table.NetdiskSearchSiteConfigure
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = netdSev.GetSearchSiteConfigureService().SaveConfigure(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "NETDISK_CONFIGURE_SAVE", fmt.Sprintf("保存网盘资源配置失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_CONFIGURE_SAVE", fmt.Sprintf("保存网盘资源配置，用户 Id 为 %d", m.GetUserId(c)))
	response.Ok(c)
}

// @summary 取得网盘搜索站点配置
// @description 取得网盘搜索站点配置
// @tags netdisk
// @success 200 {object} response.Response{data=table.NetdiskSearchSiteConfigure}
// @security user
// @router /netdisk/search/site/configure [get]
func (m *SearchSiteApiRouter) GetConfigureByUserId(c *gin.Context) {
	mod, err := netdSev.GetSearchSiteConfigureService().GetByUserId(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
