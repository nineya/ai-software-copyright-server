package enum

import (
	"github.com/pkg/errors"
)

var CLIENT_TYPE = [...]string{
	"",                      // 0,未定义
	"BLOGGER_HELPER",        // 1,博主创作助手
	"TRAFFIC_TOOLBOX",       // 2,流量工具箱（小红薯文案助手）
	"OPERATION_HELPER",      // 3,公私域运营助手
	"WECHAT_TOOLBOX",        // 4,微聊宝盒
	"NETDISK_HELPER",        // 5,网盘拉新达人助手
	"COZE",                  // 6,扣子
	"NETDISK_SEARCH_WXAMP",  // 7,网盘搜索小程序
	"NETDISK_SEARCH_APP",    // 8,网盘搜索APP
	"NETDISK_SEARCH_SITE",   // 9,网盘搜索网站
	"NINEYA_NETDISK_HELPER", // 10,小玖网盘助手
}

type ClientType uint

// JsonDate反序列化
func (t *ClientType) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range CLIENT_TYPE {
		if status == value {
			*t = ClientType(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t ClientType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + CLIENT_TYPE[t] + "\""), nil
}

func ClientTypeValue(value string) (ClientType, error) {
	for i, postType := range CLIENT_TYPE {
		if postType == value {
			return ClientType(i), nil
		}
	}
	return ClientType(0), errors.New("未找到状态码：" + value)
}
