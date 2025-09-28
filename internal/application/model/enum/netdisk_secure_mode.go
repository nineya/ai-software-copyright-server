package enum

import (
	"github.com/pkg/errors"
)

var NETDISK_SECURE_MODE = [...]string{
	"",       // 0,未定义
	"NORMAL", // 1,正常模式
	"SECURE", // 2,安全模式
	"AUDIT",  // 3,审核模式
}

type NetdiskSecureMode uint

// JsonDate反序列化
func (t *NetdiskSecureMode) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range NETDISK_SECURE_MODE {
		if status == value {
			*t = NetdiskSecureMode(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t NetdiskSecureMode) MarshalJSON() ([]byte, error) {
	return []byte("\"" + NETDISK_SECURE_MODE[t] + "\""), nil
}

func NetdiskSecureModeValue(value string) (NetdiskSecureMode, error) {
	for i, postType := range NETDISK_SECURE_MODE {
		if postType == value {
			return NetdiskSecureMode(i), nil
		}
	}
	return NetdiskSecureMode(0), errors.New("未找到状态码：" + value)
}
