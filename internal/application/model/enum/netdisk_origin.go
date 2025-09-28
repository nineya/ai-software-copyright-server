package enum

import (
	"github.com/pkg/errors"
)

var NETDISK_ORIGIN = [...]string{
	"",               // 0, 未定义
	"UPLOAD",         // 1,用户上传
	"COLLECT",        // 2,搜索采集
	"SHORT_LINK",     // 3,创建短链
	"AUTO_SYNC",      // 4,自动同步
	"HELPER_OPERATE", // 5,助手操作
}

type NetdiskOrigin uint

// JsonDate反序列化
func (t *NetdiskOrigin) UnmarshalJSON(data []byte) error {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range NETDISK_ORIGIN {
		if status == value {
			*t = NetdiskOrigin(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t NetdiskOrigin) MarshalJSON() ([]byte, error) {
	return []byte("\"" + NETDISK_ORIGIN[t] + "\""), nil
}

func NetdiskOriginValue(value string) (NetdiskOrigin, error) {
	for i, postType := range NETDISK_ORIGIN {
		if postType == value {
			return NetdiskOrigin(i), nil
		}
	}
	return NetdiskOrigin(0), errors.New("未找到状态码：" + value)
}
