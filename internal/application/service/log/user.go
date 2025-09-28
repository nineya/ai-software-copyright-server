package log

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type UserService struct {
	service.UserCrudService[table.UserLog]
}

var onceUser = sync.Once{}
var userService *UserService

// 获取单例
func GetUserService() *UserService {
	onceUser.Do(func() {
		userService = new(UserService)
		userService.Db = global.DB
	})
	return userService
}

func (s *UserService) DeleteByUserId(adminId int64) error {
	_, err := s.WhereUserSession(adminId).Delete(&table.UserLog{})
	return err
}
