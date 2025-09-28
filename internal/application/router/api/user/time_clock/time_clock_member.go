package time_clock

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	tcSev "ai-software-copyright-server/internal/application/service/time_clock"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type TimeClockMemberApiRouter struct {
	api.BaseApi
}

func (m *TimeClockMemberApiRouter) InitTimeClockMemberApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("timeClock/member")
	m.Router = router
	router.POST("", m.Create)
	router.DELETE(":id", m.DeleteById)
	router.PUT("audit", m.AuditById)
	router.GET("list", m.GetByPage)
}

// @summary 创建打卡成员
// @description 创建打卡成员
// @tags timeClock
// @accept json
// @param param body table.TimeClock true "创建打卡成员"
// @success 200 {object} response.Response{data=[]table.TimeClockMember}
// @security user
// @router /timeClock/member [post]
func (m *TimeClockMemberApiRouter) Create(c *gin.Context) {
	var param table.TimeClockMember
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	param.Status = enum.TimeClockMemberStatus(2)
	mod, err := tcSev.GetTimeClockMemberService().Create(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "TIME_CLOCK_MEMBER_CREATE", fmt.Sprintf("创建打卡成员失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "TIME_CLOCK_MEMBER_CREATE", fmt.Sprintf("创建打卡成员 %s", param.UserId))
	response.OkWithData(mod, c)
}

// @summary 审核成员
// @description 审核成员
// @tags timeClock
// @param id path int64 true "打卡id"
// @success 200 {object} response.Response
// @security user
// @router /timeClock/member/audit [put]
func (m *TimeClockMemberApiRouter) AuditById(c *gin.Context) {
	var param table.TimeClockMember
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	if err = tcSev.GetTimeClockMemberService().AuditById(m.GetUserId(c), param); err != nil {
		m.UserLog(c, "TIME_CLOCK_MEMBER_AUDIT", fmt.Sprintf("审核 Id 为 %d 的打卡成员失败，原因：%s", param.Id, err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "TIME_CLOCK_MEMBER_AUDIT", fmt.Sprintf("审核 Id 为 %d 的打卡成员", param.Id))
	response.Ok(c)
}

// @summary 删除打卡成员
// @description 删除打卡成员
// @tags timeClock
// @param id path int64 true "打卡id"
// @success 200 {object} response.Response
// @security user
// @router /timeClock/member/{id} [delete]
func (m *TimeClockMemberApiRouter) DeleteById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	if err = tcSev.GetTimeClockMemberService().DeleteById(m.GetUserId(c), id); err != nil {
		m.UserLog(c, "TIME_CLOCK_MEMBER_DELETE", fmt.Sprintf("删除 Id 为 %d 的打卡失败，原因：%s", id, err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "TIME_CLOCK_MEMBER_DELETE", fmt.Sprintf("删除 Id 为 %d 的打卡", id))
	response.Ok(c)
}

// @summary 列表分页查询打卡成员列表
// @description 列表分页查询打卡成员列表
// @tags timeClock
// @param param query request.QueryPageParam true "分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.TimeClockMember}}
// @security user
// @router /timeClock/member/list [get]
func (m *TimeClockMemberApiRouter) GetByPage(c *gin.Context) {
	var param request.PageableParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	page, err := tcSev.GetTimeClockMemberService().GetByPage(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}
