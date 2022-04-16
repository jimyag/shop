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
