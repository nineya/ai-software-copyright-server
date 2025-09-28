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

type TimeClockApiRouter struct {
	api.BaseApi
}

func (m *TimeClockApiRouter) InitTimeClockApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("timeClock")
	m.Router = router
	router.POST("", m.Create)
	router.DELETE(":id", m.DeleteById)
	router.PUT(":id", m.UpdateById)
	router.GET("list", m.GetByPage)
	router.GET("my", m.GetListByMemberId)
	router.GET("members", m.GetMembersById)
	router.GET("records", m.GetRecordsById)
	router.GET(":id/myInfo", m.GetMyInfoById)
}

// @summary 创建打卡
// @description 创建打卡
// @tags timeClock
// @accept json
// @param param body table.TimeClock true "创建打卡"
// @success 200 {object} response.Response{data=[]table.TimeClock}
// @security user
// @router /timeClock [post]
func (m *TimeClockApiRouter) Create(c *gin.Context) {
	var param table.TimeClock
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := tcSev.GetTimeClockService().Create(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "TIME_CLOCK_CREATE", fmt.Sprintf("创建打卡失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "TIME_CLOCK_CREATE", fmt.Sprintf("创建打卡 %s", param.Name))
	response.OkWithData(mod, c)
}

// @summary 删除打卡
// @description 删除打卡
// @tags timeClock
// @param id path int64 true "打卡id"
// @success 200 {object} response.Response
// @security user
// @router /timeClock/{id} [delete]
func (m *TimeClockApiRouter) DeleteById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	if err = tcSev.GetTimeClockService().DeleteById(m.GetUserId(c), id); err != nil {
		m.UserLog(c, "TIME_CLOCK_DELETE", fmt.Sprintf("删除 Id 为 %d 的打卡失败，原因：%s", id, err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "TIME_CLOCK_DELETE", fmt.Sprintf("删除 Id 为 %d 的打卡", id))
	response.Ok(c)
}

// @summary 更新打卡信息
// @description 更新打卡信息
// @tags timeClock
// @accept json
// @param id path int64 true "打卡id"
// @param param body table.TimeClock true "打卡信息"
// @success 200 {object} response.Response
// @security user
// @router /timeClock/{id}  [put]
func (m *TimeClockApiRouter) UpdateById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	var param table.TimeClock
	err = c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	if err = tcSev.GetTimeClockService().UpdateById(m.GetUserId(c), id, param); err != nil {
		m.UserLog(c, "TIME_CLOCK_UPDATE", fmt.Sprintf("更新网盘资源 %s 失败，原因：%s", param.Name, err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "TIME_CLOCK_UPDATE", fmt.Sprintf("更新网盘资源 %s，资源 Id 为 %d", param.Name, id))
	response.Ok(c)
}

// @summary 列表分页查询打卡列表
// @description 列表分页查询打卡列表
// @tags timeClock
// @param param query request.QueryPageParam true "分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.TimeClock}}
// @security user
// @router /timeClock/list [get]
func (m *TimeClockApiRouter) GetByPage(c *gin.Context) {
	var param request.PageableParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	page, err := tcSev.GetTimeClockService().GetByPage(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}

// @summary 列表分页查询我的打卡列表
// @description 列表分页查询我的打卡列表
// @tags timeClock
// @param param query request.QueryPageParam true "分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.TimeClock}}
// @security user
// @router /timeClock/my [get]
func (m *TimeClockApiRouter) GetListByMemberId(c *gin.Context) {
	var param request.PageableParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	page, err := tcSev.GetTimeClockService().GetListByMemberId(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}

// @summary 列表分页查询打卡成员列表
// @description 列表分页查询打卡成员列表
// @tags timeClock
// @param param query request.QueryPageParam true "分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.TimeClockMember}}
// @security user
// @router /timeClock/members [get]
func (m *TimeClockApiRouter) GetMembersById(c *gin.Context) {
	var param request.TimeClockQueryPageParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	page, err := tcSev.GetTimeClockMemberService().GetListByClockId(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}

// @summary 查询我的打卡信息
// @description 查询我的打卡信息
// @tags timeClock
// @param param query request.QueryPageParam true "分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=table.TimeClock}}
// @security user
// @router /timeClock/{id}/member/myInfo [get]
func (m *TimeClockApiRouter) GetMyInfoById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	page, err := tcSev.GetTimeClockMemberService().GetMyInfoById(m.GetUserId(c), id)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}

// @summary 查询我的打卡记录
// @description 查询我的打卡记录
// @tags timeClock
// @param param query request.QueryPageParam true "分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=table.TimeClockRecord}}
// @security user
// @router /timeClock/records [get]
func (m *TimeClockApiRouter) GetRecordsById(c *gin.Context) {
	var param request.TimeClockQueryPageParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	page, err := tcSev.GetTimeClockRecordService().GetListByClockId(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}
