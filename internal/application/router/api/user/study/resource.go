package study

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	stuSev "ai-software-copyright-server/internal/application/service/study"
	"github.com/gin-gonic/gin"
)

type ResourceApiRouter struct {
	api.BaseApi
}

func (m *ResourceApiRouter) InitResourceApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("study/resource")
	m.Router = router
	router.GET("list", m.GetByPage)
}

// @summary 列表分页查询学习资源
// @description 列表分页查询学习资源
// @tags study
// @param param query request.StudyResourceQueryPageParam true "学习资源列表的分页查询信息"
// @success 200 {object} response.Response{data=response.PageResponse{content=[]table.StudyResource}}
// @security user
// @router /study/resource/list [get]
func (m *ResourceApiRouter) GetByPage(c *gin.Context) {
	var param request.StudyResourceQueryPageParam
	err := c.ShouldBindQuery(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	page, err := stuSev.GetResourceService().GetByPage(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(page, c)
}
