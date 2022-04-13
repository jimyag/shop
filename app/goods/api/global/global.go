package global

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/jimyag/shop/app/goods/api/config"
	"github.com/jimyag/shop/common/model"
	"github.com/jimyag/shop/common/proto"
	"github.com/jimyag/shop/common/utils/paseto"
)

var (
	Logger         *zap.Logger             // logger
	ConfigCenter   *model.ConfigCenterInfo // 配置中心的配置
	RemoteConfig   *config.ServerInfo      // 配置中心中的配置
	GoodsSrvClient proto.GoodsClient       // goods grpc 服务的客户端
	Trans          ut.Translator           // 公共的翻译
	Validate       *validator.Validate     // 公共的validate
	PasetoMaker    *paseto.PasetoMaker     // paseto 的maker
)
