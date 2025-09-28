package user

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"sync"
	"xorm.io/xorm"
)

type AuthService struct {
	service.BaseService
}

var onceAuth = sync.Once{}
var authService *AuthService

// 获取单例
func GetAuthService() *AuthService {
	onceAuth.Do(func() {
		authService = new(AuthService)
		authService.Db = global.DB
	})
	return authService
}

func (s *AuthService) Login(param request.UserLoginParam) (*common.Token, error) {
	var user table.User
	exist, err := s.Db.Where("phone = ?", param.Phone).Get(&user)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("用户不存在或已被禁用")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), utils.Md5ByBytes(param.Password))
	if err != nil {
		global.LOG.Error(fmt.Sprintf("登录密码不正确：%s, %s, %+v", param.Password, user.Password, err))
		return nil, errors.New("登录密码不正确")
	}
	return utils.AuthToken(user.Id, global.User)
}

func (s *AuthService) Register(param table.User) error {
	exist, err := s.Db.Get(table.User{Phone: param.Phone})
	if err != nil {
		return errors.Wrap(err, "查询手机号失败")
	}
	if exist {
		return errors.New("该手机号已注册")
	}
	exist, err = s.Db.Get(table.User{Email: param.Email})
	if err != nil {
		return errors.Wrap(err, "查询邮箱失败")
	}
	if exist {
		return errors.New("该邮箱已注册")
	}

	// 新增用户
	mod := &table.User{Nickname: param.Nickname, Phone: param.Phone, Email: param.Email, Inviter: param.Inviter}
	hashPwd, err := bcrypt.GenerateFromPassword(utils.Md5ByBytes(param.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err, "密码不符合要求")
	}
	mod.Password = string(hashPwd)
	return s.DbTransaction(func(session *xorm.Session) error {
		_, err = session.Insert(mod)
		if err != nil {
			return err
		}
		// 奖励自己
		myRewardCredits := table.CreditsChange{
			Type:          enum.CreditsChangeType(5),
			ChangeCredits: 50,
			Remark:        "首次使用，平台赠送50个积分",
		}
		// 判断邀请码是否有效
		inviter := &table.User{InviteCode: param.Inviter}
		if param.Inviter != "" {
			exist, _ = session.Get(inviter)
			if !exist {
				return errors.New("邀请码不存在")
			}
			myRewardCredits.ChangeCredits = 100
			myRewardCredits.Remark = "您的好友邀请了您，并送了您100个积分"
		}
		// 奖励用户自己
		user, err := GetUserService().ChangeCreditsRunning(mod.Id, session, myRewardCredits)
		if err != nil {
			return err
		}
		mod = user
		// 添加邀请人的邀请奖励
		if param.Inviter != "" {
			inviterRewardCredits := request.UserInviterCreditsParam{
				Inviter:       mod.Inviter,
				Type:          enum.InviteType(1),
				RewardCredits: 50,
				Remark:        fmt.Sprintf("邀请新用户（%s）奖励50个积分", utils.MaskContent(mod.Nickname)),
			}
			err = GetUserService().InviterCreditsRunning(mod.Id, session, inviterRewardCredits)
			if err != nil {
				return err
			}
		}
		// 生成邀请码
		mod.InviteCode = strconv.FormatInt(mod.Id+10000, 10)
		_, err = session.ID(mod.Id).Update(mod)
		return err
	})
}
