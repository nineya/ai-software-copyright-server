package response

import (
	"ai-software-copyright-server/internal/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HtmlResponse struct {
	Context *gin.Context
	Now     time.Time
	Title   string
	Data    any
}

func GenerateHtmlResult(c *gin.Context) HtmlResponse {
	return HtmlResponse{
		Context: c,
		Now:     time.Now(),
	}
}

func (h *HtmlResponse) result(httpStatus int, name string) {
	h.Context.HTML(httpStatus, name, h)
}

func (h *HtmlResponse) Ok(name string) {
	h.result(http.StatusOK, name)
}

func (h *HtmlResponse) OkWithData(name string, data any) {
	h.Data = data
	h.result(http.StatusOK, name)
}

func (h *HtmlResponse) OkWithContentType(name, contentType string) {
	h.Context.Header("Content-Type", contentType)
	h.result(http.StatusOK, name)
}

func (h *HtmlResponse) Error(err error) {
	h.ErrorWithTitle("服务器错误，请稍后重试", err)
}

func (h *HtmlResponse) ErrorWithTitle(title string, err error) {
	h.Data = ErrorResponse{
		Status:  500,
		Title:   title,
		Message: err.Error(),
	}
	global.LOG.Error(fmt.Sprintf("Service error stack: %+v", err))
	h.result(http.StatusInternalServerError, "internal/error/error.html")
}

func (h *HtmlResponse) NotFund() {
	h.Data = ErrorResponse{
		Status:  404,
		Title:   "访问的内容不存在",
		Message: "Not Found.",
	}
	h.result(http.StatusNotFound, "internal/error/error.html")
}
