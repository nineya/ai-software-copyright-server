package netdisk

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	netdSev "ai-software-copyright-server/internal/application/service/netdisk"
	"ai-software-copyright-server/internal/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"unicode/utf8"
)

type HelperApiRouter struct {
	api.BaseApi
}

func (m *HelperApiRouter) InitHelperApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("netdisk/helper")
	m.Router = router
	router.POST("configure/save", m.SaveConfigure)
	router.POST("sendRequest", m.SendRequest)
	router.GET("configure", m.GetConfigureByUserId)
	router.GET("client", m.GetClientByUserId)
}

// @summary 保存网盘工具人配置
// @description 保存网盘工具人配置
// @tags netdisk
// @accept json
// @param param body table.NetdiskHelperConfigure true "网盘资源配置信息"
// @success 200 {object} response.Response
// @security user
// @router /netdisk/helper/configure/save [post]
func (m *HelperApiRouter) SaveConfigure(c *gin.Context) {
	var param table.NetdiskHelperConfigure
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = netdSev.GetHelperConfigureService().SaveConfigure(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "NETDISK_CONFIGURE_SAVE", fmt.Sprintf("保存网盘助手配置失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "NETDISK_CONFIGURE_SAVE", fmt.Sprintf("保存网盘助手配置，用户 Id 为 %d", m.GetUserId(c)))
	response.Ok(c)
}

// @summary 发送网盘助手信息
// @description 发送网盘助手信息
// @tags client
// @accept json
// @param param body common.SocketMessage true "发送网盘助手信息"
// @success 200 {object} response.Response{data=string}
// @security user
// @router /netdisk/helper/sendRequest [post]
func (m *HelperApiRouter) SendRequest(c *gin.Context) {
	var param common.SocketMessage
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}

	var data any
	existResult, err := global.SOCKET.SendMessage(m.GetUserId(c), param, &data)
	// 字符串超长截断
	message := param.Data
	wordCount := utf8.RuneCountInString(message)
	if wordCount > 10000 {
		message = string([]rune(message)[:10000]) + "..."
	}
	if err != nil {
		m.UserLog(c, "NETDISK_HELPER_SEND_REQUEST", fmt.Sprintf("发送网盘助手请求失败：userId = %d, type = %d, message = %s，原因：%s", m.GetUserId(c), param.Type, message, err.Error()))
		response.FailWithError(err, c)
		return
	}
	if !existResult {
		data = "请求成功"
	}
	m.UserLog(c, "NETDISK_HELPER_SEND_REQUEST", fmt.Sprintf("发送网盘助手请求成功：userId = %d, type = %d, message = %s", m.GetUserId(c), param.Type, message))
	response.OkWithData(data, c)
}

// @summary 取得网盘工具人配置
// @description 取得网盘工具人配置
// @tags netdisk
// @success 200 {object} response.Response{data=table.NetdiskHelperConfigure}
// @security user
// @router /netdisk/helper/configure [get]
func (m *HelperApiRouter) GetConfigureByUserId(c *gin.Context) {
	mod, err := netdSev.GetHelperConfigureService().GetByUserId(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}

// @summary 取得网盘工具人客户端信息
// @description 取得网盘工具人客户端信息
// @tags netdisk
// @success 200 {object} response.Response{data=response.NetdiskHelperClientResponse}
// @security user
// @router /netdisk/helper/client [get]
func (m *HelperApiRouter) GetClientByUserId(c *gin.Context) {
	client := global.SOCKET.GetClient(m.GetUserId(c))
	result := response.NetdiskHelperClientResponse{}
	if client != nil {
		result.Online = true
		result.Version = client.Version
		latestVersion := global.NetdiskHelperUpdateNotes[0]
		if client.Version != latestVersion.Version {
			result.IsUpdate = true
			result.NewVersion = latestVersion.Version
			result.NewVersionDownloadUrl = latestVersion.DownloadUrl
			result.NewVersionDescription = latestVersion.Description
		}
	}
	response.OkWithData(result, c)
}
