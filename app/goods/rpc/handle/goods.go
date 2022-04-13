package handler

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jimyag/shop/app/goods/rpc/global"
	"github.com/jimyag/shop/app/goods/rpc/model"
	"github.com/jimyag/shop/common/proto"
)

//
// GoodsServer
//  @Description: Goods 的 server
//
type GoodsServer struct {
	model.Store
}

func NewGoodsServer(store model.Store) *GoodsServer {
	return &GoodsServer{
		store,
	}
}

//
// CreateGoods
//  @Description: 创建商品
//  @receiver server
//  @param ctx
//  @param req
//  @return *proto.GoodsInfo
//  @return error
//
func (server *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodRequest) (*proto.GoodsInfo, error) {
	// 判断这个商品是否存在
	_, err := server.GetGoodsByName(ctx, req.Name)
	// 如果有那么就为nil
	if err == nil {
		return &proto.GoodsInfo{}, status.Error(codes.AlreadyExists, "商品已存在")
	}
	// 剩下有错误的，并且错误不是没有找到
	if !errors.Is(err, sql.ErrNoRows) {
		global.Logger.Error(err.Error())
		return &proto.GoodsInfo{}, status.Error(codes.Internal, "内部错误")
	}

	arg := model.CreateGoodsParams{
		Name:  req.Name,
		Price: float64(req.Price),
	}
	goods, err := server.Store.CreateGoods(ctx, arg)
	if err != nil {
		global.Logger.Error(err.Error())
		return &proto.GoodsInfo{}, status.Error(codes.Internal, "内部错误")
	}

	response := proto.GoodsInfo{
		Id:    int32(goods.ID),
		Name:  goods.Name,
		Price: float32(goods.Price),
	}
	return &response, nil
}

//
// UpdateGoods
//  @Description: 更新商品的信息
//  @receiver server
//  @param ctx
//  @param req
//  @return *proto.GoodsInfo
//  @return error
//
func (server *GoodsServer) UpdateGoods(ctx context.Context, req *proto.GoodsInfo) (*proto.GoodsInfo, error) {
	// 查找商品是否存在
	goods, err := server.GetGoodsByID(ctx, int64(req.Id))
	// 剩下有错误的，并且错误不是没有找到
	if errors.Is(err, sql.ErrNoRows) {
		return &proto.GoodsInfo{}, status.Error(codes.NotFound, "没有找到该商品")
	} else if err != nil {
		global.Logger.Error(err.Error())
		return &proto.GoodsInfo{}, status.Error(codes.Internal, "内部错误")
	}

	arg := model.UpdateGoodsParams{
		UpdatedAt: time.Now(),
		Name:      req.Name,
		Price:     float64(req.Price),
		ID:        int64(req.Id),
	}

	goods, err = server.Store.UpdateGoods(ctx, arg)
	if err != nil {
		global.Logger.Error(err.Error())
		return &proto.GoodsInfo{}, status.Error(codes.Internal, "内部错误")
	}

	response := proto.GoodsInfo{
		Id:    int32(goods.ID),
		Name:  goods.Name,
		Price: float32(goods.Price),
	}
	return &response, nil
}

//
// GetGoods
//  @Description: 使用goods id 拿到goods信息
//  @receiver server
//  @param ctx
//  @param req
//  @return *proto.GoodsInfo
//  @return error
//
func (server *GoodsServer) GetGoods(ctx context.Context, req *proto.GoodID) (*proto.GoodsInfo, error) {
	response := proto.GoodsInfo{}
	goods, err := server.Store.GetGoodsByID(ctx, int64(req.Id))
	if errors.Is(err, sql.ErrNoRows) {
		return &proto.GoodsInfo{}, status.Error(codes.NotFound, "没有找到该商品")
	}
	if err != nil {
		global.Logger.Error(err.Error())
		return &proto.GoodsInfo{}, status.Error(codes.Internal, "内部错误")
	}

	response.Id = int32(goods.ID)
	response.Price = float32(goods.Price)
	response.Name = goods.Name
	return &response, nil
}

//
// DeleteGoods
//  @Description: 删除商品信息
//  @receiver server
//  @param ctx
//  @param req
//  @return *proto.Empty
//  @return error
//
func (server *GoodsServer) DeleteGoods(ctx context.Context, req *proto.GoodsInfo) (*proto.Empty, error) {
	_, err := server.Store.GetGoodsByID(ctx, int64(req.Id))
	if errors.Is(err, sql.ErrNoRows) {
		return &proto.Empty{}, status.Error(codes.NotFound, "没有找到该商品")
	} else if err != nil {
		global.Logger.Error(err.Error())
		return &proto.Empty{}, status.Error(codes.Internal, "内部错误")
	}
	_, err = server.Store.DeleteGoods(ctx, model.DeleteGoodsParams{
		DeletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		ID:        int64(req.Id),
	})
	if err != nil {
		return &proto.Empty{}, status.Error(codes.Internal, "内部错误")
	}
	return &proto.Empty{}, nil
}

//
// GetGoodsBatchInfo
//  @Description: 批量获取商品的信息
//  @receiver server
//  @param ctx
//  @param req
//  @return *proto.ManyGoodsInfos
//  @return error
//
func (server *GoodsServer) GetGoodsBatchInfo(ctx context.Context, req *proto.ManyGoodsID) (*proto.ManyGoodsInfos, error) {
	rsp := proto.ManyGoodsInfos{
		Data: make([]*proto.GoodsInfo, 0),
	}
	var err error
	err = server.Store.ExecTx(ctx, func(queries *model.Queries) error {
		for _, d := range req.GoodsIDs {
			goods, getGoodErr := queries.GetGoodsByID(ctx, int64(d.GetId()))
			if getGoodErr != nil {
				return getGoodErr
			}
			info := proto.GoodsInfo{
				Id:    int32(goods.ID),
				Name:  goods.Name,
				Price: float32(goods.Price),
			}
			rsp.Data = append(rsp.GetData(), &info)
		}
		return nil
	})

	if err != nil {
		return &rsp, status.Error(codes.Unknown, err.Error())
	}
	return &rsp, nil
}
