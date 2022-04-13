package middlewares

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/jimyag/shop/app/user/api/global"
	"github.com/jimyag/shop/common/model"
	"github.com/jimyag/shop/common/utils/paseto"
)

func Paseto() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenHeader := context.Request.Header.Get("Authorization")
		if tokenHeader == "" {
			model.FailWithMsg("token 无效", context)
			context.Abort()
			return
		}
		check := strings.SplitN(tokenHeader, " ", 2)
		if len(check) != 2 && check[0] != "Bearer" {
			model.FailWithMsg("token 格式错误", context)
			context.Abort()
			return
		}
		payload, err := global.PasetoMaker.VerifyToken(check[1])
		if errors.Is(err, paseto.ErrInvalidToken) {
			model.FailWithMsg("token 格式错误", context)
			context.Abort()
			return
		} else if errors.Is(err, paseto.ErrExpiredToken) {
			model.FailWithMsg("token 过期", context)
			context.Abort()
			return
		}
		context.Set("payload", payload)
	}
}
