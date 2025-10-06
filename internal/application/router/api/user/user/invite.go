package user

import (
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"github.com/gin-gonic/gin"
)

type InviteApiRouter struct {
	api.BaseApi
}

func (m *InviteApiRouter) InitInviteApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("user/invite")
	m.Router = router
	router.GET("all", m.GetAll)
	router.GET("statistic", m.Statistic)
}

// @summary 取得全部邀请记录
// @description 取得全部邀请记录
// @tags user
// @accept json
// @success 200 {object} response.Response{data=[]table.InviteRecord}
// @security user
// @router /user/invite/all [get]
func (m *InviteApiRouter) GetAll(c *gin.Context) {
	data, err := userSev.GetInviteRecordService().GetAll(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(data, c)
}

// @summary 取得邀请统计
// @description 取得邀请统计
// @tags user
// @accept json
// @success 200 {object} response.Response{data=table.InviteStatistic}
// @security user
// @router /user/invite/statistic [get]
func (m *InviteApiRouter) Statistic(c *gin.Context) {
	mod, err := userSev.GetInviteRecordService().Statistic(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
