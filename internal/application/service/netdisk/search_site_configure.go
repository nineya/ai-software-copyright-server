package netdisk

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type SearchSiteConfigureService struct {
	service.UserCrudService[table.NetdiskSearchSiteConfigure]
}

var onceSearchSiteConfigure = sync.Once{}
var searchSiteConfigureService *SearchSiteConfigureService

// 获取单例
func GetSearchSiteConfigureService() *SearchSiteConfigureService {
	onceSearchSiteConfigure.Do(func() {
		searchSiteConfigureService = new(SearchSiteConfigureService)
		searchSiteConfigureService.Db = global.DB
	})
	return searchSiteConfigureService
}

// 分页查询列表
func (s *SearchSiteConfigureService) SaveConfigure(userId int64, param table.NetdiskSearchSiteConfigure) error {
	mod := &table.NetdiskSearchSiteConfigure{UserId: userId}
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

func (s *SearchSiteConfigureService) GetByUserId(userId int64) (table.NetdiskSearchSiteConfigure, error) {
	mod := &table.NetdiskSearchSiteConfigure{UserId: userId}
	_, err := s.WhereUserSession(userId).Get(mod)
	return *mod, err
}
