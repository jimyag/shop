一个微服务商城的实现

此项目使用gRPC从0到1实现一个完整的微服务的商城项目。主要用到的技术栈有：gin、postgresql、paseto、sqlc、migrate、docker、consul、jaeger、protobuf。

项目中一共涉及到：

1. 用户服务
2. 商品服务
3. 库存服务
4. 订单和购物车服服务



### todo

- [x] 用户服务
  - [x] 用户的RPC服务
  - [x] 用户的HTTP服务
- [x] 商品服务
  - [x] 商品的RPC服务
  - [x] 商品的HTTP服务
- [x] 库存服务
  - [x] 库存的RPC服务
- [ ] 订单服务
  - [x] 订单的RPC服务
  - [ ] 订单的HTTP服务
- [x] 使用consul实现服务注册
- [x] 使用consul实现服务发现
- [x] 使用consul实现配置中心
- [x] 实现同个服务多个实例的负载均衡
- [x] 使用jaeger实现链路追踪
- [x] 使用redis实现分布式锁

 
