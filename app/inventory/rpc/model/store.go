package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jimyag/shop/app/inventory/rpc/global"
	"github.com/jimyag/shop/common/proto"
)

type Store interface {
	SetInvTx(ctx context.Context, arg CreateInventoryParams) (Inventory, error)
	SellTx(ctx context.Context, arg *proto.SellInfo) error
	RollBackTx(ctx context.Context, arg *proto.SellInfo) error
	Querier
}

type SqlStore struct {
	*Queries
	db *sql.DB
}

func NewSQLStore(db *sql.DB) Store {
	return &SqlStore{
		Queries: New(db),
		db:      db,
	}
}

func (store *SqlStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

func (store *SqlStore) SetInvTx(ctx context.Context, arg CreateInventoryParams) (Inventory, error) {
	var inventory Inventory
	var err error
	err = store.execTx(ctx, func(queries *Queries) error {
		inventory, err = queries.GetInventoryByGoodsID(ctx, arg.Goods)
		if err != nil {
			if err == sql.ErrNoRows {
				// 没有找到
				inventory, err = queries.CreateInventory(ctx, arg)
				return nil
			} else {
				global.Logger.Error("", zap.Error(err))
				return status.Error(codes.Internal, "内部错误")
			}
		}
		updateArg := UpdateInventoryParams{
			UpdatedAt: time.Now(),
			Goods:     inventory.Goods,
			Counts:    arg.Sticks,
		}
		inventory, err = queries.UpdateInventory(ctx, updateArg)
		return err
	})
	return inventory, err
}

func (store *SqlStore) SellTx(ctx context.Context, arg *proto.SellInfo) error {
	// 本地事务  要不都卖，要不都不卖
	// 拿到所有的商品，
	// 		判断是否有库存
	// 		判断库存是否够
	// 		扣减库存 - 库存 会出现数据不一致的问题
	err := store.execTx(ctx, func(queries *Queries) error {
		var inventory Inventory
		var err error
		for _, info := range arg.GetGoodsInfo() {
			inventory, err = queries.GetInventoryByGoodsID(ctx, info.GoodsId)
			if err != nil {
				if err == sql.ErrNoRows {
					return status.Error(codes.NotFound, "没有该货物")
				} else {
					return status.Error(codes.Internal, "内部错误")
				}
			}
			if inventory.Sticks < info.Num {
				return status.Error(codes.InvalidArgument, "货物不够")
			}
			updateArg := UpdateInventoryParams{}
			updateArg.Goods = info.GoodsId
			updateArg.Counts = -info.Num // 这边应该时负数
			updateArg.UpdatedAt = time.Now()
			inventory, err = queries.UpdateInventory(ctx, updateArg)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (store *SqlStore) RollBackTx(ctx context.Context, arg *proto.SellInfo) error {
	err := store.execTx(ctx, func(queries *Queries) error {
		var err error
		for _, info := range arg.GetGoodsInfo() {
			_, err = queries.GetInventoryByGoodsID(ctx, info.GoodsId)
			if err != nil {
				if err == sql.ErrNoRows {
					return status.Error(codes.NotFound, "没有该货物")
				} else {
					return status.Error(codes.Internal, "内部错误")
				}
			}
			updateArg := UpdateInventoryParams{}
			updateArg.Goods = info.GoodsId
			updateArg.Counts = info.Num // 这边应该时正数
			updateArg.UpdatedAt = time.Now()
			_, err = queries.UpdateInventory(ctx, updateArg)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
