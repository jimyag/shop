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
		GoodsId: 1,
		Num:     9000,
	}
	inventory, err := inventoryClient.SetInv(context.Background(), &in)
	require.NoError(t, err)
	require.NotNil(t, inventory)
}

func TestInvDetail(t *testing.T) {
	in := proto.GoodInvInfo{
		GoodsId: 1,
		Num:     10,
	}
	inventory, err := inventoryClient.InvDetail(context.Background(), &in)
	require.NoError(t, err)
	require.NotNil(t, inventory)
	require.Equal(t, in.Num, inventory.Num)
	require.Equal(t, in.GoodsId, inventory.GoodsId)
}

func TestSell(t *testing.T) {
	in := proto.GoodInvInfo{
		GoodsId: 5,
		Num:     1,
	}
	inventory, err := inventoryClient.SetInv(context.Background(), &in)
	require.NoError(t, err)
	require.NotNil(t, inventory)

	ins := proto.SellInfo{
		GoodsInfo: []*proto.GoodInvInfo{
			{
				GoodsId: 1,
				Num:     10,
			},
			{
				GoodsId: 4,
				Num:     100,
			},
		},
	}
	inventory, err = inventoryClient.Sell(context.Background(), &ins)
	require.Error(t, err)
	require.Nil(t, inventory)

}

func TestInventoryServer_Sell(t *testing.T) {

	n := 100
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(wg *sync.WaitGroup) {
			ins := proto.SellInfo{
				GoodsInfo: []*proto.GoodInvInfo{
					{
						GoodsId: 1,
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
		},
	}
	inventory, err := inventoryClient.Rollback(context.Background(), &in)
	require.NoError(t, err)
	require.NotNil(t, inventory)

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
	inventory, err = inventoryClient.Rollback(context.Background(), &in)
	require.Error(t, err)
	require.Nil(t, inventory)
}
