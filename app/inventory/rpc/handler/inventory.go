package handler

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jimyag/shop/app/inventory/rpc/global"
	"github.com/jimyag/shop/app/inventory/rpc/model"
	"github.com/jimyag/shop/common/proto"
)

type InventoryServer struct {
	model.Store
}

func NewInventoryServer(store model.Store) *InventoryServer {
	return &InventoryServer{store}
}

//
// SetInv
//  @Description: 如果没有就添加，如果有就更新
//  @receiver i
//  @param ctx
//  @param req
//  @return *proto.Empty
//  @return error
//  todo 在创建库存的时候应该保证商品已经存在
//
func (i *InventoryServer) SetInv(ctx context.Context, req *proto.GoodInvInfo) (*proto.Empty, error) {
	// 查找库存
	inventory, err := i.Store.GetInventoryByGoodsID(ctx, req.GoodsId)
	if err != nil {
		// 如果不是没有找到
		if !errors.Is(err, sql.ErrNoRows) {
			global.Logger.Error(err.Error())
			return &proto.Empty{}, status.Error(codes.Internal, "内部错误")
		}

		// 如果没有找到那么就要创建
		err = i.ExecTx(ctx, func(queries *model.Queries) error {
			_, err = queries.CreateInventory(ctx, model.CreateInventoryParams{
				GoodsID: req.GoodsId,
				Sticks:  req.Num,
				Version: 0,
			})
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return &proto.Empty{}, status.Error(codes.Unknown, err.Error())
		}
		return &proto.Empty{}, nil
	}

	// 更新
	_, err = i.UpdateInventory(ctx, model.UpdateInventoryParams{
		UpdatedAt: time.Now(),
		GoodsID:   inventory.GoodsID,
		Counts:    req.Num,
	})
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, "内部错误")
	}
	return &proto.Empty{}, nil
}

//
// InvDetail
//  @Description:  获得详情
//  @receiver i
//  @param ctx
//  @param req
//  @return *proto.GoodInvInfo
//  @return error
//
func (i *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodInvInfo) (*proto.GoodInvInfo, error) {
	// 查询
	res, err := i.GetInventoryByGoodsID(ctx, req.GoodsId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &proto.GoodInvInfo{}, status.Error(codes.NotFound, "货物不存在")
		}

		return &proto.GoodInvInfo{}, status.Error(codes.Internal, "内部错误")
	}

	rsp := proto.GoodInvInfo{
		GoodsId: res.GoodsID,
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

	// 拿到所有的商品，
	err := i.ExecTx(ctx, func(queries *model.Queries) error {
		for _, info := range req.GoodsInfo {
			//判断是否有库存
			inventory, err := queries.GetInventoryByGoodsID(ctx, info.GoodsId)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return status.Error(codes.NotFound, "货物不存在")
				}
				return status.Error(codes.Internal, "内部错误")
			}
			//判断库存是否够
			if inventory.Sticks < info.Num {
				return status.Error(codes.ResourceExhausted, "货物不足")
			}

			//扣减库存 - 库存 会出现数据不一致的问题
			_, err = queries.UpdateInventory(ctx, model.UpdateInventoryParams{
				UpdatedAt: time.Now(),
				GoodsID:   inventory.GoodsID,
				Counts:    -info.Num,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return &proto.Empty{}, err
	}

	return &proto.Empty{}, nil
}
func (i *InventoryServer) Rollback(ctx context.Context, req *proto.SellInfo) (*proto.Empty, error) {
	// 订单超时归还
	// 订单创建失败
	// 收到取消归还
	// 批量归还
	err := i.ExecTx(ctx, func(queries *model.Queries) error {
		for _, info := range req.GoodsInfo {
			//判断是否有库存
			inventory, err := queries.GetInventoryByGoodsID(ctx, info.GoodsId)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return status.Error(codes.NotFound, "货物不存在")
				}
				return status.Error(codes.Internal, "内部错误")
			}

			//增加库存 - 库存 会出现数据不一致的问题
			_, err = queries.UpdateInventory(ctx, model.UpdateInventoryParams{
				UpdatedAt: time.Now(),
				GoodsID:   inventory.GoodsID,
				Counts:    info.Num,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return &proto.Empty{}, err
	}

	return &proto.Empty{}, nil
}
