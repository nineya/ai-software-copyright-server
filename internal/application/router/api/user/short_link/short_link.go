package short_link

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	slSev "ai-software-copyright-server/internal/application/service/short_link"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ShortLinkApiRouter struct {
	api.BaseApi
}

func (m *ShortLinkApiRouter) InitShortLinkApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("shortLink")
	m.Router = router
	router.POST("createCloudDisk", m.CreateCloudDisk)
	router.POST("redirect", m.Redirect)
}

// @summary Create cloud disk link Info
// @description Create cloud disk link Info
// @tags shortLink
// @accept json
// @param param body request.ShortLinkCreateCloudDiskParam true "Cloud disk link Info"
// @success 200 {object} response.Response{data=response.UserBuyContentResponse}
// @security user
// @router /shortLink/createCloudDisk [post]
// Deprecated
func (m *ShortLinkApiRouter) CreateCloudDisk(c *gin.Context) {
	var param request.ShortLinkCreateCloudDiskParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}

	mod, err := slSev.GetNetdiskService().Create(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "NETDISK_SHORT_LINK_CREATE", fmt.Sprintf("创建网盘短链失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_SHORT_LINK_CREATE", fmt.Sprintf("创建网盘短链，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}

// @summary 重定向网盘短链
// @description 重定向网盘短链
// @tags shortLink
// @accept json
// @param param body request.ShortLinkRedirectParam true "Redirect link Info"
// @success 200 {object} response.Response{data=response.ShortLinkRedirectResponse}
// @security user
// @router /shortLink/redirect [post]
// Deprecated
func (m *ShortLinkApiRouter) Redirect(c *gin.Context) {
	var param request.ShortLinkRedirectParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}

	mod, err := slSev.GetNetdiskService().Redirect(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "NETDISK_SHORT_LINK_REDIRECT", fmt.Sprintf("网盘短链重定向失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_SHORT_LINK_REDIRECT", fmt.Sprintf("网盘短链重定向，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}
