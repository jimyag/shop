package handler

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jimyag/shop/app/order/rpc/global"
	"github.com/jimyag/shop/app/order/rpc/model"
	"github.com/jimyag/shop/common/proto"
)

//
// OrderServer
//  @Description: order 的server
//
type OrderServer struct {
	Store model.Store
}

//
// NewOrderServer
//  @Description: 创建 order server
//  @param store
//  @return *OrderServer
//
func NewOrderServer(store model.Store) *OrderServer {
	return &OrderServer{Store: store}
}

//
// CartItemList
//  @Description: 获取用户的购物车列表
//  @receiver server
//  @param ctx
//  @param req
//  @return *proto.CartItemListResponse
//  @return error
//
func (server *OrderServer) CartItemList(ctx context.Context, req *proto.CartItemListRequest) (*proto.CartItemListResponse, error) {
	cartList, err := server.Store.GetCartListByUid(ctx, req.GetUid())
	// 没数据不用关心
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		global.Logger.Error(err.Error())
		return &proto.CartItemListResponse{}, status.Error(codes.Internal, "未知错误")
	}
	response := proto.CartItemListResponse{}
	response.Total = int32(len(cartList))
	cartResponse := make([]*proto.ShopCartInfoResponse, 0)
	for _, v := range cartList {
		infoResponse := proto.ShopCartInfoResponse{
			Id:      int32(v.ID),
			UserID:  v.UserID,
			GoodsID: v.GoodsID,
			Nums:    v.Nums,
			Checked: v.Checked,
		}
		cartResponse = append(cartResponse, &infoResponse)
	}
	response.Data = cartResponse
	return &response, nil
}

//
// CreateCartItem
//  @Description: 将商品添加到购物车， 有两种 原本没有这件商品，这个商品添加了合并
//  @receiver server
//  @param ctx
//  @param req
//  @return *proto.ShopCartInfoResponse
//  @return error
//
func (server *OrderServer) CreateCartItem(ctx context.Context, req *proto.CreateCartItemRequest) (*proto.ShopCartInfoResponse, error) {
	// 查询是否有了
	getCartDetailByUIDAndGoodsIDParams := model.GetCartDetailByUIDAndGoodsIDParams{
		GoodsID: req.GoodsID,
		UserID:  req.UserID,
	}
	shoppingCart, err := server.Store.GetCartDetailByUIDAndGoodsID(ctx, getCartDetailByUIDAndGoodsIDParams)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			global.Logger.Error("查询购物车失败", zap.Error(err))
			return &proto.ShopCartInfoResponse{}, status.Error(codes.Internal, "内部错误")
		}
		// 没有的话就新建
		cartInfo, err := server.Store.CreateCart(ctx, model.CreateCartParams{UserID: req.UserID, GoodsID: req.GoodsID, Nums: req.Nums})
		if err != nil {
			global.Logger.Error("创建购物车记录失败", zap.Error(err))
			return &proto.ShopCartInfoResponse{}, status.Error(codes.Internal, "未知错误")
		}
		resp := proto.ShopCartInfoResponse{
			Id:      int32(cartInfo.ID),
			UserID:  cartInfo.UserID,
			GoodsID: cartInfo.GoodsID,
			Nums:    cartInfo.Nums,
		}
		return &resp, nil
	}

	// 如果已经有了就把数量加上去
	cartInfo, err := server.Store.UpdateCartItem(ctx, model.UpdateCartItemParams{
		UpdatedAt: time.Now(),
		Nums:      shoppingCart.Nums + req.Nums,
		Checked:   req.Checked,
		UserID:    req.UserID,
		GoodsID:   req.GoodsID,
	})
	if err != nil {
		global.Logger.Error(err.Error())
		return &proto.ShopCartInfoResponse{}, status.Error(codes.Internal, "未知错误")
	}
	resp := proto.ShopCartInfoResponse{
		Id:      int32(cartInfo.ID),
		UserID:  cartInfo.UserID,
		GoodsID: cartInfo.GoodsID,
		Nums:    cartInfo.Nums,
		Checked: cartInfo.Checked,
	}
	return &resp, nil
}

// DeleteCartItems 删除购物车中的某条信息
func (server *OrderServer) DeleteCartItems(ctx context.Context, req *proto.DeleteCartItemsRequest) (*proto.Empty, error) {
	arg := model.DeleteCartItemParams{
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UserID:    req.UserID,
		GoodsID:   req.GoodsID,
	}
	_, err := server.Store.DeleteCartItem(ctx, arg)
	if err == sql.ErrNoRows {
		return &proto.Empty{}, status.Error(codes.NotFound, "没有该条记录")
	} else if err != nil {
		global.Logger.Error(err.Error())
		return &proto.Empty{}, status.Error(codes.Internal, "未知错误")
	}
	return &proto.Empty{}, nil
}

