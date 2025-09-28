package netdisk

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/service"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type CollectService struct {
	service.UserCrudService[table.NetdiskResource]
	Client http.Client
}

var onceCollect = sync.Once{}
var collectService *CollectService

// 获取单例
func GetCollectService() *CollectService {
	onceCollect.Do(func() {
		collectService = new(CollectService)
		collectService.Db = global.DB
		// 创建一个忽略证书验证的Transport
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		collectService.Client = http.Client{Transport: tr, Timeout: 6 * time.Second}
	})
	return collectService
}

func (s *CollectService) Collect(param request.NetdiskCollectParam) []table.NetdiskResource {
	keyword := strings.TrimSpace(param.Keyword)
	isQuark := false
	isBaidu := false
	for _, item := range param.Types {
		switch item {
		case enum.NetdiskType(2):
			isQuark = true
		case enum.NetdiskType(4):
			isBaidu = true
		}
	}
	resultCount := 5 // 默认全部选中的消息数量
	if !isQuark {    // 仅有百度时的消息数量
		resultCount = 3
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	resultChan := make(chan []table.NetdiskResource, resultCount) // 收集结果
	defer func() {
		cancel()
		close(resultChan)
	}()

	// 执行采集任务
	go s.Pansosuo(keyword, ctx, resultChan, isQuark, isBaidu)
	go s.Melost(keyword, ctx, resultChan, isQuark, isBaidu)
	go s.Woxiangsou(keyword, ctx, resultChan, isQuark, isBaidu)
	// 仅支持夸克
	if isQuark {
		go s.Kuake8(keyword, ctx, resultChan)
		go s.Funletu(keyword, ctx, resultChan)
		go s.WWW17x2(keyword, ctx, resultChan)
		go s.Yunso(keyword, ctx, resultChan)
	}

	// 整合采集结果
	resources := make([]table.NetdiskResource, 0)
	for i := 0; i < resultCount; i++ {
		select {
		case res := <-resultChan:
			resources = append(resources, res...)
		case <-ctx.Done(): // 防止死等
			return s.handleResources(keyword, resources)
		}
	}
	return s.handleResources(keyword, resources)
}

// https://www.pansosuo.com/ 采集，支持夸克、百度
func (s *CollectService) Pansosuo(keyword string, ctx context.Context, resultChan chan []table.NetdiskResource, isQuark, isBaidu bool) {
	// 指定请求参数
	param := make(map[string]any, 0)
	param["name"] = keyword
	bytesData, _ := json.Marshal(param)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.pansosuo.com/api/sources/aipan-search", bytes.NewBuffer(bytesData))
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.pansosuo.com/）: %+v", err))
		return
	}
	s.SetCommonHeader(req.Header)
	req.Header.Set("origin", "https://www.pansosuo.com")
	req.Header.Set("referer", "https://www.pansosuo.com/search")

	// 发起请求
	resp, err := s.Client.Do(req)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.pansosuo.com/）: %+v", err))
		return
	}
	defer resp.Body.Close()

	// 获取请求结果
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.pansosuo.com/）: %+v", err))
		return
	}
	global.LOG.Info("采集结果（https://www.pansosuo.com/）：", zap.Any("param", param), zap.String("result", string(content)))

	// 解析结果结构体
	result := struct {
		List []struct {
			Name  string `json:"name"`
			Links []struct {
				Pwd  string `json:"pwd"`
				Link string `json:"link"`
			} `json:"links"`
		} `json:"list"`
	}{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.pansosuo.com/）: %+v", err))
		return
	}
	// 取得列表数据
	if result.List == nil {
		return
	}
	resources := make([]table.NetdiskResource, 0)
	for _, item := range result.List {
		if item.Links == nil {
			continue
		}
		for _, link := range item.Links {
			if (isQuark && strings.Contains(link.Link, "https://pan.quark.cn/s/")) ||
				(isBaidu && strings.Contains(link.Link, "https://pan.baidu.com/s/")) {
				resources = append(resources, table.NetdiskResource{
					Name:           item.Name,
					TargetUrl:      link.Link,
					ShareTargetUrl: link.Link,
					SharePwd:       link.Pwd,
					Type:           utils.TransformNetdiskType(link.Link),
					Origin:         enum.NetdiskOrigin(2),
					Status:         enum.NetdiskStatus(1),
				})
			}
		}
	}
	// 如果还没结束，传入参数
	if err := ctx.Err(); err == nil {
		resultChan <- resources
	}
}

