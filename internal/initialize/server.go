package initialize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"tool-server/internal/global"
)

func RunServer(router *gin.Engine) {
	server := initServer(router)
	global.LOG.Error(server.ListenAndServe().Error())
}

func initServer(router *gin.Engine) *http.Server {
	address := fmt.Sprintf(":%d", global.CONFIG.Server.Port)
	global.LOG.Info("Listening and serving HTTP on " + address)
	global.LOG.Info("The working path of tool-server is " + global.WORK_DIR)
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
