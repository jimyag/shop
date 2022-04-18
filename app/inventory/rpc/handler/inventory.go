package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
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

		sellDetail := model.StockSellDetail{
			OrderID: req.OrderId,
			Status:  1, // 默认表示已经扣减了
		}
		details := make([]model.GoodsDetail, 0)
		for _, info := range req.GoodsInfo {
			details = append(details, model.GoodsDetail{
				GoodsID: info.GoodsId,
				Nums:    info.Num,
			})
			// 分布式锁
			mutex := global.RedSync.NewMutex(fmt.Sprintf("goods_%d", info.GoodsId))
			if err := mutex.Lock(); err != nil {
				global.Logger.Error("", zap.Error(err))
				return status.Error(codes.Internal, "内部错误-获取分布式锁失败")
			}

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

			if ok, err := mutex.Unlock(); !ok || err != nil {
				return status.Error(codes.Internal, "内部错误-释放分布式锁失败")
			}
		}
		sellDetail.Detail = details
		_, err := queries.CreateSellDetail(ctx, model.CreateSellDetailParams{
			OrderID: sellDetail.OrderID,
			Status:  sellDetail.Status,
			Detail:  sellDetail.Detail,
		})
		if err != nil {
			return err
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

func AutoRollBack(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderID int64 `json:"order_id"`
	}
	for _, msg := range msgs {
		// 既然要归还库存，就应该直到每件商品应该归还多少， 这时候出现 重复归还的问题
		// 这个接口应该保证幂等性，不能因为消息的重复发送而导致一个订单的库存归还多次，没有扣减的库存不能归还。
		// 新建一张表，记录了详细的订单扣减细节，以及归还的情况
		var orderInfo OrderInfo
		err := json.Unmarshal(msg.Body, &orderInfo)
		if err != nil {
			global.Logger.Error("JSON 解析失败", zap.Error(err))
			// 根据业务来，如果赶紧时自己代码问题就用
			//return consumer.ConsumeRetryLater,nil
			// 否则就直接忽略这个消息
			return consumer.ConsumeSuccess, nil
		}
		// 将inv的库存加回去，同时将sell status 变为2
		// todo
		_, err = global.DB.Begin()
		if err != nil {
			global.Logger.Error("获得事务失败", zap.Error(err))
			return consumer.ConsumeRetryLater, nil
		}

		//	将状态变为2
	}
	return consumer.ConsumeSuccess, nil
}
