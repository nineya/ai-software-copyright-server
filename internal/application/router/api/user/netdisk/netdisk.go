package netdisk

import (
	"ai-software-copyright-server/internal/application/router/api"
	"github.com/gin-gonic/gin"
)

type NetdiskApiRouter struct {
	api.BaseApi
}

func (m *NetdiskApiRouter) InitNetdiskApiRouter(Router *gin.RouterGroup) {
	router := Router.Group("netdisk")
	m.Router = router
}
