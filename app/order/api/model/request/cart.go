package request

//
// GetCartListRequest
//  @Description: 获得用户购物车列表的参数
//
type GetCartListRequest struct {
	UserId int32 `json:"user_id" validate:"required,min=1" label:"用户ID"`
}

//
// CreateCartItemRequest
//  @Description: 创建购物车项的参数
//
type CreateCartItemRequest struct {
	UserId  int32 `json:"user_id" validate:"required,min=1" label:"用户ID"`
	GoodsId int32 `json:"goods_id" validate:"required,min=1" label:"商品ID"`
	Nums    int32 `json:"nums" validate:"required,min=1" label:"数量"`
}

//
// DeleteShopCartRequest
//  @Description: 删除购物车项的参数
//
type DeleteShopCartRequest struct {
	UserId  int32 `json:"user_id" validate:"required,min=1" label:"用户ID"`
	GoodsId int32 `json:"goods_id" validate:"required,min=1" label:"商品ID"`
}

type UpdateShopCartRequest struct {
	UserId  int32 `json:"user_id" validate:"required,min=1" label:"用户ID"`
	GoodsId int32 `json:"goods_id" validate:"required,min=1" label:"商品ID"`
	Nums    int32 `json:"nums" validate:"required,min=1" label:"数量"`
	Checker bool  `json:"checker" validate:"required" label:"是否选中"`
}
