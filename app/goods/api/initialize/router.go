package initialize

import (
	"github.com/gin-gonic/gin"

	router2 "github.com/jimyag/shop/app/goods/api/router"
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

	// goods 的路由
	goodsRouter := router.Group("/goods/v1")
	router2.GoodsRouter(goodsRouter)
	return router
}
