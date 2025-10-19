package software_copyright

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	scSev "ai-software-copyright-server/internal/application/service/software_copyright"
	"fmt"
	"github.com/gin-gonic/gin"
)

type SoftwareCopyrightApiRouter struct {
	api.BaseApi
}

func (m *SoftwareCopyrightApiRouter) InitSoftwareCopyrightApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("softwareCopyright")
	m.Router = router
	router.POST("trigger", m.Trigger)
}

// @summary 触发软著生成任务
// @description 触发软著生成任务
// @tags softwareCopyright
// @accept json
// @param param body table.SoftwareCopyright true "创建软著申请信息"
// @success 200 {object} response.Response{data=[]table.SoftwareCopyright}
// @security user
// @router /softwareCopyright [post]
func (m *SoftwareCopyrightApiRouter) Trigger(c *gin.Context) {
	var param request.SCTriggerParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	err = scSev.GetSoftwareCopyrightService().TriggerGenerate(param.Id)
	if err != nil {
		m.AdminLog(c, "SOFTWARE_COPYRIGHT_TRIGGER", fmt.Sprintf("重新触发软著生成任务 %d 失败，原因：%s", param.Id, err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.AdminLog(c, "SOFTWARE_COPYRIGHT_TRIGGER", fmt.Sprintf("重新触发软著生成任务 %d", param.Id))
	response.Ok(c)
}
