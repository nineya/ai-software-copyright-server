package public

import (
	"ai-software-copyright-server/internal/application/param/response"
	difyPlugin "ai-software-copyright-server/internal/application/plugin/dify"
	"ai-software-copyright-server/internal/application/router/api"
	"github.com/gin-gonic/gin"
)

type DifyApiRouter struct {
	api.BaseApi
}

func (m *DifyApiRouter) InitDifyApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("dify")
	router.POST("chatMessage", m.ChatMessage)
}

// @summary 发送对话消息
// @description 发送对话消息
// @tags public,dify
// @accept json
// @param param body dify.DifyChatMessageParam true "对话消息参数"
// @success 200 {object} response.Response{data=string}
// @router /public/dify/chatMessage [post]
func (m *DifyApiRouter) ChatMessage(c *gin.Context) {
	var param difyPlugin.DifyChatMessageParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	message, err := difyPlugin.GetDifyPlugin().SendChat("app-kPGnBkdf9bSG850c5kgCS3SC", param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(message, c)
}
