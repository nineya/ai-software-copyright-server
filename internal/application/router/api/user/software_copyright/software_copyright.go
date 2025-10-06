package software_copyright

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	scSev "ai-software-copyright-server/internal/application/service/software_copyright"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type SoftwareCopyrightApiRouter struct {
	api.BaseApi
}

func (m *SoftwareCopyrightApiRouter) InitSoftwareCopyrightApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("softwareCopyright")
	m.Router = router
	router.POST("", m.Create)
	router.GET("list", m.GetByPage)
	router.GET(":id", m.GetById)
	router.GET("statistic", m.Statistic)
}

// @summary 创建软著申请
// @description 创建软著申请
// @tags softwareCopyright
// @accept json
// @param param body table.SoftwareCopyright true "创建软著申请"
// @success 200 {object} response.Response{data=[]table.SoftwareCopyright}
// @security user
// @router /softwareCopyright [post]
func (m *SoftwareCopyrightApiRouter) Create(c *gin.Context) {
	var param table.SoftwareCopyright
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := scSev.GetSoftwareCopyrightService().Create(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "SOFTWARE_COPYRIGHT_CREATE", fmt.Sprintf("创建软著任务失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "SOFTWARE_COPYRIGHT_CREATE", fmt.Sprintf("创建软著任务 %s", param.Name))
	response.OkWithData(mod, c)
}

// @summary 查询软著申请
// @description 查询软著申请
// @tags softwareCopyright
// @param id path int64 true "软著申请id"
// @success 200 {object} response.Response{data=table.SoftwareCopyright}
// @security admin
// @router /softwareCopyright/{id} [get]
func (m *SoftwareCopyrightApiRouter) GetById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	mod, err := scSev.GetSoftwareCopyrightService().GetById(m.GetUserId(c), id)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}

// @summary 列表分页查询软著申请
// @description 列表分页查询软著申请
// @tags softwareCopyright
// @param param query request.QueryPageParam true "软著申请列表的分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.SoftwareCopyright}}
// @security user
// @router /softwareCopyright/list [get]
func (m *SoftwareCopyrightApiRouter) GetByPage(c *gin.Context) {
	var param request.QueryPageParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	page, err := scSev.GetSoftwareCopyrightService().GetByPage(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}

// @summary 软著申请数量统计
// @description 软著申请数量统计
// @tags softwareCopyright
// @accept json
// @success 200 {object} response.Response{data=table.SoftwareCopyrightStatistic}
// @security user
// @router /softwareCopyright/statistic [get]
func (m *SoftwareCopyrightApiRouter) Statistic(c *gin.Context) {
	mod, err := scSev.GetSoftwareCopyrightService().Statistic(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
