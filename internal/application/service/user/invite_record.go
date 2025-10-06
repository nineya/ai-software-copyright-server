package user

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type InviteRecordService struct {
	service.UserCrudService[table.InviteRecord]
}

var onceInviteRecord = sync.Once{}
var inviteRecordService *InviteRecordService

// 获取单例
func GetInviteRecordService() *InviteRecordService {
	onceInviteRecord.Do(func() {
		inviteRecordService = new(InviteRecordService)
		inviteRecordService.Db = global.DB
	})
	return inviteRecordService
}

func (s *InviteRecordService) Statistic(userId int64) (*table.InviteStatistic, error) {
	mod := &table.InviteStatistic{}
	_, err := s.WhereUserSession(userId).Select("COUNT(*) AS total_count, sum(CASE WHEN type = 1 THEN reward_coin END) AS invite_credits, count(CASE WHEN type = 1 and create_time >= DATE_SUB(NOW(), INTERVAL 1 MONTH) THEN 1 END) AS month_count").Get(mod)
	return mod, err
}
