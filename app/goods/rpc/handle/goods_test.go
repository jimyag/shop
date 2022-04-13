package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jimyag/shop/common/proto"
	"github.com/jimyag/shop/common/utils/test_util"
)

func createGoods(t *testing.T) *proto.GoodsInfo {
	in := proto.CreateGoodRequest{
		Name:  test_util.RandomString(20),
		Price: test_util.RandomPrice(),
	}
	goods, err := goodsClient.CreateGoods(context.Background(), &in)
	require.NoError(t, err)
	require.NotNil(t, goods)
	require.Equal(t, in.Price, goods.Price)
	require.Equal(t, in.Name, goods.Name)
	return goods
}

func TestGoodsServer_CreateGoods(t *testing.T) {

	goods := createGoods(t)
	in := proto.CreateGoodRequest{
		Name:  goods.Name,
		Price: goods.Price,
	}
	_, err := goodsClient.CreateGoods(context.Background(), &in)
	require.Error(t, err)
}

func TestGoodsServer_GetGoods(t *testing.T) {
	goods := createGoods(t)
	getGoods, err := goodsClient.GetGoods(context.Background(), &proto.GoodID{Id: goods.GetId()})
	require.NoError(t, err)
	require.Equal(t, goods, getGoods)
}

func TestGoodsServer_UpdateGoods(t *testing.T) {
	goods := createGoods(t)
	arg := proto.GoodsInfo{
		Id:    goods.Id,
		Name:  goods.Name,
		Price: test_util.RandomPrice(),
	}
	getGoods, err := goodsClient.UpdateGoods(context.Background(), &arg)
	require.NoError(t, err)
	require.Equal(t, arg.Id, getGoods.Id)
	require.Equal(t, arg.Name, getGoods.Name)
	require.Equal(t, arg.Price, getGoods.Price)
	arg.Id += 1
	getGoods, err = goodsClient.UpdateGoods(context.Background(), &arg)
	require.Error(t, err)
}

func TestGoodsServer_DeleteGoods(t *testing.T) {

}
