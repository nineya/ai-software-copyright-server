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
	router.POST("use", m.Use)
}

// @summary 使用Cdkey
// @description 使用Cdkey
// @tags cdkey
// @accept json
// @param param body request.CdkeyCreateParam true "创建Cdkey信息"
// @success 200 {object} response.Response{data=string}
// @security user
// @router /cdkey [post]
func (m *CdkeyApiRouter) Use(c *gin.Context) {
	var param request.CdkeyUseParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := cdkeySev.GetCdkeyService().Use(m.GetUserId(c), param.Cdkey)
	if err != nil {
		m.UserLog(c, "CDKEY_USE", fmt.Sprintf("使用Cdkey失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "CDKEY_USE", fmt.Sprintf("使用Cdkey %s 成功，添加 %d 积分", param.Cdkey, mod.NyCredits))
	response.OkWithData(mod, c)
}
