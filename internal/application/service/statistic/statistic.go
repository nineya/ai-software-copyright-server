package statistic

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"github.com/gin-gonic/gin"
	"sync"
)

type StatisticService struct {
	service.AdminCrudService[table.Statistic]
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

func (s *StatisticService) Create(c *gin.Context) error {
	statistic := table.Statistic{
		Url:        c.Request.RequestURI,
		IpAddress:  c.ClientIP(),
		Referrer:   c.Request.Referer(),
		Origin:     utils.GetHost(c.Request.Referer()),
		HttpStatus: c.Writer.Status(),
		UserAgent:  c.Request.UserAgent(),
	}
	_, err := global.DB.Insert(&statistic)
	if err != nil {
		global.LOG.Sugar().Warnf("新增访问记录失败: %+v", err)
	}
	return err
}

// 清理请求统计表两个月以前的数据
func (s *StatisticService) ClearStatistic() (int64, error) {
	return s.Db.Where("create_time < DATE_SUB(NOW(), INTERVAL 2 MONTH)").Delete(&table.Statistic{})
}
