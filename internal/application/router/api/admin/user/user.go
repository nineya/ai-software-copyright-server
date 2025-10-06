package user

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserApiRouter struct {
	api.BaseApi
}

var store = utils.NewCaptchaStore()

func (m *UserApiRouter) InitUserApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("user")
	m.Router = router
	router.POST("auditShare", m.AuditShare)
	router.POST("addCredits", m.AddCredits)
}

// @summary 审核分享
// @description 审核分享
// @tags public,user
// @accept json
// @param param body table.ShareRecord true "分享审核信息"
// @success 200 {object} response.Response{data=[]table.ShareRecord}
// @router /user/auditShare [post]
func (m *UserApiRouter) AuditShare(c *gin.Context) {
	var param table.ShareRecord
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := userSev.GetShareRecordService().Audit(param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}

// @summary Add credits
// @description Add credits
// @tags public,user
// @accept json
// @param param body table.UserRewardCreditsParam true "Add credits information"
// @success 200 {object} response.Response{data=[]table.User}
// @router /user/addCredits [post]
func (m *UserApiRouter) AddCredits(c *gin.Context) {
	var param request.UserAddCreditsParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := userSev.GetUserService().AddCredits(param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
