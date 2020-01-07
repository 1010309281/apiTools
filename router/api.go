package router

import (
	"apiTools/controlers"
)

// 初始化api路由
func initApiRouter() {
	apiGroup := Router.Group("/api")
	// whois router
	{
		apiGroup.GET("/whoisquery", controlers.WhoisQuery)
		apiGroup.POST("/whoisquery", controlers.WhoisQuery)
	}
	// short router
	{
		// 长链接转换为短链接
		apiGroup.GET("/toshorturl", controlers.ShortToShortUrl)
		apiGroup.POST("/toshorturl", controlers.ShortToShortUrl)
		// 短链接解析回长链接
		apiGroup.GET("/parseshorturl", controlers.ShortParseShortUrl)
		apiGroup.POST("/parseshorturl", controlers.ShortParseShortUrl)
	}
}