// https://kuake8.com/ 采集，支持夸克
func (s *CollectService) Kuake8(keyword string, ctx context.Context, resultChan chan []table.NetdiskResource) {
	// 指定请求参数
	param := make(map[string]any, 0)
	param["q"] = keyword
	param["page"] = 1
	param["size"] = 50
	param["exact"] = true
	bytesData, _ := json.Marshal(param)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://kuake8.com/v1/search/disk", bytes.NewBuffer(bytesData))
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://kuake8.com/）: %+v", err))
		return
	}
	s.SetCommonHeader(req.Header)
	req.Header.Set("origin", "https://kuake8.com")
	req.Header.Set("referer", "https://kuake8.com/search")

	// 发起请求
	resp, err := s.Client.Do(req)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://kuake8.com/）: %+v", err))
		return
	}
	defer resp.Body.Close()

	// 获取请求结果
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://kuake8.com/）: %+v", err))
		return
	}
	global.LOG.Info("采集结果（https://kuake8.com/）：", zap.Any("param", param), zap.String("result", string(content)))

	// 解析结果结构体
	result := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			List []struct {
				DiskName string `json:"disk_name"`
				DiskPass string `json:"disk_pass"`
				Link     string `json:"link"`
			} `json:"list"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://kuake8.com/）: %+v", err))
		return
	}
	if result.Code != 200 {
		global.LOG.Sugar().Error("采集任务执行失败（https://kuake8.com/）：" + result.Msg)
		return
	}
	// 取得列表数据
	if result.Data.List == nil {
		return
	}
	resources := make([]table.NetdiskResource, 0)
	for _, item := range result.Data.List {
		if strings.Contains(item.Link, "https://pan.quark.cn/s/") {
			resources = append(resources, table.NetdiskResource{
				Name:           strings.ReplaceAll(strings.ReplaceAll(item.DiskName, "<em>", ""), "</em>", ""),
				TargetUrl:      item.Link,
				ShareTargetUrl: item.Link,
				SharePwd:       item.DiskPass,
				Type:           enum.NetdiskType(2),
				Origin:         enum.NetdiskOrigin(2),
				Status:         enum.NetdiskStatus(1),
			})
		}
	}
	// 如果还没结束，传入参数
	if err := ctx.Err(); err == nil {
		resultChan <- resources
	}
}

// https://pan.funletu.com/ 趣盘搜采集，支持夸克
func (s *CollectService) Funletu(keyword string, ctx context.Context, resultChan chan []table.NetdiskResource) {
	// 指定请求参数
	param := make(map[string]any, 0)
	param["style"] = "get"
	param["datasrc"] = "search"
	queryParam := make(map[string]any, 0)
	queryParam["courseid"] = 1
	queryParam["searchtext"] = keyword
	param["query"] = queryParam
	pageParam := make(map[string]any, 0)
	pageParam["pageIndex"] = 1
	pageParam["pageSize"] = 50
	param["page"] = pageParam
	orderParam := make(map[string]any, 0)
	orderParam["prop"] = "views"
	orderParam["order"] = "desc"
	param["order"] = orderParam
	bytesData, _ := json.Marshal(param)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://v.funletu.com/search", bytes.NewBuffer(bytesData))
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://pan.funletu.com/）: %+v", err))
		return
	}
	s.SetCommonHeader(req.Header)
	req.Header.Set("origin", "https://pan.funletu.com")
	req.Header.Set("referer", "https://pan.funletu.com/")

	// 发起请求
	resp, err := s.Client.Do(req)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://pan.funletu.com/）: %+v", err))
		return
	}
	defer resp.Body.Close()

	// 获取请求结果
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://pan.funletu.com/）: %+v", err))
		return
	}
	global.LOG.Info("采集结果（https://pan.funletu.com/）：", zap.Any("param", param), zap.String("result", string(content)))

	// 解析结果结构体
	result := struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    []struct {
			Title   string `json:"title"`
			Url     string `json:"url"`
			Extcode string `json:"extcode"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://pan.funletu.com/）: %+v", err))
		return
	}
	if result.Status != 200 {
		global.LOG.Sugar().Error("采集任务执行失败（https://pan.funletu.com/）：" + result.Message)
		return
	}
	// 取得列表数据
	if result.Data == nil {
		return
	}
	resources := make([]table.NetdiskResource, 0)
	for _, item := range result.Data {
		if strings.Contains(item.Url, "https://pan.quark.cn/s/") {
			resources = append(resources, table.NetdiskResource{
				Name:           strings.ReplaceAll(strings.ReplaceAll(item.Title, "<em>", ""), "</em>", ""),
				TargetUrl:      item.Url,
				ShareTargetUrl: item.Url,
				SharePwd:       item.Extcode,
				Type:           enum.NetdiskType(2),
				Origin:         enum.NetdiskOrigin(2),
				Status:         enum.NetdiskStatus(1),
			})
		}
	}
	// 如果还没结束，传入参数
	if err := ctx.Err(); err == nil {
		resultChan <- resources
	}
}

