package netdisk

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type ResourceSearchService struct {
	service.UserCrudService[table.NetdiskResourceSearch]
}

var onceResourceSearch = sync.Once{}
var resourceSearchService *ResourceSearchService

// 获取单例
func GetResourceSearchService() *ResourceSearchService {
	onceResourceSearch.Do(func() {
		resourceSearchService = new(ResourceSearchService)
		resourceSearchService.Db = global.DB
	})
	return resourceSearchService
}

// 管理员后台分页查询列表
func (s *ResourceSearchService) GetByPage(userId int64, param request.QueryPageParam) (*response.PageResponse, error) {
	session := s.WhereUserSession(userId).Desc("create_time")
	if param.Keyword != "" {
		session.And("keyword like concat('%',?,'%')", param.Keyword)
	}
	list := make([]table.NetdiskResourceSearch, 0)
	return s.HandlePageable(param.PageableParam, &list, session)
}
