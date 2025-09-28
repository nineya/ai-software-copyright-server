package baidu

import (
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

type BaiduPlugin struct {
	Client http.Client
}

var onceBaidu = sync.Once{}
var baiduPlugin *BaiduPlugin

// 获取单例
func GetBaiduPlugin() *BaiduPlugin {
	onceBaidu.Do(func() {
		baiduPlugin = new(BaiduPlugin)
		baiduPlugin.Client = http.Client{}
	})
	return baiduPlugin
}

// 转存资源,返回内容：info=转存信息,self=是不是自己转存自己,error
func (p *BaiduPlugin) Dump(cookie, bdstoken, shareUrl, pwd, toPath string) (*response.NetdiskBaiduTransferResponse, bool, error) {
	if bdstoken == "" {
		bdstoken, _ = p.GetBdstoken(cookie)
	}
	shareToken, passcode := p.GetShareTokenAndPasscode(shareUrl)
	// 使用传入的密码
	if passcode == "" {
		passcode = pwd
	}
	if shareToken == "" {
		// 链接有问题，只更新检查时间
		return nil, false, errors.New("错误的分享链接：" + shareUrl)
	}
	// 验证提取码
	cookies, _, err := p.VerifyPassCodeApi(cookie, shareToken, passcode)
	if err != nil {
		return nil, false, errors.New("验证提取码失败：" + shareUrl)
	}
	// 组合提取码的cookie信息
	cookie = p.GetCookie(cookies) + cookie
	// 获取资源信息
	info, err := p.GetShareInfoByHtmlApi(cookie, shareUrl)
	if err != nil {
		return nil, false, err
	}
	if len(info.FileList) == 0 {
		return nil, false, errors.New("未查询到分享资源")
	}
	// 获取需要转存的文件信息
	fsIds := utils.ListTransform(info.FileList, func(item response.NetdiskBaiduShareInfoByHtmlFileItem) int64 {
		return item.FsId
	})
	// 如果包含多个文件，则创建一个子目录用于存储
	if len(fsIds) > 1 {
		dirName := strconv.FormatInt(time.Now().Unix(), 10) + "-" + info.FileList[0].ServerFilename
		// 字符串截取
		if utf8.RuneCountInString(dirName) > 225 {
			dirName = string([]rune(dirName)[:225])
		}
		newDirInfo, err := p.NewDirApi(cookie, bdstoken, toPath, dirName)
		if err != nil {
			return nil, false, err
		}
		toPath = newDirInfo.Path
	}
	// 转存
	result, err := p.TransferApi(cookie, bdstoken, info.ShareId, info.ShareUk, fsIds, toPath)
	return result, false, err
}

// 转存并夹带私货，并分享
func (p *BaiduPlugin) DumpAndShare(bdstoken, shareUrl, passcode string) (string, error) {
	baiduConfig := global.CONFIG.Plugin.Baidu
	if bdstoken == "" {
		bdstoken, _ = p.GetBdstoken(baiduConfig.Cookie)
	}
	// 转存文件
	dumpInfo, self, err := p.Dump(baiduConfig.Cookie, bdstoken, shareUrl, passcode, baiduConfig.ToPath)
	// 如果是自己转存自己，直接返回原链接
	if self {
		return shareUrl, nil
	}
	if err != nil {
		return "", err
	}
	// 保存附加广告文件，不考虑转存是否成功
	if baiduConfig.AppendShare != "" {
		err = p.CopyFileApi(baiduConfig.Cookie, bdstoken, baiduConfig.AppendShare, dumpInfo.Extra.List[0].To, "")
		if err != nil {
			global.LOG.Error(fmt.Sprintf("转存附加文件失败: %+v", err))
		}
	}
	// 分享文件
	fsIds := utils.ListTransform(dumpInfo.Extra.List, func(item response.NetdiskBaiduTransferExtraItem) int64 {
		return item.ToFsId
	})
	shareInfo, err := p.ShareApi(baiduConfig.Cookie, bdstoken, baiduConfig.Pwd, fsIds)
	if err != nil {
		return "", err
	}
	// 获取分享链接
	return shareInfo.ShareUrl, nil
}

// 获取目录所有内容
func (p *BaiduPlugin) GetDirContent(cookie, path string) ([]response.NetdiskBaiduListItem, error) {
	list := make([]response.NetdiskBaiduListItem, 0)
	page := 1
	for true {
		detail, next, err := p.ListApi(cookie, path, page)
		if err != nil {
			return nil, err
		}
		list = append(list, detail.List...)
		if !next {
			break
		}
		time.Sleep(3 * time.Second)
		page++
	}
	return list, nil
}

// 遍历重命名
func (p *BaiduPlugin) BatchRenameFile(cookie, bdstoken, path string, replaceFunc func(string) string, replaceText string, recursive bool) {
	if bdstoken == "" {
		bdstoken, _ = p.GetBdstoken(cookie)
	}
	list, err := p.GetDirContent(cookie, path)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("批量重命名文件，获取资源列表失败：%+v", err))
		return
	}
	for _, item := range list {
		oldName := item.ServerFilename
		newName := replaceFunc(item.ServerFilename)
		if oldName != newName {
			global.LOG.Info(fmt.Sprintf("批量重命名文件，%s => %s", oldName, newName))
			err = p.RenameFileApi(cookie, bdstoken, item.FsId, item.Path, newName)
			if err != nil {
				global.LOG.Error(fmt.Sprintf("批量重命名文件，重命名文件失败：%+v", err))
			}
			time.Sleep(3 * time.Second)
		}
		// 如果是目录，且需要递归，递归重命名
		if item.IsDir == 1 && recursive {
			p.BatchRenameFile(cookie, bdstoken, item.Path, replaceFunc, replaceText, recursive)
		}
	}
}

