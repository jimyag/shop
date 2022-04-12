package initialize

import (
	"log"

	"go.uber.org/zap"

	"github.com/jimyag/shop/app/user/api/global"
)

//
// InitLogger
//  @Description: 初始化日志
//
func InitLogger() {
	var err error
	global.Logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("初始化 logger 失败 :%s \n", err.Error())
	}
	global.Logger.Info("初始化 logger 成功......")
}
