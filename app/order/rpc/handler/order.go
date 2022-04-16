package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jimyag/shop/app/order/rpc/global"
	"github.com/jimyag/shop/app/order/rpc/model"
	"github.com/jimyag/shop/app/order/rpc/tools/generate"
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
	if cartList == nil {
		return &proto.CartItemListResponse{}, status.Error(codes.NotFound, "没有记录")
	} else if err != nil {
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

//
// DeleteCartItems
//  @Description: 删除购物车中的某条信息
//  @receiver server
//  @param ctx
//  @param req
//  @return *proto.Empty
//  @return error
//
func (server *OrderServer) DeleteCartItems(ctx context.Context, req *proto.DeleteCartItemsRequest) (*proto.Empty, error) {
	// 查询是否存在
	getCartDetailByUIDAndGoodsIDParams := model.GetCartDetailByUIDAndGoodsIDParams{
		GoodsID: req.GoodsID,
		UserID:  req.UserID,
	}
	_, err := server.Store.GetCartDetailByUIDAndGoodsID(ctx, getCartDetailByUIDAndGoodsIDParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.Empty{}, status.Error(codes.NotFound, "此条记录不存在")
		}

		global.Logger.Error("获得购物车记录失败", zap.Error(err))
		return &proto.Empty{}, status.Error(codes.Internal, "未知错误")

	}
	arg := model.DeleteCartItemParams{
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UserID:    req.UserID,
		GoodsID:   req.GoodsID,
	}
	_, err = server.Store.DeleteCartItem(ctx, arg)
	if err == sql.ErrNoRows {
		return &proto.Empty{}, status.Error(codes.NotFound, "没有该条记录")
	} else if err != nil {
		global.Logger.Error(err.Error())
		return &proto.Empty{}, status.Error(codes.Internal, "未知错误")
	}
	return &proto.Empty{}, nil
}

//
// UpdateCartItem
//  @Description: 更新购物车记录，更新数量和选中状态
//  @receiver server
//  @param ctx
//  @param req
//  @return *proto.Empty
//  @return error
//
func (server *OrderServer) UpdateCartItem(ctx context.Context, req *proto.UpdateCartItemRequest) (*proto.Empty, error) {
	cartDetail := model.GetCartDetailByUIDAndGoodsIDParams{
		UserID:  req.UserID,
		GoodsID: req.GoodsID,
	}
	// 拿到之前的条目信息
	_, err := server.Store.GetCartDetailByUIDAndGoodsID(ctx, cartDetail)
	if err == sql.ErrNoRows {
		return &proto.Empty{}, status.Error(codes.NotFound, "没有找到该条目")
	} else if err != nil {
		global.Logger.Error(err.Error())
		return &proto.Empty{}, status.Error(codes.Internal, "内部错误")
	}

	updateArg := model.UpdateCartItemParams{
		UserID:    req.UserID,
		GoodsID:   req.GoodsID,
		UpdatedAt: time.Now(),
	}

	// 状态传过来是什么就是什么
	updateArg.Checked = req.Checked
	// 只有当传过来的值大于0的时候才能更新
	if req.Nums > 0 {
		updateArg.Nums = req.Nums
	}

	_, err = server.Store.UpdateCartItem(ctx, updateArg)
	if err != nil {
		global.Logger.Error(err.Error())
		return &proto.Empty{}, status.Error(codes.Internal, "内部错误")
	}
	return &proto.Empty{}, nil
}

type OrderListener struct {
	Code        codes.Code
	Detail      string
	OrderID     int64
	OrderAmount float32
	server      *OrderServer
	ctx         context.Context
}

