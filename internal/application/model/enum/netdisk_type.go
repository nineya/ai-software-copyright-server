package enum

import (
	"github.com/pkg/errors"
)

var NETDISK_TYPE = [...]string{
	"",       // 0, 未定义
	"OTHER",  // 1,其他
	"QUARK",  // 2,夸克
	"XUNLEI", // 3,迅雷
	"BAIDU",  // 4,百度
	"UC",     // 5,UC
	"CY139",  // 6,移动云盘
	"P123",   // 7,123网盘
	"WUKONG", // 8,悟空网盘
	"KUAITU", // 9,快兔网盘
}

type NetdiskType uint

// JsonDate反序列化
func (t *NetdiskType) UnmarshalJSON(data []byte) error {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range NETDISK_TYPE {
		if status == value {
			*t = NetdiskType(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t NetdiskType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + NETDISK_TYPE[t] + "\""), nil
}

func NetdiskTypeValue(value string) (NetdiskType, error) {
	for i, postType := range NETDISK_TYPE {
		if postType == value {
			return NetdiskType(i), nil
		}
	}
	return NetdiskType(0), errors.New("未找到状态码：" + value)
}
