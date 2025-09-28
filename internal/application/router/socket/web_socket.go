package socket

import (
	"ai-software-copyright-server/internal/application/model/common"
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	netdSev "ai-software-copyright-server/internal/application/service/netdisk"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

type WebSocket struct {
	Clients        map[int64]*common.SocketClient
	upgrader       websocket.Upgrader
	awaitResultMap map[string]chan common.SocketMessage
}

var onceWebSocket = sync.Once{}
var webSocket *WebSocket

// 获取单例
func GetWebSocket() *WebSocket {
	onceWebSocket.Do(func() {
		webSocket = new(WebSocket)
		webSocket.Clients = make(map[int64]*common.SocketClient, 0)
		webSocket.upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return r.Header.Get("Access-Key") != ""
			},
		}
		webSocket.awaitResultMap = make(map[string]chan common.SocketMessage, 0)
	})
	return webSocket
}

func (s *WebSocket) NewHandler(c *gin.Context) {
	var user *table.User
	var configure table.NetdiskHelperConfigure
	var version string
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// 从accessKey中获取用户信息
			accessKey := c.GetHeader("Access-Key")
			if accessKey == "" {
				global.LOG.Error("Access-Key 不存在")
				return false
			}
			version = c.GetHeader("Client-Version")
			var err error
			user, err = userSev.GetUserService().GetByAccessKey(accessKey)
			if err != nil || user.Id == 0 {
				global.LOG.Error(fmt.Sprintf("获取用户信息失败（%s）: %+v", accessKey, err))
				return false
			}
			if s.Clients[user.Id] != nil {
				sendErr := s.Clients[user.Id].WriteMessage(enum.SocketMessageType(1), "有其他客户端尝试登录，可能您的AccessKey已泄露！")
				// 发送失败，认为这个客户端已经失效，允许继续下一步骤
				if sendErr != nil {
					global.LOG.Error(fmt.Sprintf("重复连接客户端，发送客户端通知消息失败（UserId = %d）：%+v", user.Id, sendErr))
				} else {
					global.LOG.Error(fmt.Sprintf("客户端已连接，不能重复连接（UserId = %d）", user.Id))
					return false
				}
			}
			configure, err = netdSev.GetHelperConfigureService().GetByUserId(user.Id)
			if err != nil {
				global.LOG.Error(fmt.Sprintf("获取网盘配置信息失败（UserId = %d）: %+v", user.Id, err))
				return false
			}
			if configure.ExpireTime == nil || configure.ExpireTime.Before(time.Now()) {
				global.LOG.Error(fmt.Sprintf("网盘助手服务已过期：%d", user.Id))
				return false
			}
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("协议升级失败: %+v", err))
		return
	}
	// 保存
	client := &common.SocketClient{Conn: conn, Version: version}
	s.Clients[user.Id] = client
	defer func() {
		_ = conn.Close()
		delete(s.Clients, user.Id)
		global.LOG.Warn(fmt.Sprintf("WebSocket（%d）客户端断开连接, %+v", user.Id, recover()))
		go func() {
			// 暂停一分钟，如果没有重新上线，就发邮件通知
			time.Sleep(1 * time.Minute)
			// 已经重新上线了，不用发通知
			if s.GetClient(user.Id) != nil {
				return
			}
			err = userSev.GetUserService().SendMail(user.Id, "异常通知：网盘助手客户端离线", "网盘助手客户端离线，将导致网盘资源无法正常转存，微信群无法正常进行资源搜索，请尽快检查！")
			if err != nil {
				global.LOG.Error(fmt.Sprintf("WebSocket（%d）客户端断开连接, 发送邮件通知失败, %+v", user.Id, err))
			}
		}()
	}()
	err = client.WriteMessage(enum.SocketMessageType(9), user)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("发送用户信息失败: %+v", err))
		return
	}
	err = client.WriteMessage(enum.SocketMessageType(4), configure)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("发送配置信息失败: %+v", err))
		return
	}
	global.LOG.Info(fmt.Sprintf("新增一个WebSocket（%d）V%s: 已发送配置信息", user.Id, client.Version))
	// 不是最新版
	// TODO 1.1.3版本之后才支持发送版本通知
	if utils.VersionCode(version) >= utils.VersionCode("1.1.3") && version != global.NetdiskHelperUpdateNotes[0].Version {
		note := global.NetdiskHelperUpdateNotes[0]
		versionCode := utils.VersionCode(version)
		for i := 1; i < len(global.NetdiskHelperUpdateNotes); i++ {
			v := global.NetdiskHelperUpdateNotes[i]
			// 如果这个版本比当前版本新，把更新说明一起加进来
			if utils.VersionCode(v.Version) > versionCode {
				note.Description += "\n" + v.Description
			}
		}
		note.Description = utils.ListJoin(strings.Split(note.Description, "\n"), "\n", func(i int, item string) string {
			return fmt.Sprintf("%d. %s", i+1, item)
		})
		noteErr := client.WriteMessage(enum.SocketMessageType(18), note)
		if noteErr != nil {
			global.LOG.Error(fmt.Sprintf("发送版本更新提示失败: %+v", err))
		}
	}
	// 获取信息
	for {
		message := common.SocketMessage{}
		err = conn.ReadJSON(&message)
		if err != nil {
			panic(err)
		}
		go s.HandlerMessage(conn, message)
	}
}

