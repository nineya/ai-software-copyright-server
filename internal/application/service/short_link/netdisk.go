package short_link

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	netdSev "ai-software-copyright-server/internal/application/service/netdisk"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

type NetdiskService struct {
	service.UserCrudService[table.ShortLink]
}

var onceNetdisk = sync.Once{}
var netdiskService *NetdiskService

// 获取单例
func GetNetdiskService() *NetdiskService {
	onceNetdisk.Do(func() {
		netdiskService = new(NetdiskService)
		netdiskService.Db = global.DB
	})
	return netdiskService
}

// 先判断目标短链是不是网盘，然后判断短链是用户所有，再通过目标短链取得目标资源，目标资源不存在则创建。
// 去除短链绑定的原资源信息，更新短链，更新新的目标资源的短链信息
func (s *NetdiskService) Redirect(userId int64, param request.ShortLinkRedirectParam) (*response.ShortLinkRedirectResponse, error) {
	expenseCredits := 20
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}

	// 取得目标链接
	targetUrl := param.TargetUrl
	reg, err := regexp.Compile("^(https://pan.quark.cn/s/|https://pan.baidu.com/s/|https://drive.uc.cn/s/|https://pan.xunlei.com/s/|https://caiyun.139.com/m/|https://pan.wkbrowser.com/netdisk/|https://wap.diskyun.com/s/)[\\w\\-@?^=%&/~+#]+$")
	if !reg.MatchString(targetUrl) {
		return nil, errors.New("目标链接不是支持的网盘链接")
	}
	// 取得源短链接的别名
	alias, err := getAlias(param.SourceUrl)
	if err != nil {
		return nil, err
	}
	if alias == "" {
		return nil, errors.New("源短链地址不是本站的短链")
	}
	// 取得源短链的信息
	mod, err := GetShortLinkService().GetByAlias(alias)
	if err != nil {
		return nil, err
	}
	if mod == nil || mod.Id == 0 {
		return nil, errors.New("源短链地址不存在或已失效")
	}
	if mod.UserId != userId {
		return nil, errors.New("源短链地址不属于你")
	}
	// 通过目标链接取得目标资源
	netdisk_resource, err := netdSev.GetResourceService().GetByTargetUrl(userId, targetUrl)
	if err != nil {
		return nil, err
	}
	// 目标资源不存在，创建目标资源
	if netdisk_resource.Id == 0 {
		netdisk_resource = &table.NetdiskResource{
			UserId:    userId,
			TargetUrl: targetUrl,
			Type:      utils.TransformNetdiskType(targetUrl),
			Origin:    enum.NetdiskOrigin(3),
			Status:    enum.NetdiskStatus(1),
		}
		_, err = s.Db.Insert(netdisk_resource)
		if err != nil {
			return nil, err
		}
	}
	// 去除原网盘资源绑定的短链别名
	_, err = s.Db.NoAutoTime().Where("short_link = ?", alias).SetExpr("short_link", "").Update(table.NetdiskResource{})
	if err != nil {
		return nil, err
	}
	// 创建目标地址
	targetUrl = fmt.Sprintf("/ptm/%d.html?id=%d&token=%s", netdisk_resource.Id, userId, url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(targetUrl))))
	// 更新短链
	_, err = s.WhereUserSession(userId).And("alias = ?", alias).Update(&table.ShortLink{TargetUrl: targetUrl})
	if err != nil {
		return nil, err
	}
	// 如果目标资源短链为空，更新短链
	if netdisk_resource.ShortLink == "" {
		netdisk_resource.ShortLink = alias
		err = netdSev.GetResourceService().UpdateById(userId, netdisk_resource.Id, *netdisk_resource)
		if err != nil {
			return nil, err
		}
	}

	result := &response.ShortLinkRedirectResponse{Url: param.SourceUrl}

	// 扣款
	user, err := userSev.GetUserService().PaymentNyCredits(userId, enum.BuyType(10), expenseCredits, fmt.Sprintf("购买短链重定向服务，花费%d币", expenseCredits))
	if err != nil {
		return nil, err
	}
	result.BuyCredits = expenseCredits
	result.BalanceCredits = user.NyCredits

	return result, nil
}