// 遍历删除
func (p *BaiduPlugin) BatchDeleteFile(cookie, bdstoken, path string, matchFunc func(string) bool, recursive bool) error {
	if bdstoken == "" {
		bdstoken, _ = p.GetBdstoken(cookie)
	}
	list, err := p.GetDirContent(cookie, path)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("批量删除文件，获取资源列表失败：%+v", err))
		return err
	}
	fids := make([]string, 0)
	for _, item := range list {
		// 匹配模式为空，或者匹配成功，记录删除路径
		if matchFunc == nil || matchFunc(item.ServerFilename) {
			global.LOG.Info(fmt.Sprintf("批量删除文件，标记 %s 需要删除", item.ServerFilename))
			fids = append(fids, item.Path)
			continue
		}
		// 如果是目录，且需要递归，递归删除数据
		if item.IsDir == 1 && recursive {
			_ = p.BatchDeleteFile(cookie, bdstoken, item.Path, matchFunc, recursive)
		}
	}
	if len(fids) == 0 {
		return nil
	}
	err = p.DeleteFileApi(cookie, bdstoken, fids)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("批量删除文件，删除文件失败：%+v", err))
	}
	return err
}

// 遍历添加内容
func (p *BaiduPlugin) BatchAddFile(cookie, bdstoken, path string, matchFunc func(string) bool, appendShare string) {
	if bdstoken == "" {
		bdstoken, _ = p.GetBdstoken(cookie)
	}
	list, err := p.GetDirContent(cookie, path)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("批量添加文件，获取资源列表失败：%+v", err))
		return
	}
	for _, item := range list {
		// 如果不是目录，跳过
		if item.IsDir != 1 {
			continue
		}
		// 默认都添加广告，匹配规则不为空，且匹配不成功
		if matchFunc != nil && !matchFunc(item.ServerFilename) {
			continue
		}
		global.LOG.Info(fmt.Sprintf("批量添加文件，%s 目录需要添加", item.ServerFilename))
		err = p.CopyFileApi(cookie, bdstoken, appendShare, item.Path, "")
		if err != nil {
			global.LOG.Error(fmt.Sprintf("批量添加文件，转存附加文件失败: %+v", err))
		}
		time.Sleep(5 * time.Second)
	}
}

// 验证提取码
func (p *BaiduPlugin) VerifyPassCodeApi(cookie, shareToken, passcode string) ([]*http.Cookie, string, error) {
	requestUrl := fmt.Sprintf("https://pan.baidu.com/share/verify?t=%d&surl=%s&channel=chunlei&web=1&app_id=250528&bdstoken=&clienttype=0",
		time.Now().UnixMilli(), url.QueryEscape(shareToken))

	param := url.Values{}
	param.Set("pwd", passcode)
	param.Set("vcode", "")
	param.Set("vcode_str", "")

	cookies, content, err := p.sendRequest(cookie, http.MethodPost, requestUrl, param)
	if err != nil {
		return nil, "", err
	}

	var result response.NetdiskBaiduVerifyPassCodeResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, "", err
	}
	if result.Errno != 0 {
		return nil, "", errors.New("百度验证提取码失败：" + result.ErrMsg)
	}
	return cookies, result.Randsk, err
}

