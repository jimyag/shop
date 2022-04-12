package router

import "github.com/gin-gonic/gin"

func UserRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("")
	_ = baseRouter.Group("")

	_ = baseRouter.Group("")
}
