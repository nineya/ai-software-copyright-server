package mp

import (
	"ai-software-copyright-server/internal/application/param/request"
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type MpService struct {
	token     string // 公众号服务器配置Token
	appID     string // 公众号AppID
	appSecret string // 公众号AppSecret
}

var onceMp = sync.Once{}
var mpService *MpService

// 获取单例
func GetMpService() *MpService {
	onceMp.Do(func() {
		mpService = new(MpService)
		mpService.token = "361654768"
		mpService.appID = "wxad97d37cbbd8c8ff"
		mpService.appSecret = "f4b942fff39b976f89ac3004db0ad98f"
	})
	return mpService
}

func (s *MpService) Verify(param request.MpVerifyParam) string {
	if s.checkSignature(param) {
		return param.Echostr
	}
	return ""
}

// 校验签名有效性
func (s *MpService) checkSignature(param request.MpVerifyParam) bool {
	// 将token、timestamp、nonce按字典序排序
	strs := sort.StringSlice{s.token, param.Timestamp, param.Nonce}
	sort.Sort(strs)
	combined := strings.Join(strs, "")

	// SHA1加密
	hasher := sha1.New()
	hasher.Write([]byte(combined))
	sha := fmt.Sprintf("%x", hasher.Sum(nil))

	return sha == param.Signature
}
