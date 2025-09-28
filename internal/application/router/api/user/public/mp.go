package public

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	mpSev "ai-software-copyright-server/internal/application/service/mp"
	"ai-software-copyright-server/internal/global"
	"github.com/gin-gonic/gin"
	"io"
)

type MpApiRouter struct {
	api.BaseApi
}

func (m *MpApiRouter) InitMpApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("mp")
	m.Router = router
	router.GET("notice", m.Verify)
	router.POST("notice", m.Notice)
}

// @summary 校验服务器
// @description 校验服务器
// @tags mp
// @accept json
// @param param quary request.MpVerifyParam true "校验服务器参数"
// @success 200 {object} response.Response{data=*response.UserBuyContentResponse}
// @security user
// @router /public/mp/notice [get]
func (m *MpApiRouter) Verify(c *gin.Context) {
	var param request.MpVerifyParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod := mpSev.GetMpService().Verify(param)
	c.Writer.Write([]byte(mod))
}

// @summary 消息通知
// @description 消息通知
// @tags mp
// @accept json
// @param param quary request.MpVerifyParam true "消息通知参数"
// @success 200 {object} response.Response{data=*response.UserBuyContentResponse}
// @security user
// @router /public/mp/notice [post]
func (m *MpApiRouter) Notice(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	global.LOG.Info("收到消息：" + string(body))
}
