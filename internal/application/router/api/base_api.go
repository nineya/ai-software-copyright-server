package api

import (
	"github.com/gin-gonic/gin"
	"tool-server/internal/application/model/enum"
	"tool-server/internal/application/model/table"
	"tool-server/internal/application/param/request"
	"tool-server/internal/application/service/log"
)

type BaseApi struct {
	Router *gin.RouterGroup
}

func (b *BaseApi) GetClaims(c *gin.Context) *request.CustomClaims {
	if claims, exists := c.Get("claims"); exists {
		waitUse := claims.(*request.CustomClaims)
		return waitUse
	}
	return nil
}

func (b *BaseApi) GetUserId(c *gin.Context) int64 {
	if claims := b.GetClaims(c); claims != nil {
		return claims.UserId
	}
	return 0
}

// 添加日志
func (b *BaseApi) Log(c *gin.Context, typ string, content string) {
	b.LogBySiteAndAdminId(c, b.GetUserId(c), typ, content)
}

// 添加日志
func (b *BaseApi) LogBySiteAndAdminId(c *gin.Context, adminId int64, typ string, content string) {
	logType, _ := enum.LogTypeValue(typ)
	_, _ = log.GetLogService().Create(adminId, table.Log{
		Content:   content,
		IpAddress: c.ClientIP(),
		Type:      logType,
	})
}
