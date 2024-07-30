package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"tool-server/internal/global"
)

type Response struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

const (
	ERROR        = 500
	UNAUTHORIZED = 401
	FORBIDDEN    = 403
	NOT_FOUND    = 404
	SUCCESS      = 200
)

func Result(httpStatus int, status int, data interface{}, message string, c *gin.Context) {
	if status != SUCCESS {
		global.LOG.Error("Service error message: " + message)
	}
	c.AbortWithStatusJSON(httpStatus, Response{
		status,
		data,
		message,
	})
	c.Abort()
	return

}

func Ok(c *gin.Context) {
	Result(http.StatusOK, SUCCESS, nil, "请求成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(http.StatusOK, SUCCESS, nil, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(http.StatusOK, SUCCESS, data, "请求成功", c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(http.StatusOK, SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(http.StatusInternalServerError, ERROR, nil, "请求失败", c)
}

func UnauthorizedWithError(err error, c *gin.Context) {
	global.LOG.Error(fmt.Sprintf("Service error stack: %+v", err))
	Result(http.StatusUnauthorized, UNAUTHORIZED, nil, err.Error(), c)
}

func ForbiddenWithError(err error, c *gin.Context) {
	global.LOG.Error(fmt.Sprintf("Service error stack: %+v", err))
	Result(http.StatusForbidden, FORBIDDEN, nil, err.Error(), c)
}

func FailWithError(err error, c *gin.Context) {
	global.LOG.Error(fmt.Sprintf("Service error stack: %+v", err))
	Result(http.StatusInternalServerError, ERROR, nil, err.Error(), c)
}

func FailWithMessageAndError(message string, err error, c *gin.Context) {
	global.LOG.Error(fmt.Sprintf("Service error stack: %+v", err))
	Result(http.StatusInternalServerError, ERROR, nil, message, c)
}

func UnauthorizedWithMessage(message string, c *gin.Context) {
	Result(http.StatusUnauthorized, UNAUTHORIZED, nil, message, c)
}

func ForbiddenWithMessage(message string, c *gin.Context) {
	Result(http.StatusForbidden, FORBIDDEN, nil, message, c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(http.StatusInternalServerError, ERROR, nil, message, c)
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(http.StatusInternalServerError, ERROR, data, message, c)
}
