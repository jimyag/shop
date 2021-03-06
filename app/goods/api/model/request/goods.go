package request

//
// CreateGoods
//  @Description: 创建商品的请求
//
type CreateGoods struct {
	Name  string  `json:"name" validate:"required" label:"商品名称"`
	Price float32 `json:"price" validate:"required,min=0.1" label:"商品价格"`
}

//
// UpdateGoods
//  @Description: 更新商品信息
//
type UpdateGoods struct {
	ID    int32   `json:"id" validate:"required,min=1" label:"商品ID"`
	Name  string  `json:"name" validate:"required" label:"商品名称"`
	Price float32 `json:"price" validate:"required,min=0.1" label:"商品价格"`
}

//
// GoodsIDRequest
//  @Description: 获得商品信息
//
type GoodsIDRequest struct {
	ID int32 `json:"id" validate:"required,min=1" label:"商品ID"`
}

//
// GoodsBatchRequest
//  @Description: 批量获得商品信息
//
type GoodsBatchRequest struct {
	GoodsBatchID []GoodsIDRequest `json:"goods_batch_id" validate:"required" label:"批量查询ID"`
}
