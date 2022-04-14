package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jimyag/shop/common/proto"
)

func TestSetInv(t *testing.T) {

	in := proto.CreateOrderRequest{}
	order, err := orderClient.CreateOrder(context.Background(), &in)
	// todo
	require.NoError(t, err)
	require.NotNil(t, order)
}
