package helper

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/global"
	"encoding/json"
	"sync"
	"time"
)

type AiService struct {
}

var onceAi = sync.Once{}
var aiService *AiService

// 获取单例
func GetAiService() *AiService {
	onceAi.Do(func() {
		aiService = new(AiService)
	})
	return aiService
}

// AI对话
func (s *AiService) Chat(userId int64, param request.ZhipuAiChatParam) (*response.ZhipuAiChatResponse, error) {
	serializedData, _ := json.Marshal(param)
	var result response.ZhipuAiChatResponse
	message := common.SocketMessage{
		NeedResult: true,
		Timeout:    120 * time.Second,
		Type:       enum.SocketMessageType(10),
		Data:       string(serializedData),
	}
	existResult, err := global.SOCKET.SendMessage(userId, message, &result)
	if err != nil || !existResult {
		return nil, err
	}
	return &result, nil
}
