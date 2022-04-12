package initialize

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/jimyag/shop/app/user/api/global"
	"github.com/jimyag/shop/common/proto"
)

//
// InitGrpcClient
//  @Description: 初始化所有的grpc client
//
func InitGrpcClient() {
	initUserClient()
}

//
// initUserClient
//  @Description: 初始化 user client
//
func initUserClient() {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s://%s:%d/%s?wait=14s",
			global.ConfigCenter.Type,
			global.ConfigCenter.Host,
			global.ConfigCenter.Port,
			global.RemoteConfig.UserGrpcServer.Name,
		),
		grpc.WithInsecure(),
	)

	if err != nil {
		global.Logger.Fatal("用户服务发现错误", zap.Error(err))
	}
	global.UserSrvClient = proto.NewUserClient(conn)

	global.Logger.Info("发现用户服务......")
}
