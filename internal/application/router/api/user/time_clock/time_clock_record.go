package time_clock

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	tcSev "ai-software-copyright-server/internal/application/service/time_clock"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type TimeClockRecordApiRouter struct {
	api.BaseApi
}

func (m *TimeClockRecordApiRouter) InitTimeClockRecordApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("timeClock/record")
	m.Router = router
	router.POST("", m.Create)
	router.DELETE(":id", m.DeleteById)
	router.GET("list", m.GetByPage)
}

// @summary 创建打卡记录
// @description 创建打卡记录
// @tags timeClock
// @accept json
// @param param body table.TimeClock true "创建打卡记录"
// @success 200 {object} response.Response{data=[]table.TimeClockMember}
// @security user
// @router /timeClock/record [post]
func (m *TimeClockRecordApiRouter) Create(c *gin.Context) {
	var param table.TimeClockRecord
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := tcSev.GetTimeClockRecordService().Create(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "TIME_CLOCK_RECORD_CREATE", fmt.Sprintf("创建打卡记录失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "TIME_CLOCK_RECORD_CREATE", fmt.Sprintf("创建打卡记录 %s", param.UserId))
	response.OkWithData(mod, c)
}

// @summary 删除打卡记录
// @description 删除打卡记录
// @tags timeClock
// @param id path int64 true "打卡id"
// @success 200 {object} response.Response
// @security user
// @router /timeClock/record/{id} [delete]
func (m *TimeClockRecordApiRouter) DeleteById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	if err = tcSev.GetTimeClockRecordService().DeleteById(m.GetUserId(c), id); err != nil {
		m.UserLog(c, "TIME_CLOCK_RECORD_DELETE", fmt.Sprintf("删除 Id 为 %d 的打卡记录失败，原因：%s", id, err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "TIME_CLOCK_RECORD_DELETE", fmt.Sprintf("删除 Id 为 %d 的打卡记录", id))
	response.Ok(c)
}

// @summary 列表分页查询打卡记录列表
// @description 列表分页查询打卡记录列表
// @tags timeClock
// @param param query request.QueryPageParam true "分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.TimeClockMember}}
// @security user
// @router /timeClock/record/list [get]
func (m *TimeClockRecordApiRouter) GetByPage(c *gin.Context) {
	var param request.PageableParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	page, err := tcSev.GetTimeClockRecordService().GetByPage(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}
