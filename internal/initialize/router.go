package initialize

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	_ "tool-server/docs/admin"
	_ "tool-server/docs/content"
	"tool-server/internal/application/router/api/admin"
	"tool-server/internal/global"
	"tool-server/internal/initialize/middleware"
)

func InitRouter() *gin.Engine {
	if global.CONFIG.Server.Mode == "dev" {
		global.LOG.Info("The service is in development mode.")
		gin.SetMode(gin.DebugMode)
	} else {
		global.LOG.Info("The service is in production mode.")
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(middleware.RouterGroupCors(router))
	//router.Use(ginLogger(), ginRecovery(true))
	initSwaggerRouter(router)
	initApiRouter(router)
	global.LOG.Info("The route mapping initialization is complete.")
	return router
}

func initSwaggerRouter(router *gin.Engine) {
	swaggerGroup := router.Group("swagger")
	swaggerGroup.GET("/admin/*any", ginSwagger.WrapHandler(
		swaggerFiles.NewHandler(), func(config *ginSwagger.Config) {
			config.InstanceName = "admin"
		}))
	swaggerGroup.GET("/content/*any", ginSwagger.WrapHandler(
		swaggerFiles.NewHandler(), func(config *ginSwagger.Config) {
			config.InstanceName = "content"
		}))
}

func initApiRouter(router *gin.Engine) {
	apiGroup := router.Group("api")
	apiGroup.Use(middleware.ApiErrorHandler)
	{
		apiGroup.GET("/health", func(c *gin.Context) {
			c.String(http.StatusOK, "ok health")
		})
	}

	adminRouterGroup := admin.ApiRouterGroup
	adminGroup := apiGroup.Group("admin")
	adminPublicGroup := adminGroup.Group("public")
	{
		adminRouterGroup.Public.InitAuthApiRouter(adminPublicGroup)
	}
	privateGroup := adminGroup.Group("")
	privateGroup.Use(middleware.AdminAuth)
	{
		adminRouterGroup.Admin.InitAdminApiRouter(privateGroup)
		adminRouterGroup.Admin.InitAuthApiRouter(privateGroup)
		adminRouterGroup.Redbook.InitCookieApiRouter(privateGroup)
	}
}
