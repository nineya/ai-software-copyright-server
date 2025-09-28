package utils

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/config"
	"ai-software-copyright-server/internal/global"
	"github.com/pkg/errors"
)

func MiniProgramConfig(clientType enum.ClientType) (*config.MiniProgramItem, error) {
	var mpConfig *config.MiniProgramItem
	switch clientType {
	case enum.ClientType(1):
		mpConfig = &global.CONFIG.Weixin.MiniProgram.BloggerHelper
	case enum.ClientType(2):
		mpConfig = &global.CONFIG.Weixin.MiniProgram.TrafficToolbox
	case enum.ClientType(3):
		mpConfig = &global.CONFIG.Weixin.MiniProgram.OperationHelper
	case enum.ClientType(4):
		mpConfig = &global.CONFIG.Weixin.MiniProgram.WechatToolbox
	case enum.ClientType(5):
		mpConfig = &global.CONFIG.Weixin.MiniProgram.NetdiskHelper
	}
	if mpConfig == nil {
		return nil, errors.New("客户端错误: " + enum.CLIENT_TYPE[clientType])
	}
	if mpConfig.Appid == "" {
		return nil, errors.New("Appid 为空: " + enum.CLIENT_TYPE[clientType])
	}
	if mpConfig.Secret == "" {
		return nil, errors.New("Secret 为空: " + enum.CLIENT_TYPE[clientType])
	}
	return mpConfig, nil
}
