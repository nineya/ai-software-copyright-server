package quark

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"
)

type QuarkPlugin struct {
	Client http.Client
}

var onceQuark = sync.Once{}
var quarkPlugin *QuarkPlugin

// 获取单例
func GetQuarkPlugin() *QuarkPlugin {
	onceQuark.Do(func() {
		quarkPlugin = new(QuarkPlugin)
		quarkPlugin.Client = http.Client{}
	})
	return quarkPlugin
}

// 调用API实现转存资源,taskid,title,self,error
func (p *QuarkPlugin) Dump(shareUrl, passcode, toPdirFid string) (string, string, bool, error) {
	pwdId := p.GetPwdId(shareUrl)
	// 获取stoken
	stoken, err := p.GetStokenApi(pwdId, passcode)
	if err != nil {
		return "", "", false, err
	}
	// 获取分享链接信息
	title := ""
	fileName := ""
	fids := make([]string, 0)
	fidTokens := make([]string, 0)
	page := 1
	for true {
		detail, next, err := p.DetailApi(pwdId, stoken.Data.Stoken, page)
		if err != nil {
			return "", "", false, err
		}
		if page == 1 {
			title = detail.Share.Title
			fileName = strconv.FormatInt(time.Now().Unix(), 10) + "-" + title
			// 字符串截取
			if utf8.RuneCountInString(fileName) > 225 {
				fileName = string([]rune(fileName)[:225])
			}
		}
		fids = append(fids, utils.ListTransform(detail.List, func(item response.NetdiskQuarkDetailDataList) string { return item.Fid })...)
		fidTokens = append(fidTokens, utils.ListTransform(detail.List, func(item response.NetdiskQuarkDetailDataList) string { return item.ShareFidToken })...)
		if !next {
			break
		}
		page++
	}
	if len(fids) == 0 {
		return "", "", false, errors.New("该分享已经失效")
	}
	// 如果包含多个文件，则创建一个子目录用于存储
	if len(fids) > 1 {
		toPdirFid, err = p.NewDirApi(fileName, toPdirFid)
		if err != nil {
			return "", "", false, err
		}
	}

	// 转存
	taskId, self, err := p.SaveApi(pwdId, stoken.Data.Stoken, fids, fidTokens, toPdirFid)
	return taskId, title, self, err
}

// 调用API接口，实现转存并分享
func (p *QuarkPlugin) DumpAndShare(shareUrl, passcode, toPdirFid string) (string, error) {
	if toPdirFid == "" {
		toPdirFid = global.CONFIG.Plugin.Quark.ToPdirFid
	}
	// 转存文件
	saveTaskId, title, self, err := p.Dump(shareUrl, passcode, toPdirFid)
	// 如果是自己转存自己，直接返回原链接
	if self {
		return shareUrl, nil
	}
	if err != nil {
		return "", err
	}
	// 查询转存任务结果
	saveTaskResult, err := p.TaskApi(saveTaskId)
	if err != nil {
		return "", err
	}
	// 分享文件
	shareTaskId, err := p.ShareApi(title, saveTaskResult.SaveAs.SaveAsTopFids)
	if err != nil {
		return "", err
	}
	// 转存附加广告文件，不考虑转存是否成功
	if global.CONFIG.Plugin.Quark.AppendShare != "" {
		_, _, _, err = p.Dump(global.CONFIG.Plugin.Quark.AppendShare, "", saveTaskResult.SaveAs.SaveAsTopFids[0])
		if err != nil {
			global.LOG.Error(fmt.Sprintf("转存附加文件失败: %+v", err))
		}
	}
	// 查询分享任务结果
	shareTaskResult, err := p.TaskApi(shareTaskId)
	if err != nil {
		return "", err
	}
	// 获取分享链接
	return p.GetShareLinkApi(shareTaskResult.ShareId)
}

// 创建文件分享链接
func (p *QuarkPlugin) ShareApi(title string, fileIds []string) (string, error) {
	param := &request.NetdiskQuarkShareParam{
		Title:       title,
		FidList:     fileIds,
		ExpiredType: 1,
		UrlType:     1,
	}
	content, err := p.sendRequest(http.MethodPost, "https://drive-pc.quark.cn/1/clouddrive/share?pr=ucpro&fr=pc&uc_param_str=", param)
	if err != nil {
		return "", err
	}

	var result response.NetdiskQuarkResponse[response.NetdiskQuarkTaskIdData] // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return "", err
	}
	if result.Status != 200 || result.Code != 0 {
		return "", errors.New("夸克分享文件失败：" + result.Message)
	}
	return result.Data.TaskId, err
}

// 通过share_id获取分享链接
func (p *QuarkPlugin) GetShareLinkApi(shareId string) (string, error) {
	param := &request.NetdiskQuarkGetShareLinkParam{ShareId: shareId}
	content, err := p.sendRequest(http.MethodPost, fmt.Sprintf("https://drive-pc.quark.cn/1/clouddrive/share/password?pr=ucpro&fr=pc&uc_param_str=&__dt=%d&__t=%d", rand.Intn(60000)+1, time.Now().UnixMilli()), param)
	if err != nil {
		return "", err
	}

	var result response.NetdiskQuarkResponse[map[string]any] // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return "", err
	}
	if result.Status != 200 {
		return "", errors.New("夸克获取分享链接失败：" + result.Message)
	}
	return fmt.Sprintf("%v", result.Data["share_url"]), err
}

