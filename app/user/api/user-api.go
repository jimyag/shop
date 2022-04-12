package main

import "github.com/jimyag/shop/app/user/api/initialize"

func main() {
	// 初始化logger
	initialize.InitLogger()

	// 初始化配置中心
	initialize.InitConfigCenter()

	// 从配置中心拉取配置
	initialize.InitRemoteConfig()
}
