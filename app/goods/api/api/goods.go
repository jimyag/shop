package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/jimyag/shop/app/goods/api/global"
	"github.com/jimyag/shop/app/goods/api/model/request"
	"github.com/jimyag/shop/common/model"
	"github.com/jimyag/shop/common/proto"
	"github.com/jimyag/shop/common/utils/handle_grpc_error"
	"github.com/jimyag/shop/common/utils/paseto"
	"github.com/jimyag/shop/common/utils/validate"
)

//
// CreateGoods
//  @Description: 创建商品
//  todo   只能是管理员创建 ,目前是谁都可以创建
//  @param ctx
//
func CreateGoods(ctx *gin.Context) {
	payload, err := paseto.GetPayloadFormCtx(ctx)
	if err != nil {
		model.FailWithMsg("权限不足", ctx)
		return
	}

	if payload.Role == 1 {
		model.FailWithMsg("权限不足", ctx)
		return
	}

	createGoodsRequest := request.CreateGoods{}
	_ = ctx.ShouldBindJSON(&createGoodsRequest)
	msg, err := validate.Validate(createGoodsRequest, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}
	in := proto.CreateGoodRequest{
		Name:  createGoodsRequest.Name,
		Price: createGoodsRequest.Price,
	}
	goodsInfo, err := global.GoodsSrvClient.CreateGoods(ctx, &in)
	if err != nil {
		global.Logger.Error("创建用户失败", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	model.OkWithData(goodsInfo, ctx)
}

//
// UpdateGoodsInfo
//  @Description: 更新商品的信息
// todo 只能是管理员更新
//  @param ctx
//
func UpdateGoodsInfo(ctx *gin.Context) {
	// todo
}

//
// GetGoods
//  @Description: 获得商品信息
//  @param ctx
//
func GetGoods(ctx *gin.Context) {
	// todo
}

//
// DeleteGoods
//  @Description: 删除商品
// todo 只能是管理员删除
//  @param ctx
//
func DeleteGoods(ctx *gin.Context) {
	//	todo
}

//
// GetBatchGoods
//  @Description: 批量获得商品信息
//  @param ctx
//
func GetBatchGoods(ctx *gin.Context) {
	// todo
}
