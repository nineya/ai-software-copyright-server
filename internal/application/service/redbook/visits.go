package redbook

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type VisitsService struct {
	service.AdminCrudService[table.RedbookVisitsTask]
}

var onceVisits = sync.Once{}
var visitsService *VisitsService

// 获取单例
func GetVisitsService() *VisitsService {
	onceVisits.Do(func() {
		visitsService = new(VisitsService)
		visitsService.Db = global.DB
	})
	return visitsService
}

func (s *VisitsService) UpdateIncreaseById(siteId int64, id int64) error {
	_, err := s.WhereAndOmitAdminSession(siteId).ID(id).Incr("current_count", 1).NoAutoTime().Update(&table.RedbookVisitsTask{})
	return err
}

func (s *VisitsService) UpdateStatusById(adminId int64, id int64, status enum.TaskStatus) error {
	_, err := s.WhereAndOmitAdminSession(adminId).ID(id).Update(&table.RedbookVisitsTask{Status: status})
	return err
}
