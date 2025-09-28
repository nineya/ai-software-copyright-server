package middleware

import (
	modelErrors "ai-software-copyright-server/internal/application/model/errors"
	"ai-software-copyright-server/internal/application/param/response"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// 统一API异常处理
func ApiErrorHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			switch err.(type) {
			case modelErrors.ForbiddenError:
				response.ForbiddenWithError(err.(modelErrors.ForbiddenError), c)
			case string:
				response.FailWithError(errors.New(err.(string)), c)
			default:
				response.FailWithError(err.(error), c)
			}
		}
	}()
	c.Next()
}
