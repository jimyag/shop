package request

//
// GetCartListRequest
//  @Description: 获得用户购物车列表的参数
//
type GetCartListRequest struct {
	UserId int32 `json:"userId" validate:"required,min=1" label:"用户ID"`
}
