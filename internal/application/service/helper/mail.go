package helper

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/global"
	"encoding/json"
	"sync"
)

type MailService struct {
}

var onceMail = sync.Once{}
var mailService *MailService

// 获取单例
func GetMailService() *MailService {
	onceMail.Do(func() {
		mailService = new(MailService)
	})
	return mailService
}

// AI对话
func (s *MailService) SendHtml(userId int64, param request.MailParam) error {
	serializedData, _ := json.Marshal(param)
	message := common.SocketMessage{
		NeedResult: false,
		Type:       enum.SocketMessageType(11),
		Data:       string(serializedData),
	}
	_, err := global.SOCKET.SendMessage(userId, message, nil)
	return err
}
