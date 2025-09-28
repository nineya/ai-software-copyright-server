package short_link

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ShortLinkService struct {
	service.UserCrudService[table.ShortLink]
}

var onceShortLink = sync.Once{}
var shortLinkService *ShortLinkService

// 获取单例
func GetShortLinkService() *ShortLinkService {
	onceShortLink.Do(func() {
		shortLinkService = new(ShortLinkService)
		shortLinkService.Db = global.DB
	})
	return shortLinkService
}

func (s *ShortLinkService) Redirect(userId int64, param request.ShortLinkRedirectParam) (*response.ShortLinkRedirectResponse, error) {
	expenseCredits := 20
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}

	// 取得源链接的别名
	alias, err := getAlias(param.SourceUrl)
	if err != nil {
		return nil, err
	}
	if alias == "" {
		return nil, errors.New("源短链地址不是本站的工具短链")
	}
	// 取得源短链的信息
	mod, err := s.GetByAlias(alias)
	if err != nil {
		return nil, err
	}
	if mod == nil || mod.Id == 0 {
		return nil, errors.New("源短链地址不存在或已失效")
	}
	if mod.UserId != userId {
		return nil, errors.New("源短链地址不属于你")
	}

	// 取得目标链接
	targetUrl, err := s.getIndeedUrl(param.TargetUrl)
	if err != nil {
		return nil, err
	}
	// 更新短链
	_, err = s.WhereUserSession(userId).And("alias = ?", alias).Update(&table.ShortLink{TargetUrl: targetUrl})
	if err != nil {
		return nil, err
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

func (s *ShortLinkService) Create(userId int64, targetUrl string) (*table.ShortLink, error) {
	targetUrl, err := s.getIndeedUrl(targetUrl)
	if err != nil {
		return nil, err
	}
	mod := &table.ShortLink{
		UserId:    userId,
		Alias:     strings.TrimSpace(fmt.Sprintf("%32s", strconv.FormatInt(time.Now().UnixMilli(), 32))),
		TargetUrl: targetUrl,
	}
	_, err = s.Db.Insert(mod)
	return mod, err
}

func (s *ShortLinkService) GetByAlias(alias string) (*table.ShortLink, error) {
	mod := &table.ShortLink{Alias: alias}
	_, err := s.Db.Get(mod)
	return mod, err
}

func (s *ShortLinkService) GetByTargetUrl(userId int64, targetUrl string) (*table.ShortLink, error) {
	mod := &table.ShortLink{UserId: userId, TargetUrl: targetUrl}
	_, err := s.Db.Get(mod)
	return mod, err
}

func (s *ShortLinkService) UpdateVisitsIncreaseById(id int64, userAgent string) error {
	// 爬虫请求不处理，直接返回
	if global.BotReg.MatchString(userAgent) {
		return nil
	}
	_, err := s.Db.ID(id).Incr("visits", 1).NoAutoTime().Update(&table.ShortLink{})
	return err
}

// 取得实际的目标地址
func (s *ShortLinkService) getIndeedUrl(url string) (string, error) {
	alias, err := getAlias(url)
	if err != nil {
		return "", err
	}
	if alias == "" {
		return url, nil
	}
	mod, err := s.GetByAlias(alias)
	if err != nil {
		return "", err
	}
	if mod.TargetUrl == "" {
		return "", errors.New("目标短链接地址不存在或已失效")
	}
	return mod.TargetUrl, nil
}

// 如果是本站的工具短链返回别名
func getAlias(url string) (string, error) {
	reg, err := regexp.Compile("^.+/s/([0-9a-zA-Z]{8,12})$")
	if err != nil {
		return "", err
	}
	if !reg.MatchString(url) {
		return "", nil
	}
	return url[strings.LastIndex(url, "/")+1:], nil
}
