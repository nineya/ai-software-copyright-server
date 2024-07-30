package log

import (
	"sync"
	"tool-server/internal/application/model/table"
	"tool-server/internal/application/service"
	"tool-server/internal/global"
)

type LogService struct {
	service.CrudService[table.Log]
}

var onceLog = sync.Once{}
var logService *LogService

// 获取单例
func GetLogService() *LogService {
	onceLog.Do(func() {
		logService = new(LogService)
		logService.Db = global.DB
	})
	return logService
}

func (s *LogService) DeleteByAdminId(adminId int64) error {
	_, err := s.WhereAdminSession(adminId).Delete(&table.Log{})
	return err
}
