package routers

import (
	"apiTools/controlers"
)

// 初始化api路由
func initApiRouter() {
	apiGroup := Router.Group("/api")
	// whois routers
	{
		apiGroup.GET("/whoisquery", controlers.WhoisQuery)
		apiGroup.POST("/whoisquery", controlers.WhoisQuery)
	}
	// short routers
	{
		// 长链接转换为短链接
		apiGroup.GET("/toshorturl", controlers.ShortToShortUrl)
		apiGroup.POST("/toshorturl", controlers.ShortToShortUrl)
		// 短链接解析回长链接
		apiGroup.GET("/parseshorturl", controlers.ShortParseShortUrl)
		apiGroup.POST("/parseshorturl", controlers.ShortParseShortUrl)
	}
	// ip query
	{
		apiGroup.GET("/ipv4query", controlers.Ipv4Query)
		apiGroup.POST("/ipv4query", controlers.Ipv4Query)
	}
	// proxy pool query
	{
		apiGroup.GET("/proxypool", controlers.ProxyPoolQuery)
		apiGroup.POST("/proxypool", controlers.ProxyPoolQuery)
	}
}
