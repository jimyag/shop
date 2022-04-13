package router

import (
	"github.com/gin-gonic/gin"

	"github.com/jimyag/shop/app/goods/api/api"
	"github.com/jimyag/shop/app/goods/api/middlewares"
)

func GoodsRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("")
	baseRouter.Use(middlewares.Tracing())
	publicRouter := baseRouter.Group("goods")
	publicRouter.Use()
	{
		publicRouter.GET("info", api.GetGoods)
		publicRouter.GET("infos", api.GetBatchGoods)
	}

	privateRouter := baseRouter.Group("goods")
	privateRouter.Use(middlewares.Paseto())
	{
		privateRouter.POST("create", api.CreateGoods)
		privateRouter.PUT("info", api.UpdateGoodsInfo)
		privateRouter.DELETE("info", api.DeleteGoods)
	}
}
