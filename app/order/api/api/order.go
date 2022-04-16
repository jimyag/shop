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

}

//
// GetOrderDetail
//  @Description:  获取个人订单详情
//
func GetOrderDetail(ctx *gin.Context) {

}

//
// GetOrderList
//  @Description: 获取个人订单列表
//  @param ctx
//
func GetOrderList(ctx *gin.Context) {
	getOrderListRequest := request.GetOrderListRequest{}
	_ = ctx.ShouldBind(&getOrderListRequest)
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
