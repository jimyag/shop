package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jimyag/shop/app/inventory/rpc/model"
	"github.com/jimyag/shop/common/proto"
)

type InventoryServer struct {
	model.Store
}

func NewInventoryServer(store model.Store) *InventoryServer {
	return &InventoryServer{store}
}

// SetInv 如果没有就添加，如果有就更新
func (i *InventoryServer) SetInv(ctx context.Context, req *proto.GoodInvInfo) (*proto.Empty, error) {
	// 查找库存

	// 数量增加

	// 创建

	// 或者更新
	arg := model.CreateInventoryParams{
		Goods:   req.GetGoodsId(),
		Sticks:  req.Num,
		Version: 1,
	}
	_, err := i.SetInvTx(ctx, arg)
	return &proto.Empty{}, err
}

// InvDetail 获得详情
func (i *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodInvInfo) (*proto.GoodInvInfo, error) {
	// 查询
	res, err := i.GetInventoryByGoodsID(ctx, req.GoodsId)
	if err != nil || res.ID == 0 {
		return nil, status.Error(codes.NotFound, "货物不存在")
	}

	rsp := proto.GoodInvInfo{
		GoodsId: res.Goods,
		Num:     res.Sticks,
	}
	// 拿到就返回
	return &rsp, nil
}

// Sell 扣减库存
func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*proto.Empty, error) {
	// 本地事务  要不都卖，要不都不卖
	// 拿到所有的商品，
	// 		判断是否有库存
	// 		判断库存是否够
	// 		扣减库存 - 库存 会出现数据不一致的问题
	// todo 分布式锁 数据不一致问题
	err := i.SellTx(ctx, req)
	return &proto.Empty{}, err
}
func (i *InventoryServer) Rollback(ctx context.Context, req *proto.SellInfo) (*proto.Empty, error) {
	// 订单超时归还
	// 订单创建失败
	// 收到取消归还
	// 批量归还
	err := i.RollBackTx(ctx, req)
	return &proto.Empty{}, err
}
