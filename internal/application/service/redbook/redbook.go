package redbook

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	mailPlugin "ai-software-copyright-server/internal/application/plugin/mail"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
	"unsafe"
)

type RedbookService struct {
	Client http.Client
}

var onceRedbook = sync.Once{}
var redbookService *RedbookService

// 获取单例
func GetRedbookService() *RedbookService {
	onceRedbook.Do(func() {
		redbookService = new(RedbookService)
		redbookService.Client = http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				//fmt.Printf("Redirecting to: %s\n", req.URL)
				// 重定向时手工拷贝cookie
				req.Header.Set("cookie", via[0].Header.Get("cookie"))
				return nil // 返回nil允许重定向
			},
		}
	})
	return redbookService
}

func (s *RedbookService) RemoveWatermark(userId int64, url string) (*response.RedbookRemoveWatermarkResponse, error) {
	expenseCredits := 20
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}

	// 获取cookie
	var html string
	cookies, err := GetCookieService().GetAllByNormal()
	if err != nil {
		return nil, err
	}
	if len(cookies) == 0 {
		return nil, errors.New("没有有效 cookie")
	}
	for i, cookie := range cookies {
		str, err := s.SendRequest(url, &cookie)
		if err != nil {
			return nil, err
		}
		if strings.Contains(*str, cookie.XhsUserId) {
			// 登录成功，把之前登录失败的设置为失效cookie，然后返回
			for j := 0; j < i; j++ {
				global.LOG.Warn(fmt.Sprintf("小红书修改 Cookie 状态为失效, Cookie = %d", cookies[j].Id))
				GetCookieService().InnerUpdateStatusById(cookies[j].Id, enum.CookieStatus(2))
			}
			html = *str
			break
		}
		global.LOG.Error(fmt.Sprintf("小红书未登录，Cookie = %d, Result = %s", cookie.Id, str))
		// 十条全部登录失效，表明是环境问题，不做处理
		if i == 10 {
			return nil, errors.New("功能暂时维护中，请稍后再试……")
		}
	}

	reg := regexp.MustCompile(`<meta\s+name=\"og:image\"\s+content=\"(//|http://|https://)([^\"]*)\">`)
	urls := reg.FindAllString(html, -1)
	for index := range urls {
		urls[index] = reg.ReplaceAllString(urls[index], "https://$2")
	}
	result := &response.RedbookRemoveWatermarkResponse{Urls: urls}

	// 扣款
	user, err := userSev.GetUserService().PaymentNyCredits(userId, enum.BuyType(4), expenseCredits, fmt.Sprintf("购买小红书去水印服务，花费%d币", expenseCredits))
	if err != nil {
		return nil, err
	}
	result.BuyCredits = expenseCredits
	result.BalanceCredits = user.NyCredits
	return result, nil
}

func (s *RedbookService) Valuation(userId int64, param request.RedbookParam) (*response.RedbookValuationResponse, error) {
	expenseCredits := 30
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}

	info, err := s.GetProfileInfo(param.Html, param.Url)
	if err != nil {
		return nil, err
	}
	userInfo := info.User.UserPageData
	result := &response.RedbookValuationResponse{
		RedbookProfileInfoUserBasicInfo: userInfo.BasicInfo,
	}
	for _, inte := range userInfo.Interactions {
		switch inte.Type {
		case "follows": // 关注
			result.FollowCount = inte.Count
		case "fans":
			result.FansCount = inte.Count
		case "interaction":
			result.InteractionCount = inte.Count
		}
	}
	followCount, err := ConversionNumber(result.FollowCount)
	if err != nil {
		return nil, err
	}
	fansCount, err := ConversionNumber(result.FansCount)
	if err != nil {
		return nil, err
	}
	interactionCount, err := ConversionNumber(result.InteractionCount)
	if err != nil {
		return nil, err
	}
	// 等于基础价20，粉丝-关注的一半，粉丝数的一半 其中的最大值，最后 * 0.2 单价
	price := 20 + math.Max(fansCount-(followCount/2.0), fansCount*0.5)*0.2
	// 再乘以赞粉比/5
	price *= math.Max(math.Min(7, interactionCount/math.Max(fansCount, 1)), 3) / 5
	for _, note := range info.User.Notes[0] {
		like, err := ConversionNumber(note.NoteCard.InteractInfo.LikedCount)
		if err != nil {
			return nil, err
		}
		switch {
		case like <= 10: // -(0~10)
			price -= 10 - float64(like)
		case like <= 50: // +(2~10)
			price += float64(like) * 0.2
		case like <= 100: // +(10~15)
			price += 10 + (float64(like)-50)*0.1
		default: // +(15~500)
			price += math.Min(15+(float64(like)-100)*0.05, 500)
		}
	}
	if price < 0 {
		price = 20
	}
	price *= 1 + ((rand.Float64() - 0.5) * 0.1)
	result.Price = fmt.Sprintf("%.2f", price)

	// 扣款
	user, err := userSev.GetUserService().PaymentNyCredits(userId, enum.BuyType(2), expenseCredits, fmt.Sprintf("购买小红书账号估值服务，花费%d币", expenseCredits))
	if err != nil {
		return nil, err
	}
	result.BuyCredits = expenseCredits
	result.BalanceCredits = user.NyCredits
	return result, err
}