// 转存资源 fsids=转存的资源列表，toPath=转存目标目录,bdstoken=不知道做啥用，可以放空
func (p *BaiduPlugin) TransferApi(cookie, bdstoken string, shareId int64, shareUk string, fsIds []int64, toPath string) (*response.NetdiskBaiduTransferResponse, error) {
	if bdstoken == "" {
		bdstoken, _ = p.GetBdstoken(cookie)
	}
	requestUrl := fmt.Sprintf("https://pan.baidu.com/share/transfer?shareid=%d&from=%s&ondup=newcopy&channel=chunlei&web=1&app_id=250528&bdstoken=%s&clienttype=0",
		shareId, shareUk, bdstoken)

	param := url.Values{}
	param.Set("fsidlist", fmt.Sprintf("[%s]", utils.ListJoin(fsIds, ",", func(index int, item int64) string {
		return strconv.FormatInt(item, 10)
	}))) // 要转存的资源列表
	param.Set("path", toPath) // 转存的目标目录

	_, content, err := p.sendRequest(cookie, http.MethodPost, requestUrl, param)
	if err != nil {
		return nil, err
	}

	var result response.NetdiskBaiduTransferResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	if result.Errno != 0 {
		return nil, errors.New("百度转存失败：" + result.ErrMsg)
	}
	return &result, err
}

// 创建目录
func (p *BaiduPlugin) NewDirApi(cookie, bdstoken, path, dirName string) (*response.NetdiskBaiduNewDirResponse, error) {
	if bdstoken == "" {
		bdstoken, _ = p.GetBdstoken(cookie)
	}
	requestUrl := fmt.Sprintf("https://pan.baidu.com/api/create?a=commit&bdstoken=%s&clienttype=0&app_id=250528&web=1", bdstoken)

	param := url.Values{}
	param.Set("path", path+"/"+dirName) // 要创建的目录
	param.Set("isdir", "1")
	param.Set("block_list", "[]")

	_, content, err := p.sendRequest(cookie, http.MethodPost, requestUrl, param)
	if err != nil {
		return nil, err
	}

	var result response.NetdiskBaiduNewDirResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	if result.Errno != 0 {
		return nil, errors.New("百度创建目录失败：" + result.ErrMsg)
	}
	return &result, err
}

// 重命名目录/文件
func (p *BaiduPlugin) RenameFileApi(cookie, bdstoken string, fsId int64, path, dirName string) error {
	if bdstoken == "" {
		bdstoken, _ = p.GetBdstoken(cookie)
	}
	requestUrl := fmt.Sprintf("https://pan.baidu.com/api/filemanager?async=2&onnest=fail&opera=rename&bdstoken=%s&clienttype=0&app_id=250528&web=1", bdstoken)

	fileList := []map[string]any{
		{
			"id":      fsId,
			"path":    path,
			"newname": dirName,
		},
	}
	// 序列化
	serializedFileList, err := json.Marshal(fileList)
	if err != nil {
		return errors.Wrap(err, "请求参数序列化失败")
	}
	param := url.Values{}
	param.Set("filelist", string(serializedFileList))

	_, content, err := p.sendRequest(cookie, http.MethodPost, requestUrl, param)
	if err != nil {
		return err
	}

	var result response.NetdiskBaiduTaskResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return err
	}
	if result.Errno != 0 {
		return errors.New("百度目录重命名失败：" + result.ErrMsg)
	}
	return nil
}

