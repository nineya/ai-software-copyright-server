package ai

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/plugin/zhipu_ai"
	"sync"
)

type AiService struct {
}

var onceWrite = sync.Once{}
var aiService *AiService

// 获取单例
func GetAiService() *AiService {
	onceWrite.Do(func() {
		aiService = new(AiService)
	})
	return aiService
}

// ai写作
func (s *AiService) Write(param request.AiWriteMessageParam) (string, error) {
	chatParam := zhipu_ai.GetDefaultChatParam()
	chatParam.Messages = []request.ZhipuAiChatMessageItem{
		{
			Role:    "system",
			Content: param.System,
		},
		{
			Role:    "user",
			Content: param.User,
		},
	}
	zhipuResult, err := zhipu_ai.GetZhipuAiPlugin().SendChat(chatParam)
	if err != nil {
		return "", err
	}
	return zhipuResult.Choices[0].Message.Content, nil
}
