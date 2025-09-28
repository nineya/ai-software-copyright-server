package initialize

import (
	"ai-software-copyright-server/internal/global"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RunServer(router *gin.Engine) {
	server := initServer(router)
	global.LOG.Error(server.ListenAndServe().Error())
}

func initServer(router *gin.Engine) *http.Server {
	address := fmt.Sprintf(":%d", global.CONFIG.Server.Port)
	global.LOG.Info("Listening and serving HTTP on " + address)
	global.LOG.Info("The working path of ai-software-copyright-server is " + global.WORK_DIR)
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    120 * time.Second,
		WriteTimeout:   120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
