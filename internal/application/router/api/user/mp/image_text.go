package mp

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	mpSev "ai-software-copyright-server/internal/application/service/mp"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ImageTextApiRouter struct {
	api.BaseApi
}

func (m *ImageTextApiRouter) InitImageTextApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("mp/imageText")
	m.Router = router
	router.POST("optimize", m.Optimize)
}

// @summary AI代写和优化图文文章
// @description AI代写和优化图文文章
// @tags mp
// @accept json
// @param param body request.MpImageTextMessageParam true "写作指示"
// @success 200 {object} response.Response{data=*response.UserBuyContentResponse}
// @security user
// @router /mp/imageText/optimize [post]
func (m *ImageTextApiRouter) Optimize(c *gin.Context) {
	var param request.MpImageTextMessageParam
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
	mod, err := mpSev.GetImageTextService().Optimize(m.GetUserId(c), param.Message)
	if err != nil {
		m.UserLog(c, "MP_IMAGE_TEXT_OPTIMIZE", fmt.Sprintf("公众号图文优化失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "MP_IMAGE_TEXT_OPTIMIZE", fmt.Sprintf("公众号图文优化，花费：%d，剩余：%d", mod.BuyCredits, mod.BalanceCredits))
	response.OkWithData(mod, c)
}