func (s *NetdiskService) Create(userId int64, param request.ShortLinkCreateCloudDiskParam) (*response.UserBuyContentResponse, error) {
	//expenseCredits := 20
	expenseCredits := 10
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}

	configure, err := netdSev.GetShortLinkConfigureService().GetByUserId(userId)
	if err != nil {
		return nil, err
	}
	host := global.Host
	if configure.CustomExpireTime != nil && configure.CustomExpireTime.After(time.Now()) && configure.CustomHost != "" {
		host = configure.CustomHost
	}

	reg, err := regexp.Compile("(https://pan.quark.cn/s/|https://pan.baidu.com/s/|https://drive.uc.cn/s/|https://pan.xunlei.com/s/|https://caiyun.139.com/m/|https://pan.wkbrowser.com/netdisk/|https://wap.diskyun.com/s/)[\\w\\-@?^=%&/~+#]+")
	content := param.Content

	replaceFunc := func(u string) string {
		netdisk_resource, err := netdSev.GetResourceService().GetByTargetUrl(userId, u)
		if err != nil {
			return u
		}
		// 资源不存在，创建网盘资源
		if netdisk_resource.Id == 0 {
			netdisk_resource = &table.NetdiskResource{
				UserId:    userId,
				TargetUrl: u,
				Type:      utils.TransformNetdiskType(u),
				Origin:    enum.NetdiskOrigin(3),
				Status:    enum.NetdiskStatus(1),
			}
			_, err = s.Db.Insert(netdisk_resource)
			if err != nil {
				return u
			}
		}
		if netdisk_resource.ShortLink == "" {
			encodedData := fmt.Sprintf("/ptm/%d.html?id=%d&token=%s", netdisk_resource.Id, userId, url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(u))))
			short_link, err := GetShortLinkService().Create(userId, encodedData)
			// 如果插入数据库失败，就不替换链接
			if err != nil {
				return u
			}
			netdisk_resource.ShortLink = short_link.Alias
			err = netdSev.GetResourceService().UpdateById(userId, netdisk_resource.Id, *netdisk_resource)
			if err != nil {
				return u
			}
		}
		return host + "/s/" + netdisk_resource.ShortLink
	}
	content = reg.ReplaceAllStringFunc(content, replaceFunc)
	result := &response.UserBuyContentResponse{Content: content}

	// 扣款
	//user, err := userSev.GetUserService().PaymentNyCredits(userId, enum.BuyType(8), expenseCredits, fmt.Sprintf("购买网盘链接转换服务，花费%d币", expenseCredits))
	user, err := userSev.GetUserService().PaymentNyCredits(userId, enum.BuyType(8), expenseCredits, "网盘转链限时特惠，本次操作五折，仅花费10币！")
	if err != nil {
		return nil, err
	}
	result.BuyCredits = expenseCredits
	result.BalanceCredits = user.NyCredits
	result.BuyMessage = "限时特惠，本次操作五折，仅花费10币！"

	return result, nil
}

