package enum

import (
	"github.com/pkg/errors"
)

var PICTURE_ORIGIN_TYPE = [...]string{
	"",       // 0, 未定义
	"CAMERA", // 1,拍照
	"ALBUM",  // 2,从相册选取
	"CHAT",   // 3,从聊天记录选取
}

type PictureOriginType uint

// JsonDate反序列化
func (t *PictureOriginType) UnmarshalJSON(data []byte) (err error) {
	value := string(data)
	value = value[1 : len(value)-1]
	for i, status := range PICTURE_ORIGIN_TYPE {
		if status == value {
			*t = PictureOriginType(i)
			return nil
		}
	}
	return errors.New("未找到状态码：" + value)
}

// JsonDate序列化
func (t PictureOriginType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + PICTURE_ORIGIN_TYPE[t] + "\""), nil
}

func PictureOriginTypeValue(value string) (PictureOriginType, error) {
	for i, postType := range PICTURE_ORIGIN_TYPE {
		if postType == value {
			return PictureOriginType(i), nil
		}
	}
	return PictureOriginType(0), errors.New("未找到状态码：" + value)
}
