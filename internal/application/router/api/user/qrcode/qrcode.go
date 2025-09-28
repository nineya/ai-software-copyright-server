package qrcode

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	qrcodeSev "ai-software-copyright-server/internal/application/service/qrcode"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type QrcodeApiRouter struct {
	api.BaseApi
}

func (m *QrcodeApiRouter) InitQrcodeApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("qrcode")
	m.Router = router
	router.POST("build", m.Build)
	router.POST("loose", m.Create)
	router.POST("loose/:id/image", m.AddImageById)
	router.DELETE("loose/:id", m.DeleteById)
	router.DELETE("loose/:id/image", m.DeleteImageById)
	router.PUT("loose/:id", m.UpdateById)
	router.GET("loose/list", m.GetByPage)
}

// @summary Build qrcode
// @description Build qrcode
// @tags qrcode
// @accept json
// @param param body request.QrcodeBuildParam true "Qrcode content Info"
// @success 200 {object} response.Response{data=string}
// @security user
// @router /qrcode/build [post]
func (m *QrcodeApiRouter) Build(c *gin.Context) {
	var param request.QrcodeBuildParam
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

	mod, err := qrcodeSev.GetQrcodeService().Build(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "QRCODE_BUILD", fmt.Sprintf("二维码生成失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "QRCODE_BUILD", fmt.Sprintf("二维码生成，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}

// @summary 创建二维码活码
// @description 创建二维码活码
// @tags qrcode
// @accept json
// @param param body table.NetdiskResource true "活码信息"
// @success 200 {object} response.Response{data=[]table.Qrcode}
// @security user
// @router /qrcode/loose [post]
func (m *QrcodeApiRouter) Create(c *gin.Context) {
	var param table.Qrcode
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	param.Alias = strings.TrimSpace(fmt.Sprintf("%32s", strconv.FormatInt(time.Now().UnixMilli(), 32)))
	mod, err := qrcodeSev.GetQrcodeService().Create(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "QRCODE_LOOSE_CREATE", fmt.Sprintf("创建活码失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "QRCODE_LOOSE_CREATE", fmt.Sprintf("创建活码 %s", param.Title))
	response.OkWithData(mod, c)
}

// @summary 上传活码图片
// @description 上传活码图片
// @tags qrcode
// @accept x-www-form-urlencoded
// @param file formData file true "图片文件流"
// @success 200 {object} response.Response{data=response.QrcodeAddImageResponse}
// @security user
// @router /qrcode/loose/{id}/image [post]
func (m *QrcodeApiRouter) AddImageById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	file, err := c.FormFile("image")
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := qrcodeSev.GetQrcodeService().AddImageById(m.GetUserId(c), id, file)
	if err != nil {
		m.UserLog(c, "QRCODE_LOOSE_ADD_IMAGE", fmt.Sprintf("活码添加图片失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "QRCODE_LOOSE_ADD_IMAGE", fmt.Sprintf("活码添加图片，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}

// @summary 删除活码
// @description 删除活码
// @tags qrcode
// @param id path int64 true "活码id"
// @success 200 {object} response.Response
// @security user
// @router /qrcode/loose/{id} [delete]
func (m *QrcodeApiRouter) DeleteById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	if err = qrcodeSev.GetQrcodeService().DeleteById(m.GetUserId(c), id); err != nil {
		m.UserLog(c, "QRCODE_LOOSE_DELETE", fmt.Sprintf("删除 Id 为 %d 的活码失败，原因：%s", id, err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "QRCODE_LOOSE_DELETE", fmt.Sprintf("删除 Id 为 %d 的活码", id))
	response.Ok(c)
}

// @summary 删除活码图片
// @description 删除活码图片
// @tags qrcode
// @accept json
// @param id path int64 true "活码id"
// @param param body table.QrcodeDeleteImageParam true "活码信息"
// @success 200 {object} response.Response
// @security user
// @router /qrcode/loose/{id}/image [delete]
func (m *QrcodeApiRouter) DeleteImageById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	var param request.QrcodeDeleteImageParam
	err = c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = qrcodeSev.GetQrcodeService().DeleteImageById(m.GetUserId(c), id, param.TargetUrl)
	if err != nil {
		m.UserLog(c, "QRCODE_LOOSE_DELETE_IMAGE", fmt.Sprintf("活码删除图片失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "QRCODE_LOOSE_DELETE_IMAGE", fmt.Sprintf("活码删除图片，图片地址：%s", param.TargetUrl))
	response.Ok(c)
}

// @summary 更新活码信息
// @description 更新活码信息
// @tags qrcode
// @accept json
// @param id path int64 true "活码id"
// @param param body table.Qrcode true "活码信息"
// @success 200 {object} response.Response
// @security user
// @router /qrcode/loose/{id} [put]
func (m *QrcodeApiRouter) UpdateById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	var param table.Qrcode
	err = c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	if err = qrcodeSev.GetQrcodeService().UpdateById(m.GetUserId(c), id, param); err != nil {
		m.UserLog(c, "QRCODE_LOOSE_UPDATE", fmt.Sprintf("更新活码 %s 失败，原因：%s", param.Title, err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "QRCODE_LOOSE_UPDATE", fmt.Sprintf("更新活码  %s，活码 Id 为 %d", param.Title, id))
	response.Ok(c)
}

// @summary 列表分页查询活码
// @description 列表分页查询活码
// @tags qrcode
// @param param query request.QueryPageParam true "分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.Qrcode}}
// @security user
// @router /qrcode/loose/list [get]
func (m *QrcodeApiRouter) GetByPage(c *gin.Context) {
	var param request.PageableParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	page, err := qrcodeSev.GetQrcodeService().GetByPage(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}