func (s *NetdiskService) Statistic(userId, id int64) (*response.ShortLinkStatisticResponse, error) {
	expenseCredits := 20
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}

	mod, err := GetShortLinkService().GetById(userId, id)
	if err != nil {
		return nil, err
	}
	if mod == nil || mod.Id == 0 {
		return nil, errors.New("短链不存在或已失效")
	}

	// 今天访问数据
	today := table.ShortLinkStatistic{}
	_, err = s.Db.Select("count(ip_address) as pv, count(DISTINCT ip_address) as uv, SUM(@is_mobile := user_agent REGEXP 'phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone') as mobile, COUNT(DISTINCT CASE WHEN @is_mobile THEN ip_address END) AS mobile_uv, COUNT(DISTINCT CASE WHEN NOT @is_mobile THEN ip_address END) AS pc_uv").
		Where("url = ? and create_time >= CURDATE() and user_agent not regexp 'okhttp|Go-http-client|python-requests|spider|bot'", "/s/"+mod.Alias).
		Get(&today)
	if err != nil {
		return nil, err
	}

	// 今天访问来源
	todayOrigins := make([]table.ShortLinkStatistic, 0)
	err = s.Db.Select("origin, count(ip_address) as pv, count(DISTINCT ip_address) as uv, SUM(@is_mobile := user_agent REGEXP 'phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone') as mobile, COUNT(DISTINCT CASE WHEN @is_mobile THEN ip_address END) AS mobile_uv, COUNT(DISTINCT CASE WHEN NOT @is_mobile THEN ip_address END) AS pc_uv").
		Where("url = ? and create_time >= CURDATE() and user_agent not regexp 'okhttp|Go-http-client|python-requests|spider|bot'", "/s/"+mod.Alias).
		GroupBy("origin").Desc("pv").Find(&todayOrigins)
	if err != nil {
		return nil, err
	}

	total := table.ShortLinkStatistic{}
	_, err = s.Db.Select("count(ip_address) as pv, count(DISTINCT ip_address) as uv, SUM(@is_mobile := user_agent REGEXP 'phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone') as mobile, COUNT(DISTINCT CASE WHEN @is_mobile THEN ip_address END) AS mobile_uv, COUNT(DISTINCT CASE WHEN NOT @is_mobile THEN ip_address END) AS pc_uv").
		Where("url = ? and create_time >= DATE_SUB(CURDATE(), INTERVAL 30 DAY) and user_agent not regexp 'okhttp|Go-http-client|python-requests|spider|bot'", "/s/"+mod.Alias).
		Get(&total)
	if err != nil {
		return nil, err
	}

	totalOrigins := make([]table.ShortLinkStatistic, 0)
	err = s.Db.Select("origin, count(ip_address) as pv, count(DISTINCT ip_address) as uv, SUM(@is_mobile := user_agent REGEXP 'phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone') as mobile, COUNT(DISTINCT CASE WHEN @is_mobile THEN ip_address END) AS mobile_uv, COUNT(DISTINCT CASE WHEN NOT @is_mobile THEN ip_address END) AS pc_uv").
		Where("url = ? and create_time >= DATE_SUB(CURDATE(), INTERVAL 30 DAY) and user_agent not regexp 'okhttp|Go-http-client|python-requests|spider|bot'", "/s/"+mod.Alias).
		GroupBy("origin").Desc("pv").Find(&totalOrigins)
	if err != nil {
		return nil, err
	}

	days := make([]table.ShortLinkStatistic, 0)
	err = s.Db.Select("date_format( create_time, '%Y-%m-%d' ) AS date, count(ip_address) as pv, count(DISTINCT ip_address) as uv, SUM(@is_mobile := user_agent REGEXP 'phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone') as mobile, COUNT(DISTINCT CASE WHEN @is_mobile THEN ip_address END) AS mobile_uv, COUNT(DISTINCT CASE WHEN NOT @is_mobile THEN ip_address END) AS pc_uv").
		Where("url = ? and create_time >= DATE_SUB(CURDATE(), INTERVAL 30 DAY) and user_agent not regexp 'okhttp|Go-http-client|python-requests|spider|bot'", "/s/"+mod.Alias).
		GroupBy("date").Desc("date").Find(&days)
	if err != nil {
		return nil, err
	}

	result := &response.ShortLinkStatisticResponse{Today: today, TodayOrigins: todayOrigins, Total: total, TotalOrigins: totalOrigins, Days: days}

	// 扣款
	user, err := userSev.GetUserService().PaymentNyCredits(userId, enum.BuyType(17), expenseCredits, fmt.Sprintf("购买短链分析服务，花费%d币", expenseCredits))
	if err != nil {
		return nil, err
	}
	result.BuyCredits = expenseCredits
	result.BalanceCredits = user.NyCredits

	return result, err
}

// 简单的今日访问短链数据
func (s *NetdiskService) TodayVisits(userId int64) (table.ShortLinkStatistic, error) {
	// 今天访问数据
	today := table.ShortLinkStatistic{}
	_, err := s.Db.Select("count(ip_address) as pv, count(DISTINCT ip_address) as uv, SUM(@is_mobile := user_agent REGEXP 'phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone') as mobile, COUNT(DISTINCT CASE WHEN @is_mobile THEN ip_address END) AS mobile_uv, COUNT(DISTINCT CASE WHEN NOT @is_mobile THEN ip_address END) AS pc_uv").
		Where("url in (select concat('/s/', alias) from short_link where user_id = ?) and create_time >= CURDATE() and user_agent not regexp 'okhttp|Go-http-client|python-requests|spider|bot'", userId).
		Get(&today)
	return today, err
}

