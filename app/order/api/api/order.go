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
// CreateOrder
//  @Description: 创建订单
//  @param ctx*gin.Context
//
func CreateOrder(ctx *gin.Context) {
	createOrderRequest := request.CreateOrderRequest{}
	_ = ctx.ShouldBindJSON(&createOrderRequest)
	msg, err := validate.Validate(createOrderRequest, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}
	orderInfo, err := global.OrderSrvClient.CreateOrder(ctx, &proto.CreateOrderRequest{
		UserID:  createOrderRequest.UserID,
		Address: createOrderRequest.Address,
		Mobile:  createOrderRequest.Mobile,
		Name:    createOrderRequest.Name,
		Post:    createOrderRequest.Post,
	})

	if err != nil {
		model.FailWithMsg(err.Error(), ctx)
		return
	}
	model.OkWithData(orderInfo, ctx)
}

//
// GetOrderDetail
//  @Description:  获取个人订单详情
//  todo 验证是否是当前用户的订单
//
func GetOrderDetail(ctx *gin.Context) {
	orderDetailRequest := &request.GetOrderDetailRequest{}
	_ = ctx.ShouldBindJSON(&orderDetailRequest)
	msg, err := validate.Validate(orderDetailRequest, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}
	rsp, err := global.OrderSrvClient.GetOrderDetail(ctx, &proto.GetOrderDetailRequest{
		OrderID: orderDetailRequest.OrderID,
	})
	if err != nil {
		model.FailWithMsg(err.Error(), ctx)
		return
	}
	model.OkWithData(rsp, ctx)

}

//
// GetOrderList
//  @Description: 获取个人订单列表
//  @param ctx
//
func GetOrderList(ctx *gin.Context) {
	getOrderListRequest := request.GetOrderListRequest{}
	_ = ctx.ShouldBindJSON(&getOrderListRequest)
	msg, err := validate.Validate(getOrderListRequest, global.Validate, global.Trans)
	if err != nil {
		model.FailWithMsg(msg, ctx)
		return
	}

	rsp, err := global.OrderSrvClient.GetOrderList(ctx, &proto.GetOrderListRequest{
		UserID:   getOrderListRequest.UserID,
		PageNum:  getOrderListRequest.PageNum,
		PageSize: getOrderListRequest.PageSize,
	})
	if err != nil {
		model.FailWithMsg(err.Error(), ctx)
		return
	}
	model.OkWithData(rsp, ctx)
}
