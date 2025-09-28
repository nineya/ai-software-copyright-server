package user

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"strconv"
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

func (s *InviteRecordService) GetInviteInfo(userId int64) (*table.InviteInfo, error) {
	mod := &table.InviteInfo{}
	results, err := s.Db.QueryString(`select * from 
(select count(*) invite_num, sum(reward_credits) invite_credits from invite_record where user_id = ? and type = 1) a,
(select sum(reward_credits) active_credits from invite_record where user_id = ? and type = 2) b`, userId, userId)
	if err != nil {
		return nil, err
	}
	mod.InviteNum, _ = strconv.Atoi(results[0]["invite_num"])
	mod.InviteCredits, _ = strconv.Atoi(results[0]["invite_credits"])
	mod.ActiveCredits, _ = strconv.Atoi(results[0]["active_credits"])
	return mod, err
}
