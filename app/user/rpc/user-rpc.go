package main

import "github.com/jimyag/shop/app/user/rpc/initialize"

func main() {
	// 初始化 logger
	initialize.InitLogger()

	// 初始化配置中心的配置信息
	initialize.InitConfigCenter()

	// 拉取远程配置中心中的配置
	initialize.InitRemoteConfig()

	// 初始化 database
	initialize.InitDataBase()
}
