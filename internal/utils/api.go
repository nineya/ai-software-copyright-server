package utils

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"
)

func GetClientType(c *gin.Context) enum.ClientType {
	clientTypeStr := c.Request.Header.Get("Client-Type")
	// TODO 适配旧版本的网盘搜索小程序的客户端类型
	if clientTypeStr == "NETDISK_SEARCH" {
		clientTypeStr = "NETDISK_SEARCH_WXAMP"
	}
	clientType, err := enum.ClientTypeValue(clientTypeStr)
	if err != nil {
		PanicErr(errors.Wrap(err, "系统异常"))
	}
	return clientType
}

func GetHeaderUserId(c *gin.Context) int64 {
	userIdStr := c.Request.Header.Get("User-Id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		PanicErr(errors.Wrap(err, "参数错误"))
	}
	return userId
}
