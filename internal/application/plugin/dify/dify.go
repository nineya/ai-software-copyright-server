package dify

import (
	"ai-software-copyright-server/internal/global"
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"sync"
)

type DifyPlugin struct {
	Client http.Client
}

var onceDify = sync.Once{}
var difyPlugin *DifyPlugin

// 获取单例
func GetDifyPlugin() *DifyPlugin {
	onceDify.Do(func() {
		difyPlugin = new(DifyPlugin)
		difyPlugin.Client = http.Client{}
	})
	return difyPlugin
}

func (p *DifyPlugin) ConversationRename(apiKey, conversationId string, param DifyConversationRenameParam) (*DifyConversationRenameResponse, error) {
	content, err := p.sendRequest("/conversations/"+conversationId+"/name", apiKey, param)
	if err != nil {
		return nil, err
	}

	var result DifyConversationRenameResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	return &result, err
}

func (p *DifyPlugin) SendChat(apiKey string, param DifyChatMessageParam) (*DifyChatMessageResponse, error) {
	param.ResponseMode = "blocking"
	content, err := p.sendRequest("/chat-messages", apiKey, param)
	if err != nil {
		return nil, err
	}

	var result DifyChatMessageResponse // 反序列化JSON到结构体
	err = json.Unmarshal(content, &result)
	return &result, err
}

func (p *DifyPlugin) SendSSEChat(apiKey string, param DifyChatMessageParam) (string, string, error) {
	return p.SendSSEChatAndRetry(apiKey, param, 0)
}

func (p *DifyPlugin) SendSSEChatAndRetry(apiKey string, param DifyChatMessageParam, retryCount int) (string, string, error) {
	param.ResponseMode = "streaming"

	resultText := ""
	conversationId := ""
	err := p.sendSSERequest("/chat-messages", apiKey, param, func(content string) error {
		content = strings.TrimPrefix(content, "data:")

		var result DifyChatMessageSSEResponse // 反序列化JSON到结构体
		err := json.Unmarshal([]byte(content), &result)
		if err != nil {
			return err
		}
		switch result.Event {
		case "message":
			resultText += result.Answer
		case "error": // 执行出错
			global.LOG.Error("Dify SSE执行失败：", zap.Any("apiKey", apiKey), zap.String("result", content))
			return errors.New(result.Message)
		}
		if conversationId == "" {
			conversationId = result.ConversationId
		}
		return nil
	})
	if err != nil && retryCount > 0 {
		global.LOG.Error("Dify SSE执行失败，准备开始重试：", zap.Any("apiKey", apiKey), zap.Int("retryCount", retryCount))
		return p.SendSSEChatAndRetry(apiKey, param, retryCount-1)
	}
	return resultText, conversationId, err
}

func (p *DifyPlugin) sendRequest(url, apiKey string, param any) ([]byte, error) {
	var body io.Reader
	if param != nil {
		bytesData, _ := json.Marshal(param)
		body = bytes.NewBuffer(bytesData)
	}
	req, err := http.NewRequest("POST", global.CONFIG.Plugin.Dify.Host+url, body)
	if err != nil {
		return nil, err
	}

	// 设置头部
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := p.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	global.LOG.Info("Dify请求：", zap.String("url", url), zap.Any("apiKey", apiKey), zap.Any("param", param), zap.String("result", string(result)))
	return result, nil
}

func (p *DifyPlugin) sendSSERequest(url, apiKey string, param any, event func(bytes string) error) error {
	var body io.Reader
	if param != nil {
		bytesData, _ := json.Marshal(param)
		body = bytes.NewBuffer(bytesData)
	}
	req, err := http.NewRequest("POST", global.CONFIG.Plugin.Dify.Host+url, body)
	if err != nil {
		return err
	}

	// 设置头部
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	resp, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	global.LOG.Info("Dify SSE请求发起：", zap.String("url", url), zap.Any("apiKey", apiKey), zap.Any("param", param))

	// 读取 SSE 流
	reader := bufio.NewReader(resp.Body)
	for {
		result, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		if !strings.HasPrefix(result, "data:") {
			continue
		}
		//global.LOG.Debug("Dify SSE请求结果片段：", zap.String("url", url), zap.String("result", result))
		err = event(result)
		if err != nil {
			global.LOG.Error("Dify SSE请求Event执行失败：", zap.String("url", url), zap.Any("apiKey", apiKey), zap.Any("result", result), zap.Error(err))
			return err
		}
	}
}

func (p *DifyPlugin) handleRequest(method, url, apiKey string, param any) (*http.Response, error) {
	var body io.Reader
	if param != nil {
		bytesData, _ := json.Marshal(param)
		body = bytes.NewBuffer(bytesData)
	}
	req, err := http.NewRequest(method, global.CONFIG.Plugin.Dify.Host+url, body)
	if err != nil {
		return nil, err
	}

	// 设置头部
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return p.Client.Do(req)
}
