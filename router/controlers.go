package router

import "apiTools/controlers"

// 初始化路由
func initControlers() {
	Router.GET("/", controlers.ApiIndex)
	Router.GET("/docs/:apiName", controlers.ApiDocs)
	Router.GET("/about", controlers.ApiAbout)
}
