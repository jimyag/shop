package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jimyag/shop/common/proto"
)

//
//  TestOrderServer_CreateCartItem
//  @Description: 测试创建购物车记录
//  @param t
//
func TestOrderServer_CreateCartItem(t *testing.T) {

	shopCartItem, err := orderClient.CreateCartItem(context.Background(), &proto.CreateCartItemRequest{
		UserID:  116,
		GoodsID: 5,
		Nums:    10,
		Checked: true,
	})
	require.NoError(t, err)
	require.NotNil(t, shopCartItem)
}

func TestGetCartList(t *testing.T) {
	list, err := orderClient.CartItemList(context.Background(), &proto.CartItemListRequest{Uid: 116})
	require.NoError(t, err)
	require.True(t, len(list.Data) > 0)
	t.Log(list)

	list, err = orderClient.CartItemList(context.Background(), &proto.CartItemListRequest{Uid: 117})
	require.Error(t, err)
	require.Nil(t, list)
}

func TestOrderServer_UpdateCartItem(t *testing.T) {
	_, err := orderClient.UpdateCartItem(context.Background(), &proto.UpdateCartItemRequest{
		UserID:  116,
		GoodsID: 5,
		Nums:    1,
		Checked: true,
	})
	require.NoError(t, err)

	_, err = orderClient.UpdateCartItem(context.Background(), &proto.UpdateCartItemRequest{
		UserID:  116,
		GoodsID: 9999999,
		Nums:    1,
		Checked: false,
	})
	require.Error(t, err)
}

func TestOrderServer_CreateOrder(t *testing.T) {
	order, err := orderClient.CreateOrder(context.Background(), &proto.CreateOrderRequest{
		UserID:  116,
		Address: "中国上海",
		Mobile:  "18522222222",
		Name:    "jimyag",
		Post:    "201314",
	})
	require.NoError(t, err)
	require.NotNil(t, order)
	t.Logf("%v", order)
}

func TestOrderServer_GetOrderDetail(t *testing.T) {
	rsp, err := orderClient.GetOrderDetail(context.Background(),
		&proto.GetOrderDetailRequest{
			OrderID: 20224151781711695,
		})
	require.NoError(t, err)
	require.NotEmpty(t, rsp)
	t.Logf("%v", rsp)
}

func TestOrderServer_GetOrderList(t *testing.T) {
	rsp, err := orderClient.GetOrderList(
		context.Background(),
		&proto.GetOrderListRequest{
			UserID:   116,
			PageSize: 10,
			PageNum:  1})
	require.NoError(t, err)
	require.True(t, len(rsp.Data) > 0)
	t.Logf("%v", rsp)
}

//
//  TestOrderServer_UpdateOrderStatus
//  @Description: 测试更新订单状态 // 1 待支付 2 成功 3 超时关闭
//  @param t
//
func TestOrderServer_UpdateOrderStatus(t *testing.T) {
	order, err := orderClient.UpdateOrderStatus(context.Background(), &proto.OrderInfo{
		UserID:  116,
		OrderID: 20224151781711695,
		PayType: "已超时",
		Status:  3,
	})
	require.NoError(t, err)
	t.Logf("%v", order)
}
