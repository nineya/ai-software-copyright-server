package cdkey

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	netdSev "ai-software-copyright-server/internal/application/service/netdisk"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"time"
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
			AdminId:           adminId,
			Cdkey:             "NY_" + key,
			CreditsNum:        param.CreditsNum,
			HelperStandardDay: param.HelperStandardDay,
			HelperWechatDay:   param.HelperWechatDay,
			TotalCount:        param.Count,
			SurplusCount:      param.Count,
			ExpireTime:        param.ExpireTime,
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
			UserId:            userId,
			Cdkey:             mod.Cdkey,
			CreditsNum:        mod.CreditsNum,
			HelperStandardDay: mod.HelperStandardDay,
			HelperWechatDay:   mod.HelperWechatDay,
		})
		if err != nil {
			return err
		}
		// 添加积分
		if mod.CreditsNum > 0 {
			myRewardCredits := table.CreditsChange{
				Type:          enum.CreditsChangeType(4),
				ChangeCredits: mod.CreditsNum,
				Remark:        fmt.Sprintf("核销Cdkey（%s），添加%d积分", mod.Cdkey, mod.CreditsNum),
			}
			result.NyCredits = mod.CreditsNum
			_, err = userSev.GetUserService().ChangeCreditsRunning(userId, session, myRewardCredits)
			if err != nil {
				return err
			}
		}
		now := time.Now()
		// 赠送网盘助手时间
		if mod.HelperStandardDay > 0 || mod.HelperWechatDay > 0 {
			config, err := netdSev.GetHelperConfigureService().GetByUserId(userId)
			if err != nil {
				return errors.Wrap(err, "获取网盘助手信息失败")
			}
			// 未获取到助手配置，需要保存
			if config.Id == 0 {
				err = netdSev.GetHelperConfigureService().SaveConfigure(userId, config)
				if err != nil {
					return errors.Wrap(err, "保存网盘助手配置失败")
				}
			}
			session = s.AddWhereUser(userId, session).NoAutoTime()
			// 增加标准版时间
			if mod.HelperStandardDay > 0 {
				if config.ExpireTime == nil || config.ExpireTime.Before(now) {
					session.SetExpr("expire_time", now.AddDate(0, 0, mod.HelperStandardDay))
				} else {
					session.SetExpr("expire_time", config.ExpireTime.AddDate(0, 0, mod.HelperStandardDay))
				}
			}
			// 增加微信工具人版时间
			if mod.HelperWechatDay > 0 {
				if config.WechatExpireTime == nil || config.WechatExpireTime.Before(now) {
					session.SetExpr("wechat_expire_time", now.AddDate(0, 0, mod.HelperWechatDay))
				} else {
					session.SetExpr("wechat_expire_time", config.WechatExpireTime.AddDate(0, 0, mod.HelperWechatDay))
				}
			}
			_, err = session.Update(table.NetdiskHelperConfigure{})
			if err != nil {
				return errors.Wrap(err, "更新网盘助手信息失败")
			}
		}
		if mod.Remark != "" {
			result.Remark = "核销成功，" + mod.Remark
		}
		return nil
	})
	return result, err
}
