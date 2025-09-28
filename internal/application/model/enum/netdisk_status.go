package enum

import "errors"

// NORMAL：正常
// FAILURE：失效
var NETDISK_STATUS = [...]string{
	"",          // 0,未定义
	"NORMAL",    // 1,正常
	"CONCEAL",   // 2,隐藏
	"CANCEL",    // 3,取消
	"VIOLATION", // 4,违规
	"FAILURE",   // 5,失效
}

type NetdiskStatus uint

// JsonDate反序列化
func (t *NetdiskStatus) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range NETDISK_STATUS {
		if status == value {
			*t = NetdiskStatus(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t NetdiskStatus) MarshalJSON() ([]byte, error) {
	return []byte("\"" + NETDISK_STATUS[t] + "\""), nil
}

func NetdiskStatusValue(value string) (NetdiskStatus, error) {
	for i, status := range NETDISK_STATUS {
		if status == value {
			return NetdiskStatus(i), nil
		}
	}
	return NetdiskStatus(0), errors.New("未找到状态码：" + value)
}
