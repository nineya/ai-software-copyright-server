package user

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	wxSev "ai-software-copyright-server/internal/application/service/weixin"
	"ai-software-copyright-server/internal/global"
	"errors"
	"sync"
)

type ClientInfoService struct {
	service.UserCrudService[table.ClientInfo]
}

var onceClientInfo = sync.Once{}
var clientInfoService *ClientInfoService

// 获取单例
func GetClientInfoService() *ClientInfoService {
	onceClientInfo.Do(func() {
		clientInfoService = new(ClientInfoService)
		clientInfoService.Db = global.DB
	})
	return clientInfoService
}

func (s *ClientInfoService) MsgSecCheck(userId int64, clientType enum.ClientType, msg string) error {
	if clientType != 1 && clientType != 2 && clientType != 3 && clientType != 4 && clientType != 5 {
		return nil
	}
	clientInfo, err := s.GetClientInfo(userId, clientType)
	if err != nil || clientInfo.Id == 0 {
		return errors.New("获取客户端信息失败")
	}
	return wxSev.GetMiniProgramService().MsgSecCheck(clientType, clientInfo.WxOpenid, msg)
}

func (s *ClientInfoService) GetClientInfo(userId int64, clientType enum.ClientType) (table.ClientInfo, error) {
	mod := table.ClientInfo{UserId: userId, Type: clientType}
	_, err := s.WhereUserSession(userId).Get(&mod)
	return mod, err
}

func (s *ClientInfoService) GetInviteCode(userId int64, clientType enum.ClientType) (string, error) {
	mod := &table.ClientInfo{UserId: userId, Type: clientType}
	exist, err := s.WhereUserSession(userId).Get(mod)
	if err != nil || !exist {
		return "", err
	}
	// 邀请码已生成，直接返回
	if mod.QrCodeUrl != "" {
		return mod.QrCodeUrl, nil
	}
	// 邀请码未生成，生成邀请码
	user, err := GetUserService().GetById(userId)
	if err != nil {
		return "", err
	}
	upload, err := wxSev.GetMiniProgramService().SaveInviteCode(clientType, user.InviteCode, "user")
	if err != nil {
		return "", err
	}
	_, _ = s.Db.ID(mod.Id).Update(table.ClientInfo{QrCodeUrl: upload.Url})
	return upload.Url, nil
}
