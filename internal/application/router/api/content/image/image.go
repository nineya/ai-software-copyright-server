package image

import (
	"ai-software-copyright-server/internal/application/param/response"
	"ai-software-copyright-server/internal/application/router/api"
	attaSev "ai-software-copyright-server/internal/application/service/attachment"
	"fmt"
	"github.com/gin-gonic/gin"
)

type ImageApiRouter struct {
	api.BaseApi
}

func (m *ImageApiRouter) InitImageApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("image")
	router.POST("upload/:bucket", m.Upload)
}

// @summary Upload image attachment
// @description Upload image attachment
// @tags image
// @accept x-www-form-urlencoded
// @param file formData file true "Image file stream"
// @param bucket path string true "Attachment team"
// @success 200 {object} response.Response{data=table.Attachment}
// @security content
// @router /image/upload/{bucket} [post]
func (m *ImageApiRouter) Upload(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	bucket := c.Param("bucket")
	mod, err := attaSev.GetImageService().Upload(file, bucket)
	if err != nil {
		response.FailWithError(err, c)
		return
	}
	m.AdminLog(c, "UPLOAD_IMAGE", fmt.Sprintf("上传图片 %s 到路径：%s", file.Filename, mod.Url))
	response.OkWithData(mod, c)
}
