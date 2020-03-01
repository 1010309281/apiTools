package cmd

import (
	"apiTools/apps/crontab"
	"apiTools/apps/proxyPool"
	"apiTools/libs/config"
	"apiTools/libs/logger"
	"apiTools/modles"
	"apiTools/routers"
	"fmt"
)

func InitServer() (err error) {
	// 初始化配置文件
	err = config.InitConfig()
	if err != nil {
		return
	}

	// 初始化Logger
	err = logger.InitLogger()
	if err != nil {
		return
	}

	// 初始化redis连接池
	err = modles.InitRedis()
	if err != nil {
		return
	}

	// 初始化mysql连接
	err = modles.InitMysql()
	if err != nil {
		return
	}

	// 初始化api配置
	err = initApiConfig()
	if err != nil {
		return
	}

	// 初始化路由引擎
	routers.InitRouter()

	//初始化apps
	err = runApps()
	if err != nil {
		return
	}

	// 启动服务
	err = routers.Router.Run(fmt.Sprintf(":%s", config.GetString("web::port")))
	if err != nil {
		return
	}
	return
}

// 初始化api配置
func initApiConfig() (err error) {
	// 初始化api docs json数据
	err = modles.InitApiDocsJsonData()
	if err != nil {
		return
	}

	// 初始化whois服务器列表
	err = modles.InitWhoisServers()
	if err != nil {
		return
	}

	// 初始化ipv4db数据库信息
	err = modles.InitIp4DB()
	if err != nil {
		return
	}
	return
}

// 运行app程序
func runApps() (err error) {
	//启动proxy app
	err = proxyPool.RunProxyPoolApp()
	if err != nil {
		logger.Echo.Errorf("run proxy pool app fail, Error Msg: %v", err)
		err = nil
	}
	// 启动定时任务app
	err = crontab.RunCrontabApp()
	if err != nil {
		logger.Echo.Errorf("run crontab apps fail, Error Msg: %v", err)
		err = nil
	}
	return
}