// 处理消息
func (s *WebSocket) HandlerMessage(conn *websocket.Conn, message common.SocketMessage) {
	startTime := time.Now()
	var result any
	// 捕获异常
	defer func() {
		r := recover()
		if r != nil {
			global.LOG.Error(fmt.Sprintf("WebSocket消息处理异常（%d）:%s\n%+v", message.Type, message.Data, r))
		}
		// 不需要结果通知
		if !message.NeedResult {
			return
		}
		// 超时也不用通知了
		if time.Now().Sub(startTime) >= message.Timeout {
			global.LOG.Error(fmt.Sprintf("WebSocket发送结果通知超时（%d）：%+v", message.Type, message.Data))
			return
		}
		// 开始处理结果通知
		resultMsg := common.SocketMessage{TaskId: message.TaskId}
		if r != nil {
			resultMsg.Type = enum.SocketMessageType(3)
			if resErr, ok := r.(error); ok {
				resultMsg.Data = "错误：" + resErr.Error()
			} else {
				resultMsg.Data = "未知异常"
			}
		} else {
			resultMsg.Type = enum.SocketMessageType(2)
			// 序列化
			serializedResult, err := json.Marshal(result)
			if err != nil {
				global.LOG.Error(fmt.Sprintf("WebSocket未能处理结果通知（%d）：%+v", message.Type, err))
				return
			}
			resultMsg.Data = string(serializedResult)
		}
		// 发送通知
		resultErr := conn.WriteJSON(resultMsg)
		if resultErr != nil {
			global.LOG.Error(fmt.Sprintf("Socket未能发送结果通知（%d）：%+v", message.Type, resultErr))
		}
	}()
	switch message.Type {
	case enum.SocketMessageType(2), enum.SocketMessageType(3): // 结果通知
		channel := s.awaitResultMap[message.TaskId]
		if channel != nil {
			channel <- message
		}
	default:
		global.LOG.Info(fmt.Sprintf("WebSocket未能处理的消息类型（%s）: %s", message.Type, message.Data))
	}
}

// 发送一个消息
func (s *WebSocket) SendMessage(userId int64, param common.SocketMessage, result any) (bool, error) {
	conn := s.GetClient(userId)
	if conn == nil {
		return false, errors.New("未搭建网盘助手或助手已经离线，搭建请联系微信：nineyaccz")
	}
	lowVersion := ""
	switch param.Type {
	case 12, 13, 14, 15, 16, 17:
		lowVersion = "1.1.0"
	case 19, 20, 21, 22, 23:
		lowVersion = "1.1.5"
	}
	if utils.VersionCode(conn.Version) < utils.VersionCode(lowVersion) {
		return false, errors.New(fmt.Sprintf("该功能需要%s或以上版本的网盘助手，请先升级网盘助手版本", lowVersion))
	}
	if param.NeedResult {
		param.TaskId = uuid.New().String()
		if param.Timeout <= 0 {
			param.Timeout = 8 * time.Second
		}
		err := conn.WriteJSON(param)
		if err != nil {
			return false, err
		}
		// 开始准备等待接收消息
		channel := make(chan common.SocketMessage)
		s.awaitResultMap[param.TaskId] = channel
		defer func() {
			delete(s.awaitResultMap, param.TaskId)
			close(channel)
		}()
		select {
		case message := <-channel: // 获取执行结果
			if message.Type == enum.SocketMessageType(3) {
				global.LOG.Info(fmt.Sprintf("WebSocket收到异常结果通知（%d）: TasId = %s, %s", message.Type, message.TaskId, message.Data))
				return false, errors.New(message.Data)
			}
			jsonUnmarshal(message.Data, result)
			return true, nil
		case <-time.After(param.Timeout): // 超时了
			return false, errors.New("等待超时，未能获取结果")
		}
	}
	return false, conn.WriteJSON(param)
}

func (s *WebSocket) GetClient(userId int64) *common.SocketClient {
	return s.Clients[userId]
}

func jsonUnmarshal(data string, v any) {
	if strings.HasPrefix(data, "\"") {
		var temp string
		err := json.Unmarshal([]byte(data), &temp)
		if err != nil {
			global.LOG.Error(fmt.Sprintf("反序列化失败: %+v", err))
			panic(errors.New("反序列化失败：" + string(data)))
		}
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(temp))
		return
	}
	err := json.Unmarshal([]byte(data), v)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("反序列化失败: %+v", err))
		panic(errors.New("反序列化失败：" + string(data)))
	}
}
