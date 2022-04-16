package router

import (
	"github.com/gin-gonic/gin"

	"github.com/jimyag/shop/app/order/api/api"
	"github.com/jimyag/shop/app/order/api/middlewares"
)

func OrderRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("")
	baseRouter.Use(middlewares.Tracing())
	publicRouter := baseRouter.Group("order")
	publicRouter.Use()

	privateRouter := baseRouter.Group("order")
	privateRouter.Use(middlewares.Paseto())
	{
		privateRouter.POST("create", api.CreateOrder)    // 创建订单
		privateRouter.GET("info", api.GetOrderDetail)    // 获得订单详情
		privateRouter.GET("infos", api.GetOrderList)     // 获得个人订单列表
		privateRouter.PUT("update", api.UpdateOrderInfo) // 更新订单
	}
}