// https://www.melost.cn/ 影盘社，支持夸克、百度
func (s *CollectService) Melost(keyword string, ctx context.Context, resultChan chan []table.NetdiskResource, isQuark, isBaidu bool) {
	// 指定请求参数
	param := make(map[string]any, 0)
	param["q"] = keyword
	param["page"] = 1
	param["size"] = 50
	param["exact"] = true
	param["type"] = ""

	if isQuark && !isBaidu {
		param["type"] = "QUARK"
	} else if !isQuark && isBaidu {
		param["type"] = "BDY"
	}

	bytesData, _ := json.Marshal(param)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://www.melost.cn/v1/search/disk", bytes.NewBuffer(bytesData))
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.melost.cn/）: %+v", err))
		return
	}
	s.SetCommonHeader(req.Header)
	req.Header.Set("origin", "https://www.melost.cn")
	req.Header.Set("referer", "https://www.melost.cn/search")

	// 发起请求
	resp, err := s.Client.Do(req)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://melost.cn/）: %+v", err))
		return
	}
	defer resp.Body.Close()

	// 获取请求结果
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://melost.cn/）: %+v", err))
		return
	}
	global.LOG.Info("采集结果（https://melost.cn/）：", zap.Any("param", param), zap.String("result", string(content)))

	// 解析结果结构体
	result := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			List []struct {
				DiskName string `json:"disk_name"`
				DiskPass string `json:"disk_pass"`
				Link     string `json:"link"`
			} `json:"list"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://melost.cn/）: %+v", err))
		return
	}
	if result.Code != 200 {
		global.LOG.Sugar().Error("采集任务执行失败（https://melost.cn/）：" + result.Msg)
		return
	}
	// 取得列表数据
	if result.Data.List == nil {
		return
	}
	resources := make([]table.NetdiskResource, 0)
	for _, item := range result.Data.List {
		if (isQuark && strings.Contains(item.Link, "https://pan.quark.cn/s/")) ||
			(isBaidu && strings.Contains(item.Link, "https://pan.baidu.com/s/")) {
			resources = append(resources, table.NetdiskResource{
				Name:           strings.ReplaceAll(strings.ReplaceAll(item.DiskName, "<em>", ""), "</em>", ""),
				TargetUrl:      item.Link,
				ShareTargetUrl: item.Link,
				SharePwd:       item.DiskPass,
				Type:           utils.TransformNetdiskType(item.Link),
				Origin:         enum.NetdiskOrigin(2),
				Status:         enum.NetdiskStatus(1),
			})
		}
	}
	// 如果还没结束，传入参数
	if err := ctx.Err(); err == nil {
		resultChan <- resources
	}
}

// https://www.17x2.cn/ 夸克短剧，支持夸克
func (s *CollectService) WWW17x2(keyword string, ctx context.Context, resultChan chan []table.NetdiskResource) {
	params := url.Values{}
	params.Add("text", keyword)
	reqUrl := "https://www.17x2.cn/api.php?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.17x2.cn/）: %+v", err))
		return
	}
	s.SetCommonHeader(req.Header)
	req.Header.Set("origin", "https://www.17x2.cn/")
	req.Header.Set("referer", "https://www.17x2.cn/")

	// 发起请求
	resp, err := s.Client.Do(req)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.17x2.cn/）: %+v", err))
		return
	}
	defer resp.Body.Close()

	// 获取请求结果
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.17x2.cn/）: %+v", err))
		return
	}
	global.LOG.Info("采集结果（https://www.17x2.cn/）：", zap.Any("url", reqUrl), zap.String("result", string(content)))

	// 解析结果结构体
	var result []struct {
		PlayName string `json:"play_name"`
		PlayUrl  string `json:"play_url"`
	}
	err = json.Unmarshal(content, &result)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.17x2.cn/）: %+v", err))
		return
	}
	// 取得列表数据
	resources := make([]table.NetdiskResource, 0)
	for _, item := range result {
		if strings.Contains(item.PlayUrl, "https://pan.quark.cn/s/") {
			resources = append(resources, table.NetdiskResource{
				Name:           item.PlayName,
				TargetUrl:      item.PlayUrl,
				ShareTargetUrl: item.PlayUrl,
				SharePwd:       "",
				Type:           enum.NetdiskType(2),
				Origin:         enum.NetdiskOrigin(2),
				Status:         enum.NetdiskStatus(1),
			})
		}
	}
	// 如果还没结束，传入参数
	if err := ctx.Err(); err == nil {
		resultChan <- resources
	}
}

// https://www.woxiangsou.com/ 口袋云，支持夸克、百度
func (s *CollectService) Woxiangsou(keyword string, ctx context.Context, resultChan chan []table.NetdiskResource, isQuark, isBaidu bool) {
	params := url.Values{}
	params.Add("keyword", keyword)
	params.Add("deviceId", uuid.New().String())
	params.Add("sourceType", "3")
	params.Add("timeFilter", "4")
	params.Add("offset", "0")
	params.Add("limit", "50")
	params.Add("channelId", "4")
	reqUrl := "https://www.woxiangsou.com/api/v1/resInfo/search?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.woxiangsou.com/）: %+v", err))
		return
	}
	s.SetCommonHeader(req.Header)
	req.Header.Set("origin", "https://www.woxiangsou.com/")
	req.Header.Set("referer", "https://www.woxiangsou.com/")

	// 发起请求
	resp, err := s.Client.Do(req)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.woxiangsou.com/）: %+v", err))
		return
	}
	defer resp.Body.Close()

	// 获取请求结果
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.woxiangsou.com/）: %+v", err))
		return
	}
	global.LOG.Info("采集结果（https://www.woxiangsou.com/）：", zap.Any("url", reqUrl), zap.String("result", string(content)))

	// 解析结果结构体
	result := struct {
		ResInfos []struct {
			Code string `json:"code"`
			Url  string `json:"url"`
			Name string `json:"name"`
		} `json:"resInfos"`
	}{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.woxiangsou.com/）: %+v", err))
		return
	}
	if result.ResInfos == nil || len(result.ResInfos) == 0 {
		global.LOG.Error("采集任务（https://www.woxiangsou.com/: 结果为空")
		return
	}
	// 取得列表数据
	resources := make([]table.NetdiskResource, 0)
	for _, item := range result.ResInfos {
		parsedURL, _ := url.Parse(item.Url)
		if parsedURL == nil {
			continue
		}
		netdiskUrl := parsedURL.Query().Get("url")
		if (isQuark && strings.Contains(netdiskUrl, "https://pan.quark.cn/s/")) ||
			(isBaidu && strings.Contains(netdiskUrl, "https://pan.baidu.com/s/")) {
			resources = append(resources, table.NetdiskResource{
				Name:           item.Name,
				TargetUrl:      netdiskUrl,
				ShareTargetUrl: netdiskUrl,
				SharePwd:       item.Code,
				Type:           utils.TransformNetdiskType(netdiskUrl),
				Origin:         enum.NetdiskOrigin(2),
				Status:         enum.NetdiskStatus(1),
			})
		}
	}
	// 如果还没结束，传入参数
	if err := ctx.Err(); err == nil {
		resultChan <- resources
	}
}

// https://www.yunso.net/ 小云搜索
func (s *CollectService) Yunso(keyword string, ctx context.Context, resultChan chan []table.NetdiskResource) {
	params := url.Values{}
	params.Add("wd", keyword)
	params.Add("mode", "90002")
	params.Add("stype", "20500")
	params.Add("scope_content", "0")
	params.Add("page", "1")
	params.Add("limit", "50")
	reqUrl := "https://www.yunso.net/api/validate/searchX2?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.yunso.net/）: %+v", err))
		return
	}
	s.SetCommonHeader(req.Header)
	req.Header.Set("origin", "https://www.yunso.net/")
	req.Header.Set("referer", "https://www.yunso.net/")
	req.Header.Set("user-agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 16_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.6 Mobile/15E148 Safari/604.1")

	// 发起请求
	resp, err := s.Client.Do(req)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.yunso.net/）: %+v", err))
		return
	}
	defer resp.Body.Close()

	// 获取请求结果
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.yunso.net/）: %+v", err))
		return
	}
	global.LOG.Info("采集结果（https://www.yunso.net/）：", zap.Any("url", reqUrl), zap.String("result", string(content)))

	// 解析结果结构体
	result := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data string `json:"data"`
	}{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("采集任务执行失败（https://www.yunso.net/）: %+v", err))
		return
	}
	if result.Code != 0 {
		global.LOG.Error("采集任务执行失败（https://www.yunso.net/）: " + result.Msg)
		return
	}

	// 取得列表数据
	resources := make([]table.NetdiskResource, 0)
	reg, _ := regexp.Compile("(https://pan.quark.cn/s/[a-z0-9]+)[^>]+>([^<]+)<")
	matches := reg.FindAllStringSubmatch(result.Data, -1)
	for _, item := range matches {
		if len(item) != 3 {
			continue
		}
		quarkUrl := item[1]
		if !strings.Contains(quarkUrl, "https://pan.quark.cn/s/") {
			continue
		}
		fileName := strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(item[2], "?", ""), "标签：", ""), "名称：", ""))
		// 字符串截取
		if strings.HasPrefix(fileName, "简介：") || utf8.RuneCountInString(fileName) > 32 {
			continue
		}
		resources = append(resources, table.NetdiskResource{
			Name:           fileName,
			TargetUrl:      quarkUrl,
			ShareTargetUrl: quarkUrl,
			SharePwd:       "",
			Type:           enum.NetdiskType(2),
			Origin:         enum.NetdiskOrigin(2),
			Status:         enum.NetdiskStatus(1),
		})
	}
	// 如果还没结束，传入参数
	if err := ctx.Err(); err == nil {
		resultChan <- resources
	}
}

// 设置公共请求头
func (s *CollectService) SetCommonHeader(header http.Header) {
	//设置请求头
	header.Set("sec-ch-ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	header.Set("accept", "application/json, text/plain, */*")
	header.Set("content-type", "application/json")
	header.Set("sec-ch-ua-mobile", "?0")
	header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	header.Set("sec-ch-ua-platform", "macOS")
	header.Set("sec-fetch-site", "same-site")
	header.Set("sec-fetch-mode", "cors")
	header.Set("sec-fetch-dest", "empty")
	//header.Set("accept-encoding", "gzip, deflate, br")
	header.Set("accept-language", "zh-CN,zh;q=0.9")
}

// 处理资源
func (s *CollectService) handleResources(keyword string, resources []table.NetdiskResource) []table.NetdiskResource {

	// 数组去重
	results := make([]common.Similarity[table.NetdiskResource], 0)
	for i := range resources {
		flag := true
		for j := range results {
			if resources[i].Name == results[j].Data.Name || resources[i].TargetUrl == results[j].Data.TargetUrl {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			score := jaccardSimilarity(keyword, resources[i].Name)
			if score > 0 {
				results = append(results, common.Similarity[table.NetdiskResource]{Data: resources[i], Score: score})
			}
		}
	}
	// 按分值减序
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	return utils.ListTransform(results, func(item common.Similarity[table.NetdiskResource]) table.NetdiskResource {
		return item.Data
	})
	// 打乱数组顺序
	//utils.ListShuffle(resources)
}

// 比较字符串相似度
func jaccardSimilarity(s1, s2 string) float64 {
	set1 := make(map[rune]bool)
	set2 := make(map[rune]bool)

	for _, c := range s1 {
		set1[c] = true
	}

	for _, c := range s2 {
		set2[c] = true
	}

	// s2和s1相同字符数
	intersection := 0
	for c := range set1 {
		if set2[c] {
			intersection++
		}
	}

	// s1+s2-相同字符数=总字符数
	union := len(set1) + len(set2) - intersection
	// s2比s1长，根据长度增加权重
	if len(set2) > len(set1) && strings.Contains(strings.ToLower(s2), strings.ToLower(s1)) || len(s1) > 6 {
		intersection += (len(set2) - len(set1)) / 2
	}

	return float64(intersection) / float64(union)
}
