package initialize

import (
	"github.com/gin-gonic/gin"

	healthRouter "github.com/jimyag/shop/common/utils/router"
)

//
// InitRouter
//  @Description: 初始化router
//  @return *gin.Engine
//
func InitRouter() *gin.Engine {
	router := gin.Default()
	rootGroup := router.Group("")

	// 初始健康检查的路由
	healthRouter.HealthRouter(rootGroup)

	return router
}
