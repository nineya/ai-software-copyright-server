package user

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type UserApiRouter struct {
	api.BaseApi
}

func (m *UserApiRouter) InitUserApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("user")
	m.Router = router
	router.POST("rewardAd", m.RewardAd)
	router.POST("rewardGoods", m.RewardGoods)
	router.PUT("accessKey", m.UpdateAccessKey)
	router.PUT("updateInfo", m.UpdateUserInfo)
	router.GET("getInviteInfo", m.GetInviteInfo)
	router.GET("getInviteCode", m.GetInviteCode)
	router.GET("profiles", m.GetProfiles)
}

// @summary Reward ad
// @description Reward ad
// @tags user
// @accept json
// @success 200 {object} response.Response{data=table.User}
// @security user
// @router /user/rewardAd [post]
func (m *UserApiRouter) RewardAd(c *gin.Context) {
	var param request.UserRewardNyCreditsParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := userSev.GetUserService().ChangeNyCredits(m.GetUserId(c), table.CreditsChange{Type: enum.CreditsChangeType(2), ChangeCredits: param.RewardCredits, Remark: param.Remark})
	if err != nil {
		m.UserLog(c, "USER_REWARD_AD", fmt.Sprintf("用户观看激励广告失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "USER_REWARD_AD", fmt.Sprintf("用户 %s 观看激励广告，赠送 %d 个积分", mod.Nickname, param.RewardCredits))
	response.OkWithData(mod, c)
}

// @summary Reward goods
// @description Reward goods
// @tags user
// @accept json
// @success 200 {object} response.Response{data=table.User}
// @security user
// @router /user/rewardGoods [post]
func (m *UserApiRouter) RewardGoods(c *gin.Context) {
	var param request.UserRewardGoodsParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	mod, err := userSev.GetUserService().PaymentNyCredits(m.GetUserId(c), enum.BuyType(15), 0, fmt.Sprintf("购买激励物品《%s》", param.Name))
	if err != nil {
		m.UserLog(c, "USER_REWARD_GOODS", fmt.Sprintf("用户观看激励广告失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "USER_REWARD_AD", fmt.Sprintf("用户 %s 观看激励广告，获赠物品《%s》", mod.Nickname, param.Name))
	response.OkWithData(mod, c)
}

// @summary 更新AccessKey
// @description 更新AccessKey
// @tags user
// @accept json
// @success 200 {object} response.Response{data=string}
// @security user
// @router /user/accessKey [put]
func (m *UserApiRouter) UpdateAccessKey(c *gin.Context) {
	mod, err := userSev.GetUserService().UpdateAccessKey(m.GetUserId(c))
	if err != nil {
		m.UserLog(c, "USER_ACCESS_KEY_UPDATE", fmt.Sprintf("用户AccessKey更新失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "USER_ACCESS_KEY_UPDATE", "用户AccessKey更新")
	response.OkWithData(mod, c)
}

// @summary 更新用户信息
// @description 更新用户信息
// @tags user
// @accept json
// @success 200 {object} response.Response{data=*response.UserRewardResponse}
// @security user
// @router /user/updateInfo [put]
func (m *UserApiRouter) UpdateUserInfo(c *gin.Context) {
	var param request.UserUpdateInfoParam
	err := c.ShouldBindJSON(&param)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	if param.Phone != "" && !utils.CheckPhone(param.Phone) {
		response.FailWithError(errors.New("手机号格式错误"), c)
		return
	}
	if param.Email != "" && !utils.CheckEmail(param.Email) {
		response.FailWithError(errors.New("邮箱格式错误"), c)
		return
	}
	mod, err := userSev.GetUserService().UpdateUserInfo(m.GetUserId(c), param)
	if err != nil {
		m.UserLog(c, "USER_INFO_UPDATE", fmt.Sprintf("用户信息更新失败，原因：%s", err.Error()))
		response.FailWithError(err, c)
		return
	}
	m.UserLog(c, "USER_INFO_UPDATE", "用户信息更新")
	response.OkWithData(mod, c)
}

// @summary Get user invite Info
// @description Get user invite Info
// @tags user
// @accept json
// @success 200 {object} response.Response{data=table.InviteInfo}
// @security user
// @router /user/getInviteInfo [post]
func (m *UserApiRouter) GetInviteInfo(c *gin.Context) {
	mod, err := userSev.GetInviteRecordService().GetInviteInfo(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(mod, c)
}

// @summary Get invite code url
// @description Get invite code url
// @tags user
// @accept json
// @success 200 {object} response.Response{data=string}
// @security user
// @router /user/inviteCode [get]
func (m *UserApiRouter) GetInviteCode(c *gin.Context) {
	url, err := userSev.GetClientInfoService().GetInviteCode(m.GetUserId(c), m.GetClientType(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(url, c)
}

// @summary Get user profiles
// @description Get user profiles
// @tags user
// @Produce json
// @success 200 {object} response.Response{data=table.User}
// @Security admin
// @router /user/profiles [get]
func (m *UserApiRouter) GetProfiles(c *gin.Context) {
	user, err := userSev.GetUserService().GetById(m.GetUserId(c))
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	response.OkWithData(user, c)
}
