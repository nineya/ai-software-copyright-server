package feed

import (
	"ai-software-copyright-server/internal/application/param/response"
	"github.com/gin-gonic/gin"
)

func Robots(c *gin.Context) {
	htmlResponse := response.GenerateHtmlResult(c)
	htmlResponse.OkWithContentType("internal/feed/robots.tmpl", "text/plain; charset=utf-8")
}
