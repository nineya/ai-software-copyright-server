package utils

import (
	"net/url"
	"strings"
)

func GetHost(rawUrl string) string {
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
