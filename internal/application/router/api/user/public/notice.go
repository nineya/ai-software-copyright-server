package public

import (
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	notSev "ai-software-copyright-server/internal/application/service/notice"
	"github.com/gin-gonic/gin"
)

type NoticeApiRouter struct {
	api.BaseApi
}

func (m *NoticeApiRouter) InitNoticeApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("notice")
	router.GET("platform", m.GetPlatform)
}

// @summary Get platform notice
// @description Get platform notice
// @tags public,notice
// @accept json
// @success 200 {object} response.Response{data=table.Notice}
// @router /public/notice/platform [get]
func (m *NoticeApiRouter) GetPlatform(c *gin.Context) {
	mod, err := notSev.GetNoticeService().GetPlatform(m.GetClientType(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
