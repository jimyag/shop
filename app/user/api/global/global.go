package global

import (
	"go.uber.org/zap"

	"github.com/jimyag/shop/app/user/api/config"
	"github.com/jimyag/shop/common/model"
	"github.com/jimyag/shop/common/proto"
)

var (
	Logger        *zap.Logger             // logger
	ConfigCenter  *model.ConfigCenterInfo // 配置中心的配置
	RemoteConfig  *config.ServerInfo      // 配置中心中的配置
	UserSrvClient proto.UserClient        // user grpc 服务的客户端
)