func (server *OrderServer) UpdateCartItem(ctx context.Context, req *proto.UpdateCartItemRequest) (*proto.Empty, error) {
	cartDetail := model.GetCartDetailByUIDAndGoodsIDParams{
		UserID:  req.UserID,
		GoodsID: req.GoodsID,
	}
	// 拿到之前的条目信息
	oldCart, err := server.Store.GetCartDetailByUIDAndGoodsID(ctx, cartDetail)
	if err == sql.ErrNoRows {
		return &proto.Empty{}, status.Error(codes.NotFound, "没有找到该条目")
	} else if err != nil {
		global.Logger.Error(err.Error())
		return &proto.Empty{}, status.Error(codes.Internal, "内部错误")
	}

	updateArg := model.UpdateCartItemParams{
		UpdatedAt: time.Now(),
	}
	// 如果之前的和现在的不同，就使用现在的
	if req.Checked != oldCart.Checked {
		updateArg.Checked = req.Checked
	}

	if req.Nums != oldCart.Nums {
		updateArg.Nums = req.Nums
	}

	_, err = server.Store.UpdateCartItem(ctx, updateArg)
	if err != nil {
		global.Logger.Error(err.Error())
		return &proto.Empty{}, status.Error(codes.Internal, "内部错误")
	}
	return &proto.Empty{}, nil
}
func (server *OrderServer) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.OrderInfo, error) {
	// 从购物车中拿到选中的商品
	getCheckedCart := model.GetCartListCheckedParams{
		UserID:  req.UserID,
		Checked: true,
	}
	_, err := server.Store.GetCartListChecked(ctx, getCheckedCart)
	if errors.Is(err, sql.ErrNoRows) {
		return &proto.OrderInfo{}, status.Error(codes.InvalidArgument, "没有选中的商品")
	} else if err != nil {
		return &proto.OrderInfo{}, status.Error(codes.Internal, "内部错误")
	}

	// todo
	// 批量获得goods的信息
	return nil, status.Errorf(codes.Unimplemented, "method CreateOrder not implemented")
}

// GetOrderList 从第一页开始获得
func (server *OrderServer) GetOrderList(ctx context.Context, req *proto.GetOrderListRequest) (*proto.GetOrderListResponse, error) {
	arg := model.GetOrderListParams{}
	arg.UserID = req.UserID
	arg.Limit = req.PageSize
	arg.Offset = (req.PageNum - 1) * req.PageSize
	orderList, err := server.Store.GetOrderList(ctx, arg)
	if err == sql.ErrNoRows {
		return &proto.GetOrderListResponse{}, status.Error(codes.NotFound, "没有订单")
	} else if err != nil {
		global.Logger.Error(err.Error())
		return &proto.GetOrderListResponse{}, status.Error(codes.Internal, "未知错误")
	}
	response := proto.GetOrderListResponse{}
	response.Total = int32(len(orderList))
	responseDatas := make([]*proto.OrderInfo, 0)
	for _, v := range orderList {
		orderInfo := proto.OrderInfo{
			Id:      int32(v.ID),
			UserID:  v.UserID,
			OrderID: v.OrderID,
			PayType: v.PayType.String,
			Status:  int32(v.Status),
			Total:   float32(v.OrderMount.Float64),
			Post:    v.Post,
			Address: v.Post,
			Name:    v.SignerName,
			Mobile:  v.SignerMobile,
		}
		responseDatas = append(responseDatas, &orderInfo)
	}
	response.Data = responseDatas

	return &response, nil
}
func (server *OrderServer) GetOrderDetail(ctx context.Context, req *proto.GetOrderDetailRequest) (*proto.OrderInfo, error) {
	orderInfo, err := server.Store.GetOrderDetail(ctx, model.GetOrderDetailParams{UserID: req.UserID, OrderID: req.OrderID})
	if err == sql.ErrNoRows {
		return &proto.OrderInfo{}, status.Error(codes.NotFound, "没有找到该订单")
	} else if err != nil {
		global.Logger.Error(err.Error())
		return &proto.OrderInfo{}, status.Error(codes.Internal, "内部错误")
	}
	response := proto.OrderInfo{
		Id:      int32(orderInfo.ID),
		UserID:  orderInfo.UserID,
		OrderID: orderInfo.OrderID,
		PayType: orderInfo.PayType.String,
		Status:  int32(orderInfo.Status),
		Total:   float32(orderInfo.OrderMount.Float64),
		Post:    orderInfo.Post,
		Address: orderInfo.Post,
		Name:    orderInfo.SignerName,
		Mobile:  orderInfo.SignerMobile,
	}
	return &response, nil
}
func (server *OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderInfo) (*proto.OrderInfo, error) {
	_, err := server.Store.GetOrderDetail(ctx, model.GetOrderDetailParams{UserID: req.UserID, OrderID: req.OrderID})
	if err == sql.ErrNoRows {
		return nil, status.Error(codes.NotFound, "没有订单")
	} else if err != nil {
		global.Logger.Error(err.Error())
		return nil, status.Error(codes.Internal, "未知错误")
	}
	arg := model.UpdateOrderParams{
		UpdatedAt: time.Now(),
	}
	// 更新支付
	if req.PayType != "" {
		arg.PayType = sql.NullString{String: req.PayType, Valid: true}
		arg.PayTime = sql.NullTime{Time: time.Now(), Valid: true}
		arg.Status = int64(req.Status)
		arg.OrderID = req.OrderID
		orderInfo, err := server.Store.UpdateOrder(ctx, arg)
		if err != nil {
			global.Logger.Error(err.Error())
			return &proto.OrderInfo{}, status.Error(codes.Internal, "未知错误")
		}
		rsp := proto.OrderInfo{
			Id:      int32(orderInfo.ID),
			UserID:  orderInfo.UserID,
			OrderID: orderInfo.OrderID,
			PayType: orderInfo.PayType.String,
			Status:  int32(orderInfo.Status),
			Total:   float32(orderInfo.OrderMount.Float64),
			Post:    orderInfo.Post,
			Address: orderInfo.Post,
			Name:    orderInfo.SignerName,
			Mobile:  orderInfo.SignerMobile,
		}
		return &rsp, nil
	}
	return &proto.OrderInfo{}, status.Error(codes.InvalidArgument, "参数无效")
}