func (s *RedbookService) Weight(userId int64, param request.RedbookParam) (*response.RedbookWeightResponse, error) {
	expenseCredits := 40
	// 预检余额
	_, err := userSev.GetUserService().GetAndCheckBalance(userId, expenseCredits)
	if err != nil {
		return nil, err
	}

	info, err := s.GetProfileInfo(param.Html, param.Url)
	if err != nil {
		return nil, err
	}
	userInfo := info.User.UserPageData
	userBasicInfo := userInfo.BasicInfo
	result := &response.RedbookWeightResponse{
		Info: response.RedbookWeightInfoResponse{RedbookProfileInfoUserBasicInfo: userBasicInfo},
	}

	var follows string     // 关注
	var fans string        // 粉丝
	var interaction string //赞收
	for _, inte := range userInfo.Interactions {
		switch inte.Type {
		case "follows": // 关注
			follows = inte.Count
			result.Info.FollowCount = inte.Count
		case "fans":
			fans = inte.Count
			result.Info.FansCount = inte.Count
		case "interaction":
			interaction = inte.Count
			result.Info.InteractionCount = inte.Count
		}
	}
	followCount, err := ConversionNumber(follows)
	if err != nil {
		return nil, err
	}
	fansCount, err := ConversionNumber(fans)
	if err != nil {
		return nil, err
	}
	interactionCount, err := ConversionNumber(interaction)
	if err != nil {
		return nil, err
	}
	interFansRate := interactionCount / math.Max(fansCount, 1) // 赞粉比
	fansFollowRate := fansCount / math.Max(followCount, 1)     // 粉关比

	noteCount := len(info.User.Notes[0])
	hotCount := 0
	likes := utils.ListTransform(info.User.Notes[0], func(item response.RedbookProfileInfoNoteItem) float64 {
		like, _ := ConversionNumber(item.NoteCard.InteractInfo.LikedCount)
		if like > 1000 {
			hotCount++
		}
		return like
	})
	hotRate := float64(hotCount) * 100 / math.Max(float64(noteCount), 1)
	// 排序
	sort.Float64s(likes)
	// 用中位数算最低下限和流量池
	var floorlLike float64
	var flow float64
	if noteCount > 0 {
		if noteCount%2 == 0 {
			floorlLike = (likes[(noteCount-1)/2] + likes[noteCount/2]) / 10
			flow = (likes[(noteCount-1)/2] + likes[noteCount/2]) * 14
		} else {
			floorlLike = likes[noteCount/2] / 5
			flow = likes[noteCount/2] * 28
		}
	}
	// 限流数量
	var restrictCount int
	for _, like := range likes {
		if like < floorlLike {
			restrictCount++
		}
	}
	if restrictCount > 1 {
		restrictCount--
	}

	// 得分
	score := 60.0

	// 检查用户名
	result.Nickname = response.RedbookWeightItemResponse{
		Value: userBasicInfo.Nickname,
	}
	if regexp.MustCompile(`^小红薯[-_a-zA-Z0-9]+$`).MatchString(userBasicInfo.Nickname) {
		result.Nickname.Hint = "建议更换默认昵称"
		result.Nickname.Level = enum.HintLevel(4)
		score -= 5
	} else if regexp.MustCompile(`^[-_a-zA-Z0-9]+$`).MatchString(userBasicInfo.Nickname) {
		result.Nickname.Hint = "建议优化"
		result.Nickname.Level = enum.HintLevel(4)
		score -= 5
	} else {
		result.Nickname.Hint = "名字正常"
		result.Nickname.Level = enum.HintLevel(2)
		score += 2
	}

	// 粉丝
	result.Fans = response.RedbookWeightItemResponse{
		Value: fans,
	}
	switch {
	case fansCount < 200:
		result.Fans.Hint = "粉丝较低"
		result.Fans.Level = enum.HintLevel(4)
	case fansCount < 2000:
		result.Fans.Hint = "粉丝一般"
		result.Fans.Level = enum.HintLevel(3)
		score += fansCount / 500
	case fansCount < 5000:
		result.Fans.Hint = "粉丝活跃"
		result.Fans.Level = enum.HintLevel(2)
		score += 4 + (fansCount / 1000)
	default:
		result.Fans.Hint = "粉丝优质"
		result.Fans.Level = enum.HintLevel(2)
		score += 10 + (fansCount / 10000)
	}

	// 赞粉数
	result.Interaction = response.RedbookWeightItemResponse{
		Value: interaction,
	}
	switch {
	case interFansRate < 1:
		result.Interaction.Hint = fmt.Sprintf("赞粉比低: %.2f", interFansRate)
		result.Interaction.Level = enum.HintLevel(4)
		score *= 0.7
	case interFansRate < 5:
		result.Interaction.Hint = fmt.Sprintf("赞粉比一般: %.2f", interFansRate)
		result.Interaction.Level = enum.HintLevel(3)
		score *= 0.9
	default:
		result.Interaction.Hint = fmt.Sprintf("赞粉比高: %.2f", interFansRate)
		result.Interaction.Level = enum.HintLevel(2)
		score *= 1.05
	}

	// 关注人数
	result.Follow = response.RedbookWeightItemResponse{
		Value: follows,
	}
	switch {
	case fansFollowRate < 5:
		result.Follow.Hint = fmt.Sprintf("粉关比低: %.2f", fansFollowRate)
		result.Follow.Level = enum.HintLevel(4)
		score *= 0.9
	case fansFollowRate < 20:
		result.Follow.Hint = fmt.Sprintf("粉关比一般: %.2f", fansFollowRate)
		result.Follow.Level = enum.HintLevel(3)
	default:
		result.Follow.Hint = fmt.Sprintf("粉关比高: %.2f", fansFollowRate)
		result.Follow.Level = enum.HintLevel(2)
		score *= 1.05
	}

	// 热门
	result.Hot = response.RedbookWeightItemResponse{
		Value: strconv.Itoa(hotCount),
		Hint:  fmt.Sprintf("近期热门率: %.1f%%", hotRate),
	}
	switch {
	case hotRate < 10:
		result.Hot.Level = enum.HintLevel(4)
	case fansFollowRate < 40:
		result.Hot.Level = enum.HintLevel(3)
		score *= 1.05
	default:
		result.Hot.Level = enum.HintLevel(2)
		score *= 1.1
	}

	// 限流
	result.Restrict = response.RedbookWeightItemResponse{
		Value: strconv.Itoa(restrictCount),
	}
	switch {
	case restrictCount > 1:
		result.Restrict.Hint = "建议优化"
		result.Restrict.Level = enum.HintLevel(4)
		score *= 0.9
	default:
		result.Restrict.Hint = "无需优化"
		result.Restrict.Level = enum.HintLevel(2)
	}

	// 简介
	result.Desc = response.RedbookWeightItemResponse{
		Value: userBasicInfo.Desc,
	}
	if userBasicInfo.Desc == "" || userBasicInfo.Desc == "还没有简介" {
		result.Desc.Hint = "没有简介"
		result.Desc.Level = enum.HintLevel(4)
		score *= 0.9
	} else if utf8.RuneCountInString(userBasicInfo.Desc) < 20 {
		result.Desc.Hint = "简介内容少"
		result.Desc.Level = enum.HintLevel(3)
		score *= 0.95
	} else if regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+.[a-zA-Z]{2,}`).MatchString(userBasicInfo.Desc) {
		result.Desc.Hint = "包含引流内容"
		result.Desc.Level = enum.HintLevel(4)
		score *= 0.9
	} else {
		result.Desc.Hint = "简介正常"
		result.Desc.Level = enum.HintLevel(2)
	}

	// 流量
	switch {
	case flow < 200:
		result.Flow = response.RedbookWeightItemResponse{
			Value: "0-200",
			Hint:  "百人流量池，流量差",
			Level: enum.HintLevel(4),
		}
	case flow < 900:
		result.Flow = response.RedbookWeightItemResponse{
			Value: "200-900",
			Hint:  "百人流量池，流量较差",
			Level: enum.HintLevel(3),
		}
	case flow < 5000:
		result.Flow = response.RedbookWeightItemResponse{
			Value: "900-5000",
			Hint:  "千人流量池，流量一般",
			Level: enum.HintLevel(3),
		}
	case flow < 9000:
		result.Flow = response.RedbookWeightItemResponse{
			Value: "5000-9000",
			Hint:  "千人流量池，流量较好",
			Level: enum.HintLevel(2),
		}
	case flow < 50000:
		result.Flow = response.RedbookWeightItemResponse{
			Value: "9000-5W",
			Hint:  "万人流量池，流量较好",
			Level: enum.HintLevel(2),
		}
	case flow < 90000:
		result.Flow = response.RedbookWeightItemResponse{
			Value: "5W-9W",
			Hint:  "万人流量池，流量好",
			Level: enum.HintLevel(2),
		}
	case flow < 500000:
		result.Flow = response.RedbookWeightItemResponse{
			Value: "9W-50W",
			Hint:  "十万流量池，流量好",
			Level: enum.HintLevel(2),
		}
	case flow < 1000000:
		result.Flow = response.RedbookWeightItemResponse{
			Value: "50W-100W",
			Hint:  "十万流量池，流量好",
			Level: enum.HintLevel(2),
		}
	case flow < 5000000:
		result.Flow = response.RedbookWeightItemResponse{
			Value: "100W-500W",
			Hint:  "百万流量池，流量极好",
			Level: enum.HintLevel(2),
		}
	case flow < 10000000:
		result.Flow = response.RedbookWeightItemResponse{
			Value: "500W-1000W",
			Hint:  "百万流量池，流量极好",
			Level: enum.HintLevel(2),
		}
	case flow < 50000000:
		result.Flow = response.RedbookWeightItemResponse{
			Value: "500W-1000W",
			Hint:  "百万流量池，流量极好",
			Level: enum.HintLevel(2),
		}
	case flow < 50000000:
		result.Flow = response.RedbookWeightItemResponse{
			Value: "1000W-5000W",
			Hint:  "千万流量池，流量极好",
			Level: enum.HintLevel(2),
		}
	default:
		result.Flow = response.RedbookWeightItemResponse{
			Value: "5000W+",
			Hint:  "千万流量池，流量极好",
			Level: enum.HintLevel(2),
		}
	}

	// 发帖时间
	result.PostTime = response.RedbookWeightItemResponse{
		Value: "18:00-20:00",
		Hint:  "建议发帖时间",
		Level: enum.HintLevel(1),
	}

	// 直播时间
	result.LiveTime = response.RedbookWeightItemResponse{
		Value: "19:00-21:00",
		Hint:  "建议直播时间",
		Level: enum.HintLevel(1),
	}

	if score > 95 {
		score = 95 + (score / 100)
	}
	if userInfo.UserAccountStatus.Type == 1 {
		score = 5
	}
	result.Score = int(score)
	switch {
	case score < 30: // 5 - 10
		result.ScoreMsg = fmt.Sprintf("已经超越 %.2f%% 的博主", 5+(score*5/30))
	case score < 40: // 10-40
		result.ScoreMsg = fmt.Sprintf("已经超越 %.2f%% 的博主", 10+(score*30/40))
	case score < 50: // 40-70
		result.ScoreMsg = fmt.Sprintf("已经超越 %.2f%% 的博主", 40+(score*30/50))
	case score < 60: // 70-80
		result.ScoreMsg = fmt.Sprintf("已经超越 %.2f%% 的博主", 70+(score*10/60))
	case score < 70: // 80-88
		result.ScoreMsg = fmt.Sprintf("已经超越 %.2f%% 的博主", 80+(score*8/70))
	case score < 80: // 88-94
		result.ScoreMsg = fmt.Sprintf("已经超越 %.2f%% 的博主", 88+(score*6/80))
	case score < 90: // 94-97
		result.ScoreMsg = fmt.Sprintf("已经超越 %.2f%% 的博主", 94+(score*3/90))
	default: // 97-99.8
		result.ScoreMsg = fmt.Sprintf("已经超越 %.2f%% 的博主", 97+(score*2.98/100))
	}

	// 质量
	switch {
	case userInfo.UserAccountStatus.Type != 0:
		result.Quality = response.RedbookWeightItemResponse{
			Value: userInfo.UserAccountStatus.Toast,
			Hint:  "建议换号或申述",
			Level: enum.HintLevel(4),
		}
	case score < 40:
		result.Quality = response.RedbookWeightItemResponse{
			Value: "僵尸号",
			Hint:  "建议注销或重新养号",
			Level: enum.HintLevel(4),
		}
	case score < 50:
		result.Quality = response.RedbookWeightItemResponse{
			Value: "低权重号",
			Hint:  "建议增发视频笔记",
			Level: enum.HintLevel(3),
		}
	case score < 70:
		result.Quality = response.RedbookWeightItemResponse{
			Value: "普通号",
			Hint:  "建议增加活跃度",
			Level: enum.HintLevel(3),
		}
	case score < 80:
		result.Quality = response.RedbookWeightItemResponse{
			Value: "中热号",
			Hint:  "建议加强用户停留",
			Level: enum.HintLevel(3),
		}
	case score < 90:
		result.Quality = response.RedbookWeightItemResponse{
			Value: "优品号",
			Hint:  "建议加强用户互动",
			Level: enum.HintLevel(2),
		}
	default:
		result.Quality = response.RedbookWeightItemResponse{
			Value: "优质号",
			Hint:  "继续保持",
			Level: enum.HintLevel(2),
		}
	}

	// 扣款
	user, err := userSev.GetUserService().PaymentNyCredits(userId, enum.BuyType(3), expenseCredits, fmt.Sprintf("购买小红书账号权重检测服务，花费%d币", expenseCredits))
	if err != nil {
		return nil, err
	}
	result.BuyCredits = expenseCredits
	result.BalanceCredits = user.NyCredits
	return result, err
}

func (s *RedbookService) GetProfileInfo(html, url string) (*response.RedbookProfileInfoResponse, error) {
	// 定义要替换的正则表达式
	regex := regexp.MustCompile(`^[\s\S]*__INITIAL_STATE__\s*=\s*(\{.+\})[\s\S]*`)

	var result response.RedbookProfileInfoResponse // 反序列化JSON到结构体

	if html == "" {
		re := regexp.MustCompile(`(http|https)://(www.xiaohongshu.com|xhslink.com)[A-Za-z0-9_\-+.:?&@=/%#,;]*`)
		url = re.FindString(url)
		if url == "" {
			return nil, errors.New("错误的博主主页链接")
		}
		// 获取cookie
		cookies, err := GetCookieService().GetAllByNormal()
		if err != nil {
			return nil, err
		}
		if len(cookies) == 0 {
			return nil, errors.New("没有有效 cookie")
		}
		for i, cookie := range cookies {
			str, err := s.SendRequest(url, &cookie)
			if err != nil {
				return nil, err
			}
			jsonStr := regex.ReplaceAllString(*str, "$1")
			jsonStr = strings.ReplaceAll(jsonStr, "undefined", "null")
			err = json.Unmarshal([]byte(jsonStr), &result)
			if err != nil {
				return nil, err
			}
			if result.User.LoggedIn {
				// 登录成功，把之前登录失败的设置为失效cookie，然后返回
				for j := 0; j < i; j++ {
					global.LOG.Warn(fmt.Sprintf("小红书修改 Cookie 状态为失效, Cookie = %d", cookies[j].Id))
					mailErr := mailPlugin.GetMailPlugin().SendHtmlMail(request.MailParam{
						To:      global.CONFIG.Plugin.Mail.AdminMail,
						Subject: "小红书Cookie被设置为失效，请检查！",
						Content: fmt.Sprintf("小红书Cookie（%s）被设置为失效，请检查！", cookies[j].Nickname),
					})
					if mailErr != nil {
						global.LOG.Warn(fmt.Sprintf("小红书提示邮件发送失败：%+v", mailErr))
					}
					_ = GetCookieService().InnerUpdateStatusById(cookies[j].Id, enum.CookieStatus(2))
				}
				return &result, nil
			}
			global.LOG.Error(fmt.Sprintf("小红书未登录，Cookie = %d, Result = %s", cookie.Id, jsonStr))
			// 十条全部登录失效，表明是环境问题，不做处理
			if i == 10 {
				return nil, errors.New("功能暂时维护中，请稍后再试……")
			}
		}
		mailErr := mailPlugin.GetMailPlugin().SendHtmlMail(request.MailParam{
			To:      global.CONFIG.Plugin.Mail.AdminMail,
			Subject: "小红书所有Cookie皆已失效，请检查！",
			Content: "小红书已经没有可用Cookie，请检查Cookie情况！",
		})
		if mailErr != nil {
			global.LOG.Warn(fmt.Sprintf("小红书提示邮件发送失败：%+v", mailErr))
		}
		// 登录成功，把之前登录失败的设置为失效cookie，然后返回
		for _, cookie := range cookies {
			global.LOG.Warn(fmt.Sprintf("小红书修改 Cookie 状态为失效, Cookie = %d", cookie.Id))
			_ = GetCookieService().InnerUpdateStatusById(cookie.Id, enum.CookieStatus(2))
		}
	}

	jsonStr := regex.ReplaceAllString(html, "$1")
	jsonStr = strings.ReplaceAll(jsonStr, "undefined", "null")
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	// 用户提交的html有问题，改用cookie试试
	if !result.User.LoggedIn {
		return s.GetProfileInfo("", url)
	}
	return &result, nil
}

