package mail

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/global"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
	"mime"
	"path"
	"sync"
)

type MailPlugin struct {
	MailDialer *gomail.Dialer
	From       string
}

var onceMail = sync.Once{}
var mailService *MailPlugin

// 获取单例
func GetMailPlugin() *MailPlugin {
	onceMail.Do(func() {
		mailService = new(MailPlugin)
		mailService.MailDialer = &gomail.Dialer{
			Host:     global.CONFIG.Plugin.Mail.Host,
			Port:     global.CONFIG.Plugin.Mail.Port,
			Username: global.CONFIG.Plugin.Mail.Username,
			Password: global.CONFIG.Plugin.Mail.Password,
			SSL:      true,
		}
		mailService.From = fmt.Sprintf("%s<%s>", mime.QEncoding.Encode("UTF-8", path.Base(global.CONFIG.Plugin.Mail.From)), global.CONFIG.Plugin.Mail.Username)
	})
	return mailService
}

func (s *MailPlugin) SendHtmlMail(param request.MailParam) error {
	if s.MailDialer == nil {
		return errors.New("您没有配置邮箱信息")
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.From)
	msg.SetHeader("To", param.To)
	msg.SetHeader("Subject", param.Subject)
	msg.SetBody("text/html", param.Content)
	return s.MailDialer.DialAndSend(msg)
}
