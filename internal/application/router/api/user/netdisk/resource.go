package netdisk

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	netdSev "ai-software-copyright-server/internal/application/service/netdisk"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type ResourceApiRouter struct {
	api.BaseApi
}

func (m *ResourceApiRouter) InitResourceApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("netdisk/resource")
	m.Router = router
	router.POST("", m.Create)
	router.POST("import", m.Import)
	router.POST("save", m.Save)
	router.DELETE(":id", m.DeleteById)
	router.DELETE("batch", m.DeleteInBatch)
	router.DELETE("clear", m.Clear)
	router.PUT(":id", m.UpdateById)
	router.PUT("status/:status", m.UpdateStatusInBatch)
	router.GET("deleteSearch", m.GetDeleteSearch)
	router.GET("list", m.GetByPage)
	router.GET("search", m.Search)
	router.GET("search/list", m.GetSearchByPage)
}

// @summary 创建网盘资源
// @description 创建网盘资源
// @tags netdisk
// @accept json
// @param param body table.NetdiskResource true "网盘资源信息"
// @success 200 {object} response.Response{data=[]table.NetdiskResource}
// @security user
// @router /netdisk/resource [post]
func (m *ResourceApiRouter) Create(c *gin.Context) {
	var param table.NetdiskResource
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := netdSev.GetResourceService().Create(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "NETDISK_RESOURCE_CREATE", fmt.Sprintf("创建网盘资源失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_RESOURCE_CREATE", fmt.Sprintf("创建网盘资源 %s，花费：%d，剩余：%d", param.Name, mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}

// @summary 批量导入网盘资源
// @description 批量导入网盘资源
// @tags netdisk
// @accept json
// @param file formData file true "网盘资源导入文件"
// @success 200 {object} response.Response{}
// @security user
// @router /netdisk/resource/import [post]
func (m *ResourceApiRouter) Import(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := netdSev.GetResourceService().Import(m.GetUserId(c), file)
	if err != nil {
		m.UserLog(c, "NETDISK_RESOURCE_IMPORT", fmt.Sprintf("批量导入网盘资源失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_RESOURCE_IMPORT", fmt.Sprintf("批量导入网盘资源，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}

// @summary 转存网盘资源
// @description 转存网盘资源，给扣子用
// @tags netdisk
// @accept json
// @param param body request.NetdiskResourceSaveParam true "网盘资源信息"
// @success 200 {object} response.Response{data=table.NetdiskResource}
// @security user
// @router /netdisk/resource/save [post]
func (m *ResourceApiRouter) Save(c *gin.Context) {
	var param request.NetdiskResourceSaveParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	if param.ShareTargetUrl != "" && global.SOCKET.GetClient(m.GetUserId(c)) == nil {
		response.FailWithMessage("转存资源需要网盘助手客户端支持，未搭建网盘助手或助手已经离线，搭建请联系微信：nineyaccz", c)
		return
	}
	mod, err := netdSev.GetResourceService().Save(m.GetUserId(c), param)
	if err != nil {
		//m.UserLog(c, "NETDISK_RESOURCE_SAVE", fmt.Sprintf("用户 %d 转存网盘资源失败，原因：%s", m.GetUserId(c), err.Error()))
		response.FailWithError(err, c)
		return
	}
	//m.UserLog(c, "NETDISK_RESOURCE_SAVE", fmt.Sprintf("用户 %d 转存网盘资源，转存后分享链接：%s", m.GetUserId(c), mod))
	response.OkWithData(mod, c)
}

// @summary 删除网盘资源
// @description 删除网盘资源
// @tags netdisk
// @param id path int64 true "资源id"
// @success 200 {object} response.Response
// @security admin
// @router /netdisk/resource/{id} [delete]
func (m *ResourceApiRouter) DeleteById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	if err = netdSev.GetResourceService().DeleteById(m.GetUserId(c), id); err != nil {
		m.UserLog(c, "NETDISK_RESOURCE_DELETE", fmt.Sprintf("删除 Id 为 %d 的网盘资源失败，原因：%s", id, err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_RESOURCE_DELETE", fmt.Sprintf("删除 Id 为 %d 的网盘资源", id))
	response.Ok(c)
}

// @summary 批量删除网盘资源
// @description 批量删除网盘资源
// @tags netdisk
// @accept json
// @param param body []int64 true "资源id列表"
// @success 200 {object} response.Response
// @security user
// @router /netdisk/resource/batch [delete]
func (m *ResourceApiRouter) DeleteInBatch(c *gin.Context) {
	param := make([]int64, 0)
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	if err = netdSev.GetResourceService().DeleteInBatch(m.GetUserId(c), param); err != nil {
		m.UserLog(c, "NETDISK_RESOURCE_DELETE", fmt.Sprintf("批量删除 Id 为 %v 的网盘资源失败，原因：%s", param, err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_RESOURCE_DELETE", fmt.Sprintf("批量删除 Id 为 %v 的网盘资源", param))
	response.Ok(c)
}

// @summary 清空用户全部资源
// @description 清空用户全部资源
// @tags netdisk
// @accept json
// @success 200 {object} response.Response
// @security user
// @router /netdisk/resource/clear [delete]
func (m *ResourceApiRouter) Clear(c *gin.Context) {
	origin, err := enum.NetdiskOriginValue(c.Query("origin"))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	status, err := enum.NetdiskStatusValue(c.Query("status"))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	msg, err := netdSev.GetResourceService().Clear(m.GetUserId(c), origin, status)
	if err != nil {
		m.UserLog(c, "NETDISK_RESOURCE_DELETE", fmt.Sprintf("清除用户 %d 的网盘资源失败，原因：%s", m.GetUserId(c), err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_RESOURCE_DELETE", fmt.Sprintf("清除用户 %d 的网盘资源", m.GetUserId(c)))
	response.OkWithData(msg, c)
}

// @summary 更新网盘资源信息
// @description 更新网盘资源信息
// @tags netdisk
// @accept json
// @param id path int64 true "资源id"
// @param param body table.NetdiskResource true "资源信息"
// @success 200 {object} response.Response
// @security user
// @router /netdisk/resource/{id} [put]
func (m *ResourceApiRouter) UpdateById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	var param table.NetdiskResource
	err = c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	param.Type = utils.TransformNetdiskType(param.TargetUrl)
	if err = netdSev.GetResourceService().UpdateById(m.GetUserId(c), id, param); err != nil {
		m.UserLog(c, "NETDISK_RESOURCE_UPDATE", fmt.Sprintf("更新网盘资源 %s 失败，原因：%s", param.Name, err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_RESOURCE_UPDATE", fmt.Sprintf("更新网盘资源 %s，资源 Id 为 %d", param.Name, id))
	response.Ok(c)
}

// @summary 批量更新网盘资源状态
// @description 批量更新网盘资源状态
// @tags netdisk
// @accept json
// @param param body request.NetdiskResourceUpdateInBatchParam true "批量更新状态信息参数"
// @success 200 {object} response.Response
// @security admin
// @router /netdisk/resource/status/{status} [put]
func (m *ResourceApiRouter) UpdateStatusInBatch(c *gin.Context) {
	param := make([]int64, 0)
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	status, err := enum.NetdiskStatusValue(c.Param("status"))
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	if err = netdSev.GetResourceService().UpdateStatusInBatch(m.GetUserId(c), param, status); err != nil {
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_RESOURCE_UPDATE", fmt.Sprintf("批量更新 Id 为 %v 的网盘资源信息", param))
	response.Ok(c)
}

// @summary 取得超过指定时间需要删除的夸克资源
// @description 取得超过指定时间需要删除的夸克资源
// @tags netdisk
// @success 200 {object} response.Response{data=[]table.NetdiskResource}
// @security user
// @router /netdisk/resource/deleteSearch [get]
func (m *ResourceApiRouter) GetDeleteSearch(c *gin.Context) {
	time, err := strconv.Atoi(c.Query("time"))
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	typ, err := enum.NetdiskTypeValue(c.Query("type"))
	if typ == 0 || err != nil {
		// TODO 1.1.3版本前，默认夸克网盘，兼容旧版本
		typ = enum.NetdiskType(2)
	}
	response.OkWithData(netdSev.GetResourceService().GetDeleteSearchResource(m.GetUserId(c), typ, time), c)
}

// @summary 列表分页查询网盘资源
// @description 列表分页查询网盘资源
// @tags netdisk
// @param param query request.NetdiskResourceQueryPageParam true "资源列表的分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.NetdiskResource}}
// @security user
// @router /netdisk/resource/list [get]
func (m *ResourceApiRouter) GetByPage(c *gin.Context) {
	var param request.NetdiskResourceQueryPageParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	page, err := netdSev.GetResourceService().GetByPage(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}

// @summary 列表分页搜索网盘资源
// @description 列表分页搜索网盘资源，给扣子插件和网盘助手用
// @tags netdisk
// @param param query request.QueryPageParam true "资源列表的分页搜索信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.NetdiskResource}}
// @security user
// @router /netdisk/resource/search [get]
func (m *ResourceApiRouter) Search(c *gin.Context) {
	var param request.NetdiskResourceSearchParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	if len(param.CollectTypes) > 0 && global.SOCKET.GetClient(m.GetUserId(c)) == nil {
		response.FailWithMessage("资源采集需要网盘助手客户端支持，未搭建网盘助手或助手已经离线，搭建请联系微信：nineyaccz", c)
		return
	}
	page, err := netdSev.GetResourceService().Search(m.GetUserId(c), m.GetClientType(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}

// @summary 列表分页查询网盘资源搜索记录
// @description 列表分页查询网盘资源搜索记录
// @tags netdisk
// @param param query request.QueryPageParam true "搜索记录列表的分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.NetdiskResourceSearch}}
// @security user
// @router /netdisk/resource/search/list [get]
func (m *ResourceApiRouter) GetSearchByPage(c *gin.Context) {
	var param request.QueryPageParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	page, err := netdSev.GetResourceSearchService().GetByPage(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}
