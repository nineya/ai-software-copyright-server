package public

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	aiSev "ai-software-copyright-server/internal/application/service/ai"
	"github.com/gin-gonic/gin"
)

type AiApiRouter struct {
	api.BaseApi
}

func (m *AiApiRouter) InitAiApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("ai")
	router.POST("write", m.Write)
}

// @summary Ai写作
// @description Ai写作
// @tags public,ai
// @accept json
// @param param body request.AiWriteMessageParam true "Ai写作命令"
// @success 200 {object} response.Response{data=string}
// @router /public/ai/write [post]
func (m *AiApiRouter) Write(c *gin.Context) {
	var param request.AiWriteMessageParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	message, err := aiSev.GetAiService().Write(param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(message, c)
}
