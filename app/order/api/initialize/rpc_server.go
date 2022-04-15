package initialize

import (
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/jimyag/shop/app/order/api/global"
	"github.com/jimyag/shop/common/proto"
	"github.com/jimyag/shop/common/utils/otgrpc"
)

//
// InitGrpcClient
//  @Description: 初始化所有的grpc client
//
func InitGrpcClient() {
	initOrderClient()
	initGoodsClient()
}

//
// initUserClient
//  @Description: 初始化 order client
//
func initOrderClient() {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s://%s:%d/%s?wait=14s",
			global.ConfigCenter.Type,
			global.ConfigCenter.Host,
			global.ConfigCenter.Port,
			global.RemoteConfig.OrderGrpcServer.Name,
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
		global.Logger.Fatal("订单服务发现错误", zap.Error(err))
	}
	global.OrderSrvClient = proto.NewOrderClient(conn)

	global.Logger.Info("发现订单服务......")
}

//
// initUserClient
//  @Description: 初始化 goods client
//
func initGoodsClient() {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s://%s:%d/%s?wait=14s",
			global.ConfigCenter.Type,
			global.ConfigCenter.Host,
			global.ConfigCenter.Port,
			global.RemoteConfig.GoodsGrpcServer.Name,
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
		global.Logger.Fatal("商品服务发现错误", zap.Error(err))
	}
	global.GoodsSrvClient = proto.NewGoodsClient(conn)

	global.Logger.Info("发现商品服务......")
}
