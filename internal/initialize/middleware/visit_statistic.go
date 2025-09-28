package middleware

import (
	"ai-software-copyright-server/internal/application/model/table"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/utils"
	"github.com/gin-gonic/gin"
)

// 访问统计
func VisitStatisticHandler(c *gin.Context) {
	defer func() {
		go func() {
			statistic := table.Statistic{
				Url:        c.Request.RequestURI,
				IpAddress:  c.ClientIP(),
				Referrer:   c.Request.Referer(),
				Origin:     utils.GetHost(c.Request.Referer()),
				HttpStatus: c.Writer.Status(),
				UserAgent:  c.Request.UserAgent(),
			}
			_, err := global.DB.Insert(&statistic)
			if err != nil {
				global.LOG.Sugar().Errorf("统计数据写入数据库失败： %+v", statistic, err)
			}
		}()
	}()
	c.Next()
}