// 取得stoken
func (p *QuarkPlugin) GetStokenApi(pwdId, passcode string) (*response.NetdiskQuarkResponse[response.NetdiskQuarkGetStokenData], error) {
	param := &request.NetdiskQuarkGetStokenParam{PwdId: pwdId, Passcode: passcode}
	content, err := p.sendRequest(http.MethodPost, fmt.Sprintf("https://drive-pc.quark.cn/1/clouddrive/share/sharepage/token?pr=ucpro&fr=pc&uc_param_str=&__dt=%d&__t=%d", rand.Intn(60000)+1, time.Now().UnixMilli()), param)
	if err != nil {
		return nil, err
	}

	var result response.NetdiskQuarkResponse[response.NetdiskQuarkGetStokenData] // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	if result.Status != 200 {
		return &result, errors.New("夸克获取Stoken失败：" + result.Message)
	}
	return &result, nil
}

// 查询分享的详情，页数从1开始
func (p *QuarkPlugin) DetailApi(pwdId, stoken string, page int) (*response.NetdiskQuarkDetailData, bool, error) {
	requestUrl := fmt.Sprintf("https://drive-h.quark.cn/1/clouddrive/share/sharepage/detail?pr=ucpro&fr=pc&uc_param_str=&pwd_id=%s&stoken=%s&pdir_fid=0&force=0&_page=%d&_size=50&_fetch_banner=1&_fetch_share=1&_fetch_total=1&_sort=file_type:asc,updated_at:desc&__dt=%d&__t=%d",
		pwdId, url.QueryEscape(stoken), page, rand.Intn(60000)+1, time.Now().UnixMilli())

	content, err := p.sendRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, false, err
	}

	var result response.NetdiskQuarkResponse[response.NetdiskQuarkDetailData] // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, false, err
	}
	if result.Status != 200 {
		return nil, false, errors.New("夸克获取分享链接信息失败：" + result.Message)
	}
	return &result.Data, result.Data.Share.FileNum > page*50, nil
}

// 创建保存分享链接文件的任务
func (p *QuarkPlugin) SaveApi(pwdId, stoken string, fids, fidTokens []string, toPdirFid string) (string, bool, error) {
	param := &request.NetdiskQuarkSaveParam{
		FidList:      fids,
		FidTokenList: fidTokens,
		PdirFid:      "0",
		PwdId:        pwdId,
		Scene:        "link",
		Stoken:       stoken,
		ToPdirFid:    toPdirFid,
	}
	content, err := p.sendRequest(http.MethodPost, fmt.Sprintf("https://drive-pc.quark.cn/1/clouddrive/share/sharepage/save?pr=ucpro&fr=pc&uc_param_str=&__dt=%d&__t=%d", rand.Intn(60000)+1, time.Now().UnixMilli()), param)
	if err != nil {
		return "", false, err
	}

	var result response.NetdiskQuarkResponse[response.NetdiskQuarkTaskIdData] // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return "", false, err
	}
	if result.Status != 200 {
		if result.Code == 41017 { // 自己转存自己的文件
			return "", true, errors.New("用户禁止转存自己的分享")
		}
		return "", false, errors.New("夸克保存文件失败：" + result.Message)
	}
	return result.Data.TaskId, false, nil
}

// 创建目录
func (p *QuarkPlugin) NewDirApi(fileName, pdirFid string) (string, error) {
	param := &request.NetdiskQuarkNewDirParam{DirInitLock: false, DirPath: "", FileName: fileName, PdirFid: pdirFid}
	content, err := p.sendRequest(http.MethodPost, "https://drive-pc.quark.cn/1/clouddrive/file?pr=ucpro&fr=pc&uc_param_str=", param)
	if err != nil {
		return "", err
	}

	var result response.NetdiskQuarkResponse[map[string]any] // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return "", err
	}
	if result.Status != 200 {
		return "", errors.New("夸克创建目录失败：" + result.Message)
	}
	return fmt.Sprintf("%v", result.Data["fid"]), err

}

// 查询任务
func (p *QuarkPlugin) TaskApi(taskId string) (*response.NetdiskQuarkTaskData, error) {
	for i := 0; i < 30; i++ {
		time.Sleep(600 * time.Millisecond)
		requestUrl := fmt.Sprintf("https://drive-pc.quark.cn/1/clouddrive/task?pr=ucpro&fr=pc&uc_param_str=&task_id=%s&retry_index=%d&__dt=%d&__t=%d",
			taskId, i, rand.Intn(60000)+1, time.Now().UnixMilli())

		content, err := p.sendRequest(http.MethodGet, requestUrl, nil)
		if err != nil {
			return nil, err
		}

		var result response.NetdiskQuarkResponse[response.NetdiskQuarkTaskData] // 反序列化JSON到结构体
		err = json.Unmarshal(content, &result)
		if err != nil {
			continue
		}
		if result.Status != 200 {
			return &result.Data, errors.New("夸克任务失败：" + result.Message)
		}
		if result.Data.Status != 0 {
			return &result.Data, nil
		}
	}
	return nil, errors.New("获取任务结果失败：" + taskId)
}

func (p *QuarkPlugin) GetPwdId(shareUrl string) string {
	compileRegex := regexp.MustCompile("/s/([0-9a-z]+)")
	matchArr := compileRegex.FindStringSubmatch(shareUrl)
	if len(matchArr) > 0 {
		return matchArr[len(matchArr)-1]
	}
	return ""
}

func (p *QuarkPlugin) sendRequest(method string, url string, param any) ([]byte, error) {
	var body io.Reader
	if param != nil {
		bytesData, _ := json.Marshal(param)
		body = bytes.NewBuffer(bytesData)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	//设置请求头
	req.Header.Set("sec-ch-ua", "\"Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\", \"Google Chrome\";v=\"126\"")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("origin", "https://pan.quark.cn")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", "https://pan.quark.cn/")
	//req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("cookie", global.CONFIG.Plugin.Quark.Cookie)

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	global.LOG.Info("夸克请求：", zap.String("url", url), zap.Any("param", param), zap.String("result", string(result)))
	return result, nil
}
