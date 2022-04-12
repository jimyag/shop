package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// HealthRouter
//  @Description: health 的路由，包含一个 根/health GET 的请求
//  @param router
//
func HealthRouter(router *gin.RouterGroup) {
	publicRouter := router.Group("")
	publicRouter.GET("/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, "Agent alive and reachable")
	})

}
