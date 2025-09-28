package user

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	wxSev "ai-software-copyright-server/internal/application/service/weixin"
	"ai-software-copyright-server/internal/global"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"sync"
	"time"
	"xorm.io/xorm"
)

type CreditsOrderService struct {
	service.UserCrudService[table.CreditsOrder]
}

var onceCreditsOrder = sync.Once{}
var creditsOrderService *CreditsOrderService

// 获取单例
func GetCreditsOrderService() *CreditsOrderService {
	onceCreditsOrder.Do(func() {
		creditsOrderService = new(CreditsOrderService)
		creditsOrderService.Db = global.DB
	})
	return creditsOrderService
}

func (s *CreditsOrderService) CreateOrder(userId int64, clientType enum.ClientType, creditsPriceId int64) (*response.CreditsCreateOrderResponse, error) {
	creditsPrice := &table.CreditsPrice{}
	exist, err := s.Db.ID(creditsPriceId).Get(creditsPrice)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("商品信息不存在")
	}
	clientInfo := &table.ClientInfo{UserId: userId, Type: clientType}
	exist, err = s.Db.Get(clientInfo)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, errors.New("状态异常，请删除小程序重新登录")
	}
	order := &table.CreditsOrder{
		UserId:      userId,
		TradeNo:     fmt.Sprintf("POINTS_%d_%d_%d", userId, creditsPrice.Credits, time.Now().UnixMilli()),
		ClientType:  clientType,
		WxOpenid:    clientInfo.WxOpenid,
		Description: fmt.Sprintf("购买%d个积分", creditsPrice.Credits),
		CreditsNum:  creditsPrice.Credits,
		OrderAmount: creditsPrice.Price,
		Status:      enum.OrderStatus(1),
	}
	appletParams, err := wxSev.GetWechatPayService().CreateOrder(order)
	if err != nil {
		return nil, err
	}
	_, err = s.Db.Insert(order)
	return &response.CreditsCreateOrderResponse{CreditsOrder: order, AppletParams: appletParams}, err
}

// 支付回调通知
func (s *CreditsOrderService) PayNotify(request *http.Request) error {
	result, err := wxSev.GetWechatPayService().PayNotify(request)
	if err != nil {
		return err
	}
	global.LOG.Sugar().Infof("收到微信支付通知: %+v", result)
	// 查询订单
	creditsOrder := &table.CreditsOrder{TradeNo: result.OutTradeNo}
	_, err = s.Db.Get(creditsOrder)
	if err != nil {
		return err
	}
	// 已经是成功状态，不再进行其他操作
	if creditsOrder.Status == enum.OrderStatus(2) {
		return nil
	}
	// 更新订单状态和币金额
	err = s.DbTransaction(func(session *xorm.Session) error {
		_, err := session.ID(creditsOrder.Id).Update(&table.CreditsOrder{Status: enum.OrderStatus(2)})
		if err != nil {
			return err
		}
		_, err = GetUserService().ChangeCreditsRunning(creditsOrder.UserId, session, table.CreditsChange{
			Type:          enum.CreditsChangeType(4),
			ChangeCredits: creditsOrder.CreditsNum,
			Remark:        creditsOrder.Description,
		})
		return err
	})
	return err
}

// 分页查询订单列表
func (s *CreditsOrderService) GetByPage(userId int64, param request.PageableParam) (*response.PageResponse, error) {
	session := s.WhereUserSession(userId).And("status = ?", enum.OrderStatus(2)).Desc("create_time")
	list := make([]table.CreditsOrder, 0)
	return s.HandlePageable(param, &list, session)
}