// 复制目录/文件,path=源文件路径
func (p *BaiduPlugin) CopyFileApi(cookie, bdstoken, path, toDirPath, toName string) error {
	if bdstoken == "" {
		bdstoken, _ = p.GetBdstoken(cookie)
	}
	requestUrl := fmt.Sprintf("https://pan.baidu.com/api/filemanager?async=2&onnest=fail&opera=copy&bdstoken=%s&clienttype=0&app_id=250528&web=1", bdstoken)

	if toName == "" {
		// 取得要添加的文件名
		if lastIndex := strings.LastIndex(path, "/"); lastIndex > -1 {
			toName = path[lastIndex+1:]
		}
	}
	fileList := []map[string]any{
		{
			"path":    path,
			"dest":    toDirPath,
			"newname": toName,
		},
	}
	// 序列化
	serializedFileList, err := json.Marshal(fileList)
	if err != nil {
		return errors.Wrap(err, "请求参数序列化失败")
	}
	param := url.Values{}
	param.Set("filelist", string(serializedFileList))

	_, content, err := p.sendRequest(cookie, http.MethodPost, requestUrl, param)
	if err != nil {
		return err
	}

	var result response.NetdiskBaiduTaskResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return err
	}
	if result.Errno != 0 {
		return errors.New("百度文件复制失败：" + result.ErrMsg)
	}
	return nil
}

// 删除目录/文件
func (p *BaiduPlugin) DeleteFileApi(cookie, bdstoken string, paths []string) error {
	if bdstoken == "" {
		bdstoken, _ = p.GetBdstoken(cookie)
	}
	requestUrl := fmt.Sprintf("https://pan.baidu.com/api/filemanager?async=2&onnest=fail&opera=delete&bdstoken=%s&newVerify=1&clienttype=0&app_id=250528&web=1", bdstoken)

	// 序列化
	serializedFileList, err := json.Marshal(paths)
	if err != nil {
		return errors.Wrap(err, "请求参数序列化失败")
	}
	param := url.Values{}
	param.Set("filelist", string(serializedFileList))

	_, content, err := p.sendRequest(cookie, http.MethodPost, requestUrl, param)
	if err != nil {
		return err
	}

	var result response.NetdiskBaiduTaskResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return err
	}
	if result.Errno != 0 {
		return errors.New("百度目录删除失败：" + result.ErrMsg)
	}
	return nil
}

// 创建文件分享链接
func (p *BaiduPlugin) ShareApi(cookie, bdstoken, pwd string, fsIds []int64) (*response.NetdiskBaiduShareResponse, error) {
	if bdstoken == "" {
		bdstoken, _ = p.GetBdstoken(cookie)
	}
	pwd = RandomPassCode(pwd)
	requestUrl := fmt.Sprintf("https://pan.baidu.com/share/pset?channel=chunlei&bdstoken=%s&clienttype=0&app_id=250528&web=1", bdstoken)

	param := url.Values{}
	param.Set("is_knowledge", "0")
	param.Set("public", "0")
	param.Set("period", "0")
	param.Set("pwd", pwd)
	param.Set("eflag_disable", "true")
	param.Set("linkOrQrcode", "link")
	param.Set("channel_list", "[]")
	param.Set("schannel", "4")
	param.Set("fid_list", fmt.Sprintf("[%s]", utils.ListJoin(fsIds, ",", func(index int, item int64) string {
		return strconv.FormatInt(item, 10)
	}))) // 分享的目标目录

	_, content, err := p.sendRequest(cookie, http.MethodPost, requestUrl, param)
	if err != nil {
		return nil, err
	}

	var result response.NetdiskBaiduShareResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	if result.Errno != 0 {
		return nil, errors.New("百度分享文件失败：" + result.ErrMsg)
	}
	result.ShareUrl = result.Link + "?pwd=" + pwd
	return &result, err
}

// 通过html网页取得分享的信息
func (p *BaiduPlugin) GetShareInfoByHtmlApi(cookie, shareUrl string) (*response.NetdiskBaiduShareInfoByHtmlResponse, error) {
	_, content, err := p.sendRequest(cookie, http.MethodGet, shareUrl, nil)
	if err != nil {
		return nil, err
	}

	var result response.NetdiskBaiduShareInfoByHtmlResponse // 反序列化JSON到结构体
	compileRegex := regexp.MustCompile("locals\\.mset\\((\\{.*?\\})\\);")
	matchArr := compileRegex.FindStringSubmatch(string(content))
	if len(matchArr) > 1 {
		err = json.Unmarshal([]byte(matchArr[1]), &result)
		if err != nil {
			return nil, errors.Wrap(err, "解析分享信息失败")
		}
	} else {
		return nil, errors.New("解析分享信息失败，未找到信息，cookie：" + cookie)
	}
	if result.Errno != 0 {
		return &result, errors.New("获取分享信息失败：" + result.ErrMsg)
	}
	return &result, err
}

