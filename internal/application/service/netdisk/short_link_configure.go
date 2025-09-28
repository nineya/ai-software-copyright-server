package netdisk

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type ShortLinkConfigureService struct {
	service.UserCrudService[table.NetdiskShortLinkConfigure]
}

var onceShortLinkConfigure = sync.Once{}
var shortLinkConfigureService *ShortLinkConfigureService

// 获取单例
func GetShortLinkConfigureService() *ShortLinkConfigureService {
	onceShortLinkConfigure.Do(func() {
		shortLinkConfigureService = new(ShortLinkConfigureService)
		shortLinkConfigureService.Db = global.DB
	})
	return shortLinkConfigureService
}

// 分页查询列表
func (s *ShortLinkConfigureService) SaveConfigure(userId int64, param table.NetdiskShortLinkConfigure) error {
	mod := &table.NetdiskShortLinkConfigure{UserId: userId}
	exist, err := s.Db.Get(mod)
	if err != nil {
		return err
	}
	param.UserId = userId
	if exist {
		_, err = s.WhereUserSession(userId).AllCols().Update(&param)
	} else {
		_, err = s.Db.Insert(&param)
	}
	return err
}

func (s *ShortLinkConfigureService) GetByUserId(userId int64) (table.NetdiskShortLinkConfigure, error) {
	mod := &table.NetdiskShortLinkConfigure{UserId: userId}
	_, err := s.WhereUserSession(userId).Get(mod)
	return *mod, err
}
