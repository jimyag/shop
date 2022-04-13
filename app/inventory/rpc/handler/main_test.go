package handler

import (
	"log"
	"os"
	"testing"

	"google.golang.org/grpc"

	"github.com/jimyag/shop/common/proto"
)

var (
	inventoryClient proto.InventoryClient
)

const (
	target = "192.168.0.2:58560"
)

func TestMain(m *testing.M) {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot dial %s :%v\n", target, err)
	}
	inventoryClient = proto.NewInventoryClient(conn)

	log.Printf("dial %s success....\n", target)
	os.Exit(m.Run())
}
