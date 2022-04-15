package router

import (
	"github.com/gin-gonic/gin"

	"github.com/jimyag/shop/app/order/api/api"
	"github.com/jimyag/shop/app/order/api/middlewares"
)

func ShopCartRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("")
	baseRouter.Use(middlewares.Tracing())
	publicRouter := baseRouter.Group("cart")
	publicRouter.Use()

	privateRouter := baseRouter.Group("cart")
	privateRouter.Use(middlewares.Paseto())
	{
		privateRouter.POST("create", api.CreateShopCart)       // 添加商品到购物车记录
		privateRouter.GET("list", api.GetShopCartList)         // 获得购物车列表
		privateRouter.DELETE("remove", api.DeleteShopCartItem) // 删除购物车的某条记录
		privateRouter.PUT("update", api.UpdateShopCartItem)    // 更新的某条记录
	}
}
