package time_clock

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"fmt"
	"sync"
)

type TimeClockService struct {
	service.UserCrudService[table.TimeClock]
}

var onceTimeClock = sync.Once{}
var timeClockService *TimeClockService

// 获取单例
func GetTimeClockService() *TimeClockService {
	onceTimeClock.Do(func() {
		timeClockService = new(TimeClockService)
		timeClockService.Db = global.DB
	})
	return timeClockService
}

func (s *TimeClockService) Create(userId int64, param table.TimeClock) (*response.TimeClockCreateResponse, error) {
	expenseCredits := 80
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}

	param.SetUserId(userId)
	_, err = s.Db.Insert(&param)
	if err != nil {
		return nil, err
	}

	result := &response.TimeClockCreateResponse{Data: param}
	// 扣款
	user, err := userSev.GetUserService().PaymentNyCredits(userId, enum.BuyType(19), expenseCredits, fmt.Sprintf("购打卡活动发起服务，花费%d币", expenseCredits))
	if err != nil {
		return nil, err
	}
	result.BuyCredits = expenseCredits
	result.BalanceCredits = user.NyCredits

	return result, nil
}

// 分页查询
func (s *TimeClockService) GetByPage(userId int64, param request.PageableParam) (*response.PageResponse, error) {
	session := s.WhereUserSession(userId).Desc("create_time")
	list := make([]table.TimeClock, 0)
	return s.HandlePageable(param, &list, session)
}

// 我的打卡分页查询
func (s *TimeClockService) GetListByMemberId(userId int64, param request.PageableParam) (*response.PageResponse, error) {
	session := s.Db.NoCache().Where("id in (select clock_id from time_clock_member where user_id = ?)", userId).Desc("end_time")
	list := make([]table.TimeClock, 0)
	return s.HandlePageable(param, &list, session)
}

// 通过id取得打卡
func (s *TimeClockService) GetByOnlyId(id int64) (*table.TimeClock, error) {
	mod := &table.TimeClock{}
	_, err := s.Db.ID(id).Get(mod)
	return mod, err
}
