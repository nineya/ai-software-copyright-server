package zhipu_ai

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/global"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

type ZhipuAiPlugin struct {
	Client  http.Client
	ChatUrl string
	ViewUrl string
}

var onceZhipuAi = sync.Once{}
var zhipuAiPlugin *ZhipuAiPlugin

// 获取单例
func GetZhipuAiPlugin() *ZhipuAiPlugin {
	onceZhipuAi.Do(func() {
		zhipuAiPlugin = new(ZhipuAiPlugin)
		zhipuAiPlugin.Client = http.Client{}
		zhipuAiPlugin.ChatUrl = "https://open.bigmodel.cn/api/paas/v4/chat/completions"
		zhipuAiPlugin.ViewUrl = "https://open.bigmodel.cn/api/paas/v4/images/generations"
	})
	return zhipuAiPlugin
}

func GetDefaultChatParam() request.ZhipuAiChatParam {
	return request.ZhipuAiChatParam{
		Model:     "glm-4-flash-250414",
		MaxTokens: 4095,
		Tools: []request.ZhipuAiChatToolItem{
			{
				Type: "web_search",
				WebSearch: &request.ZhipuAiChatToolWebSearch{
					Enable: true,
				},
			},
		},
	}
}

func (p *ZhipuAiPlugin) SendView(param request.ZhipuAiViewParam) (*response.ZhipuAiViewResponse, error) {
	content, err := p.sendRequest(param, p.ViewUrl)
	if err != nil {
		return nil, err
	}

	var result response.ZhipuAiViewResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	return &result, err
}

func (p *ZhipuAiPlugin) SendChat(param request.ZhipuAiChatParam) (*response.ZhipuAiChatResponse, error) {
	content, err := p.sendRequest(param, p.ChatUrl)
	if err != nil {
		return nil, err
	}

	//str := (*string)(unsafe.Pointer(&content)) //转化为string,优化内存
	//fmt.Println(*str)

	var result response.ZhipuAiChatResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	return &result, err
}

func (p *ZhipuAiPlugin) sendRequest(param any, url string) ([]byte, error) {
	bytesData, _ := json.Marshal(param)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bytesData))
	if err != nil {
		return nil, err
	}
	//设置请求头
	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+global.CONFIG.Plugin.ZhipuAi.Apikey)

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
