package public

import (
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	qrcodeSev "ai-software-copyright-server/internal/application/service/qrcode"
	"github.com/gin-gonic/gin"
)

type QrcodeApiRouter struct {
	api.BaseApi
}

func (m *QrcodeApiRouter) InitQrcodeApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("qrcode")
	m.Router = router
	router.GET("loose/:alias", m.GetByAlias)
}

// @summary 取得活码
// @description 取得活码
// @tags qrcode
// @accept json
// @param id path int64 true "活码别名"
// @success 200 {object} response.Response{data=table.Qrcode}
// @security user
// @router /qrcode/loose/{alias} [get]
func (m *QrcodeApiRouter) GetByAlias(c *gin.Context) {
	alias := c.Param("alias")
	if alias == "" {
		response.ForbiddenWithMessage("活码不存在", c)
		return
	}
	mod, err := qrcodeSev.GetQrcodeService().GetByAlias(alias)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
