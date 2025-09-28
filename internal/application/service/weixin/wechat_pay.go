package weixin

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"net/http"
	"sync"
	"time"
)

type WechatPayService struct {
	service.BaseService
}

var onceWechatPay = sync.Once{}
var wechatPayService *WechatPayService

// 获取单例
func GetWechatPayService() *WechatPayService {
	onceWechatPay.Do(func() {
		wechatPayService = new(WechatPayService)
		wechatPayService.Db = global.DB
	})
	return wechatPayService
}

// 创建订单
func (s *WechatPayService) CreateOrder(order *table.CreditsOrder) (*wechat.AppletParams, error) {
	price, err := decimal.NewFromString(order.OrderAmount)
	if err != nil {
		return nil, err
	}
	price = price.Mul(decimal.NewFromInt(100))

	mpConfig, err := utils.MiniProgramConfig(order.ClientType)
	if err != nil {
		return nil, err
	}
	merchant := global.CONFIG.Weixin.Merchant
	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("appid", mpConfig.Appid).
		Set("mchid", merchant.Mchid).
		Set("description", order.Description).
		Set("out_trade_no", order.TradeNo).
		Set("time_expire", expire).
		Set("notify_url", merchant.NotifyUrl).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", price.IntPart()).
				Set("currency", "CNY")
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", order.WxOpenid)
		})
	wxRsp, err := global.WECHAT_PAY.V3TransactionJsapi(context.Background(), bm)
	if err != nil {
		return nil, err
	}
	fmt.Println(wxRsp.Code)
	if wxRsp.Code != 0 {
		return nil, errors.New(wxRsp.Error)
	}
	// 生成小程序的pay sign
	return global.WECHAT_PAY.PaySignOfApplet(mpConfig.Appid, wxRsp.Response.PrepayId)
}

// 支付回调通知
func (s *WechatPayService) PayNotify(request *http.Request) (*wechat.V3DecryptPayResult, error) {
	notifyReq, err := wechat.V3ParseNotify(request)
	if err != nil {
		return nil, err
	}
	// 获取微信平台证书
	certMap := global.WECHAT_PAY.WxPublicKeyMap()
	// 验证异步通知的签名
	err = notifyReq.VerifySignByPKMap(certMap)
	if err != nil {
		return nil, err
	}
	// 普通支付通知解密
	return notifyReq.DecryptPayCipherText(global.CONFIG.Weixin.Merchant.ApiV3Key)
}

// 退款回调通知
func (s *WechatPayService) RefundNotify(c *gin.Context) error {
	notifyReq, err := wechat.V3ParseNotify(c.Request)
	if err != nil {
		return err
	}
	// 获取微信平台证书
	certMap := global.WECHAT_PAY.WxPublicKeyMap()
	// 验证异步通知的签名
	err = notifyReq.VerifySignByPKMap(certMap)
	if err != nil {
		return err
	}
	// 退款通知解密
	result, err := notifyReq.DecryptRefundCipherText(global.CONFIG.Weixin.Merchant.ApiV3Key)
	if err != nil {
		c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.FAIL, Message: "内容解密失败"})
		return err
	}
	print(result.OutTradeNo)
	c.JSON(http.StatusOK, &wechat.V3NotifyRsp{Code: gopay.SUCCESS, Message: "成功"})
	return nil
}
