package admin

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type AdminService struct {
	service.BaseService
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

func (s *AdminService) GetById(id int64) (table.Admin, error) {
	mod := &table.Admin{}
	_, err := s.Db.ID(id).Get(mod)
	return *mod, err
}
