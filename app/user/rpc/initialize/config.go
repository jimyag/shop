package initialize

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/jimyag/shop/app/user/rpc/global"
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
