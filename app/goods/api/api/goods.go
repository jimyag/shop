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
//  @param ctx
//
func UpdateGoodsInfo(ctx *gin.Context) {
	payload, err := paseto.GetPayloadFormCtx(ctx)
	if err != nil {
		model.FailWithMsg("权限不足", ctx)
		return
	}

	if payload.Role == 1 {
		model.FailWithMsg("权限不足", ctx)
		return
	}

	arg := request.UpdateGoods{}
	_ = ctx.ShouldBindJSON(&arg)
	msg, err := validate.Validate(&arg, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	in := proto.GoodsInfo{
		Id:    arg.ID,
		Name:  arg.Name,
		Price: arg.Price,
	}
	goodsInfo, err := global.GoodsSrvClient.UpdateGoods(ctx, &in)
	if err != nil {
		global.Logger.Error("更新商品信息失败", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	model.OkWithData(goodsInfo, ctx)
}

//
// GetGoods
//  @Description: 获得商品信息
//  @param ctx
//
func GetGoods(ctx *gin.Context) {
	arg := request.GoodsIDRequest{}
	_ = ctx.ShouldBindJSON(&arg)
	msg, err := validate.Validate(&arg, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	goodsInfo, err := global.GoodsSrvClient.GetGoods(ctx, &proto.GoodID{Id: arg.ID})
	if err != nil {
		global.Logger.Error("获得商品信息失败", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	model.OkWithData(goodsInfo, ctx)
}

//
// DeleteGoods
//  @Description: 删除商品
//  @param ctx
//
func DeleteGoods(ctx *gin.Context) {
	payload, err := paseto.GetPayloadFormCtx(ctx)
	if err != nil {
		model.FailWithMsg("权限不足", ctx)
		return
	}

	if payload.Role == 1 {
		model.FailWithMsg("权限不足", ctx)
		return
	}

	goodsID := request.GoodsIDRequest{}
	_ = ctx.ShouldBindJSON(&goodsID)

	msg, err := validate.Validate(goodsID, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	_, err = global.GoodsSrvClient.DeleteGoods(ctx, &proto.GoodsInfo{Id: goodsID.ID})
	if err != nil {
		global.Logger.Error("删除商品信息失败", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	model.OkWithMsg("成功删除商品", ctx)
}

//
// GetBatchGoods
//  @Description: 批量获得商品信息
//  @param ctx
//
func GetBatchGoods(ctx *gin.Context) {
	arg := request.GoodsBatchRequest{}
	_ = ctx.ShouldBindJSON(&arg)
	msg, err := validate.Validate(&arg, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	in := proto.ManyGoodsID{
		GoodsIDs: make([]*proto.GoodID, 0),
	}
	for _, idRequest := range arg.GoodsBatchID {
		goodId := proto.GoodID{Id: idRequest.ID}
		in.GoodsIDs = append(in.GoodsIDs, &goodId)
	}
	goodsInfos, err := global.GoodsSrvClient.GetGoodsBatchInfo(ctx, &in)
	if err != nil {
		global.Logger.Error("批量获得商品信息失败", zap.Error(err))
		handle_grpc_error.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	model.OkWithData(goodsInfos, ctx)
}
