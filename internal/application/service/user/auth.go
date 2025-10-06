package user

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/service"
	attaSev "ai-software-copyright-server/internal/application/service/attachment"
	wechatSev "ai-software-copyright-server/internal/application/service/wechat"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
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

func (s *AuthService) Register(param request.UserInfoParam) (*table.User, error) {
	if param.Phone != "" {
		exist, err := s.Db.Get(&table.User{Phone: &param.Phone})
		if err != nil {
			return nil, errors.Wrap(err, "查询手机号失败")
		}
		if exist {
			return nil, errors.New("该手机号已注册")
		}
	}
	if param.Email != "" {
		exist, err := s.Db.Get(&table.User{Email: &param.Email})
		if err != nil {
			return nil, errors.Wrap(err, "查询邮箱失败")
		}
		if exist {
			return nil, errors.New("该邮箱已注册")
		}
	}

	// 新增用户
	mod := &table.User{Nickname: param.Nickname, Phone: &param.Phone, Email: &param.Email, WxUnionid: param.WxUnionid, WxOpenid: param.WxOpenid, Inviter: param.Inviter}
	if param.Password != "" {
		hashPwd, err := bcrypt.GenerateFromPassword(utils.Md5ByBytes(param.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.Wrap(err, "密码不符合要求")
		}
		mod.Password = string(hashPwd)
	}
	err := s.DbTransaction(func(session *xorm.Session) error {
		_, err := session.Insert(mod)
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
			exist, _ := session.Get(inviter)
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
				RewardCredits: 25,
				Remark:        fmt.Sprintf("邀请新用户（%s）奖励25个积分", utils.MaskContent(mod.Nickname)),
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
	return mod, err
}

// 微信授权登录
func (s *AuthService) Authorization(code, inviter string) (*common.Token, error) {
	if code == "" {
		return nil, errors.New("登录状态错误，请刷新")
	}
	// 通过code获取Unionid
	tokenResult, err := wechatSev.GetWechatService().Oauth2AccessToken(global.CONFIG.Weixin.Site.Main.Appid, global.CONFIG.Weixin.Site.Main.Secret, code)
	if err != nil {
		return nil, errors.Wrap(err, "系统繁忙，请稍后重试")
	}
	if tokenResult.AccessToken == "" || tokenResult.Unionid == "" || tokenResult.Openid == "" {
		return nil, errors.New("系统繁忙，请稍后重试")
	}
	// 取得用户信息
	user, err := GetUserService().GetByWxUnionid(tokenResult.Unionid)
	if err != nil {
		return nil, errors.Wrap(err, "获取用户信息失败")
	}
	// 用户未注册，注册用户
	if user.Id == 0 {
		// 获取用户信息
		userInfoResult, err := wechatSev.GetWechatService().UserInfo(tokenResult.AccessToken, tokenResult.Openid)
		if err != nil {
			return nil, errors.Wrap(err, "系统繁忙，请稍后重试")
		}
		if userInfoResult.Nickname == "" || userInfoResult.HeadImgUrl == "" {
			return nil, errors.New("系统繁忙，请稍后重试")
		}
		userInfo := request.UserInfoParam{
			Nickname:  userInfoResult.Nickname,
			WxOpenid:  tokenResult.Openid,
			WxUnionid: tokenResult.Unionid,
			Inviter:   inviter,
		}
		// 存储用户头像
		resp, err := http.Get(userInfoResult.HeadImgUrl)
		defer resp.Body.Close()
		codeBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			global.LOG.Error(fmt.Sprintf("下载用户头像失败：%+v", err))
		} else {
			imageResponse, err := attaSev.GetImageService().UploadByBytes(codeBytes, ".png", "user")
			if err != nil {
				global.LOG.Error(fmt.Sprintf("存储用户头像失败：%+v", err))
			} else {
				userInfo.Avatar = imageResponse.Url
			}
		}
		// 注册用户
		user, err = s.Register(userInfo)
		if err != nil {
			return nil, errors.Wrap(err, "用户注册失败，请重试")
		}
	}
	// 生成登录token
	return utils.AuthToken(user.Id, global.User)
}
