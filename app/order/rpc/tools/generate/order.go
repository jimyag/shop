package generate

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

//
// GenerateOrderID
//  @Description: 订单号的生成
//  年月日时分秒+uid+随机数
//  @param userID
//
func GenerateOrderID(userID int32) int64 {
	now := time.Now()
	rand.Seed(now.UnixNano())
	orderSn := fmt.Sprintf("%d%d%d%d%d%d%d%d",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
		userID,
		rand.Intn(90)+10,
	)

	orderID, _ := strconv.ParseInt(orderSn, 10, 64)
	return orderID
}
