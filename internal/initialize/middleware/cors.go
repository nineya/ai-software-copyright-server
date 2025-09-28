package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 处理跨域请求,支持options访问
func RouterGroupCors(group *gin.Engine) gin.HandlerFunc {
	// 因为没有全局注册，为了避免post请求的options请求404，必须加上这段路由
	// 如果是全局注册，则可不需要这段内容
	group.OPTIONS("/*options_support", Cors)
	return Cors
}

func Cors(c *gin.Context) {
	method := c.Request.Method
	//requestUrl := c.Request.RequestURI
	origin := c.Request.Header.Get("Origin")
	c.Header("Access-Control-Allow-Origin", origin)
	c.Header("Access-Control-Allow-Headers", "Content-Type, AccessToken, X-CSRF-Token, Admin-Authorization, User-Authorization, User-Access-Key, Authorization, Token, X-Token, X-User-Id")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PUT")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	// 放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}

	// 处理请求
	c.Next()
}
