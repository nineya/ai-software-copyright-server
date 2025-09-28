package netdisk

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/response"
	netdSev "ai-software-copyright-server/internal/application/service/netdisk"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
	"time"
)

// PC转移动流量
func PcToMobile(c *gin.Context) {
	htmlResponse := response.GenerateHtmlResult(c)
	// 先取得资源
	resource := &table.NetdiskResource{}
	if idStr, exist := c.Params.Get("id"); exist { // 通过资源id获取
		id, err := strconv.ParseInt(utils.RemoveSuffix(idStr, ".html"), 10, 64)
		if err == nil && id > 0 {
			resource, err = netdSev.GetResourceService().GetByOnlyId(id)
			if err != nil {
				global.LOG.Warn(fmt.Sprintf("查询网盘资源失败: %+v", err))
			}
		}
	}
	// 没有取得资源，通过解析参数获取资源
	if resource.Id == 0 {
		// 解析用户id
		userId, _ := strconv.ParseInt(c.Query("id"), 10, 64)
		resource.UserId = userId
		// 解析目标地址
		decodeToken, _ := url.QueryUnescape(c.Query("token"))
		tokenBytes, _ := base64.StdEncoding.DecodeString(decodeToken)
		resource.TargetUrl = string(tokenBytes)
		if resource.TargetUrl == "" {
			resource.TargetUrl = global.Host
		}
		resource.Type = utils.TransformNetdiskType(resource.TargetUrl)
		// 解析资源名称
		decodeName, _ := url.QueryUnescape(c.Query("name"))
		nameBytes, _ := base64.StdEncoding.DecodeString(decodeName)
		resource.Name = string(nameBytes)
	}
	if resource.Name == "" {
		resource.Name = "网盘热门资源2000T"
	}

	// 取得短链配置
	headerUserId, _ := strconv.ParseInt(c.Request.Header.Get("User-Id"), 10, 64)
	var configure table.NetdiskShortLinkConfigure
	if resource.UserId > 0 {
		if headerUserId > 0 && headerUserId != resource.UserId {
			response.FailWithError(errors.New("该资源不存在"), c)
			return
		}
		configure, _ = netdSev.GetShortLinkConfigureService().GetByUserId(resource.UserId)
	} else if headerUserId > 0 {
		configure, _ = netdSev.GetShortLinkConfigureService().GetByUserId(headerUserId)
	}

	// 取得自定义配置
	host := global.Host
	favicon := "https://blog.nineya.com/upload/2023/04/favicon.ico"
	custom := false
	// 如果配置id为空，就不判断是不是定制
	if configure.Id > 0 {
		//if configure.CustomHost != "" && configure.CustomHost != global.Host { // 这个办法会导致不设置CustomHost就不会检测自定义域名
		// 只有header有user-id的访问，才认为是定制版的短链访问
		if headerUserId != 0 {
			if configure.CustomExpireTime == nil || configure.CustomExpireTime.Before(time.Now()) {
				response.FailWithError(errors.New("短链定制版服务已过期"), c)
				return
			}
			host = configure.CustomHost
			if configure.CustomFavicon != "" {
				favicon = configure.CustomFavicon
			}
			custom = true
		}
	}
	// 返回参数
	htmlResponse.OkWithData("feature/netdisk/ptm.html", gin.H{
		"Custom":    custom,
		"Host":      host,
		"Favicon":   favicon,
		"Type":      enum.NETDISK_TYPE[resource.Type],
		"Resource":  resource,
		"Configure": configure,
		"Failure":   resource.Status != 0 && resource.Status != 1 && resource.Status != 2,
	})

	//if idStr, exist := c.Params.Get("id"); exist {
	//	id, err := strconv.ParseInt(utils.RemoveSuffix(idStr, ".html"), 10, 64)
	//	if err == nil && id > 0 {
	//		resource, err := netdSev.GetResourceService().GetByOnlyId(id)
	//		configure, _ := netdSev.GetConfigureService().GetByUserId(resource.UserId)
	//		failure := resource.Status != 1 && resource.Status != 2
	//		failureText := ""
	//		if failure {
	//			if err == nil {
	//				failureText = configure.ShortLinkFailureText
	//			}
	//		}
	//		name := resource.Name
	//		if name == "" {
	//			name = "网盘热门资源2000T"
	//		}
	//		if err == nil && resource.TargetUrl != "" {
	//			htmlResponse.OkWithData("feature/netdisk/ptm.html", gin.H{
	//				"Resource":    resource,
	//				"Name":        name,
	//				"Tips":        configure.ShortLinkTips,
	//				"More":        configure.ShortLinkMoreResource,
	//				"Type":        enum.NETDISK_TYPE[resource.Type],
	//				"Failure":     failure,
	//				"FailureText": failureText,
	//			})
	//			return
	//		}
	//		global.LOG.Warn(fmt.Sprintf("查询网盘资源失败: %+v", err))
	//	}
	//}
	//
	//// 解析目标地址
	//decodeToken, _ := url.QueryUnescape(c.Query("token"))
	//tokenBytes, _ := base64.StdEncoding.DecodeString(decodeToken)
	//token := string(tokenBytes)
	//if token == "" {
	//	token = "https://tool.nineya.com"
	//}
	//// 解析用户id
	//headerUserId, _ := strconv.ParseInt(c.Query("id"), 10, 64)
	//var configure table.NetdiskConfigure
	//if headerUserId > 0 {
	//	configure, _ = netdSev.GetConfigureService().GetByUserId(headerUserId)
	//}
	//// 解析资源名称
	//decodeName, _ := url.QueryUnescape(c.Query("name"))
	//nameBytes, _ := base64.StdEncoding.DecodeString(decodeName)
	//name := string(nameBytes)
	//if name == "" {
	//	name = "网盘热门资源2000T"
	//}
	//
	//htmlResponse.OkWithData("feature/netdisk/ptm.html", gin.H{
	//	"Resource": table.NetdiskResource{TargetUrl: token},
	//	"Name":     name,
	//	"Tips":     configure.ShortLinkTips,
	//	"More":     configure.ShortLinkMoreResource,
	//	"Type":     enum.NETDISK_TYPE[utils.TransformNetdiskType(token)],
	//	"Failure":  false,
	//})
}
