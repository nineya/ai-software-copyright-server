package redbook

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	rbSev "ai-software-copyright-server/internal/application/service/redbook"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type VisitsApiRouter struct {
	api.BaseApi
}

func (m *VisitsApiRouter) InitVisitsApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("redbook/visits")
	router.POST("", m.Create)
	router.PUT(":id/status/:status", m.UpdateStatusById)
	router.PUT(":id/increase", m.UpdateIncreaseById)
	router.GET("all", m.GetAll)
}

// @summary Create redbook visits task
// @description Create redbook visits task
// @tags redbook
// @accept json
// @param param body table.RedbookVisitsTask true "Redbook visits task information"
// @success 200 {object} response.Response{data=[]table.RedbookVisitsTask}
// @security admin
// @router /redbook/visits [post]
func (m *VisitsApiRouter) Create(c *gin.Context) {
	var param table.RedbookVisitsTask
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := rbSev.GetVisitsService().Create(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	m.AdminLog(c, "CREATED_LINK", fmt.Sprintf("创建友链 %s，友链 Id 为 %d", mod.Name, mod.Id))
	response.OkWithData(mod, c)
}

// @summary Increase redbook visits task current count by id
// @description Increase redbook visits task current count by id
// @tags redbook
// @accept json
// @param id path int64 true "Redbook visits task id"
// @success 200 {object} response.Response{}
// @security admin
// @router /redbook/visits/{id}/increase [put]
func (m *VisitsApiRouter) UpdateIncreaseById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	err = rbSev.GetVisitsService().UpdateIncreaseById(m.GetUserId(c), id)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.Ok(c)
}

// @summary Update redbook visits task status by id
// @description Update redbook visits task status by id
// @tags redbook
// @accept json
// @param id path int64 true "Redbook visits task id"
// @param status path string true "Redbook visits task new status"
// @success 200 {object} response.Response{}
// @security admin
// @router /redbook/visits/{id}/status/{status} [put]
func (m *VisitsApiRouter) UpdateStatusById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	status, err := enum.TaskStatusValue(c.Param("status"))
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	err = rbSev.GetVisitsService().UpdateStatusById(m.GetUserId(c), id, status)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	m.AdminLog(c, "REDBOOK_UPDATE_COOKIE", fmt.Sprintf("修改 Id 为 %d 的小红书浏览量任务状态为 %s", id, enum.TASK_STATUS[status]))
	response.Ok(c)
}

// @summary Get all redbook visits task
// @description Get all redbook visits task
// @tags redbook
// @accept json
// @success 200 {object} response.Response{data=[]table.RedbookVisitsTask}
// @security admin
// @router /redbook/visits/all [get]
func (m *VisitsApiRouter) GetAll(c *gin.Context) {
	mod, err := rbSev.GetVisitsService().GetAll(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
