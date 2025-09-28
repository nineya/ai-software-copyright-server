package user

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	mailPlugin "ai-software-copyright-server/internal/application/plugin/mail"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"sync"
	"xorm.io/xorm"
)

type UserService struct {
	service.BaseService
}

var onceUser = sync.Once{}
var userService *UserService

// 获取单例
func GetUserService() *UserService {
	onceUser.Do(func() {
		userService = new(UserService)
		userService.Db = global.DB
	})
	return userService
}

func (s *UserService) SendMail(userId int64, title, content string) error {
	mod, err := s.GetById(userId)
	if err != nil {
		return errors.Wrap(err, "发送邮件失败，获取用户信息失败")
	}
	if mod.Email == nil || *mod.Email == "" {
		global.LOG.Warn("未配置邮箱，无法发送邮件")
		return nil
	}
	return mailPlugin.GetMailPlugin().SendHtmlMail(request.MailParam{
		To:      *mod.Email,
		Subject: title,
		Content: content,
	})
}

func (s *UserService) UpdateUserInfo(userId int64, param request.UserUpdateInfoParam) (*response.UserRewardResponse, error) {
	mod, err := s.GetById(userId)
	if err != nil {
		return nil, err
	}
	updateMod := &table.User{}
	first := false
	if param.Phone != "" {
		if mod.Phone == nil {
			first = true
			if mod.Password == "" && param.Password == "" {
				param.Password = param.Phone[len(param.Phone)-8 : len(param.Phone)]
			}
		}
		if mod.Phone == nil || *mod.Phone != param.Phone {
			exist, err := s.Db.Get(&table.User{Phone: &param.Phone})
			if err != nil {
				return nil, err
			}
			if exist {
				return nil, errors.New("该手机号已绑定其他账号")
			}
			updateMod.Phone = &param.Phone
		}
	}
	if param.Email != "" {
		if mod.Email == nil || *mod.Email != param.Email {
			exist, err := s.Db.Get(&table.User{Email: &param.Email})
			if err != nil {
				return nil, err
			}
			if exist {
				return nil, errors.New("该邮箱已绑定其他账号")
			}
			updateMod.Email = &param.Email
		}
	}
	if param.Password != "" {
		hashPwd, err := bcrypt.GenerateFromPassword(utils.Md5ByBytes(param.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updateMod.Password = string(hashPwd)
	}

	result := &response.UserRewardResponse{}

	err = s.DbTransaction(func(session *xorm.Session) error {
		_, err = session.ID(userId).Update(updateMod)
		if err != nil {
			return err
		}
		if first {
			// 奖励自己
			myRewardCredits := table.CreditsChange{
				Type:          enum.CreditsChangeType(5),
				ChangeCredits: 100,
				Remark:        "首次维护个人信息，送100积分",
			}
			user, err := GetUserService().ChangeCreditsRunning(mod.Id, session, myRewardCredits)
			if err != nil {
				return err
			}
			result.Message = myRewardCredits.Remark
			result.RewardCredits = myRewardCredits.ChangeCredits
			result.BalanceCredits = user.Credits
		}
		return nil
	})
	return result, err
}

// 购物减少积分
func (s *UserService) PaymentNyCredits(userId int64, typ enum.BuyType, expenseCredits int, remark string) (*table.User, error) {
	mod := &table.User{}
	err := s.DbTransaction(func(session *xorm.Session) error {
		// 执行金额变动和记录
		user, err := s.ChangeCreditsRunning(userId, session, table.CreditsChange{Type: enum.CreditsChangeType(1), ChangeCredits: expenseCredits * -1, Remark: remark})
		if err != nil {
			return err
		}
		mod = user
		// 记录购买信息
		_, err = session.Insert(table.Buy{
			UserId:     userId,
			Type:       typ,
			PayCredits: expenseCredits,
			Remark:     remark,
		})
		return err
	})
	return mod, err
}

// 添加积分
func (s *UserService) AddNyCredits(param request.UserAddNyCreditsParam) ([]table.User, error) {
	list := make([]table.User, len(param.InviteCodes))
	creditsChange := table.CreditsChange{Type: param.Type, ChangeCredits: param.AddCredits, Remark: param.Remark}
	err := s.DbTransaction(func(session *xorm.Session) error {
		for i, code := range param.InviteCodes {
			mod := &table.User{InviteCode: code}
			exist, err := s.Db.Get(mod)
			if !exist {
				return errors.New("邀请码不存在：" + code)
			}
			user, err := s.ChangeCreditsRunning(mod.Id, session, creditsChange)
			if err != nil {
				return err
			}
			list[i] = *user
		}
		return nil
	})
	return list, err
}

// 修改积分
func (s *UserService) ChangeNyCredits(userId int64, param table.CreditsChange) (*table.User, error) {
	var mod *table.User
	err := s.DbTransaction(func(session *xorm.Session) error {
		user, err := s.ChangeCreditsRunning(userId, session, param)
		mod = user
		return err
	})
	return mod, err
}

// 购物减少积分
func (s *UserService) PaymentNyCreditsRunning(userId int64, session *xorm.Session, typ enum.BuyType, expenseCredits int, remark string) (*table.User, error) {
	// 执行金额变动和记录
	mod, err := s.ChangeCreditsRunning(userId, session, table.CreditsChange{Type: enum.CreditsChangeType(1), ChangeCredits: expenseCredits * -1, Remark: remark})
	if err != nil {
		return nil, err
	}
	// 记录购买信息
	_, err = session.Insert(table.Buy{
		UserId:     userId,
		Type:       typ,
		PayCredits: expenseCredits,
		Remark:     remark,
	})
	return mod, err
}

// 修改积分执行实体，并添加积分修改记录
func (s *UserService) ChangeCreditsRunning(userId int64, session *xorm.Session, param table.CreditsChange) (*table.User, error) {
	mod := &table.User{}
	// 取得原始用户信息
	_, err := session.ID(userId).Get(mod)
	if err != nil {
		return nil, err
	}
	// 如果金额为0，则不用记录后面的余额变动
	if param.ChangeCredits == 0 {
		return mod, nil
	}
	session.ID(userId)
	// 需要判断余额
	if param.ChangeCredits < 0 {
		session.Where("credits >= ?", param.ChangeCredits*-1)
	}
	// 修改金额
	num, err := session.Incr("credits", param.ChangeCredits).Update(&table.User{})
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, errors.New("积分不足，邀请好友可赠送积分")
	}
	// 创建余额变动实体
	param.UserId = userId
	param.OriginCredits = mod.Credits
	// 重新查询用户信息
	mod = &table.User{}
	_, err = session.ID(userId).Get(mod)
	if err != nil {
		return nil, err
	}
	// 记录积分变动
	param.BalanceCredits = mod.Credits
	_, err = session.Insert(param)
	return mod, err
}

// 邀请获得积分执行实体
func (s *UserService) InviterCreditsRunning(userId int64, session *xorm.Session, param request.UserInviterCreditsParam) error {
	if param.Inviter == "" {
		return nil
	}
	// 取得邀请人信息
	inviter := &table.User{InviteCode: param.Inviter}
	_, err := session.Get(inviter)
	if err != nil {
		return err
	}
	if inviter.Id == 0 {
		return nil
	}
	// 执行金额变动和记录
	inviter, err = s.ChangeCreditsRunning(inviter.Id, session, table.CreditsChange{Type: enum.CreditsChangeType(3), ChangeCredits: param.RewardCredits, Remark: param.Remark})
	if err != nil {
		return err
	}
	// 记录邀请信息
	_, err = session.Insert(table.InviteRecord{
		UserId:        inviter.Id,
		InviteeId:     userId,
		Type:          param.Type,
		RewardCredits: param.RewardCredits,
		Remark:        param.Remark,
	})
	return err
}

func (s *UserService) GetById(id int64) (*table.User, error) {
	mod := &table.User{}
	_, err := s.Db.ID(id).Get(mod)
	return mod, err
}

func (s *UserService) GetAndCheckBalance(userId int64, credits int) (*table.User, error) {
	user, err := s.GetById(userId)
	if err != nil {
		return nil, err
	}
	if user.Credits < credits {
		return nil, errors.New("积分不足，邀请好友可获赠积分")
	}
	return user, nil
}
