package software_copyright

import (
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type DownloadApiRouter struct {
	api.BaseApi
}

func (m *DownloadApiRouter) InitDownloadApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("softwareCopyright/download")
	m.Router = router
	router.GET(":id/:filename", m.Download)
}

// @summary 下载著作权文件
// @description 下载著作权文件
// @tags softwareCopyright
// @accept json
// @produce octet-stream
// @param id path int64 true "软著申请id"
// @param filename path string true "文件名称"
// @success 200 {file} byte "文件内容"
// @security user
// @router /softwareCopyright/download/{id}/{filename} [get]
func (m *DownloadApiRouter) Download(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}
	filename := strings.TrimSpace(c.Param("filename"))
	if filename == "" {
		response.FailWithMessageAndError("参数获取失败", err, c)
		return
	}

	downloadPath := filepath.Join(utils.GetSoftwareCopyrightPath(id), filename)
	// 检查文件是否存在
	if _, err = os.Stat(downloadPath); os.IsNotExist(err) {
		response.FailWithMessageAndError("文件不存在", err, c)
		return
	}

	downloadId := uuid.NewString()
	err = global.CACHE.SetCache("SC_FILE_"+downloadId, downloadPath, time.Minute)
	if err != nil {
		response.FailWithMessageAndError("获取文件下载链接失败", err, c)
		return
	}
	response.OkWithData(global.Host+"/download/softwareCopyright/"+downloadId, c)
}
