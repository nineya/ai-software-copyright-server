package enum

import (
	"github.com/pkg/errors"
)

var SOCKET_MESSAGE_TYPE = [...]string{
	"",                                  // 0, 未定义
	"SEND_MESSAGE",                      // 1, 发送消息
	"RESULT",                            // 2, 结果通知
	"ERROR",                             // 3, 异常通知
	"UPDATE_CONFIG",                     // 4, 更新配置
	"WECHAT_SEND_MESSAGE",               // 5,【微信工具人】群发微信消息
	"NETDISK_RESOURCE_COLLECT",          // 6,【网盘助手】资源采集
	"NETDISK_RESOURCE_SEARCH_SAVE",      // 7,【网盘助手】网盘资源搜索转存
	"NETDISK_RESOURCE_TEXT_SAVE",        // 8,【网盘助手】夸克通过文案内容批量转存
	"USER_INFO",                         // 9, 同步用户信息
	"AI_CHAT",                           // 10,【网盘助手】AI对话
	"MAIL_SEND",                         // 11,【网盘助手】发送邮件
	"NETDISK_RESOURCE_ACCOUNT_TRANSFER", // 12,【网盘助手】网盘资源迁移
	"NETDISK_RESOURCE_FILE_LIST",        // 13,【网盘助手】网盘文件列表
	"NETDISK_RESOURCE_BATCH_SHARE",      // 14,【网盘助手】批量分享
	"NETDISK_RESOURCE_BATCH_RENAME",     // 15,【网盘助手】批量重命名
	"NETDISK_RESOURCE_BATCH_DELETE_AD",  // 16,【网盘助手】批量去广告
	"NETDISK_RESOURCE_BATCH_ADD_AD",     // 17,【网盘助手】批量加广告
	"VERSION_NOTE",                      // 18, 发送版本通知
	"NETDISK_RESOURCE_BATCH_SAVE",       // 19,【网盘助手】批量转存
	"NETDISK_RESOURCE_TASK_LIST",        // 20,【网盘助手】任务列表
	"NETDISK_RESOURCE_TASK_LOG",         // 21,【网盘助手】任务日志
	"NETDISK_RESOURCE_TASK_RESULT",      // 22,【网盘助手】任务结果
	"NETDISK_RESOURCE_TASK_DELETE",      // 23,【网盘助手】删除任务
	"TEST3",                             // 24, 测试
	"TEST3",                             // 25, 测试
	"TEST3",                             // 26, 测试
	"WECHAT_SEND_TEXT",                  // 27,【微信工具人】群发微信文本通知 TODO 临时过渡
}

type SocketMessageType int

// JsonDate反序列化
func (t *SocketMessageType) UnmarshalJSON(data []byte) error {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range SOCKET_MESSAGE_TYPE {
		if status == value {
			*t = SocketMessageType(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t SocketMessageType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + SOCKET_MESSAGE_TYPE[t] + "\""), nil
}

func SocketMessageTypeValue(value string) (SocketMessageType, error) {
	for i, postType := range SOCKET_MESSAGE_TYPE {
		if postType == value {
			return SocketMessageType(i), nil
		}
	}
	return SocketMessageType(0), errors.New("未找到状态码：" + value)
}
