package middleware

import (
	"github.com/gin-gonic/gin"
	"tool-server/internal/application/model/errors"
	"tool-server/internal/application/param/response"
)

// 统一API异常处理
func ApiErrorHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case errors.ForbiddenError:
				response.ForbiddenWithError(err.(errors.ForbiddenError), c)
			default:
				response.FailWithError(err.(error), c)
			}
		}
	}()
	c.Next()
}
