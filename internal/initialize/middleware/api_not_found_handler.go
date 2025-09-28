package middleware

import (
	"ai-software-copyright-server/internal/application/param/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 统一API异常处理
func ApiNotFoundHandler(c *gin.Context) {
	c.JSON(http.StatusNotFound, response.Response{
		response.NOT_FOUND,
		nil,
		"接口不存在",
	})
	c.Abort()
}
