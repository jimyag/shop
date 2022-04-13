package main

import "github.com/jimyag/shop/common/utils/initialize"

const (
	consulAddress = "http://localhost:8500"
	debugPath     = "shop/goods/rpc/debug.yaml"
	releasePath   = "shop/goods/rpc/release.yaml"
	fileName      = "remote-config.yaml"
)

//
// main
//  @Description: 用来初始化远程配置文件
//
func main() {
	initialize.LocalConfigToRemoteConfig(consulAddress, debugPath, releasePath, fileName)
}
