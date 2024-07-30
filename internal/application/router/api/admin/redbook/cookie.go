package redbook

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"tool-server/internal/application/model/table"
	"tool-server/internal/application/param/response"
	"tool-server/internal/application/router/api"
	rbSev "tool-server/internal/application/service/redbook"
)

type CookieApiRouter struct {
	api.BaseApi
}

func (m *CookieApiRouter) InitCookieApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("redbook/cookie")
	m.Router = router
	router.POST("batch", m.CreateInBatch)
	router.GET("all", m.GetAll)
}

// @summary Create redbook cookie in batches
// @description Create redbook cookie in batches
// @tags redbook
// @accept json
// @param param body []table.RedbookCookie true "Redbooke cookie information"
// @success 200 {object} response.Response
// @security admin
// @router /redbook/cookie/batch [post]
func (m *CookieApiRouter) CreateInBatch(c *gin.Context) {
	var param []table.RedbookCookie
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	if err = rbSev.GetCookieService().CreateInBatch(m.GetUserId(c), param); err != nil {
		response.FailWithError(err, c)
		return
	}
	m.Log(c, "CREATED_PHOTO", fmt.Sprintf("批量创建小红书 Cookie %d 条", len(param)))
	response.Ok(c)
}

// @summary Get all redbook cookie
// @description Get all redbook cookie
// @tags redbook
// @accept json
// @success 200 {object} response.Response{data=[]table.RedbookCookie}
// @security admin
// @router /redbook/cookie/all [get]
func (m *CookieApiRouter) GetAll(c *gin.Context) {
	mod, err := rbSev.GetCookieService().GetAll(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
