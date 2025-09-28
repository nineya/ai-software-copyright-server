package credits

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"fmt"
	"github.com/gin-gonic/gin"
)

type CreditsOrderApiRouter struct {
	api.BaseApi
}

func (m *CreditsOrderApiRouter) InitCreditsOrderApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("credits/order")
	router.POST("create", m.Create)
	router.GET("list", m.GetByPage)
}

// @summary 创建付款订单
// @description 创建付款订单
// @tags credits
// @accept json
// @success 200 {object} response.Response{data=table.CreditsOrder}
// @security user
// @router /credits/order/create [post]
func (m *CreditsOrderApiRouter) Create(c *gin.Context) {
	var param request.CreditsCreateOrderParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := userSev.GetCreditsOrderService().CreateOrder(m.GetUserId(c), m.GetClientType(c), param.CreditsPriceId)
	if err != nil {
		m.UserLog(c, "POINTS_ORDER_CREATE", fmt.Sprintf("创建付款订单失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "POINTS_ORDER_CREATE", fmt.Sprintf("创建付款订单：%s，金额：%s, 积分数量：%d", mod.TradeNo, mod.OrderAmount, mod.CreditsNum))
	response.OkWithData(mod, c)
}

// @summary 列表分页查询充值订单
// @description 列表分页查询充值订单
// @tags form
// @param param query request.PageableParam true "订单列表的分页信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.FormData}}
// @security user
// @router /credits/order/list [get]
func (m *CreditsOrderApiRouter) GetByPage(c *gin.Context) {
	var param request.PageableParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	page, err := userSev.GetCreditsOrderService().GetByPage(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}
