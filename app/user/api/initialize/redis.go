package initialize

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"golang.org/x/net/context"

	"github.com/jimyag/shop/app/user/api/global"
)

//
// InitRedis
//  @Description: 初始化所有的redis配置
//
func InitRedis() {
	allRedis := global.AllRedis{}
	allRedis.CreateUser = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%d",
			global.RemoteConfig.Redis.Host,
			global.RemoteConfig.Redis.Port,
		),
		DB: 0,
	})
	_, err := allRedis.CreateUser.Ping(context.Background()).Result()
	if err != nil {
		global.Logger.Fatal("初始化创建用户的邮箱验证码redis失败", zap.Error(err))
	}
	global.Redis = &allRedis
}