func (s *RedbookService) GetCookie() (*table.RedbookCookie, error) {
	// 获取cookie
	cookies, err := GetCookieService().GetAllByNormal()
	if err != nil {
		return nil, err
	}
	for _, ck := range cookies {
		if ck.Status == enum.CookieStatus(1) {
			return &ck, nil
		}
	}
	return nil, nil
}

func (s *RedbookService) SendRequest(url string, cookie *table.RedbookCookie) (*string, error) {
	// 发起请求
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	//设置请求头
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-encoding", "gzip, deflate, br, zstd")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6,zh-TW;q=0.5")
	req.Header.Set("content-type", "text/html; charset=utf-8")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("sec-ch-ua", "\"Chromium\";v=\"128\", \"Not;A=Brand\";v=\"24\", \"Google Chrome\";v=\"128\"")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36")
	if cookie != nil {
		req.Header.Set("cookie", cookie.Cookie)
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 如果gzip压缩了，需要特别处理
	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		//处理gzip响应流
		reader, _ = gzip.NewReader(resp.Body)
		defer reader.Close()
	default:
		reader = resp.Body
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	str := (*string)(unsafe.Pointer(&content)) //转化为string,优化内存
	return str, nil
}

func ConversionNumber(count string) (float64, error) {
	if count == "" {
		return 0, nil
	}
	if strings.HasSuffix(count, "万") {
		value, err := strconv.ParseFloat(strings.ReplaceAll(count, "万", ""), 64)
		if err != nil {
			return 0, err
		}
		return value * 10000.0, nil
	}
	return strconv.ParseFloat(count, 64)
}
