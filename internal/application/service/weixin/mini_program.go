package weixin

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	attaSev "ai-software-copyright-server/internal/application/service/attachment"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"
)

type MiniProgramService struct {
	service.BaseService
}

var onceMiniProgram = sync.Once{}
var miniProgramService *MiniProgramService

// 获取单例
func GetMiniProgramService() *MiniProgramService {
	onceMiniProgram.Do(func() {
		miniProgramService = new(MiniProgramService)
		miniProgramService.Db = global.DB
	})
	return miniProgramService
}

// 通过授权码获取用户id
func (s *MiniProgramService) WeixinSessionLogin(clientType enum.ClientType, code string) (*response.WeixinSessionResponse, error) {
	mpConfig, err := utils.MiniProgramConfig(clientType)
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		mpConfig.Appid,
		mpConfig.Secret,
		code))
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	global.LOG.Info("微信小程序：通过授权码获取用户信息", zap.String("result", string(content)))

	var result response.WeixinSessionResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	if result.Errcode != 0 || result.Openid == "" {
		return nil, errors.New("登录状态失效")
	}
	if result.Unionid == "" {
		return nil, errors.New("该应用未绑定开放平台")
	}
	return &result, nil
}

// 获取小程序token
func (s *MiniProgramService) WeixinCgiBinAccessToken(clientType enum.ClientType) (string, error) {
	mpConfig, err := utils.MiniProgramConfig(clientType)
	if err != nil {
		return "", err
	}
	resp, err := http.Get(fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?appid=%s&secret=%s&grant_type=client_credential",
		mpConfig.Appid,
		mpConfig.Secret))
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result response.WeixinCgiBinAccessTokenResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return "", err
	}
	return result.AccessToken, nil
}

// 保存小程序邀请码
func (s *MiniProgramService) SaveInviteCode(clientType enum.ClientType, inviteCode, bucket string) (*response.UploadImageResponse, error) {
	accessToken, err := s.WeixinCgiBinAccessToken(clientType)
	if err != nil {
		return nil, err
	}

	// 设置参数
	param := &request.WeixinCodeUnlimitParam{
		Page:      "pages/index/index",
		Scene:     "inviter=" + inviteCode,
		CheckPath: true,
	}
	if global.CONFIG.Server.Mode == "prod" {
		param.EnvVersion = "release"
	} else {
		param.EnvVersion = "trial"
	}
	bytesData, _ := json.Marshal(param)
	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s", accessToken)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bytesData))
	if err != nil {
		return nil, err
	}

	// 发起请求
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	codeBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return attaSev.GetImageService().UploadByBytes(codeBytes, ".png", bucket)
}

// 文本内容安全识别，flase表示不通过
func (s *MiniProgramService) MsgSecCheck(clientType enum.ClientType, openid, msg string) error {
	accessToken, err := s.WeixinCgiBinAccessToken(clientType)
	if err != nil {
		return err
	}

	// 设置参数
	param := &request.WeixinMsgSecCheckParam{
		Content: msg,
		Version: 2,
		Scene:   2,
		Openid:  openid,
	}
	bytesData, _ := json.Marshal(param)
	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/msg_sec_check?access_token=%s", accessToken)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bytesData))
	if err != nil {
		return err
	}

	// 发起请求
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result response.WeixinMsgSecCheckResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return err
	}
	if result.Errcode != 0 {
		return errors.New(result.Errmsg)
	}
	if result.Result.Label != 100 && result.Result.Suggest == "risky" {
		switch result.Result.Label {
		case 10001:
			return errors.New("包含广告内容，禁止发布")
		case 20001:
			return errors.New("包含时政内容，禁止发布")
		case 20002:
			return errors.New("包含色情内容，禁止发布")
		case 20003:
			return errors.New("包含辱骂内容，禁止发布")
		case 20006:
			return errors.New("包含违法犯罪内容，禁止发布")
		case 20008:
			return errors.New("包含欺诈内容，禁止发布")
		case 20012:
			return errors.New("包含低俗内容，禁止发布")
		case 20013:
			return errors.New("包含版权内容，禁止发布")
		case 21000:
			return errors.New("包含风险内容，禁止发布")
		}
	}
	return nil
}
