package initialize

import (
	"log"

	"go.uber.org/zap"

	"github.com/jimyag/shop/app/inventory/rpc/global"
)

func InitLogger() {
	var err error
	global.Logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("初始化 logger 失败 :%v\n", err)
	}

	global.Logger.Info("初始化 logger 成功.....")
}
