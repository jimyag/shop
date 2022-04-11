package initialize

import (
	"fmt"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"go.uber.org/zap"

	"github.com/jimyag/shop/app/user/rpc/global"
	"github.com/jimyag/shop/common/utils/env"
)

// InitConfigCenter 初始化配置中心的配置
func InitConfigCenter() {
	configCenterPath := "config-center.yaml"

	v := viper.New()
	v.SetConfigFile(configCenterPath)
	if err := v.ReadInConfig(); err != nil {
		global.Logger.Fatal("加载配置文件失败", zap.Error(err))
	}
	global.Logger.Info("加载配置文件成功......")

	if err := v.Unmarshal(&global.ConfigCenter); err != nil {
		global.Logger.Fatal("解析配置文件失败", zap.Error(err))
	}
	global.Logger.Info("解析配置文件成功", zap.Any("content", global.ConfigCenter))

}

// InitRemoteConfig 拉取远程配置中心的配置
func InitRemoteConfig() {
	v := viper.New()

	path := global.ConfigCenter.ReleasePath
	if debug := env.GetEnvBool(global.ConfigCenter.EnvName); debug {
		path = global.ConfigCenter.DebugPath
	}

	v.SetConfigType(global.ConfigCenter.FileType)
	err := v.AddRemoteProvider(
		global.ConfigCenter.Type,
		fmt.Sprintf("%s:%d",
			global.ConfigCenter.Host,
			global.ConfigCenter.Port),
		path,
	)
	if err != nil {
		global.Logger.Fatal("拉取远程配置文件失败", zap.Error(err))
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		global.Logger.Fatal("读取远程配置文件失败", zap.Error(err))
	}

	err = v.Unmarshal(&global.RemoteConfig)
	if err != nil {
		global.Logger.Fatal("解析远程配置文件失败", zap.Error(err))
	}
	global.Logger.Info("成功加载远程配置文件......", zap.Any("content", global.RemoteConfig))
}
