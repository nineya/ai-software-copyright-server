package middleware

import (
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/global"
	"fmt"
	"github.com/gin-gonic/gin"
)

func AdminAuth(c *gin.Context) {
	token := c.Request.Header.Get("Admin-Authorization")
	if token == "" {
		response.UnauthorizedWithMessage("登录状态已失效", c)
		return
	}
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
}
