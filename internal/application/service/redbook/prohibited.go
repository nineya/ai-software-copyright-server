package redbook

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/service"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"strings"
	"sync"
)

type ProhibitedService struct {
	service.BaseService
}

var onceProhibited = sync.Once{}
var prohibitedService *ProhibitedService

// 获取单例
func GetProhibitedService() *ProhibitedService {
	onceProhibited.Do(func() {
		prohibitedService = new(ProhibitedService)
		prohibitedService.Db = global.DB
	})
	return prohibitedService
}

func (s *ProhibitedService) Save(param table.RedbookProhibited) error {
	oldMod := &table.RedbookProhibited{Type: param.Type, UserId: param.UserId}
	exist, err := s.Db.Get(oldMod)
	if err != nil {
		return err
	}
	if exist {
		_, err = s.Db.Where("id = ?", oldMod.Id).Update(&table.RedbookProhibited{Words: param.Words})
	} else {
		_, err = s.Db.Insert(&param)
	}
	return err
}

func (s *ProhibitedService) UpdateCustomProhibited(userId int64, param table.RedbookProhibited) error {
	param.SetUserId(userId)
	param.Type = enum.ProhibitedType(3)
	param.Words = utils.MapGetKeys(utils.ListToMap(
		utils.ListFilter(param.Words, func(item string) bool { return item != "" }),
		func(item string) string { return item },
	))
	return s.Save(param)
}

func (s *ProhibitedService) Detection(userId int64, param request.RedbookProhibitedDetectionParam) (*response.RedbookProhibitedDetectionResponse, error) {
	expenseCredits := 4
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}

	content := param.Content
	// 违禁词
	prohibited := &table.RedbookProhibited{}
	_, err = s.Db.Where("type = ?", enum.ProhibitedType(1)).Get(prohibited)
	if err != nil {
		return nil, err
	}
	// 敏感词
	sensitive := &table.RedbookProhibited{}
	_, err = s.Db.Where("type = ?", enum.ProhibitedType(2)).Get(sensitive)
	if err != nil {
		return nil, err
	}
	// 自定义词汇
	custom := &table.RedbookProhibited{}
	customExist, err := s.WhereUserSession(userId).Where("type = ?", enum.ProhibitedType(3)).Get(custom)
	if err != nil {
		return nil, err
	}

	result := &response.RedbookProhibitedDetectionResponse{}

	contentLength := len(content)
	for _, word := range prohibited.Words {
		content = strings.ReplaceAll(content, word, fmt.Sprintf("<span class=\"prohibited\">%s</span>", word))
	}
	// 每替换一次增加32个字符
	result.ProhibitedCount = (len(content) - contentLength) / 32

	contentLength = len(content)
	for _, word := range sensitive.Words {
		content = strings.ReplaceAll(content, word, fmt.Sprintf("<span class=\"sensitive\">%s</span>", word))
	}
	// 每替换一次增加31个字符
	result.SensitiveCount = (len(content) - contentLength) / 31

	if customExist {
		contentLength = len(content)
		for _, word := range custom.Words {
			content = strings.ReplaceAll(content, word, fmt.Sprintf("<span class=\"user\">%s</span>", word))
		}
		// 每替换一次增加26个字符
		result.CustomCount = (len(content) - contentLength) / 26
	}

	result.Content = content

	// 扣款
	user, err := userSev.GetUserService().PaymentNyCredits(userId, enum.BuyType(1), expenseCredits, fmt.Sprintf("购买小红书敏感词检测服务，花费%d币", expenseCredits))
	if err != nil {
		return nil, err
	}
	result.BuyCredits = expenseCredits
	result.BalanceCredits = user.NyCredits
	return result, err
}

func (s *ProhibitedService) GetCustomProhibited(userId int64) (*table.RedbookProhibited, error) {
	mod := &table.RedbookProhibited{Type: enum.ProhibitedType(3), UserId: userId}
	_, err := s.WhereUserSession(userId).Get(mod)
	return mod, err
}
