package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/jimyag/shop/app/order/rpc/global"
	"github.com/jimyag/shop/common/proto"
	"github.com/jimyag/shop/common/utils/otgrpc"
)

//
// InitGrpcClient
//  @Description: 初始化所有的grpc client
//
func InitGrpcClient() {
	initGoodsClient()
	initInventoryClient()
}

//
// initGoodsClient
//  @Description: 初始化 user client
//
func initGoodsClient() {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s://%s:%d/%s?wait=14s",
			global.ConfigCenter.Type,
			global.ConfigCenter.Host,
			global.ConfigCenter.Port,
			global.RemoteConfig.ThirdServer.GoodsGrpcServer.Name,
		),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(
				opentracing.GlobalTracer(),
			),
		),
	)

	if err != nil {
		global.Logger.Fatal("goods服务发现错误", zap.Error(err))
	}
	global.GoodsClient = proto.NewGoodsClient(conn)
	global.Logger.Info("发现goods服务......")
}

//
// initInventoryClient
//  @Description: 初始化 user client
//
func initInventoryClient() {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s://%s:%d/%s?wait=14s",
			global.ConfigCenter.Type,
			global.ConfigCenter.Host,
			global.ConfigCenter.Port,
			global.RemoteConfig.ThirdServer.InventoryGrpcServer.Name,
		),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(
				opentracing.GlobalTracer(),
			),
		),
	)

	if err != nil {
		global.Logger.Fatal("Inventory服务发现错误", zap.Error(err))
	}
	global.InventoryClient = proto.NewInventoryClient(conn)
	global.Logger.Info("发现Inventory服务......")
}
