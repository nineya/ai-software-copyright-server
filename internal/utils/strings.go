package utils

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"crypto/rand"
	"encoding/base64"
	"math"
	"net/url"
	"regexp"
	"strings"
)

func GetHost(rawUrl string) string {
	if rawUrl == "" {
		return ""
	}
	newUrl, _ := url.Parse(rawUrl)
	return newUrl.Host
}

func RemoveSuffix(str string, suffix string) string {
	if strings.HasSuffix(str, suffix) {
		return str[:len(str)-len(suffix)]
	} else {
		return str // 不需要修改直接返回原始字符串
	}
}

func MaskContent(content string) string {
	contentRune := []rune(content)
	switch len(contentRune) {
	case 1:
		return "*"
	case 2:
		return string(contentRune[0:1]) + "*"
	default:
		maskSize := int(math.Round(float64(len(contentRune)) / 3))
		var str string
		for i := 0; i < maskSize; i++ {
			str += "*"
		}
		return string(contentRune[0:maskSize]) + str + string(contentRune[maskSize+maskSize:])
	}
}

// 判断是否为手机号的函数
func CheckPhone(phone string) bool {
	phoneRegexp := regexp.MustCompile(`^1[0-9]{10,13}$`)
	return phoneRegexp.MatchString(phone)
}

// 判断是否为邮箱的函数
func CheckEmail(email string) bool {
	emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegexp.MatchString(email)
}

// 字符串转换
func TransformNetdiskType(url string) enum.NetdiskType {
	switch {
	case strings.Contains(url, "https://pan.quark.cn/s/"):
		return enum.NetdiskType(2)
	case strings.Contains(url, "https://pan.xunlei.com/s/"):
		return enum.NetdiskType(3)
	case strings.Contains(url, "https://pan.baidu.com/s/"):
		return enum.NetdiskType(4)
	case strings.Contains(url, "https://drive.uc.cn/s/"):
		return enum.NetdiskType(5)
	case strings.Contains(url, "https://caiyun.139.com/m/"):
		return enum.NetdiskType(6)
	case strings.Contains(url, "https://www.123684.com/s/"):
		return enum.NetdiskType(7)
	case strings.Contains(url, "https://pan.wkbrowser.com/netdisk/"):
		return enum.NetdiskType(8)
	case strings.Contains(url, "https://wap.diskyun.com/s/"):
		return enum.NetdiskType(9)
	default:
		return enum.NetdiskType(1)
	}
}

// 字符串转换
func TransformNetdiskName(typ enum.NetdiskType) string {
	switch typ {
	case enum.NetdiskType(2):
		return "夸克网盘"
	case enum.NetdiskType(3):
		return "迅雷网盘"
	case enum.NetdiskType(4):
		return "百度网盘"
	case enum.NetdiskType(5):
		return "UC网盘"
	case enum.NetdiskType(6):
		return "移动云盘"
	case enum.NetdiskType(7):
		return "123网盘"
	case enum.NetdiskType(8):
		return "悟空网盘"
	case enum.NetdiskType(9):
		return "快兔网盘"
	default:
		return "网络来源"
	}
}

// 生成随机字符串
func GenerateRandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}
