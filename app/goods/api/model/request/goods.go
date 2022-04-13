package request

//
// CreateGoods
//  @Description: 创建商品的请求
//
type CreateGoods struct {
	Name  string  `json:"name" validate:"required" label:"商品名称"`
	Price float32 `json:"price" validate:"required,min=0.1" label:"商品价格"`
}
