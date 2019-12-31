package router

import (
	"apiTools/controlers"
	"apiTools/router/middlreware"
)

// 初始化api路由
func initApiRouter() {
	apiGroup := Router.Group("/api")
	// whois router
	{
		apiGroup.GET("/whoisquery", middlreware.ApiCount("whois"), controlers.WhoisQuery)
		apiGroup.POST("/whoisquery", middlreware.ApiCount("whois"), controlers.WhoisQuery)
	}
	// short router
	{
		// 长链接转换为短链接
		apiGroup.GET("/toshorturl", middlreware.ApiCount("toShort"), controlers.ShortToShortUrl)
		apiGroup.POST("/toshorturl", middlreware.ApiCount("toShort"), controlers.ShortToShortUrl)
		// 短链接解析回长链接
		apiGroup.GET("/parseshorturl", middlreware.ApiCount("parseShort"), controlers.ShortParseShortUrl)
		apiGroup.POST("/parseshorturl", middlreware.ApiCount("parseShort"), controlers.ShortParseShortUrl)
	}
}
