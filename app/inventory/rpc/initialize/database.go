package initialize

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/jimyag/shop/app/inventory/rpc/global"
)

func InitDatabase() {
	var err error
	dbSource := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		global.RemoteConfig.Postgres.Type,
		global.RemoteConfig.Postgres.User,
		global.RemoteConfig.Postgres.Password,
		global.RemoteConfig.Postgres.Host,
		global.RemoteConfig.Postgres.Port,
		global.RemoteConfig.Postgres.Database,
	)
	global.DB, err = sql.Open(global.RemoteConfig.Postgres.Type, dbSource)
	if err != nil {
		global.Logger.Fatal("cannot connect to db :", zap.Error(err))
	}
	global.DB.SetMaxOpenConns(0)
	//设置最大空闲超时
	global.DB.SetConnMaxIdleTime(time.Second)
	global.Logger.Info("connect db success....")
}
