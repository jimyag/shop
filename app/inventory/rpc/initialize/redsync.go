package initialize

import (
	"fmt"

	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"

	"github.com/jimyag/shop/app/inventory/rpc/global"
)

func InitRedSync() {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d",
			global.RemoteConfig.RedSync.Host,
			global.RemoteConfig.RedSync.Port,
		),
	})
	pool := goredis.NewPool(client)
	global.RedSync = redsync.New(pool)
	if global.RedSync == nil {
		global.Logger.Fatal("初始化分布式锁失败")
	}
	global.Logger.Info("初始化分布式锁成功......")
}
