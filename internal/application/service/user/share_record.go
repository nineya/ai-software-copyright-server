package user

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type ShareRecordService struct {
	service.UserCrudService[table.ShareRecord]
}

var onceShareRecord = sync.Once{}
var shareRecordService *ShareRecordService

// 获取单例
func GetShareRecordService() *ShareRecordService {
	onceShareRecord.Do(func() {
		shareRecordService = new(ShareRecordService)
		shareRecordService.Db = global.DB
	})
	return shareRecordService
}

func (s *ShareRecordService) Statistic(userId int64) (*table.ShareStatistic, error) {
	mod := &table.ShareStatistic{}
	_, err := s.WhereUserSession(userId).Select("sum(reward_credits) AS share_credits, count(CASE WHEN status = 1 THEN 1 END) AS await_count, count(CASE WHEN status = 2 THEN 1 END) AS pass_count").Get(mod)
	return mod, err
}
