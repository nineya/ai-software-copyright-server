package public

import (
	"ai-software-copyright-server/internal/application/router/api"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"net/http"
)

type WxNotifyApiRouter struct {
	api.BaseApi
}

func (m *WxNotifyApiRouter) InitWxNotifyApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("wxNotify")
	router.POST("pay", m.PayNotify)
}

// @summary 微信支付回调
// @description 微信支付回调
// @tags credits
// @accept json
// @success 200 {object} &wechat.V3NotifyRsp{}
// @router /public/wxNotify/pay [post]
func (m *WxNotifyApiRouter) PayNotify(c *gin.Context) {
	err := userSev.GetCreditsOrderService().PayNotify(c.Request)
	if err != nil {
		m.UserLog(c, "WX_PAY_NOTIFY", fmt.Sprintf("处理微信支付回调失败，原因：%s", err.Error()))
		c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: err.Error()})
		return
	}
	m.UserLog(c, "WX_PAY_NOTIFY", "处理微信支付回调成功")
	c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"})
}
