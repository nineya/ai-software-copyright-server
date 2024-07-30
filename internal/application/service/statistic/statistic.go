package statistic

import (
	"sync"
	"tool-server/internal/application/model/table"
	"tool-server/internal/application/service"
	"tool-server/internal/global"
)

type StatisticService struct {
	service.CrudService[table.Statistic]
}

var onceUser = sync.Once{}
var statisticService *StatisticService

// 获取单例
func GetStatisticService() *StatisticService {
	onceUser.Do(func() {
		statisticService = new(StatisticService)
		statisticService.Db = global.DB
	})
	return statisticService
}

// 清理请求统计表一个月以前的数据
func (s *StatisticService) ClearStatistic() (int64, error) {
	return s.Db.Where("create_time < DATE_SUB(NOW(), INTERVAL 1 MONTH)").Delete(&table.Statistic{})
}
