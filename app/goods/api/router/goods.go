package router

import (
	"github.com/gin-gonic/gin"

	"github.com/jimyag/shop/app/goods/api/middlewares"
)

func GoodsRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("")
	baseRouter.Use(middlewares.Tracing())
	publicRouter := baseRouter.Group("user")
	publicRouter.Use()
	{

	}

	privateRouter := baseRouter.Group("user")
	privateRouter.Use(middlewares.Paseto())
	{

	}
}
