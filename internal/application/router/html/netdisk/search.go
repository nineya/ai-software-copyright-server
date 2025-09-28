package netdisk

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	netdSev "ai-software-copyright-server/internal/application/service/netdisk"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// 首页
func Index(c *gin.Context) {
	htmlResponse := response.GenerateHtmlResult(c)
	configure, err := netdSev.GetSearchSiteConfigureService().GetByUserId(utils.GetHeaderUserId(c))
	if err != nil {
		global.LOG.Error(fmt.Sprintf("查询网盘搜索站点配置信息错误:%+v", err))
		htmlResponse.OkWithData("feature/netdisk/search/error.html", gin.H{
			"Title":     "系统繁忙",
			"Message":   "查询站点配置信息错误，请稍后再试。",
			"Configure": configure,
		})
		return
	}
	if configure.ExpireTime == nil || configure.ExpireTime.Before(time.Now()) {
		global.LOG.Error("网盘搜索站点服务已过期")
		htmlResponse.OkWithData("feature/netdisk/search/error.html", gin.H{
			"Title":     "站点服务已过期",
			"Message":   "网盘搜索站点服务已过期，请联系管理员。",
			"Configure": configure,
		})
		return
	}
	if configure.BrowserTips && strings.Contains(c.Request.UserAgent(), "MicroMessenger") {
		htmlResponse.OkWithData("feature/netdisk/search/browser_tips.html", gin.H{
			"Configure": configure,
		})
		return
	}
	param := request.NetdiskResourceSearchParam{QueryPageParam: request.QueryPageParam{PageableParam: request.PageableParam{Page: 0, Size: 15}}}
	page, err := netdSev.GetResourceService().Search(configure.UserId, enum.ClientType(9), param)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("查询最新资源信息错误:%+v", err))
	}

	htmlResponse.OkWithData("feature/netdisk/search/index.html", gin.H{
		"Page":      page,
		"Configure": configure,
	})
}

// 搜索页
func Search(c *gin.Context) {
	htmlResponse := response.GenerateHtmlResult(c)
	configure, err := netdSev.GetSearchSiteConfigureService().GetByUserId(utils.GetHeaderUserId(c))
	if err != nil {
		global.LOG.Error(fmt.Sprintf("查询网盘搜索站点配置信息错误:%+v", err))
		htmlResponse.OkWithData("feature/netdisk/search/error.html", gin.H{
			"Title":     "系统繁忙",
			"Message":   "查询站点配置信息错误，请稍后再试。",
			"Configure": configure,
		})
		return
	}
	if configure.ExpireTime == nil || configure.ExpireTime.Before(time.Now()) {
		global.LOG.Error("网盘搜索站点服务已过期")
		htmlResponse.OkWithData("feature/netdisk/search/error.html", gin.H{
			"Title":     "站点服务已过期",
			"Message":   "网盘搜索站点服务已过期，请联系管理员。",
			"Configure": configure,
		})
		return
	}
	if configure.BrowserTips && strings.Contains(c.Request.UserAgent(), "MicroMessenger") {
		htmlResponse.OkWithData("feature/netdisk/search/browser_tips.html", gin.H{
			"Configure": configure,
		})
		return
	}
	keyword := c.Query("keyword")
	if keyword == "" {
		global.LOG.Error("请求参数不能为空")
		htmlResponse.OkWithData("feature/netdisk/search/error.html", gin.H{
			"Title":     "参数错误",
			"Message":   "参数校验错误，请稍后再试。",
			"Configure": configure,
		})
		return
	}
	param := request.NetdiskResourceSearchParam{
		QueryPageParam: request.QueryPageParam{
			QueryParam:    request.QueryParam{Keyword: keyword},
			PageableParam: request.PageableParam{Page: 0, Size: 100},
		},
		CollectTypes: utils.ListTransform(configure.CollectTypes, func(item enum.NetdiskType) string {
			return enum.NETDISK_TYPE[item]
		}),
	}
	page, err := netdSev.GetResourceService().Search(configure.UserId, enum.ClientType(9), param)
	if err != nil {
		global.LOG.Error(fmt.Sprintf("查询网盘资源错误:%+v", err))
		htmlResponse.OkWithData("feature/netdisk/search/error.html", gin.H{
			"Title":     "系统繁忙",
			"Message":   "查询网盘资源失败，请稍后再试。",
			"Configure": configure,
		})
		return
	}
	list := *page.Content.(*[]table.NetdiskResource)
	for i, item := range list {
		if item.Id == 0 {
			bytes, _ := json.Marshal(map[string]string{"name": item.Name, "shareTargetUrl": item.ShareTargetUrl, "sharePwd": item.SharePwd})
			bytes, _ = utils.AesEncrypt(bytes, global.AesKey)
			list[i].ShareTargetUrl = url.QueryEscape(base64.StdEncoding.EncodeToString(bytes))
		}
	}
	htmlResponse.OkWithData("feature/netdisk/search/search.html", gin.H{
		"Keyword":   keyword,
		"Page":      page,
		"Configure": configure,
	})
}

