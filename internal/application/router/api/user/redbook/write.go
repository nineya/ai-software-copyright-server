package redbook

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	rbSev "ai-software-copyright-server/internal/application/service/redbook"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

type WriteApiRouter struct {
	api.BaseApi
}

func (m *WriteApiRouter) InitWriteApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("redbook/write")
	m.Router = router
	router.POST("title", m.Title)
	router.POST("note", m.Note)
	router.POST("planting", m.Planting)
}

// @summary Ai writing little red book title
// @description Ai writing little red book title
// @tags redbook
// @accept json
// @param param body request.RedbookWriteMessageParam true "Writing instruction"
// @success 200 {object} response.Response{data=*response.RedbookWriteTitleResponse}
// @security user
// @router /redbook/write/title [post]
func (m *WriteApiRouter) Title(c *gin.Context) {
	var param request.RedbookWriteMessageParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = userSev.GetClientInfoService().MsgSecCheck(m.GetUserId(c), utils.GetClientType(c), param.Message)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := rbSev.GetWriteService().Title(m.GetUserId(c), param.Message)
	if err != nil {
		m.UserLog(c, "REDBOOK_WRITE_TITLE", fmt.Sprintf("小红书爆款标题生成失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "REDBOOK_WRITE_TITLE", fmt.Sprintf("小红书爆款标题生成，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}

// @summary Ai writing little red book note
// @description Ai writing little red book note
// @tags redbook
// @accept json
// @param param body request.RedbookWriteMessageParam true "Writing instruction"
// @success 200 {object} response.Response{}
// @security user
// @router /redbook/write/note [post]
func (m *WriteApiRouter) Note(c *gin.Context) {
	var param request.RedbookWriteMessageParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = userSev.GetClientInfoService().MsgSecCheck(m.GetUserId(c), utils.GetClientType(c), param.Message)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := rbSev.GetWriteService().Note(m.GetUserId(c), param.Message)
	if err != nil {
		m.UserLog(c, "REDBOOK_WRITE_NOTE", fmt.Sprintf("小红书笔记帮写/优化失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "REDBOOK_WRITE_NOTE", fmt.Sprintf("小红书笔记帮写/优化，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}

// @summary Ai writing little red book planting note
// @description Ai writing little red book planting note
// @tags redbook
// @accept json
// @param param body request.RedbookWriteMessageParam true "Writing instruction"
// @success 200 {object} response.Response{data=*response.RedbookWriteMessageResponse}
// @security user
// @router /redbook/write/planting [post]
func (m *WriteApiRouter) Planting(c *gin.Context) {
	var param request.RedbookWriteMessageParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = userSev.GetClientInfoService().MsgSecCheck(m.GetUserId(c), utils.GetClientType(c), param.Message)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := rbSev.GetWriteService().Planting(m.GetUserId(c), param.Message)
	if err != nil {
		m.UserLog(c, "REDBOOK_WRITE_PLANTING", fmt.Sprintf("小红书种草笔记生成失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "REDBOOK_WRITE_PLANTING", fmt.Sprintf("小红书种草笔记生成，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}
