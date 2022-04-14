package main

import (
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/jimyag/shop/app/order/rpc/global"
	"github.com/jimyag/shop/app/order/rpc/initialize"
)

func main() {
	// 初始化 logger
	initialize.InitLogger()

	// 初始化配置中心的配置信息
	initialize.InitConfigCenter()

	// 拉取远程配置中心中的配置
	initialize.InitRemoteConfig()

	// 初始化 database
	initialize.InitDataBase()

	tracer, cl, err := initialize.InitJaeger()
	if err != nil {
		global.Logger.Fatal("创建 tracer 失败", zap.Error(err))
	}
	opentracing.SetGlobalTracer(tracer)

	cl.Close()
}
