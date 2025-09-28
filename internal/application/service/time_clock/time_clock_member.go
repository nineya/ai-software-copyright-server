package time_clock

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type TimeClockMemberService struct {
	service.UserCrudService[table.TimeClockMember]
}

var onceTimeClockMember = sync.Once{}
var timeClockMemberService *TimeClockMemberService

// 获取单例
func GetTimeClockMemberService() *TimeClockMemberService {
	onceTimeClockMember.Do(func() {
		timeClockMemberService = new(TimeClockMemberService)
		timeClockMemberService.Db = global.DB
	})
	return timeClockMemberService
}

// 审核成员id
func (s *TimeClockMemberService) AuditById(userId int64, param table.TimeClockMember) error {
	mod, err := GetTimeClockService().GetById(userId, param.ClockId)
	if err != nil {
		return err
	}
	if mod.Id == 0 {
		return errors.New("你非该打卡管理员")
	}
	param.Status = enum.TimeClockMemberStatus(1)
	now := time.Now()
	param.JoinTime = &now
	_, err = s.Db.ID(param.Id).AllCols().Update(&param)
	return err
}

// 取得当前打卡和当前成员的信息
func (s *TimeClockMemberService) GetMyInfoById(userId int64, clockId int64) (*response.TimeClockMyInfoResponse, error) {
	mod := table.TimeClock{}
	exist, err := s.Db.ID(clockId).Get(&mod)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("不存在该打卡")
	}
	member := table.TimeClockMember{ClockId: clockId}
	_, err = s.WhereUserSession(userId).Get(&member)
	if err != nil {
		return nil, err
	}
	exist, err = s.WhereUserSession(userId).Where("clock_id = ? and create_time >= CURDATE()", clockId).Exist(&table.TimeClockRecord{})
	if err != nil {
		return nil, err
	}
	count, err := s.WhereUserSession(userId).Count(&table.TimeClockRecord{ClockId: clockId})
	if err != nil {
		return nil, err
	}
	return &response.TimeClockMyInfoResponse{mod, member, exist, count}, err
}

// 分页查询
func (s *TimeClockMemberService) GetByPage(userId int64, param request.PageableParam) (*response.PageResponse, error) {
	session := s.WhereUserSession(userId).Desc("create_time")
	list := make([]table.TimeClockMember, 0)
	return s.HandlePageable(param, &list, session)
}

// 分页取得当前打卡的成员
func (s *TimeClockMemberService) GetListByClockId(userId int64, param request.TimeClockQueryPageParam) (*response.PageResponse, error) {
	mod, err := GetTimeClockService().GetById(userId, param.ClockId)
	if err != nil {
		return nil, err
	}
	if mod.Id == 0 {
		return nil, errors.New("你非该打卡管理员")
	}
	session := s.Db.Table(table.TimeClockMember{}).
		Select("*,(select count(*) from time_clock_record where clock_id = time_clock_member.clock_id and user_id = user.id) as clock_in_count").
		Join("LEFT", table.User{}, "time_clock_member.user_id = user.id").
		Where("clock_id = ?", mod.Id).
		Desc("time_clock_member.create_time")
	list := make([]response.TimeClockMemberInfoResponse, 0)
	return s.HandlePageable(param.PageableParam, &list, session)
}
