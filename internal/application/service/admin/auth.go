package admin

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"sync"
)

type AuthService struct {
	service.BaseService
}

var onceAuth = sync.Once{}
var authService *AuthService

// 获取单例
func GetAuthService() *AuthService {
	onceAuth.Do(func() {
		authService = new(AuthService)
		authService.Db = global.DB
	})
	return authService
}

func (s *AuthService) Login(param request.AdminLoginParam) (*common.Token, error) {
	var admin table.Admin
	exist, err := s.Db.Where("username = ?", param.Username).Get(&admin)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("用户不存在或已被禁用")
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), utils.Md5ByBytes(param.Password))
	if err != nil {
		return nil, err
	}
	return utils.AuthToken(admin.Id, global.Admin)
}
