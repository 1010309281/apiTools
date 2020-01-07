package router

import (
	"apiTools/controlers"
	"apiTools/libs/config"
	"apiTools/router/middlreware"
	"apiTools/utils"
	"github.com/gin-gonic/gin"
	"path/filepath"
)

var (
	Router *gin.Engine
)

// 初始化gin
func InitRouter() {
	if config.GetString("web::appMode") == "production" {
		gin.SetMode(gin.ReleaseMode)
		gin.DisableConsoleColor()
	}
	//Router = gin.Default()
	Router = gin.New()

	// 设置静态文件
	Router.LoadHTMLGlob(filepath.Join(utils.GetRootPath(), "views", "/*"))
	Router.Static("/static", filepath.Join(utils.GetRootPath(), "static"))

	// 设置全局中间件
	Router.Use(gin.Recovery())
	Router.Use(middlreware.AllowCors())
	Router.Use(middlreware.Logger())
	Router.Use(middlreware.ProApiDocs())

	// 404错误处理
	Router.NoRoute(controlers.NoRouter)
	Router.NoMethod(controlers.NoRouter)

	// 加载路由
	initApiRouter()
	initControlers()
}
