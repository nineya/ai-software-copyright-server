package public

import (
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	creditsSev "ai-software-copyright-server/internal/application/service/credits"
	"github.com/gin-gonic/gin"
)

type CreditsApiRouter struct {
	api.BaseApi
}

func (m *CreditsApiRouter) InitCreditsApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("credits")
	router.GET("allPrice", m.GetAllPrice)
}

// @summary 取得全部币售价
// @description 取得全部币售价
// @tags credits
// @accept json
// @success 200 {object} response.Response{data=[]table.CreditsPrice}
// @router /public/credits/allPrice [get]
func (m *CreditsApiRouter) GetAllPrice(c *gin.Context) {
	mod, err := creditsSev.GetCreditsPriceService().GetAll()
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
