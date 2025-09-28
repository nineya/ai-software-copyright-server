package netdisk

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"encoding/json"
	"fmt"
	"sync"
)

type HelperConfigureService struct {
	service.UserCrudService[table.NetdiskHelperConfigure]
}

var onceHelperConfigure = sync.Once{}
var helperConfigureService *HelperConfigureService

// 获取单例
func GetHelperConfigureService() *HelperConfigureService {
	onceHelperConfigure.Do(func() {
		helperConfigureService = new(HelperConfigureService)
		helperConfigureService.Db = global.DB
	})
	return helperConfigureService
}

// 分页查询列表
func (s *HelperConfigureService) SaveConfigure(userId int64, param table.NetdiskHelperConfigure) error {
	mod := &table.NetdiskHelperConfigure{UserId: userId}
	exist, err := s.Db.Get(mod)
	if err != nil {
		return err
	}
	param.UserId = userId
	if exist {
		_, err = s.WhereUserSession(userId).AllCols().Update(&param)
	} else {
		_, err = s.Db.Insert(&param)
	}
	if err != nil {
		return err
	}
	// 更新客户端的配置信息
	if global.SOCKET.GetClient(userId) != nil { // 如果是有客户端的用户，联网客户端进行查询
		serializedData, _ := json.Marshal(param)
		message := common.SocketMessage{
			Type: enum.SocketMessageType(4),
			Data: string(serializedData),
		}
		_, cliErr := global.SOCKET.SendMessage(userId, message, nil)
		// 查询错误或者没有等到结果
		if cliErr != nil {
			global.LOG.Error(fmt.Sprintf("配置同步网盘助手客户端失败: %+v", cliErr))
		}
	}
	return err
}

// 分页查询列表
func (s *HelperConfigureService) UpdateExpireTime(userId int64, param request.NetdiskHelperUpdateExpireTime) error {
	// 如果没有配置，就添加配置
	mod := &table.NetdiskHelperConfigure{UserId: userId}
	exist, err := s.Db.Get(mod)
	if err != nil {
		return err
	}
	if !exist {
		_, err = s.Db.Insert(&mod)
		if err != nil {
			return err
		}
	}
	// 去除原网盘资源绑定的短链别名
	session := s.Db.NoAutoTime().ID(mod.Id)
	if !param.ExpireTime.IsZero() {
		session.SetExpr("expire_time", param.ExpireTime)
	}
	if !param.WechatExpireTime.IsZero() {
		session.SetExpr("wechat_expire_time", param.WechatExpireTime)
	}
	_, err = session.Update("netdisk_helper_configure")
	return err
}

func (s *HelperConfigureService) GetByUserId(userId int64) (table.NetdiskHelperConfigure, error) {
	mod := &table.NetdiskHelperConfigure{UserId: userId}
	_, err := s.WhereUserSession(userId).Get(mod)
	if mod.Wechat.SearchResultTemplate == "" {
		mod.Wechat.SearchResultTemplate = "为你搜索到了如下结果：\n#{result}\n————————\n直接回复数字序号获取资源！\n搜索到#{total}条数据，当前第#{curPage}/#{totalPage}页，发送“下一页”翻页。\n用小程序搜索可获取全部结果-> #小程序://网盘搜索/9uqYNALJZlYS5Ud"
	}
	if mod.Wechat.SearchFailTemplate == "" {
		mod.Wechat.SearchFailTemplate = "没有搜索到内容，请调整搜索词\n————————\n上小程序搜索更便捷-> #小程序://网盘搜索/9uqYNALJZlYS5Ud"
	}
	if mod.Wechat.SearchResourceTemplate == "" {
		mod.Wechat.SearchResourceTemplate = "#{name}，资源链接：#{url}\n————————\n上小程序搜索更便捷-> #小程序://网盘搜索/9uqYNALJZlYS5Ud"
	}
	if mod.Wechat.SearchResourceFailTemplate == "" {
		mod.Wechat.SearchResourceFailTemplate = "获取资源失败，可以用小程序进行搜索-> #小程序://网盘搜索/9uqYNALJZlYS5Ud"
	}
	return *mod, err
}