// @summary 转存加密的网盘资源
// @description 转存加密的网盘资源
// @tags netdisk
// @accept json
// @param param body table.NetdiskResource true "网盘资源信息"
// @success 200 {object} response.Response{data=string}
// @security user
// @router /public/netdisk/resource/detail [post]
func Detail(c *gin.Context) {
	htmlResponse := response.GenerateHtmlResult(c)
	configure, err := netdSev.GetSearchSiteConfigureService().GetByUserId(utils.GetHeaderUserId(c))
	if err != nil {
		global.LOG.Error(fmt.Sprintf("查询网盘搜索站点配置信息错误:%+v", err))
		htmlResponse.OkWithData("feature/netdisk/search/error.html", gin.H{
			"Title":     "系统繁忙",
			"Message":   "查询站点配置信息错误，请稍后再试。",
			"Configure": configure,
		})
		return
	}
	if configure.ExpireTime == nil || configure.ExpireTime.Before(time.Now()) {
		global.LOG.Error("网盘搜索站点服务已过期")
		htmlResponse.OkWithData("feature/netdisk/search/error.html", gin.H{
			"Title":     "站点服务已过期",
			"Message":   "网盘搜索站点服务已过期，请联系管理员。",
			"Configure": configure,
		})
		return
	}
	if configure.BrowserTips && strings.Contains(c.Request.UserAgent(), "MicroMessenger") {
		htmlResponse.OkWithData("feature/netdisk/search/browser_tips.html", gin.H{
			"Configure": configure,
		})
		return
	}
	mod := &table.NetdiskResource{}
	if val, exist := c.Params.Get("id"); exist {
		id, err := strconv.ParseInt(utils.RemoveSuffix(val, ".html"), 10, 64)
		if err != nil {
			global.LOG.Error("参数获取失败")
			htmlResponse.OkWithData("feature/netdisk/search/error.html", gin.H{
				"Title":     "参数错误",
				"Message":   "参数校验错误，请稍后再试。",
				"Configure": configure,
			})
			return
		}
		if id > 0 {
			resource, err := netdSev.GetResourceService().GetByOnlyId(id)
			if err != nil {
				global.LOG.Error("参数获取失败")
				htmlResponse.OkWithData("feature/netdisk/search/error.html", gin.H{
					"Title":     "参数错误",
					"Message":   "参数校验错误，无法获取资源信息，请稍后再试。",
					"Configure": configure,
				})
				return
			}
			mod = resource
		}
	}
	// 还未取得参数，证明该数据还未存储
	if mod.Id == 0 {
		t, _ := url.QueryUnescape(c.Query("t"))
		if t == "" {
			global.LOG.Error("请求参数不能为空")
			htmlResponse.OkWithData("feature/netdisk/search/error.html", gin.H{
				"Title":     "参数错误",
				"Message":   "参数校验错误，请稍后再试。",
				"Configure": configure,
			})
			return
		}
		tokenBytes, _ := base64.StdEncoding.DecodeString(t)
		paramBytes, err := utils.AesEncrypt(tokenBytes, global.AesKey)
		if err != nil {
			global.LOG.Error("参数加载错误")
			htmlResponse.OkWithData("feature/netdisk/search/error.html", gin.H{
				"Title":     "参数错误",
				"Message":   "参数加载错误，请稍后再试。",
				"Configure": configure,
			})
			return
		}
		var param request.NetdiskResourceSaveParam
		err = json.Unmarshal(paramBytes, &param)
		if err != nil {
			global.LOG.Error("解析参数错误")
			htmlResponse.OkWithData("feature/netdisk/search/error.html", gin.H{
				"Title":     "参数错误",
				"Message":   "解析参数错误，请稍后再试。",
				"Configure": configure,
			})
			return
		}
		resource, err := netdSev.GetResourceService().Save(configure.UserId, param)
		if err != nil {
			global.LOG.Error(fmt.Sprintf("加载参数错误:%+v", err))
			mod = &param.NetdiskResource
		} else {
			if resource.UpdateTime == nil {
				now := time.Now()
				resource.UpdateTime = &now
			}
			mod = resource
		}
	}
	// 如果不是自己的资源，就去掉资源id，避免报错
	if mod.UserId != configure.UserId {
		mod.Id = 0
	}
	// 创建网盘短链
	jumpUrl := mod.TargetUrl
	if configure.UseShortLink && mod.TargetUrl != "" {
		host := global.Host
		slConfigure, _ := netdSev.GetShortLinkConfigureService().GetByUserId(configure.UserId)
		if slConfigure.CustomExpireTime != nil && slConfigure.CustomExpireTime.After(time.Now()) && slConfigure.CustomHost != "" {
			host = slConfigure.CustomHost
		}
		jumpUrl = fmt.Sprintf(
			"%s/ptm/%d.html?id=%d&token=%s&name=%s",
			host,
			mod.Id,
			configure.UserId,
			url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(mod.TargetUrl))),
			url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(mod.Name))),
		)
	}

	htmlResponse.OkWithData("feature/netdisk/search/detail.html", gin.H{
		"Resource":  mod,
		"JumpUrl":   jumpUrl,
		"Configure": configure,
	})
}
