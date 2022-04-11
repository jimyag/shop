package global

import (
	"database/sql"

	"go.uber.org/zap"

	remoteConfig "github.com/jimyag/shop/app/user/rpc/config"
	"github.com/jimyag/shop/common/model"
)

var (
	Logger       *zap.Logger             // logger
	RemoteConfig *remoteConfig.ALLConfig //远程配置中心里面的配置
	ConfigCenter *model.ConfigCenterInfo //配置中心的位置信息
	DB           *sql.DB                 // database
)
