package short_link

import (
	"ai-software-copyright-server/internal/application/param/response"
	slSev "ai-software-copyright-server/internal/application/service/short_link"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @summary Short link redirect
// @description Short link redirect
// @tags shortLink
// @accept x-www-form-urlencoded
// @success 200
// @security content
// @router /s/{alias} [get]
func Redirect(c *gin.Context) {
	alias := c.Param("alias")
	mod, err := slSev.GetShortLinkService().GetByAlias(alias)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	if mod.Id != 0 {
		_ = slSev.GetShortLinkService().UpdateVisitsIncreaseById(mod.Id, c.Request.UserAgent())
	}
	//// 永久重定向，后面修改链接了，对已经访问过的人不会生效
	//c.Redirect(http.StatusMovedPermanently, mod.TargetUrl)
	// 临时重定向
	c.Redirect(http.StatusFound, mod.TargetUrl)
}
