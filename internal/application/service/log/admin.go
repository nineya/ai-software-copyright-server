package log

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type AdminService struct {
	service.AdminCrudService[table.AdminLog]
}

var onceAdmin = sync.Once{}
var adminService *AdminService

// 获取单例
func GetAdminService() *AdminService {
	onceAdmin.Do(func() {
		adminService = new(AdminService)
		adminService.Db = global.DB
	})
	return adminService
}

func (s *AdminService) DeleteByAdminId(adminId int64) error {
	_, err := s.WhereAdminSession(adminId).Delete(&table.AdminLog{})
	return err
}
