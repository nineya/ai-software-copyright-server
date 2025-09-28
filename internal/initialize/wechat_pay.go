package initialize

import (
	"ai-software-copyright-server/internal/global"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/pkg/errors"
	"os"
)

func InitWechatPay() {
	merchant := global.CONFIG.Weixin.Merchant
	privateKey, err := os.ReadFile(merchant.PrivateKeyPath)
	if err != nil {
		panic(errors.Wrap(err, "初始化微信支付：读取私钥文件失败："+merchant.PrivateKeyPath))
	}
	client, err := wechat.NewClientV3(merchant.Mchid, merchant.SerialNo, merchant.ApiV3Key, string(privateKey))
	if err != nil {
		panic(errors.Wrap(err, "初始化微信支付：客户端初始化失败"))
	}
	// 启用自动同步返回验签，并定时更新微信平台API证书（开启自动验签时，无需单独设置微信平台API证书和序列号）
	err = client.AutoVerifySign()
	if err != nil {
		panic(errors.Wrap(err, "初始化微信支付：启用自动同步返回验签失败"))
	}
	// 打开Debug开关，输出日志，默认是关闭的
	client.DebugSwitch = gopay.DebugOn
	global.WECHAT_PAY = client
}
