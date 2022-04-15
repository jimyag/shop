package router

import (
	"github.com/gin-gonic/gin"

	"github.com/jimyag/shop/app/order/api/middlewares"
)

func orderRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("")
	baseRouter.Use(middlewares.Tracing())
	publicRouter := baseRouter.Group("order")
	publicRouter.Use()
	{

	}

	privateRouter := baseRouter.Group("order")
	privateRouter.Use(middlewares.Paseto())
	{

	}
}
