package flash_picture

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"xorm.io/xorm"
)

type FlashPictureService struct {
	service.UserCrudService[table.FlashPicture]
}

var onceFlashPicture = sync.Once{}
var flashPictureService *FlashPictureService

// 获取单例
func GetFlashPictureService() *FlashPictureService {
	onceFlashPicture.Do(func() {
		flashPictureService = new(FlashPictureService)
		flashPictureService.Db = global.DB
	})
	return flashPictureService
}

func (s *FlashPictureService) Browse(userId, id int64) (*response.FlashPictureBrowseResponse, error) {
	mod := &table.FlashPicture{}
	exist, err := s.Db.ID(id).Get(mod)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("该闪图分享已被撤回")
	}
	// 取得当前浏览次数
	record := &table.FlashPictureRecord{UserId: userId, FlashPictureId: id}
	_, err = s.Db.Get(record)
	if err != nil {
		return nil, err
	}
	result := &response.FlashPictureBrowseResponse{FlashPicture: mod, UseVisits: record.Visits}
	return result, err
}

func (s *FlashPictureService) Visits(userId, id int64) (*response.FlashPictureVisitsResponse, error) {
	expenseCredits := 30

	mod := &table.FlashPicture{}
	exist, err := s.Db.ID(id).Get(mod)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("该闪图分享已被撤回")
	}

	result := &response.FlashPictureVisitsResponse{FlashPicture: mod}
	err = s.DbTransaction(func(session *xorm.Session) error {
		record := &table.FlashPictureRecord{UserId: userId, FlashPictureId: id}
		exist, err := s.Db.Get(record)
		if err != nil {
			return err
		}
		if mod.TriesLimit != 0 && record.Visits >= mod.TriesLimit {
			return errors.New("该闪照浏览次数已达上限")
		}
		// 添加浏览次数
		_, err = s.Db.ID(id).Incr("visits", 1).NoAutoTime().Update(&table.FlashPicture{})
		if err != nil {
			return err
		}
		result.Visits++
		result.UseVisits = record.Visits + 1

		// 如果不是第一次，则收费
		if exist {
			// 预检余额
			_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
			if err != nil {
				return err
			}

			// 添加浏览次数
			_, err = s.Db.ID(record.Id).Incr("visits", 1).NoAutoTime().Update(&table.FlashPictureRecord{})
			if err != nil {
				return err
			}

			// 扣款
			user, err := userSev.GetUserService().PaymentCredits(userId, enum.BuyType(12), expenseCredits, fmt.Sprintf("购买浏览闪照服务，花费%d币", expenseCredits))
			if err != nil {
				return err
			}
			result.BuyCredits = expenseCredits
			result.BalanceCredits = user.Credits
			return nil
		}
		_, err = session.Insert(table.FlashPictureRecord{
			UserId:         userId,
			FlashPictureId: id,
			Visits:         1,
		})
		if err != nil {
			return err
		}
		// 查询余额
		user, err := userSev.GetUserService().GetById(userId)
		result.BalanceCredits = user.Credits
		return err
	})
	return result, err
}

func (s *FlashPictureService) GetOrigin(userId, id int64) (*response.UserBuyContentResponse, error) {
	expenseCredits := 10
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}

	mod := &table.FlashPicture{}
	exist, err := s.Db.ID(id).Get(mod)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("闪图已被撤回")
	}

	result := &response.UserBuyContentResponse{}
	switch mod.OriginType {
	case enum.PictureOriginType(1):
		result.Content = "拍照"
	case enum.PictureOriginType(2):
		result.Content = "相册选取"
	case enum.PictureOriginType(3):
		result.Content = "聊天记录选取"
	}

	// 扣款
	user, err := userSev.GetUserService().PaymentCredits(userId, enum.BuyType(11), expenseCredits, fmt.Sprintf("购买闪照查看来源服务，花费%d币", expenseCredits))
	if err != nil {
		return nil, err
	}
	result.BuyCredits = expenseCredits
	result.BalanceCredits = user.Credits

	return result, nil
}
