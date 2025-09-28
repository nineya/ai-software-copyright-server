package cdkey

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	cdkeySev "ai-software-copyright-server/internal/application/service/cdkey"
	"fmt"
	"github.com/gin-gonic/gin"
)

type CdkeyApiRouter struct {
	api.BaseApi
}

func (m *CdkeyApiRouter) InitCdkeyApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("cdkey")
	m.Router = router
	router.POST("", m.Create)
}

// @summary 创建Cdkey
// @description 创建Cdkey
// @tags cdkey
// @accept json
// @param param body request.CdkeyCreateParam true "创建Cdkey信息"
// @success 200 {object} response.Response{data=string}
// @security user
// @router /cdkey [post]
func (m *CdkeyApiRouter) Create(c *gin.Context) {
	var param request.CdkeyCreateParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := cdkeySev.GetCdkeyService().Create(m.GetUserId(c), param)
	if err != nil {
		m.AdminLog(c, "CDKEY_CREATE", fmt.Sprintf("创建Cdkey失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.AdminLog(c, "CDKEY_CREATE", fmt.Sprintf("创建Cdkey %d 条", param.CdkeyNum))
	response.OkWithData(mod, c)
}
