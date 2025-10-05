package cdkey

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"xorm.io/xorm"
)

type CdkeyService struct {
	service.UserCrudService[table.Cdkey]
}

var onceCdkey = sync.Once{}
var cdkeyService *CdkeyService

// 获取单例
func GetCdkeyService() *CdkeyService {
	onceCdkey.Do(func() {
		cdkeyService = new(CdkeyService)
		cdkeyService.Db = global.DB
	})
	return cdkeyService
}

// 创建cdkey
func (s *CdkeyService) Create(adminId int64, param request.CdkeyCreateParam) (string, error) {
	list := make([]table.Cdkey, param.CdkeyNum)
	for i := 0; i < param.CdkeyNum; i++ {
		key, err := utils.GenerateRandomString(12)
		if err != nil {
			return "", err
		}
		list[i] = table.Cdkey{
			AdminId:      adminId,
			Cdkey:        "NY_" + key,
			Credits:      param.Credits,
			TotalCount:   param.Count,
			SurplusCount: param.Count,
			ExpireTime:   param.ExpireTime,
		}
	}

	_, err := s.Db.Insert(list)
	if err != nil {
		return "", err
	}
	return utils.ListJoin(list, "\n", func(index int, item table.Cdkey) string {
		return item.Cdkey
	}), err
}

// 使用cdkey
func (s *CdkeyService) Use(userId int64, cdkey string) (response.CdkeyUseResponse, error) {
	result := response.CdkeyUseResponse{}

	err := s.DbTransaction(func(session *xorm.Session) error {
		// 检查cdkey
		mod := table.Cdkey{Cdkey: cdkey}
		exist, err := session.Where("surplus_count > 0 and (expire_time is null or expire_time > now())").Get(&mod)
		if err != nil {
			return err
		}
		if !exist {
			return errors.New("无效Cdkey")
		}
		// 检查该用户是否核销过该cdkey
		recordMod := table.CdkeyRecord{Cdkey: cdkey}
		recordExist, err := s.AddWhereUser(userId, session).Get(&recordMod)
		if err != nil {
			return err
		}
		if recordExist {
			return errors.New("您已核销过该Cdkey")
		}
		// 使用次数-1
		_, err = session.ID(mod.Id).Decr("surplus_count", 1).NoAutoTime().Update(&table.Cdkey{})
		if err != nil {
			return err
		}
		// 添加使用记录
		_, err = session.Insert(table.CdkeyRecord{
			UserId:  userId,
			Cdkey:   mod.Cdkey,
			Credits: mod.Credits,
		})
		if err != nil {
			return err
		}
		// 添加积分
		if mod.Credits > 0 {
			myRewardCredits := table.CreditsChange{
				Type:          enum.CreditsChangeType(4),
				ChangeCredits: mod.Credits,
				Remark:        fmt.Sprintf("核销Cdkey（%s），添加%d积分", mod.Cdkey, mod.Credits),
			}
			result.Credits = mod.Credits
			_, err = userSev.GetUserService().ChangeCreditsRunning(userId, session, myRewardCredits)
			if err != nil {
				return err
			}
		}
		if mod.Remark != "" {
			result.Remark = "核销成功，" + mod.Remark
		}
		return nil
	})
	return result, err
}
