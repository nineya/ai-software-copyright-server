package user

import (
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
	router.POST("addCredits", m.AddCredits)
}

// @summary Add credits
// @description Add credits
// @tags public,user
// @accept json
// @param param body table.UserRewardNyCreditsParam true "Add credits information"
// @success 200 {object} response.Response{data=[]table.User}
// @router /user/addCredits [post]
func (m *UserApiRouter) AddCredits(c *gin.Context) {
	var param request.UserAddNyCreditsParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := userSev.GetUserService().AddNyCredits(param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}
