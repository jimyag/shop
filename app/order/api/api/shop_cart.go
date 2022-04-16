package api

import (
	"github.com/gin-gonic/gin"

	"github.com/jimyag/shop/app/order/api/global"
	"github.com/jimyag/shop/app/order/api/model/request"
	"github.com/jimyag/shop/common/model"
	"github.com/jimyag/shop/common/proto"
	"github.com/jimyag/shop/common/utils/validate"
)

//
// CreateShopCart
//  @Description: 添加商品到购物车中
//  @param ctx*gin.Context
//
func CreateShopCart(ctx *gin.Context) {
	cartItemRequest := request.CreateCartItemRequest{}
	_ = ctx.ShouldBind(&cartItemRequest)
	msg, err := validate.Validate(cartItemRequest, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	// 判断商品是否存在
	_, err = global.GoodsSrvClient.GetGoods(ctx, &proto.GoodID{Id: cartItemRequest.GoodsId})
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}
	arg := &proto.CreateCartItemRequest{
		UserID:  cartItemRequest.UserId,
		GoodsID: cartItemRequest.GoodsId,
		Nums:    cartItemRequest.Nums,
	}
	rsp, err := global.OrderSrvClient.CreateCartItem(ctx, arg)
	if err != nil {
		model.FailWithMsg(err.Error(), ctx)
		return
	}
	model.OkWithData(rsp, ctx)
}

//
// GetShopCartList
//  @Description: 获取购物车列表
//
func GetShopCartList(ctx *gin.Context) {
	cartListRequest := request.GetCartListRequest{}
	_ = ctx.ShouldBindJSON(&cartListRequest)

	msg, err := validate.Validate(cartListRequest, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	arg := proto.CartItemListRequest{Uid: cartListRequest.UserId}
	rsp, err := global.OrderSrvClient.CartItemList(ctx, &arg)
	if err != nil {
		model.FailWithMsg(err.Error(), ctx)
		return
	}
	model.OkWithData(rsp.Data, ctx)
}

//
// DeleteShopCartItem
//  @Description: 删除购物车中的商品
//  @param ctx
//
func DeleteShopCartItem(ctx *gin.Context) {
	deleteShopCartRequest := request.DeleteShopCartRequest{}
	_ = ctx.ShouldBindJSON(&deleteShopCartRequest)
	msg, err := validate.Validate(deleteShopCartRequest, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}
	_, err = global.OrderSrvClient.DeleteCartItems(ctx, &proto.DeleteCartItemsRequest{
		UserID:  deleteShopCartRequest.UserId,
		GoodsID: deleteShopCartRequest.GoodsId,
	})
	if err != nil {
		model.FailWithMsg(err.Error(), ctx)
		return
	}
	model.Ok(ctx)
}

//
// UpdateShopCartItem
//  @Description: 更新购物车中的商品信息
//  @param ctx
//
func UpdateShopCartItem(ctx *gin.Context) {

	updateShopCartRequest := request.UpdateShopCartRequest{}
	_ = ctx.ShouldBindJSON(&updateShopCartRequest)
	msg, err := validate.Validate(updateShopCartRequest, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	_, err = global.OrderSrvClient.UpdateCartItem(ctx, &proto.UpdateCartItemRequest{
		UserID:  updateShopCartRequest.UserId,
		GoodsID: updateShopCartRequest.GoodsId,
		Nums:    updateShopCartRequest.Nums,
		Checked: updateShopCartRequest.Checker,
	})
	if err != nil {
		model.FailWithMsg(err.Error(), ctx)
		return
	}
	model.Ok(ctx)

}
