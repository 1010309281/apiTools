package main

import (
	"apiTools/libs/config"
	"apiTools/libs/logger"
	"apiTools/modle"
	"apiTools/router"
	"fmt"
)

func main() {
	// 初始化配置文件
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}

	// 初始化Logger
	err = logger.InitLogger()
	if err != nil {
		panic(err)
	}

	// 初始化redis连接池
	err = modles.InitRedis()
	if err != nil {
		panic(err)
	}

	// 初始化路由引擎
	router.InitRouter()

	// 启动服务
	err = router.Router.Run(fmt.Sprintf(":%s", config.GetString("web::port")))
	if err != nil {
		panic(err)
	}
}
