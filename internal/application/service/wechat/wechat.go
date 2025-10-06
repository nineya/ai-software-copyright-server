package wechat

import (
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"
)

type WechatService struct {
	service.BaseService
}

var onceWechat = sync.Once{}
var wechatService *WechatService

// 获取单例
func GetWechatService() *WechatService {
	onceWechat.Do(func() {
		wechatService = new(WechatService)
		wechatService.Db = global.DB
	})
	return wechatService
}

// 获取token和unionid
func (s *WechatService) Oauth2AccessToken(appid, secret, code string) (*response.WechatOauth2AccessTokenResponse, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", appid, secret, code)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	global.LOG.Info("向微信平台发送请求：", zap.String("url", url), zap.String("result", string(content)))
	if err != nil {
		return nil, err
	}

	result := &response.WechatOauth2AccessTokenResponse{} // 反序列化JSON到结构体
	err = json.Unmarshal(content, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 获取token和unionid
func (s *WechatService) UserInfo(accessToken, openid string) (*response.WechatUserInfoResponse, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s", accessToken, openid)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	global.LOG.Info("向微信平台发送请求：", zap.String("url", url), zap.String("result", string(content)))
	if err != nil {
		return nil, err
	}

	result := &response.WechatUserInfoResponse{} // 反序列化JSON到结构体
	err = json.Unmarshal(content, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