func NewOrderListener(server *OrderServer, ctx context.Context) *OrderListener {
	return &OrderListener{
		server: server,
		ctx:    ctx,
	}
}
func (dl *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	// 4. 从购物车中拿到选中的商品
	// 1. 商品的金额自己查询 商品服务
	// 2. 库存的扣减 库存服务
	// 3. 订单的基本信息表
	//
	// 5. 从购物车中删除已购买的记录
	// 从购物车中拿到选中的商品
	createOrderParams := model.CreateOrderParams{}
	err := json.Unmarshal(msg.Body, &createOrderParams)
	if err != nil {
		global.Logger.Error("解析消息失败", zap.Error(err))
		return primitive.RollbackMessageState
	}

	getCheckedCart := model.GetCartListCheckedParams{
		UserID:  createOrderParams.UserID,
		Checked: true,
	}
	goodsIDS := make([]*proto.GoodID, 0)
	shoppingCart, err := dl.server.Store.GetCartListChecked(dl.ctx, getCheckedCart)
	if shoppingCart == nil {
		dl.Code = codes.InvalidArgument
		dl.Detail = "购物车为空"
		return primitive.RollbackMessageState
	} else if err != nil {
		dl.Code = codes.Internal
		dl.Detail = "获取购物车失败"
		return primitive.RollbackMessageState
	}

	// 保存 商品的数量
	goodsNumMap := make(map[int32]int32)
	for _, cart := range shoppingCart {
		goodsIDS = append(goodsIDS, &proto.GoodID{Id: cart.GoodsID})
		goodsNumMap[cart.GoodsID] = cart.Nums
	}
	goodsInfos, err := global.GoodsClient.GetGoodsBatchInfo(dl.ctx, &proto.ManyGoodsID{GoodsIDs: goodsIDS})
	if err != nil {
		dl.Code = codes.Internal
		dl.Detail = "获取商品信息失败"
		return primitive.RollbackMessageState
	}

	// 订单的总金额
	var orderAmount float32
	// 订单中商品的参数
	createOrderGoodsParams := make([]*model.CreateOrderGoodsParams, 0)
	// 扣减库存 的参数
	sellInfo := proto.SellInfo{GoodsInfo: make([]*proto.GoodInvInfo, 0)}
	for _, datum := range goodsInfos.Data {
		// 求总金额
		orderAmount += datum.Price * float32(goodsNumMap[datum.Id])
		// 订单中的参数
		createOrderGoodsParams = append(createOrderGoodsParams, &model.CreateOrderGoodsParams{
			GoodsID:    datum.Id,
			GoodsName:  datum.Name,
			GoodsPrice: float64(datum.Price),
			Nums:       goodsNumMap[datum.Id],
		})
		// 扣减库存的参数
		sellInfo.GoodsInfo = append(sellInfo.GoodsInfo, &proto.GoodInvInfo{
			GoodsId: datum.Id,
			Num:     goodsNumMap[datum.Id],
		})
	}

	// 跨服务调用 扣减库存

	_, err = global.InventoryClient.Sell(dl.ctx, &sellInfo)
	if err != nil {
		// todo
		// 如果是因为网络问题，这种如何避免
		// sell 的返回逻辑 返回的状态码是否sell返回的状态码 如果是才进行rollback
		dl.Code = codes.ResourceExhausted
		dl.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
	}

	// 本地服务的事务
	err = dl.server.Store.ExecTx(dl.ctx, func(queries *model.Queries) error {
		createOrderParams.OrderMount = sql.NullFloat64{
			Float64: float64(orderAmount),
			Valid:   true,
		}
		// 保存order
		_, err = dl.server.Store.CreateOrder(dl.ctx, createOrderParams)
		if err != nil {
			dl.Code = codes.Internal
			dl.Detail = "保存订单失败"
			return err
		}
		dl.OrderAmount = orderAmount

		// 将订单id更新
		for _, good := range createOrderGoodsParams {
			good.OrderID = createOrderParams.OrderID
		}
		// 批量插入订单中的商品
		err = dl.server.Store.ExecTx(dl.ctx, func(queries *model.Queries) error {
			for _, good := range createOrderGoodsParams {
				_, err = queries.CreateOrderGoods(dl.ctx, *good)
				if err != nil {
					dl.Code = codes.Internal
					dl.Detail = "保存订单商品失败"
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}

		// 批量删除购物车中记录
		err = dl.server.Store.ExecTx(dl.ctx, func(queries *model.Queries) error {
			for _, cart := range shoppingCart {
				_, err = queries.DeleteCartItem(dl.ctx, model.DeleteCartItemParams{
					DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
					UserID:    cart.UserID,
					GoodsID:   cart.GoodsID,
				})
				if err != nil {
					dl.Code = codes.Internal
					dl.Detail = "删除购物车中商品失败"
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	// 如果有错就要把库存归还
	if err != nil {
		return primitive.CommitMessageState
	}
	return primitive.RollbackMessageState
}

func (dl *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	createOrderParams := model.CreateOrderParams{}
	err := json.Unmarshal(msg.Body, &createOrderParams)
	if err != nil {
		global.Logger.Error("解析消息失败", zap.Error(err))
		return primitive.RollbackMessageState
	}

	_, err = dl.server.GetOrderDetail(dl.ctx, &proto.GetOrderDetailRequest{OrderID: createOrderParams.OrderID})
	if err != nil {
		// 没有扣减的库存不能被归还
		return primitive.CommitMessageState
	}
	return primitive.RollbackMessageState
}

//
// CreateOrder
//  @Description: 新建订单
//  @receiver server
//  @param ctx
//  @param req
//  @return *proto.OrderInfo
//  @return error
//
func (server *OrderServer) CreateOrder(ctx context.Context, req *proto.CreateOrderRequest) (*proto.OrderInfo, error) {
	orderlistener := NewOrderListener(server, ctx)
	p, err := rocketmq.NewTransactionProducer(
		orderlistener,
		producer.WithNameServer([]string{"192.168.0.2:9876"}),
	)

	if err != nil {
		global.Logger.Error("创建生产者失败", zap.Error(err))
		return &proto.OrderInfo{}, status.Error(codes.Internal, "创建生产者失败")
	}
	err = p.Start()
	if err != nil {
		global.Logger.Error("启动生产者失败", zap.Error(err))
		return &proto.OrderInfo{}, status.Error(codes.Internal, "启动生产者失败")
	}
	topic := "order_reback"

	// 一定要在这边生成订单号
	createOrderParams := model.CreateOrderParams{
		UserID:       req.UserID,
		OrderID:      generate.GenerateOrderID(req.UserID),
		Status:       1, // 1 待支付 2 成功 3 超时关闭
		Address:      req.Address,
		SignerName:   req.Name,
		SignerMobile: req.Mobile,
		Post:         req.Post,
	}
	jsonString, err := json.Marshal(createOrderParams)
	if err != nil {
		global.Logger.Error("序列化失败", zap.Error(err))
		return &proto.OrderInfo{}, status.Error(codes.Internal, "序列化失败")
	}
	res, err := p.SendMessageInTransaction(
		ctx,
		primitive.NewMessage(
			topic,
			jsonString,
		),
	)

	if res.State == primitive.CommitMessageState {
		return &proto.OrderInfo{}, status.Error(codes.Internal, "创建订单失败")
	}

	if err != nil {
		global.Logger.Error("发送消息失败", zap.Error(err))
		return &proto.OrderInfo{}, status.Error(codes.Internal, "发送消息失败")
	}

	return &proto.OrderInfo{
		OrderID: createOrderParams.OrderID,
		Total:   orderlistener.OrderAmount,
	}, nil
}

//
// GetOrderList
//  @Description:  从第一页开始获得某个用户的订单列表
//  @receiver server
//  @param ctx
//  @param req
//  @return *proto.GetOrderListResponse
//  @return error
//  todo 分页有问题，但是不影响由于没有后台管理可以不用管
//
func (server *OrderServer) GetOrderList(ctx context.Context, req *proto.GetOrderListRequest) (*proto.GetOrderListResponse, error) {
	arg := model.GetOrderListParams{}
	arg.UserID = req.UserID
	arg.Limit = req.PageSize
	arg.Offset = (req.PageNum - 1) * req.PageSize
	orderList, err := server.Store.GetOrderList(ctx, arg)
	if err != nil {
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

//
// GetOrderDetail
//  @Description: 获得订单详情
//  @receiver server
//  @param ctx
//  @param req
//  @return *proto.OrderInfo
//  @return error
//
func (server *OrderServer) GetOrderDetail(ctx context.Context, req *proto.GetOrderDetailRequest) (*proto.OrderDetailResponse, error) {
	// 获得订单信息
	orderInfo, err := server.Store.GetOrderDetail(ctx, req.OrderID)
	if err == sql.ErrNoRows {
		return &proto.OrderDetailResponse{}, status.Error(codes.NotFound, "没有找到该订单")
	} else if err != nil {
		global.Logger.Error(err.Error())
		return &proto.OrderDetailResponse{}, status.Error(codes.Internal, "内部错误")
	}
	OrderInfo := proto.OrderInfo{
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
	response := proto.OrderDetailResponse{
		OrderInfo: &OrderInfo,
	}
	// 获得订单中包含的商品信息
	orderGoods, err := server.Store.GetOrderListByOrderID(ctx, orderInfo.OrderID)
	if err != nil {
		global.Logger.Error(err.Error())
		return &proto.OrderDetailResponse{}, status.Error(codes.Internal, "内部错误")
	} else if orderGoods == nil {
		return &proto.OrderDetailResponse{}, status.Error(codes.NotFound, "没有找到记录")
	}
	rspOrderGoods := make([]*proto.OrderGoods, 0)
	for _, good := range orderGoods {
		OrderGoods := proto.OrderGoods{
			Id:         int32(good.ID),
			OrderID:    good.OrderID,
			GoodsID:    good.GoodsID,
			GoodsName:  good.GoodsName,
			GoodsPrice: float32(good.GoodsPrice),
			GoodsNum:   good.Nums,
		}
		rspOrderGoods = append(rspOrderGoods, &OrderGoods)
	}
	response.Goods = rspOrderGoods
	return &response, nil
}

// UpdateOrderStatus
//  @Description: 更新订单状态
//  @receiver server
//  @param ctx
//  @param req
//  @return *proto.OrderInfo
//  @return error
//
func (server *OrderServer) UpdateOrderStatus(ctx context.Context, req *proto.OrderInfo) (*proto.OrderInfo, error) {
	_, err := server.Store.GetOrderDetail(ctx, req.OrderID)
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
		arg.Status = int16(req.Status)
		arg.OrderID = req.OrderID

	} else {
		arg.Status = int16(req.Status)
		arg.OrderID = req.OrderID
	}

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
