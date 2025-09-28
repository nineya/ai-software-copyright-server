package ai

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	zhipuPlugin "ai-software-copyright-server/internal/application/plugin/zhipu_ai"
	"ai-software-copyright-server/internal/application/router/api"
	helperSev "ai-software-copyright-server/internal/application/service/helper"
	"github.com/gin-gonic/gin"
)

type AiApiRouter struct {
	api.BaseApi
}

func (m *AiApiRouter) InitAiApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("ai")
	router.POST("chat", m.Chat)
}

// @summary AI对话交互
// @description AI对话交互
// @tags ai
// @accept json
// @param param body request.ZhipuAiChatParam true "对话请求内容"
// @success 200 {object} response.Response{data=*response.ZhipuAiChatResponse}
// @router /ai/chat [post]
func (m *AiApiRouter) Chat(c *gin.Context) {
	var msgParam request.ZhipuAiChatParam
	err := c.ShouldBindJSON(&msgParam)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	param := zhipuPlugin.GetDefaultChatParam()
	param.Messages = msgParam.Messages

	mod, err := helperSev.GetAiService().Chat(m.GetUserId(c), param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
