package user

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"github.com/pkg/errors"
	"sync"
	"xorm.io/xorm"
)

type ShareRecordService struct {
	service.UserCrudService[table.ShareRecord]
}

var onceShareRecord = sync.Once{}
var shareRecordService *ShareRecordService

// 获取单例
func GetShareRecordService() *ShareRecordService {
	onceShareRecord.Do(func() {
		shareRecordService = new(ShareRecordService)
		shareRecordService.Db = global.DB
	})
	return shareRecordService
}

func (s *ShareRecordService) Create(userId int64, param table.ShareRecord) error {
	return s.DbTransaction(func(session *xorm.Session) error {
		exist, err := s.AddWhereUser(userId, session).Get(&table.ShareRecord{ShareUrl: param.ShareUrl})
		if err != nil {
			return errors.Wrap(err, "查询分享链接信息失败")
		}
		if exist {
			return errors.New("该分享链接已提交")
		}
		_, err = session.Insert(&table.ShareRecord{
			UserId:        userId,
			ShareUrl:      param.ShareUrl,
			RewardCredits: param.RewardCredits,
			Status:        enum.ShareStatus(1),
		})
		return err
	})
}

// 审核分享
func (s *ShareRecordService) Audit(param table.ShareRecord) (*table.ShareRecord, error) {
	mod := &table.ShareRecord{}
	exist, err := s.WhereUserSession(param.UserId).ID(param.Id).Get(mod)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("分享记录不存在")
	}
	mod.Status = param.Status
	mod.Remark = param.Remark
	err = s.DbTransaction(func(session *xorm.Session) error {
		// 审核通过，需要奖励分享的用户
		if param.Status == enum.ShareStatus(2) && param.RewardCredits > 0 {
			mod.RewardCredits = param.RewardCredits
			rewardCredits := table.CreditsChange{
				Type:          enum.CreditsChangeType(2),
				ChangeCredits: param.RewardCredits,
				Remark:        "分享使用体验赠送积分",
			}
			// 奖励用户自己
			_, err = GetUserService().ChangeCreditsRunning(mod.Id, session, rewardCredits)
			if err != nil {
				return err
			}
		}
		// 更新分享记录
		_, err = session.ID(mod.Id).Update(mod)
		return err
	})
	return mod, err
}

func (s *ShareRecordService) GetAll(userId int64) ([]table.ShareRecord, error) {
	list := make([]table.ShareRecord, 0)
	err := s.WhereUserSession(userId).Desc("create_time").Find(&list)
	return list, err
}

func (s *ShareRecordService) Statistic(userId int64) (*table.ShareStatistic, error) {
	mod := &table.ShareStatistic{}
	_, err := s.WhereUserSession(userId).Select("sum(CASE WHEN status = 2 THEN reward_credits END) AS share_credits, count(CASE WHEN status = 1 THEN 1 END) AS await_count, count(CASE WHEN status = 2 THEN 1 END) AS pass_count").Get(mod)
	return mod, err
}