// 取得分享的信息，比html来源信息更全
func (p *BaiduPlugin) GetShareInfoApi(cookie, bdstoken, shareToken string) (*response.NetdiskBaiduShareInfoResponse, error) {
	if bdstoken == "" {
		bdstoken, _ = p.GetBdstoken(cookie)
	}
	requestUrl := fmt.Sprintf("https://pan.baidu.com/share/list?web=5&app_id=250528&desc=1&showempty=0&page=1&num=20&order=time&shorturl=%s&root=1&view_mode=1&channel=chunlei&web=1&bdstoken=%s&clienttype=0",
		shareToken, bdstoken)

	_, content, err := p.sendRequest(cookie, http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, err
	}

	var result response.NetdiskBaiduShareInfoResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, err
	}
	if result.Errno != 0 {
		return nil, errors.New("百度获取分享信息失败：" + result.ErrMsg)
	}
	return &result, nil
}

// 取得bdstoken
func (p *BaiduPlugin) GetBdstoken(cookie string) (string, error) {
	requestUrl := "https://pan.baidu.com/api/gettemplatevariable?clienttype=0&app_id=250528&web=1&fields=[%22bdstoken%22,%22token%22,%22uk%22,%22isdocuser%22,%22servertime%22]&channel=chunlei"
	_, content, err := p.sendRequest(cookie, http.MethodGet, requestUrl, nil)
	if err != nil {
		return "", err
	}

	var result response.NetdiskBaiduGetBdstokenResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return "", err
	}
	if result.Errno != 0 {
		return "", errors.New("百度获取Stoken失败：" + result.ErrMsg)
	}
	return result.Result.Bdstoken, nil
}

// 查询文件列表，页数从1开始
func (p *BaiduPlugin) ListApi(cookie, path string, page int) (*response.NetdiskBaiduListResponse, bool, error) {
	requestUrl := fmt.Sprintf("https://pan.baidu.com/api/list?clienttype=0&app_id=250528&web=1&order=time&desc=1&dir=%s&num=100&page=%d",
		url.QueryEscape(path), page)

	_, content, err := p.sendRequest(cookie, http.MethodGet, requestUrl, nil)
	if err != nil {
		return nil, false, err
	}

	var result response.NetdiskBaiduListResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	if err != nil {
		return nil, false, err
	}
	if result.Errno != 0 {
		return nil, false, errors.New("百度获取文件列表失败：" + result.ErrMsg)
	}
	return &result, len(result.List) == 100, nil
}

func (p *BaiduPlugin) sendRequest(cookie, method string, url string, param url.Values) ([]*http.Cookie, []byte, error) {
	var body io.Reader
	if param != nil {
		body = strings.NewReader(param.Encode())
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	//设置请求头
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.Header.Set("sec-ch-ua", "\"Not)A;Brand\";v=\"8\", \"Chromium\";v=\"138\", \"Google Chrome\";v=\"138\"")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua-platform", "Windows")
	req.Header.Set("origin", "https://pan.baidu.com")
	//req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-fetch-mode", "navigate")
	//req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", "https://pan.baidu.com/")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("cookie", cookie)

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	global.LOG.Info("百度网盘请求：", zap.String("url", url), zap.Any("param", param), zap.String("result", string(result)))
	return resp.Cookies(), result, nil
}

// 取得分享id和提取码
func (p *BaiduPlugin) GetShareTokenAndPasscode(shareUrl string) (shareId string, passcode string) {
	parsedURL, err := url.Parse(shareUrl)
	if err != nil {
		return
	}
	shareId = path.Base(parsedURL.Path)
	if shareId != "" {
		shareId = shareId[1:]
	}
	passcode = parsedURL.Query().Get("pwd")
	return shareId, passcode
}

// 提取cookie
func (p *BaiduPlugin) GetCookie(cookies []*http.Cookie) string {
	return utils.ListJoin(cookies, "", func(index int, item *http.Cookie) string {
		return item.Name + "=" + item.Value + ";"
	})
}

// 随机密码生成
func RandomPassCode(passCode string) string {
	if len(passCode) == 4 {
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9]{4}$`, passCode)
		if matched {
			return passCode
		}
	}
	charset := "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 4)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
