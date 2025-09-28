package api

import (
	"ai-software-copyright-server/internal/application/model/enum"
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/application/param/request"
	"ai-software-copyright-server/internal/application/service/log"
	"ai-software-copyright-server/internal/utils"
	"github.com/gin-gonic/gin"
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

func (b *BaseApi) GetClientType(c *gin.Context) enum.ClientType {
	return utils.GetClientType(c)
}

// 添加日志
func (b *BaseApi) AdminLog(c *gin.Context, typ string, content string) {
	b.LogByAdminId(c, b.GetUserId(c), typ, content)
}

func (b *BaseApi) UserLog(c *gin.Context, typ string, content string) {
	b.LogByUserId(c, b.GetUserId(c), typ, content)
}

// 添加管理员日志
func (b *BaseApi) LogByAdminId(c *gin.Context, adminId int64, typ string, content string) {
	logType, _ := enum.AdminLogTypeValue(typ)
	_, _ = log.GetAdminService().Create(adminId, table.AdminLog{
		Content:   content,
		IpAddress: c.ClientIP(),
		Type:      logType,
	})
}

// 添加用户日志
func (b *BaseApi) LogByUserId(c *gin.Context, userId int64, typ string, content string) {
	logType, _ := enum.UserLogTypeValue(typ)
	_, _ = log.GetUserService().Create(userId, table.UserLog{
		ClientType: b.GetClientType(c),
		Content:    content,
		IpAddress:  c.ClientIP(),
		Type:       logType,
	})
}