// 短链活跃数据
func (s *NetdiskService) ActiveVisits(userId int64) (table.ShortLinkActive, error) {
	// 今天访问数据
	today := table.ShortLinkActive{}
	_, err := s.WhereUserSession(userId).
		Select("count(*) as total, count((SELECT 1 FROM statistic s WHERE s.url = CONCAT('/s/', alias) AND s.create_time >= CURDATE() LIMIT 1)) as today_active, count((SELECT 1 FROM statistic s WHERE s.url = CONCAT('/s/', alias) AND s.create_time >= DATE_SUB(CURDATE(), INTERVAL 30 DAY) LIMIT 1)) as total_active").
		Get(&today)
	return today, err
}

// 所有短链总的访问数据
func (s *NetdiskService) AllStatistic(userId int64) (*response.ShortLinkStatisticResponse, error) {
	expenseCredits := 20
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}

	// 今天访问数据
	today := table.ShortLinkStatistic{}
	_, err = s.Db.Select("count(ip_address) as pv, count(DISTINCT ip_address) as uv, SUM(@is_mobile := user_agent REGEXP 'phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone') as mobile, COUNT(DISTINCT CASE WHEN @is_mobile THEN ip_address END) AS mobile_uv, COUNT(DISTINCT CASE WHEN NOT @is_mobile THEN ip_address END) AS pc_uv").
		Where("url in (select concat('/s/', alias) from short_link where user_id = ?) and create_time >= CURDATE() and user_agent not regexp 'okhttp|Go-http-client|python-requests|spider|bot'", userId).
		Get(&today)
	if err != nil {
		return nil, err
	}

	// 今天访问来源
	todayOrigins := make([]table.ShortLinkStatistic, 0)
	err = s.Db.Select("origin, count(ip_address) as pv, count(DISTINCT ip_address) as uv, SUM(@is_mobile := user_agent REGEXP 'phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone') as mobile, COUNT(DISTINCT CASE WHEN @is_mobile THEN ip_address END) AS mobile_uv, COUNT(DISTINCT CASE WHEN NOT @is_mobile THEN ip_address END) AS pc_uv").
		Where("url in (select concat('/s/', alias) from short_link where user_id = ?) and create_time >= CURDATE() and user_agent not regexp 'okhttp|Go-http-client|python-requests|spider|bot'", userId).
		GroupBy("origin").Desc("pv").Find(&todayOrigins)
	if err != nil {
		return nil, err
	}

	total := table.ShortLinkStatistic{}
	_, err = s.Db.Select("count(ip_address) as pv, count(DISTINCT ip_address) as uv, SUM(@is_mobile := user_agent REGEXP 'phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone') as mobile, COUNT(DISTINCT CASE WHEN @is_mobile THEN ip_address END) AS mobile_uv, COUNT(DISTINCT CASE WHEN NOT @is_mobile THEN ip_address END) AS pc_uv").
		Where("url in (select concat('/s/', alias) from short_link where user_id = ?) and create_time >= DATE_SUB(CURDATE(), INTERVAL 30 DAY) and user_agent not regexp 'okhttp|Go-http-client|python-requests|spider|bot'", userId).
		Get(&total)
	if err != nil {
		return nil, err
	}

	totalOrigins := make([]table.ShortLinkStatistic, 0)
	err = s.Db.Select("origin, count(ip_address) as pv, count(DISTINCT ip_address) as uv, SUM(@is_mobile := user_agent REGEXP 'phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone') as mobile, COUNT(DISTINCT CASE WHEN @is_mobile THEN ip_address END) AS mobile_uv, COUNT(DISTINCT CASE WHEN NOT @is_mobile THEN ip_address END) AS pc_uv").
		Where("url in (select concat('/s/', alias) from short_link where user_id = ?) and create_time >= DATE_SUB(CURDATE(), INTERVAL 30 DAY) and user_agent not regexp 'okhttp|Go-http-client|python-requests|spider|bot'", userId).
		GroupBy("origin").Desc("pv").Find(&totalOrigins)
	if err != nil {
		return nil, err
	}

	days := make([]table.ShortLinkStatistic, 0)
	err = s.Db.Select("date_format( create_time, '%Y-%m-%d' ) AS date, count(ip_address) as pv, count(DISTINCT ip_address) as uv, SUM(@is_mobile := user_agent REGEXP 'phone|pad|pod|iPhone|iPod|ios|iPad|Android|Mobile|BlackBerry|IEMobile|MQQBrowser|JUC|Fennec|wOSBrowser|BrowserNG|WebOS|Symbian|Windows Phone') as mobile, COUNT(DISTINCT CASE WHEN @is_mobile THEN ip_address END) AS mobile_uv, COUNT(DISTINCT CASE WHEN NOT @is_mobile THEN ip_address END) AS pc_uv").
		Where("url in (select concat('/s/', alias) from short_link where user_id = ?) and create_time >= DATE_SUB(CURDATE(), INTERVAL 30 DAY) and user_agent not regexp 'okhttp|Go-http-client|python-requests|spider|bot'", userId).
		GroupBy("date").Desc("date").Find(&days)
	if err != nil {
		return nil, err
	}

	result := &response.ShortLinkStatisticResponse{Today: today, TodayOrigins: todayOrigins, Total: total, TotalOrigins: totalOrigins, Days: days}

	// 扣款
	user, err := userSev.GetUserService().PaymentNyCredits(userId, enum.BuyType(17), expenseCredits, fmt.Sprintf("购买短链分析服务，花费%d币", expenseCredits))
	if err != nil {
		return nil, err
	}
	result.BuyCredits = expenseCredits
	result.BalanceCredits = user.NyCredits

	return result, err
}

