package initialize

import (
	"github.com/gin-gonic/gin"

	router2 "github.com/jimyag/shop/app/order/api/router"
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

	// order 的路由
	orderRouter := router.Group("/order/v1")
	router2.OrderRouter(orderRouter)
	router2.ShopCartRouter(orderRouter)
	return router
}
