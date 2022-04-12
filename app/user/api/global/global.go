package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/jimyag/shop/app/user/api/config"
	"github.com/jimyag/shop/app/user/api/tools/paseto"
	"github.com/jimyag/shop/common/model"
	"github.com/jimyag/shop/common/proto"
)

var (
	Logger        *zap.Logger             // logger
	ConfigCenter  *model.ConfigCenterInfo // 配置中心的配置
	RemoteConfig  *config.ServerInfo      // 配置中心中的配置
	UserSrvClient proto.UserClient        // user grpc 服务的客户端
	Redis         *AllRedis               // 所有的redis的配置
	Trans         ut.Translator           // 公共的翻译
	Validate      *validator.Validate     // 公共的validate
	PasetoMaker   *paseto.PasetoMaker     // paseto 的maker
)

type AllRedis struct {
	CreateUser *redis.Client
}
