package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tool-server/internal/application/param/response"
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
