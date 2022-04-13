package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/jimyag/shop/app/goods/rpc/global"
	handler "github.com/jimyag/shop/app/goods/rpc/handle"
	"github.com/jimyag/shop/app/goods/rpc/initialize"
	"github.com/jimyag/shop/app/goods/rpc/model"
	"github.com/jimyag/shop/common/proto"
	"github.com/jimyag/shop/common/utils/otgrpc"
	port2 "github.com/jimyag/shop/common/utils/port"
	"github.com/jimyag/shop/common/utils/register/consul"
	uuid2 "github.com/jimyag/shop/common/utils/uuid"
)

func main() {
	// 初始化logger
	initialize.InitLogger()

	// 初始化配置中心的配置信息
	initialize.InitConfigCenter()

	// 拉取远程配置中心中的配置
	initialize.InitRemoteConfig()

	// 初始化 database
	initialize.InitDataBase()

	// 初始化jaeger
	tracer, cl, err := initialize.InitJaeger()
	if err != nil {
		global.Logger.Fatal("初始化jaeger失败", zap.Error(err))
	}
	opentracing.SetGlobalTracer(tracer)

	// grpc 的server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(tracer),
		),
	)

	// 数据库的连接
	sqlStore := model.NewSQLStore(global.DB)

	goodsServer := handler.NewGoodsServer(sqlStore)
	proto.RegisterGoodsServer(grpcServer, goodsServer)

	// 优先使用配置的端口
	listener, err := net.Listen(
		"tcp",
		fmt.Sprintf("%s:%d",
			global.RemoteConfig.ServiceInfo.Host,
			global.RemoteConfig.ServiceInfo.Port,
		),
	)
	// 如果配置的端口不能使用
	if err != nil {
		// 拿到可用的端口
		port, err := port2.GetFreePort()
		if err == nil {
			global.RemoteConfig.ServiceInfo.Port = port
			listener, err = net.Listen(
				"tcp",
				fmt.Sprintf("%s:%d",
					global.RemoteConfig.ServiceInfo.Host,
					global.RemoteConfig.ServiceInfo.Port,
				),
			)
			if err != nil {
				global.Logger.Fatal(
					"监听端口失败",
					zap.Error(err),
					zap.String("host", global.RemoteConfig.ServiceInfo.Host),
					zap.Int("port", global.RemoteConfig.ServiceInfo.Port),
				)
			}
		} else {
			global.Logger.Fatal(
				"获得可用端口失败",
				zap.Error(err),
			)
		}

	}
	global.Logger.Info("监听端口",
		zap.String("host", global.RemoteConfig.ServiceInfo.Host),
		zap.Int("port", global.RemoteConfig.ServiceInfo.Port),
	)

	// 注册 grpc 健康检查
	grpc_health_v1.RegisterHealthServer(grpcServer, health.NewServer())

	registerClient := consul.NewRegistryRpcClient(
		global.ConfigCenter.Host,
		global.ConfigCenter.Port,
	)
	// 随机生成服务的id
	serviceID := uuid2.GetUUid()

	// 注册服务
	err = registerClient.Register(
		global.RemoteConfig.ServiceInfo.Host,
		global.RemoteConfig.ServiceInfo.Port,
		global.RemoteConfig.ServiceInfo.Name,
		nil,
		serviceID.String(),
	)
	if err != nil {
		global.Logger.Error("注册服务健康检查失败", zap.Error(err))
	}

	go func() {
		err = grpcServer.Serve(listener)
		if err != nil {
			global.Logger.Fatal("运行服务失败", zap.Error(err))
		}
	}()

	// 监听终止事件
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 注销服务
	if err = registerClient.DeRegister(serviceID.String()); err != nil {
		global.Logger.Info("服务注销失败", zap.String("serviceID", serviceID.String()))
	}
	cl.Close()

	global.Logger.Info("服务已注销", zap.String("serviceID", serviceID.String()))
}
