package enum

import (
	"github.com/pkg/errors"
)

var CLIENT_TASK_TYPE = [...]string{
	"",               // 0,未定义
	"SEND_MESSAGE",   // 1,发送消息
	"NETDISK_SEARCH", // 2,网盘资源搜索推送任务
}

type ClientTaskType uint

// JsonDate反序列化
func (t *ClientTaskType) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range CLIENT_TASK_TYPE {
		if status == value {
			*t = ClientTaskType(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t ClientTaskType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + CLIENT_TASK_TYPE[t] + "\""), nil
}

func ClientTaskTypeValue(value string) (ClientTaskType, error) {
	for i, postType := range CLIENT_TASK_TYPE {
		if postType == value {
			return ClientTaskType(i), nil
		}
	}
	return ClientTaskType(0), errors.New("未找到状态码：" + value)
}
