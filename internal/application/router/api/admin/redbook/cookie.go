package redbook

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	rbSev "ai-software-copyright-server/internal/application/service/redbook"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type CookieApiRouter struct {
	api.BaseApi
}

func (m *CookieApiRouter) InitCookieApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("redbook/cookie")
	m.Router = router
	router.POST("batch", m.CreateInBatch)
	router.PUT(":id/status/:status", m.UpdateStatusById)
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
	m.AdminLog(c, "CREATED_PHOTO", fmt.Sprintf("批量创建小红书 Cookie %d 条", len(param)))
	response.Ok(c)
}

// @summary Update redbook cookie status by id
// @description Update redbook cookie status by id
// @tags redbook
// @accept json
// @param id path int64 true "Cookie id"
// @param status path string true "Redbook cookie new status"
// @success 200 {object} response.Response{}
// @security admin
// @router /redbook/cookie/{id}/status/{status} [put]
func (m *CookieApiRouter) UpdateStatusById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	status, err := enum.CookieStatusValue(c.Param("status"))
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	err = rbSev.GetCookieService().UpdateStatusById(m.GetUserId(c), id, status)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	m.AdminLog(c, "REDBOOK_UPDATE_COOKIE", fmt.Sprintf("修改 Id 为 %d 的小红书cookie状态为 %s", id, enum.TASK_STATUS[status]))
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