// 分页查询列表
func (s *NetdiskService) GetByPage(userId int64, param request.QueryPageParam) (*response.PageResponse, error) {
	host := global.Host
	configure, _ := netdSev.GetShortLinkConfigureService().GetByUserId(userId)
	if configure.CustomExpireTime != nil && configure.CustomExpireTime.After(time.Now()) && configure.CustomHost != "" {
		host = configure.CustomHost
	}
	session := s.WhereUserSessionByTable("netdisk_resource", userId).
		Select("short_link.*, netdisk_resource.user_name, netdisk_resource.name, netdisk_resource.target_url netdisk_target_url, netdisk_resource.type, netdisk_resource.status").
		Join("INNER", "short_link", "netdisk_resource.user_id = short_link.user_id and netdisk_resource.short_link = short_link.alias").
		And("short_link != ''").Desc("update_time").Asc("id")
	if param.Keyword != "" {
		lastSlash := strings.LastIndex(param.Keyword, "/")
		if lastSlash > 0 {
			session.And("(short_link = ? or netdisk_resource.target_url like concat('%',?,'%'))", param.Keyword[lastSlash+1:], param.Keyword)
		} else {
			session.And("(netdisk_resource.target_url like concat('%',?,'%'))", param.Keyword)
		}
	}
	list := make([]table.NetdiskShortLink, 0)
	resp, err := s.HandlePageable(param.PageableParam, &list, session)
	if err != nil {
		return nil, err
	}
	for i, _ := range list {
		list[i].ShortLinkUrl = host + "/s/" + list[i].Alias
	}
	return resp, nil
}

// 分页查询列表
func (s *NetdiskService) GetByLast(count int) ([]table.NetdiskShortLink, error) {
	list := make([]table.NetdiskShortLink, 0)
	err := s.Db.
		Select("short_link.*, netdisk_resource.user_name, netdisk_resource.name, netdisk_resource.target_url netdisk_target_url, netdisk_resource.type, netdisk_resource.status").
		Join("INNER", "short_link", "netdisk_resource.user_id = short_link.user_id and netdisk_resource.short_link = short_link.alias").
		And("short_link != '' and name != '' and status = 1").Limit(count, 0).Desc("update_time").
		Find(&list)
	return list, err
}
