package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/jimyag/shop/app/user/api/global"
	"github.com/jimyag/shop/app/user/api/initialize"
	"github.com/jimyag/shop/common/utils/register/consul"
	uuid2 "github.com/jimyag/shop/common/utils/uuid"
)

func main() {
	// 初始化logger
	initialize.InitLogger()

	// 初始化配置中心
	initialize.InitConfigCenter()

	// 从配置中心拉取配置
	initialize.InitRemoteConfig()

	// 初始化 grpc 的 client
	initialize.InitGrpcClient()

	// 初始化router
	router := initialize.InitRouter()

	// 初始化 redis 的配置
	initialize.InitRedis()

	registerClient := consul.NewRegistryHttpClient(
		global.ConfigCenter.Host,
		global.ConfigCenter.Port,
	)
	serviceID := uuid2.GetUUid()

	err := registerClient.Register(
		global.RemoteConfig.ServiceInfo.Host,
		global.RemoteConfig.ServiceInfo.Port,
		global.RemoteConfig.ServiceInfo.Name,
		nil,
		serviceID.String(),
	)
	if err != nil {
		global.Logger.Error("注册用户http服务失败", zap.Error(err))
	}

	go func() {
		err = router.Run(
			fmt.Sprintf("%s:%d",
				global.RemoteConfig.ServiceInfo.Host,
				global.RemoteConfig.ServiceInfo.Port,
			),
		)
		if err != nil {
			global.Logger.Error("用户http服务启动失败", zap.Error(err))
		}
	}()
	global.Logger.Info(
		"用户http服务启动成功",
		zap.String("host", global.RemoteConfig.ServiceInfo.Host),
		zap.Int("port", global.RemoteConfig.ServiceInfo.Port),
		zap.String("name", global.RemoteConfig.ServiceInfo.Name),
	)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = registerClient.DeRegister(serviceID.String()); err != nil {
		global.Logger.Info("注销失败", zap.String("serviceID", serviceID.String()), zap.Error(err))
	}
	global.Logger.Info("服务已注销", zap.String("serviceID", serviceID.String()))
}
