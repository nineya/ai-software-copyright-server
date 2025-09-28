package study

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"sync"
)

type ResourceService struct {
	service.BaseService
}

var onceResource = sync.Once{}
var resourceService *ResourceService

// 获取单例
func GetResourceService() *ResourceService {
	onceResource.Do(func() {
		resourceService = new(ResourceService)
		resourceService.Db = global.DB
	})
	return resourceService
}

// 管理员后台分页查询列表
func (s *ResourceService) GetByPage(userId int64, param request.StudyResourceQueryPageParam) (*response.PageResponse, error) {
	session := s.Db.Desc("create_time")
	if param.Type != "" {
		studyType, err := enum.StudyTypeValue(param.Type)
		if err != nil {
			return nil, err
		}
		session = session.And("type = ?", studyType)
	}
	//if param.Status != "" {
	//	netdiskStatus, err := enum.NetdiskStatusValue(param.Status)
	//	if err != nil {
	//		return nil, err
	//	}
	//	session = session.And("status = ?", netdiskStatus)
	//}
	if param.Keyword != "" {
		session.And("name like concat('%',?,'%')", param.Keyword)
	}
	list := make([]table.StudyResource, 0)
	return s.HandlePageable(param.PageableParam, &list, session)
}
