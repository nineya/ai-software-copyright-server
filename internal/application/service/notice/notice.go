package notice

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type NoticeService struct {
	service.UserCrudService[table.Notice]
}

var onceNotice = sync.Once{}
var noticeService *NoticeService

// 获取单例
func GetNoticeService() *NoticeService {
	onceNotice.Do(func() {
		noticeService = new(NoticeService)
		noticeService.Db = global.DB
	})
	return noticeService
}

func (s *NoticeService) GetPlatform(clientType enum.ClientType) (*table.Notice, error) {
	mod := &table.Notice{ClientType: clientType}
	_, err := s.Db.Get(mod)
	return mod, err
}
