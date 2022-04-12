package env

import "github.com/spf13/viper"

//
// GetEnvBool
//  @Description: 获得指定环境变量中的bool值
//  @param env
//  @return bool
//
func GetEnvBool(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
