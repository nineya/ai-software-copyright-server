package feed

import (
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/global"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func SoftwareCopyright(c *gin.Context) {
	key := "SC_FILE_" + c.Param("id")
	downloadUrl, exits := global.CACHE.GetCache(key)
	if downloadUrl == "" || !exits {
		response.Result(http.StatusNotFound, response.NOT_FOUND, nil, "文件不存在", c)
		return
	}
	global.CACHE.DeleteCache(key)
	// 设置响应头
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(downloadUrl))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Cache-Control", "no-cache")

	// 发送文件
	c.File(downloadUrl)
}
