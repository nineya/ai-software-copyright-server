package redbook

import (
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	rbSev "ai-software-copyright-server/internal/application/service/redbook"
	"github.com/gin-gonic/gin"
)

type CookieApiRouter struct {
	api.BaseApi
}

func (m *CookieApiRouter) InitCookieApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("redbook/cookie")
	m.Router = router
	router.GET("rand", m.GetRand)
}

// @summary Rand get redbook cookie
// @description Rand get redbook cookie
// @tags redbook
// @accept json
// @success 200 {object} response.Response{data=[]table.RedbookCookie}
// @security user
// @router /redbook/cookie/rand [get]
func (m *CookieApiRouter) GetRand(c *gin.Context) {
	mod, err := rbSev.GetCookieService().GetRand()
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
