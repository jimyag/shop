package handle_grpc_error

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"

	"github.com/jimyag/shop/common/model"
)

func HandleGrpcErrorToHttp(err error, ctx *gin.Context) {
	if err == nil {
		return
	}
	if e, ok := status.FromError(err); ok {
		model.FailWithMsg(e.Message(), ctx)
	}
}
