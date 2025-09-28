package time_clock

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"github.com/pkg/errors"
	"sync"
	"time"
	"xorm.io/xorm"
)

type TimeClockRecordService struct {
	service.UserCrudService[table.TimeClockRecord]
}

var onceTimeClockRecord = sync.Once{}
var timeClockRecordService *TimeClockRecordService

// 获取单例
func GetTimeClockRecordService() *TimeClockRecordService {
	onceTimeClockRecord.Do(func() {
		timeClockRecordService = new(TimeClockRecordService)
		timeClockRecordService.Db = global.DB
	})
	return timeClockRecordService
}

func (s *TimeClockRecordService) Create(userId int64, param table.TimeClockRecord) (table.TimeClockRecord, error) {
	mod, err := GetTimeClockService().GetByOnlyId(param.ClockId)
	if err != nil {
		return param, err
	}
	if mod.Id == 0 {
		return param, errors.New("打卡不存在")
	}
	if time.Now().Before(mod.StartTime) || time.Now().After(mod.EndTime) {
		return param, errors.New("当前不在打卡时间")
	}
	param.UserId = userId
	err = s.DbTransaction(func(session *xorm.Session) error {
		exist, err := s.AddWhereUser(userId, session).Where("clock_id = ? and create_time >= CURDATE()", param.ClockId).Exist(&table.TimeClockRecord{})
		if err != nil {
			return err
		}
		if exist {
			return nil
		}
		_, err = session.Insert(&param)
		return err
	})
	return param, err
}

// 分页查询
func (s *TimeClockRecordService) GetByPage(userId int64, param request.PageableParam) (*response.PageResponse, error) {
	session := s.WhereUserSession(userId).Desc("create_time")
	list := make([]table.TimeClockRecord, 0)
	return s.HandlePageable(param, &list, session)
}

// 分页取得当前成员的打卡记录
func (s *TimeClockRecordService) GetListByClockId(userId int64, param request.TimeClockQueryPageParam) (*response.PageResponse, error) {
	session := s.WhereUserSession(userId).Where("clock_id = ?", param.ClockId).Desc("create_time")
	list := make([]table.TimeClockRecord, 0)
	return s.HandlePageable(param.PageableParam, &list, session)
}
