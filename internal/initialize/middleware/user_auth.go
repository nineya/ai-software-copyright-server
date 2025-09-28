package middleware

import (
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/param/response"
	userSev "ai-software-copyright-server/internal/application/service/user"
	"ai-software-copyright-server/internal/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func UserAuth(c *gin.Context) {
	// 小程序登录
	unionid := c.Request.Header.Get("User-Unionid")
	// TODO 临时过渡，适配工具人
	if unionid == "" {
		temp := c.Request.Header.Get("User-Access-Key")
		if strings.HasPrefix(temp, "owNRL") {
			unionid = temp
		}
	}
	// TODO 临时过度
	if unionid == "" {
		token := c.Request.Header.Get("User-Authorization")
		if len(token) < 36 {
			unionid = token
		}
	}
	if unionid != "" {
		user, err := userSev.GetUserService().GetByWxUnionid(unionid)
		if err != nil || user.Id == 0 {
			response.UnauthorizedWithMessage("登录状态已失效", c)
			return
		}
		claims := request.CustomClaims{
			UserId: user.Id,
			Type:   global.AuthToken,
		}
		c.Set("claims", &claims)
		c.Next()
		return
	}
	// accessKey登录
	accessKey := c.Request.Header.Get("User-Access-Key")
	if accessKey != "" {
		user, err := userSev.GetUserService().GetByAccessKey(accessKey)
		if err != nil || user.Id == 0 {
			response.UnauthorizedWithMessage("无效的AccessKey", c)
			return
		}
		claims := request.CustomClaims{
			UserId: user.Id,
			Type:   global.AuthToken,
		}
		c.Set("claims", &claims)
		c.Next()
		return
	}
	// 网页登录
	token := c.Request.Header.Get("User-Authorization")
	if token != "" {
		// parseToken 解析token包含的信息
		claims, err := global.JWT.ParseToken(token)
		if err != nil || claims.Type != global.AuthToken {
			response.UnauthorizedWithMessage("登录状态已失效", c)
			return
		}
		checkKey := fmt.Sprintf("%s_%s_%d_%s", global.AuthToken, claims.UserType, claims.UserId, claims.Id)
		if _, exist := global.CACHE.GetCache(checkKey); !exist {
			response.UnauthorizedWithMessage("登录状态已失效", c)
			return
		}
		c.Set("claims", claims)
		c.Next()
		return
	}
	response.UnauthorizedWithMessage("登录状态已失效", c)
}
