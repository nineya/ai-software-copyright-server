package public

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	netdSev "ai-software-copyright-server/internal/application/service/netdisk"
	"ai-software-copyright-server/internal/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

type NetdiskApiRouter struct {
	api.BaseApi
}

func (m *NetdiskApiRouter) InitNetdiskApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("netdisk")
	m.Router = router
	router.POST("resource/save", m.Save)
	router.PUT("resource/checkResult", m.UpdateCheckResult)
	router.GET("resource/:id", m.GetById)
	router.GET("resource/list", m.Search) // TODO 旧版本想废弃
	router.GET("resource/search", m.Search)
	router.GET("resource/awaitCheck", m.GetAwaitCheck)
	router.GET("configure", m.GetConfigureByUserId)
}

// @summary 转存网盘资源
// @description 转存网盘资源
// @tags netdisk
// @accept json
// @param param body table.NetdiskResource true "网盘资源信息"
// @success 200 {object} response.Response{data=string}
// @security user
// @router /public/netdisk/resource/save [post]
func (m *NetdiskApiRouter) Save(c *gin.Context) {
	var param request.NetdiskResourceSaveParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := netdSev.GetResourceService().Save(utils.GetHeaderUserId(c), param)
	if err != nil {
		//m.UserLog(c, "NETDISK_RESOURCE_SAVE", fmt.Sprintf("用户 %d 转存网盘资源失败，原因：%s", m.GetUserId(c), err.Error()))
		response.FailWithError(err, c)
		return
	}
	//m.UserLog(c, "NETDISK_RESOURCE_SAVE", fmt.Sprintf("用户 %d 转存网盘资源，转存后分享链接：%s", m.GetUserId(c), mod))
	response.OkWithData(mod.TargetUrl, c)
}

// @summary 修改检查的网盘资源状态（给客户端用，客户端协助查询和修改资源）
// @description 修改检查的网盘资源状态
// @tags netdisk
// @param param body table.NetdiskResource true "网盘资源信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.NetdiskResource}}
// @security user
// @router /public/netdisk/resource/checkResult [put]
func (m *NetdiskApiRouter) UpdateCheckResult(c *gin.Context) {
	var param table.NetdiskResource
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = netdSev.GetResourceService().UpdateCheckResult(param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.Ok(c)
}

// @summary 取得一部分待检测资源列表
// @description 取得一部分待检测资源
// @tags netdisk
// @success 200 {object} response.Response{data=[]table.NetdiskResource}
// @security user
// @router /public/netdisk/resource/awaitCheck [get]
func (m *NetdiskApiRouter) GetAwaitCheck(c *gin.Context) {
	typ, _ := enum.NetdiskTypeValue(c.DefaultQuery("type", "QUARK"))
	switch typ {
	case enum.NetdiskType(2):
		response.OkWithData(netdSev.GetResourceService().GetCheckQuarkResource(10), c)
	case enum.NetdiskType(4):
		response.OkWithData(netdSev.GetResourceService().GetCheckBaiduResource(10), c)
	default:
		response.OkWithData(netdSev.GetResourceService().GetCheckQuarkResource(10), c)
	}
}

// @summary 取得网盘资源信息
// @description 取得网盘资源信息
// @tags netdisk
// @accept json
// @param id path int64 true "资源id"
// @success 200 {object} response.Response
// @security user
// @router /netdisk/resource/{id} [get]
func (m *NetdiskApiRouter) GetById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	mod, err := netdSev.GetResourceService().GetById(utils.GetHeaderUserId(c), id)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}

// @summary 列表分页查询网盘资源
// @description 列表分页查询网盘资源，给搜索小程序和app提供服务
// @tags netdisk
// @param param query request.QueryPageParam true "资源列表的分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.NetdiskResource}}
// @security user
// @router /public/netdisk/resource/search [get]
func (m *NetdiskApiRouter) Search(c *gin.Context) {
	var param request.NetdiskResourceSearchParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	// TODO 如果来源于搜素小程序，主动再次查询网盘配置，因为旧版本搜索小程序不会主动发送配置
	if m.GetClientType(c) == enum.ClientType(7) {
		configure, err := netdSev.GetSearchWxampConfigureService().GetByUserId(utils.GetHeaderUserId(c))
		if err != nil {
			response.FailWithError(err, c)
			return
		}
		param.CollectTypes = utils.ListTransform(configure.CollectTypes, func(item enum.NetdiskType) string {
			return enum.NETDISK_TYPE[item]
		})
		param.SecureMode = enum.NETDISK_SECURE_MODE[configure.SecureMode]
	}
	page, err := netdSev.GetResourceService().Search(utils.GetHeaderUserId(c), m.GetClientType(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}

// @summary 取得网盘资源配置
// @description 取得网盘资源配置，小程序和APP中会调用这个接口获取配置
// @tags netdisk
// @success 200 {object} response.Response{data=table.NetdiskSearchWxampConfigure}
// @security user
// @router /netdisk/configure [get]
func (m *NetdiskApiRouter) GetConfigureByUserId(c *gin.Context) {
	switch m.GetClientType(c) {
	case enum.ClientType(7): //网盘搜索小程序
		mod, err := netdSev.GetSearchWxampConfigureService().GetByUserId(utils.GetHeaderUserId(c))
		if err != nil {
			response.FailWithError(err, c)
			return
		}
		// TODO 安全模式或审核模式，设置为旧版的审核，兼容旧版本
		if mod.SecureMode == enum.NetdiskSecureMode(2) || mod.SecureMode == enum.NetdiskSecureMode(3) {
			mod.IsAudit = true
		}
		// 审核模式，设置默认羊毛
		if mod.SecureMode == enum.NetdiskSecureMode(3) {
			mod.WelfareConfig = []common.WelfareLabel{
				{
					Label: "美团外卖",
					Items: []common.WelfareItem{
						{
							Title:   "[外卖]每日红包",
							Desc:    "快捷打开美团外卖小程序每日红包页面。",
							WxAppid: "wxde8ac0a21135c07d",
							WxPath:  "/index/pages/h5/h5?weburl=https%3A%2F%2Fclick.meituan.com%2Ft%3Ft%3D1%26c%3D2%26p%3DNWqXdL9zsQgZ",
						}, {
							Title:   "[外卖]美团霸王餐",
							Desc:    "快捷打开美团外卖小程序霸王餐页面。",
							WxAppid: "wxde8ac0a21135c07d",
							WxPath:  "/waimai/pages/web-view/web-view?type=REDIRECT&webviewUrl=https%3A%2F%2Foffsiteact.meituan.com%2Fact%2Fcps%2Fpromotion%3Fp%3D9bc7499724d849eeb8a3832e11af8373&utm_content=0c3bfd35279b4140b3bd8ecbc41301d6__ac118d53892c424cbd2b0689f95bafd6",
						},
					},
				},
			}
		}
		response.OkWithData(mod, c)
	case enum.ClientType(8): // 网盘搜索APP
		mod, err := netdSev.GetSearchAppConfigureService().GetByUserId(utils.GetHeaderUserId(c))
		if err != nil {
			response.FailWithError(err, c)
			return
		}
		response.OkWithData(mod, c)
	}
}
