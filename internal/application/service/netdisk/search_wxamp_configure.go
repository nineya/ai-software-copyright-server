package netdisk

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type SearchWxampConfigureService struct {
	service.UserCrudService[table.NetdiskSearchWxampConfigure]
}

var onceSearchWxampConfigure = sync.Once{}
var searchWxampConfigureService *SearchWxampConfigureService

// 获取单例
func GetSearchWxampConfigureService() *SearchWxampConfigureService {
	onceSearchWxampConfigure.Do(func() {
		searchWxampConfigureService = new(SearchWxampConfigureService)
		searchWxampConfigureService.Db = global.DB
	})
	return searchWxampConfigureService
}

// 分页查询列表
func (s *SearchWxampConfigureService) SaveConfigure(userId int64, param table.NetdiskSearchWxampConfigure) error {
	mod := &table.NetdiskSearchWxampConfigure{UserId: userId}
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

func (s *SearchWxampConfigureService) GetByUserId(userId int64) (table.NetdiskSearchWxampConfigure, error) {
	mod := &table.NetdiskSearchWxampConfigure{UserId: userId}
	_, err := s.WhereUserSession(userId).Get(mod)
	if len(mod.WelfareConfig) == 0 {
		mod.WelfareConfig = global.DefaultWelfate
	}
	return *mod, err
}
