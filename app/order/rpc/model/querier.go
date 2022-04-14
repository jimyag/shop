// Code generated by sqlc. DO NOT EDIT.

package model

import (
	"context"
)

type Querier interface {
	CreateCart(ctx context.Context, arg CreateCartParams) (ShoppingCart, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (OrderInfo, error)
	CreateOrderGoods(ctx context.Context, arg CreateOrderGoodsParams) (OrderGood, error)
	DeleteCartItem(ctx context.Context, arg DeleteCartItemParams) (ShoppingCart, error)
	GetCartDetailByUIDAndGoodsID(ctx context.Context, arg GetCartDetailByUIDAndGoodsIDParams) (ShoppingCart, error)
	GetCartListByUid(ctx context.Context, userID int32) ([]ShoppingCart, error)
	GetCartListChecked(ctx context.Context, arg GetCartListCheckedParams) ([]ShoppingCart, error)
	GetOrderDetail(ctx context.Context, arg GetOrderDetailParams) (OrderInfo, error)
	GetOrderList(ctx context.Context, arg GetOrderListParams) ([]OrderInfo, error)
	GetOrderListByOrderID(ctx context.Context, orderID int32) ([]OrderGood, error)
	UpdateCartItem(ctx context.Context, arg UpdateCartItemParams) (ShoppingCart, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (OrderInfo, error)
}

var _ Querier = (*Queries)(nil)
