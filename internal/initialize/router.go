package initialize

import (
	_ "ai-software-copyright-server/docs/admin"
	_ "ai-software-copyright-server/docs/content"
	"ai-software-copyright-server/internal/application/router/api/admin"
	"ai-software-copyright-server/internal/application/router/api/content"
	"ai-software-copyright-server/internal/application/router/api/user"
	"ai-software-copyright-server/internal/application/router/html"
	"ai-software-copyright-server/internal/global"
	"ai-software-copyright-server/internal/initialize/middleware"
	"ai-software-copyright-server/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"reflect"
	"strings"
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
	router.Static("image", utils.GetImageStorePath())
	//router.Use(ginLogger(), ginRecovery(true))
	initSwaggerRouter(router)
	initWebRouter(router)
	initApiRouter(router)
	initValidator()
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

func initWebRouter(router *gin.Engine) {
	middleware.HtmlRender(router)
	// 资源
	resourceGroup := router.Group("")
	html.InitResourceRouter(resourceGroup)
	// 页面
	htmlGroup := router.Group("")
	htmlGroup.Use(middleware.VisitStatisticHandler)
	html.InitHtmlRouter(htmlGroup)
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
		adminRouterGroup.Public.InitAiApiRouter(adminPublicGroup)
		adminRouterGroup.Public.InitAuthApiRouter(adminPublicGroup)
		adminRouterGroup.Public.InitDifyApiRouter(adminPublicGroup)
	}
	privateGroup := adminGroup.Group("")
	privateGroup.Use(middleware.AdminAuth)
	{
		adminRouterGroup.Admin.InitAdminApiRouter(privateGroup)
		adminRouterGroup.Admin.InitAuthApiRouter(privateGroup)
		adminRouterGroup.Cdkey.InitCdkeyApiRouter(privateGroup)
		adminRouterGroup.User.InitUserApiRouter(privateGroup)
	}

	contentRouterGroup := content.ApiRouterGroup
	contentGroup := apiGroup.Group("content")
	{
		contentRouterGroup.Image.InitImageApiRouter(contentGroup)
	}

	userRouterGroup := user.ApiRouterGroup
	userGroup := apiGroup.Group("user")
	userPublicGroup := userGroup.Group("public")
	{
		userRouterGroup.Public.InitAuthApiRouter(userPublicGroup)
		userRouterGroup.Public.InitCreditsApiRouter(userPublicGroup)
		userRouterGroup.Public.InitNoticeApiRouter(userPublicGroup)
		userRouterGroup.Public.InitStudyApiRouter(userPublicGroup)
		userRouterGroup.Public.InitWxNotifyApiRouter(userPublicGroup)
	}
	userPrivateGroup := userGroup.Group("")
	userPrivateGroup.Use(middleware.UserAuth)
	{
		userRouterGroup.Cdkey.InitCdkeyApiRouter(userPrivateGroup)
		userRouterGroup.Credits.InitCreditsOrderApiRouter(userPrivateGroup)

		userRouterGroup.SoftwareCopyright.InitSoftwareCopyrightApiRouter(userPrivateGroup)
		userRouterGroup.Study.InitResourceApiRouter(userPrivateGroup)
		userRouterGroup.User.InitAuthApiRouter(userPrivateGroup)
		userRouterGroup.User.InitUserApiRouter(userPrivateGroup)
	}
}

func initValidator() {
	validate, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		global.LOG.Error("验证器注册失败")
		return
	}
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		if label, ok := field.Tag.Lookup("label"); ok {
			return label
		}
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		// skip if tag key says it should be ignored
		if name == "-" {
			return ""
		}
		return name
	})
}
