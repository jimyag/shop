package request

//
// GetOrderListRequest
//  @Description:  获取订单列表
//
type GetOrderListRequest struct {
	PageNum  int32 `json:"page_num" validate:"required,min=1" label:"页码"`
	PageSize int32 `json:"page_size" validate:"required,min=1" label:"每页数量"`
	UserID   int32 `json:"user_id" validate:"required,min=1" label:"用户ID"`
}

//
// GetOrderDetailRequest
//  @Description: 获取订单详情
//
type GetOrderDetailRequest struct {
	OrderID int64 `json:"order_id" validate:"required,min=1" label:"订单ID"`
}

//
// CreateOrderRequest
//  @Description:  创建订单
//
type CreateOrderRequest struct {
	UserID  int32  `json:"user_id" validate:"required,min=1" label:"用户ID"`
	Address string `json:"address" validate:"required" label:"收货地址"`
	Mobile  string `json:"mobile" validate:"required" label:"手机号"`
	Name    string `json:"name" validate:"required" label:"收货人"`
	Post    string `json:"post" validate:"required" label:"邮编"`
}
