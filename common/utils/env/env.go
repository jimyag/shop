package env

import "github.com/spf13/viper"

func GetEnvBool(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
}
