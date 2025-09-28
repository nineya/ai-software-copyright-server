package netdisk

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type SearchAppConfigureService struct {
	service.UserCrudService[table.NetdiskSearchAppConfigure]
}

var onceSearchAppConfigure = sync.Once{}
var searchAppConfigureService *SearchAppConfigureService

// 获取单例
func GetSearchAppConfigureService() *SearchAppConfigureService {
	onceSearchAppConfigure.Do(func() {
		searchAppConfigureService = new(SearchAppConfigureService)
		searchAppConfigureService.Db = global.DB
	})
	return searchAppConfigureService
}

// 分页查询列表
func (s *SearchAppConfigureService) SaveConfigure(userId int64, param table.NetdiskSearchAppConfigure) error {
	mod := &table.NetdiskSearchAppConfigure{UserId: userId}
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

func (s *SearchAppConfigureService) GetByUserId(userId int64) (table.NetdiskSearchAppConfigure, error) {
	mod := &table.NetdiskSearchAppConfigure{UserId: userId}
	_, err := s.WhereUserSession(userId).Get(mod)
	if len(mod.WelfareConfig) == 0 {
		mod.WelfareConfig = global.DefaultWelfate
	}
	return *mod, err
}
