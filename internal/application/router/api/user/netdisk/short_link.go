package netdisk

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	netdSev "ai-software-copyright-server/internal/application/service/netdisk"
	slSev "ai-software-copyright-server/internal/application/service/short_link"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ShortLinkApiRouter struct {
	api.BaseApi
}

func (m *ShortLinkApiRouter) InitShortLinkApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("netdisk/shortLink")
	m.Router = router
	router.POST("", m.Create)
	router.POST("configure/save", m.SaveConfigure)
	router.POST("redirect", m.Redirect)
	router.POST("statistic", m.Statistic)
	router.GET("todayVisits", m.TodayVisits)
	router.GET("activeVisits", m.ActiveVisits)
	router.POST("allStatistic", m.AllStatistic)
	router.GET("list", m.GetByPage)
	router.GET("configure", m.GetConfigureByUserId)
}

// @summary Create cloud disk link Info
// @description Create cloud disk link Info
// @tags netdisk
// @accept json
// @param param body request.ShortLinkCreateCloudDiskParam true "Cloud disk link Info"
// @success 200 {object} response.Response{data=response.UserBuyContentResponse}
// @security user
// @router /netdisk/shortLink [post]
func (m *ShortLinkApiRouter) Create(c *gin.Context) {
	var param request.ShortLinkCreateCloudDiskParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = userSev.GetClientInfoService().MsgSecCheck(m.GetUserId(c), utils.GetClientType(c), param.Content)
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

// @summary 保存网盘短链配置
// @description 保存网盘短链配置
// @tags netdisk
// @accept json
// @param param body table.NetdiskSearchAppConfigure true "网盘资源配置信息"
// @success 200 {object} response.Response
// @security user
// @router /netdisk/shortLink/configure/save [post]
func (m *ShortLinkApiRouter) SaveConfigure(c *gin.Context) {
	var param table.NetdiskShortLinkConfigure
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = netdSev.GetShortLinkConfigureService().SaveConfigure(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "NETDISK_CONFIGURE_SAVE", fmt.Sprintf("保存网盘资源配置失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_CONFIGURE_SAVE", fmt.Sprintf("保存网盘资源配置，用户 Id 为 %d", m.GetUserId(c)))
	response.Ok(c)
}

// @summary 重定向网盘短链
// @description 重定向网盘短链
// @tags netdisk
// @accept json
// @param param body request.ShortLinkRedirectParam true "Redirect link Info"
// @success 200 {object} response.Response{data=response.ShortLinkRedirectResponse}
// @security user
// @router /netdisk/shortLink/redirect [post]
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

// @summary 网盘短链访问统计分析
// @description 网盘短链访问统计分析
// @tags netdisk
// @accept json
// @param id query int64 true "Short link id"
// @success 200 {object} response.Response{data=response.ShortLinkRedirectResponse}
// @security user
// @router /netdisk/shortLink/statistic [post]
func (m *ShortLinkApiRouter) Statistic(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}

	mod, err := slSev.GetNetdiskService().Statistic(m.GetUserId(c), id)
	if err != nil {
		m.UserLog(c, "NETDISK_SHORT_LINK_STATISTIC", fmt.Sprintf("网盘短链分析失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_SHORT_LINK_STATISTIC", fmt.Sprintf("网盘短链分析，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}

// @summary 今日访问数据展示
// @description 今日访问数据展示
// @tags netdisk
// @accept json
// @success 200 {object} response.Response{data=response.ShortLinkRedirectResponse}
// @security user
// @router /netdisk/shortLink/todayVisits [get]
func (m *ShortLinkApiRouter) TodayVisits(c *gin.Context) {
	mod, err := slSev.GetNetdiskService().TodayVisits(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}

// @summary 短链活跃数据
// @description 短链活跃数据
// @tags netdisk
// @accept json
// @success 200 {object} response.Response{data=response.ShortLinkRedirectResponse}
// @security user
// @router /netdisk/shortLink/activeVisits [get]
func (m *ShortLinkApiRouter) ActiveVisits(c *gin.Context) {
	mod, err := slSev.GetNetdiskService().ActiveVisits(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}

// @summary 网盘短链整体访问统计分析
// @description 网盘短链整体访问统计分析
// @tags netdisk
// @accept json
// @success 200 {object} response.Response{data=response.ShortLinkRedirectResponse}
// @security user
// @router /netdisk/shortLink/allStatistic [post]
func (m *ShortLinkApiRouter) AllStatistic(c *gin.Context) {
	mod, err := slSev.GetNetdiskService().AllStatistic(m.GetUserId(c))
	if err != nil {
		m.UserLog(c, "NETDISK_SHORT_LINK_STATISTIC", fmt.Sprintf("网盘短链分析失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_SHORT_LINK_STATISTIC", fmt.Sprintf("网盘短链分析，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}

// @summary 列表分页查询网盘短链
// @description 列表分页查询网盘短链
// @tags form
// @param param query request.NetdiskResourceQueryPageParam true "短链列表的分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.NetdiskShortLink}}
// @security user
// @router /netdisk/shortLink/list [get]
func (m *ShortLinkApiRouter) GetByPage(c *gin.Context) {
	var param request.QueryPageParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	page, err := slSev.GetNetdiskService().GetByPage(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}

// @summary 取得网盘短链配置
// @description 取得网盘短链配置
// @tags netdisk
// @success 200 {object} response.Response{data=table.NetdiskSearchAppConfigure}
// @security user
// @router /netdisk/shortLink/configure [get]
func (m *ShortLinkApiRouter) GetConfigureByUserId(c *gin.Context) {
	mod, err := netdSev.GetShortLinkConfigureService().GetByUserId(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
