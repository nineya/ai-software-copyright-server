package middleware

import (
	"github.com/gin-gonic/gin"
	"tool-server/internal/application/model/table"
	"tool-server/internal/global"
)

// 访问统计
func VisitStatisticHandler(c *gin.Context) {
	defer func() {
		go func() {
			statistic := table.Statistic{
				Url:        c.Request.RequestURI,
				IpAddress:  c.ClientIP(),
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
