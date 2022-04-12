package initialize

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/jimyag/shop/app/user/rpc/global"
)

//
// InitDataBase
//  @Description: 初始化数据库的连接
//
func InitDataBase() {
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
		global.Logger.Fatal("连接数据库失败 :", zap.Error(err))
	}
	global.Logger.Info("连接数据库成功......")
}
