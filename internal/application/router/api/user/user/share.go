package user

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ShareApiRouter struct {
	api.BaseApi
}

func (m *ShareApiRouter) InitShareApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("user/share")
	m.Router = router
	router.POST("", m.Create)
	router.GET("all", m.GetAll)
	router.GET("statistic", m.Statistic)
}

// @summary 创建分享记录
// @description 创建分享记录
// @tags user,share
// @accept json
// @param param body table.ShareRecord true "分享记录信息"
// @success 200 {object} response.Response{}
// @security user
// @router /user/share [post]
func (m *ShareApiRouter) Create(c *gin.Context) {
	var param table.ShareRecord
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = userSev.GetShareRecordService().Create(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "SHARE_CREATE", fmt.Sprintf("创建分享记录失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "SHARE_CREATE", fmt.Sprintf("创建分享记录 %s", param.ShareUrl))
	response.Ok(c)
}

// @summary 取得全部分享记录
// @description 取得全部分享记录
// @tags user,share
// @accept json
// @success 200 {object} response.Response{data=[]table.ShareRecord}
// @security user
// @router /user/share/all [get]
func (m *ShareApiRouter) GetAll(c *gin.Context) {
	data, err := userSev.GetShareRecordService().GetAll(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(data, c)
}

// @summary 取得分享统计
// @description 取得分享统计
// @tags user,share
// @accept json
// @success 200 {object} response.Response{data=table.ShareStatistic}
// @security user
// @router /user/share/statistic [get]
func (m *ShareApiRouter) Statistic(c *gin.Context) {
	mod, err := userSev.GetShareRecordService().Statistic(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
