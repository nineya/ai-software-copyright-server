package redbook

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	rbSev "ai-software-copyright-server/internal/application/service/redbook"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ProhibitedApiRouter struct {
	api.BaseApi
}

func (m *ProhibitedApiRouter) InitProhibitedApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("redbook/prohibited")
	m.Router = router
	router.POST("detection", m.Detection)
	router.POST("customProhibited", m.UpdateCustomProhibited)
	router.GET("customProhibited", m.GetCustomProhibited)
}

// @summary Detection redbook prohibited
// @description Detection redbook prohibited
// @tags redbook
// @accept json
// @param param body request.RedbookProhibitedDetectionParam true "Copywriter information"
// @success 200 {object} response.Response{}
// @security user
// @router /redbook/prohibited/detection [post]
func (m *ProhibitedApiRouter) Detection(c *gin.Context) {
	var param request.RedbookProhibitedDetectionParam
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

	mod, err := rbSev.GetProhibitedService().Detection(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "REDBOOK_PROHIBITED_DETECTION", fmt.Sprintf("小红书敏感词检测失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "REDBOOK_PROHIBITED_DETECTION", fmt.Sprintf("小红书敏感词检测，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}

// @summary Update my prohibited
// @description Update my prohibited
// @tags redbook
// @accept json
// @param param body table.RedbookProhibited true "Prohibited information"
// @success 200 {object} response.Response{}
// @security user
// @router /redbook/prohibited/customProhibited [put]
func (m *ProhibitedApiRouter) UpdateCustomProhibited(c *gin.Context) {
	var param table.RedbookProhibited
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = userSev.GetClientInfoService().MsgSecCheck(m.GetUserId(c), utils.GetClientType(c), utils.ListJoin(param.Words, ",", func(index int, item string) string {
		return item
	}))
	if err != nil {
		response.FailWithError(err, c)
		return
	}

	err = rbSev.GetProhibitedService().UpdateCustomProhibited(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "REDBOOK_PROHIBITED_CUSTOM", fmt.Sprintf("更新小红书用户自定义敏感词失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "REDBOOK_PROHIBITED_CUSTOM", fmt.Sprintf("更新小红书用户自定义敏感词，共计 %d 条", len(param.Words)))
	response.Ok(c)
}

// @summary Get my prohibited
// @description Get my prohibited
// @tags redbook
// @accept json
// @success 200 {object} response.Response{data=table.RedbookProhibited}
// @security user
// @router /redbook/prohibited/customProhibited [get]
func (m *ProhibitedApiRouter) GetCustomProhibited(c *gin.Context) {
	mod, err := rbSev.GetProhibitedService().GetCustomProhibited(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
