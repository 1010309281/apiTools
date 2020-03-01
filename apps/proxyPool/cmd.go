// 代理池主入口文件
package proxyPool

func RunProxyPoolApp() (err error) {
	// 初始化数据管道
	checkProxyJobChan = make(chan *ProxyInfo, 1024)
	checkProxyResultChan = make(chan *ProxyInfo, 1024)

	// 启动运行爬虫抓取代理
	RunSpiders()

	// 启动校验
	runCheckProxy(16)

	// 启动入库
	runSaveToDB(4)

	return
}
