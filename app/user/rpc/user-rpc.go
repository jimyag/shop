package main

import "github.com/jimyag/shop/app/user/rpc/initialize"

func main() {
	// 初始化 logger
	initialize.InitLogger()

	// 初始化配置中心的配置信息
	initialize.InitConfigCenter()
}
