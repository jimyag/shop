package handler

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jimyag/shop/common/proto"
)

func TestSetInv(t *testing.T) {
	in := proto.GoodInvInfo{
		GoodsId: 2,
		Num:     9000,
	}
	inventory, err := inventoryClient.SetInv(context.Background(), &in)
	require.NoError(t, err)
	require.NotNil(t, inventory)
}

func TestInvDetail(t *testing.T) {
	in := proto.GoodInvInfo{
		GoodsId: 1,
	}
	inventory, err := inventoryClient.InvDetail(context.Background(), &in)
	require.NoError(t, err)
	require.NotNil(t, inventory)
	require.Equal(t, in.GoodsId, inventory.GoodsId)
}

func TestSell(t *testing.T) {

	inNoErr := proto.SellInfo{
		GoodsInfo: []*proto.GoodInvInfo{
			{
				GoodsId: 1,
				Num:     10,
			},
			{
				GoodsId: 2,
				Num:     100,
			},
		},
	}
	_, err := inventoryClient.Sell(context.Background(), &inNoErr)
	require.NoError(t, err)

	// 第一件商品扣减成功了第二件商品扣减失败 回滚
	inErr := proto.SellInfo{
		GoodsInfo: []*proto.GoodInvInfo{
			{
				GoodsId: 1,
				Num:     10,
			},
			{
				GoodsId: 99999,
				Num:     100,
			},
		},
	}
	_, err = inventoryClient.Sell(context.Background(), &inErr)
	require.Error(t, err)

	// 第一件商品扣减成功了第二件商品数量不足 回滚
	inErr = proto.SellInfo{
		GoodsInfo: []*proto.GoodInvInfo{
			{
				GoodsId: 1,
				Num:     10,
			},
			{
				GoodsId: 2,
				Num:     9999999,
			},
		},
	}

	_, err = inventoryClient.Sell(context.Background(), &inErr)
	require.Error(t, err)

}

func TestInventoryServer_Sell(t *testing.T) {
	t.Parallel()
	n := 90
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(wg *sync.WaitGroup) {
			ins := proto.SellInfo{
				GoodsInfo: []*proto.GoodInvInfo{
					{
						GoodsId: 2,
						Num:     1,
					},
				},
			}
			inventory, err := inventoryClient.Sell(context.Background(), &ins)
			require.NoError(t, err)
			require.NotNil(t, inventory)
			wg.Done()
		}(&wg)

	}
	wg.Wait()

}

func TestRollBack(t *testing.T) {
	in := proto.SellInfo{
		GoodsInfo: []*proto.GoodInvInfo{
			{
				GoodsId: 1,
				Num:     10,
			},
			{
				GoodsId: 2,
				Num:     100,
			},
		},
	}
	_, err := inventoryClient.Rollback(context.Background(), &in)
	require.NoError(t, err)

	in = proto.SellInfo{
		GoodsInfo: []*proto.GoodInvInfo{
			{
				GoodsId: 1,
				Num:     10,
			},
			{
				GoodsId: 10000,
				Num:     1,
			},
		},
	}
	_, err = inventoryClient.Rollback(context.Background(), &in)
	require.Error(t, err)
}

func TestAddGoodsInventory(t *testing.T) {
	in := proto.GoodInvInfo{
		GoodsId: 2,
		Num:     9000,
	}
	inventory, err := inventoryClient.SetInv(context.Background(), &in)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	in = proto.GoodInvInfo{
		GoodsId: 3,
		Num:     9000,
	}
	inventory, err = inventoryClient.SetInv(context.Background(), &in)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	in = proto.GoodInvInfo{
		GoodsId: 4,
		Num:     9000,
	}
	inventory, err = inventoryClient.SetInv(context.Background(), &in)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	in = proto.GoodInvInfo{
		GoodsId: 5,
		Num:     9000,
	}
	inventory, err = inventoryClient.SetInv(context.Background(), &in)
	require.NoError(t, err)
	require.NotNil(t, inventory)
}
