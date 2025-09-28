package redbook

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	rbSev "ai-software-copyright-server/internal/application/service/redbook"
	"ai-software-copyright-server/internal/global"
	"fmt"
	"github.com/gin-gonic/gin"
)

type RedbookApiRouter struct {
	api.BaseApi
}

func (m *RedbookApiRouter) InitRedbookApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("redbook")
	m.Router = router
	router.POST("removeWatermark", m.RemoveWatermark)
	router.POST("valuation", m.Valuation)
	router.POST("weight", m.Weight)
}

// @summary Redbook remove watermark
// @description Redbook remove watermark
// @tags redbook
// @accept json
// @param param body request.RedbookParam true "Redbook home address"
// @success 200 {object} response.Response{data=*response.RedbookRemoveWatermarkResponse}
// @security user
// @router /redbook/removeWatermark [post]
func (m *RedbookApiRouter) RemoveWatermark(c *gin.Context) {
	var param request.RedbookParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := rbSev.GetRedbookService().RemoveWatermark(m.GetUserId(c), param.Url)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("去除图文水印失败: %+v", err))
		m.UserLog(c, "REDBOOK_REMOVE_WATERMARK", fmt.Sprintf("小红书去水印失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		//response.FailWithError(errors.New("服务维护中"), c)
		return
	}
	m.UserLog(c, "REDBOOK_REMOVE_WATERMARK", fmt.Sprintf("小红书去水印，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}

// @summary Redbook account valuation
// @description Redbook account valuation
// @tags redbook
// @accept json
// @param param body request.RedbookParam true "Redbook home address"
// @success 200 {object} response.Response{data=*response.RedbookValuationResponse}
// @security user
// @router /redbook/valuation [post]
func (m *RedbookApiRouter) Valuation(c *gin.Context) {
	var param request.RedbookParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := rbSev.GetRedbookService().Valuation(m.GetUserId(c), param)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("账号估值失败: %+v", err))
		m.UserLog(c, "REDBOOK_VALUATION", fmt.Sprintf("小红书账号估值，原因：%s", err.Error()))
		response.FailWithError(err, c)
		//response.FailWithError(errors.New("服务维护中"), c)
		return
	}
	m.UserLog(c, "REDBOOK_VALUATION", fmt.Sprintf("小红书账号估值，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}

// @summary Redbook account weight
// @description Redbook account weight
// @tags redbook
// @accept json
// @param param body request.RedbookParam true "Redbook home address"
// @success 200 {object} response.Response{data=*response.RedbookWeightResponse}
// @security user
// @router /redbook/weight [post]
func (m *RedbookApiRouter) Weight(c *gin.Context) {
	var param request.RedbookParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := rbSev.GetRedbookService().Weight(m.GetUserId(c), param)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("账号权重分析失败: %+v", err))
		m.UserLog(c, "REDBOOK_WEIGHT", fmt.Sprintf("小红书账号权重分析，原因：%s", err.Error()))
		response.FailWithError(err, c)
		//response.FailWithError(errors.New("服务维护中"), c)
		return
	}
	m.UserLog(c, "REDBOOK_WEIGHT", fmt.Sprintf("小红书账号权重分析，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}
