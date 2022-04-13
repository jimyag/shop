package router

import (
	"github.com/gin-gonic/gin"

	"github.com/jimyag/shop/app/user/api/api"
	"github.com/jimyag/shop/app/user/api/middlewares"
)

func UserRouter(router *gin.RouterGroup) {
	baseRouter := router.Group("")
	baseRouter.Use(middlewares.Tracing())
	publicRouter := baseRouter.Group("user")
	{
		// 使用邮箱和密码登录
		publicRouter.POST("login", api.PasswordLogin)
		// 注册用户之前发送的验证码
		publicRouter.POST("email/register", api.CreateUserEmail)
		// 注册用户
		publicRouter.POST("register", api.CreateUser)
	}

	privateRouter := baseRouter.Group("user")
	privateRouter.Use(middlewares.Paseto())
	{
		// 通过使用uid获得用户的信息
		privateRouter.GET("info", api.GetUserByID)
		// 通过使用email获得用户的信息
		privateRouter.GET("info_email", api.GetUserByEmail)
		// 获得用户列表
		privateRouter.GET("list", api.GetUserList)
		// 更新用户的nickname 和gender
		privateRouter.PUT("info", api.UpdateUserWithOutPassword)
		// 更新用户的密码
		privateRouter.PUT("password", api.ChangePassword)
		// 更新用户的权限
		privateRouter.PUT("role", api.ChangeRole)
	}
}
