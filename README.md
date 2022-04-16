一个微服务商城的实现

此项目使用gRPC从0到1实现一个完整的微服务的商城项目。主要用到的技术栈有：gin、postgresql、paseto、sqlc、migrate、docker、consul、jaeger、protobuf。

项目中一共涉及到：

1. 用户服务
2. 商品服务
3. 库存服务
4. 订单和购物车服服务

### 快速开始

请提前配置docker环境，可以参考[服务器环境的配置 - 步履不停 (jimyag.cn)](https://jimyag.cn/post/173a3c06/)这篇文章在linux中配置docker，[Docker基础入门 - 步履不停 (jimyag.cn)](https://jimyag.cn/post/8b63f587/#安装)在windows中配置docker环境。

参考这篇[从0到1实现完整的微服务框架-项目介绍 - 步履不停 (jimyag.cn)](https://jimyag.cn/post/5f073a52/)配置相关的环境。

#### 启动consul

```shell
docker run -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp consul consul agent -dev -client=0.0.0.0
```

#### 启动jaeger

```shell
docker run -d --name jaeger   -e COLLECTOR_ZIPKIN_HOST_PORT=:9411   -p 5775:5775/udp   -p 6831:6831/udp   -p 6832:6832/udp   -p 5778:5778   -p 16686:16686   -p 14250:14250   -p 14268:14268   -p 14269:14269   -p 9411:9411 jaegertracing/all-in-one:1.32
```

#### 配置数据库

user的数据库

```shell
docker run --name shop-user -p 35432:5432 -e POSTGRES_PASSWORD=postgres -e TZ=PRC -d postgres:14-alpine
```

inventory的数据库

```shell
docker run --name shop-inventory -p 35433:5432 -e POSTGRES_PASSWORD=postgres -e TZ=PRC -d postgres:14-alpine
```

order的数据库

```shell
docker exec -it shop-order createdb --username=postgres --owner=postgres shop
```

goods的数据库

```shell
docker run --name shop-goods -p 35435:5432 -e POSTGRES_PASSWORD=postgres -e TZ=PRC -d postgres:14-alpine
```

### 数据库迁移

user的数据库

```powershell
cd app/user/rpc
migrate -path db/migration -database "postgresql://postgres:postgres@localhost:35432/shop?sslmode=disable" -verbose up
```

inventory的数据库

```powershell
cd app/inventory/rpc
migrate -path db/migration -database "postgresql://postgres:postgres@localhost:35433/shop?sslmode=disable" -verbose up
```

order的数据库

```powershell
cd app/order/rpc
migrate -path db/migration -database "postgresql://postgres:postgres@localhost:35434/shop?sslmode=disable" -verbose up
```

goods的数据库

```powershell
cd app/goods/rpc
migrate -path db/migration -database "postgresql://postgres:postgres@localhost:35435/shop?sslmode=disable" -verbose up
```

#### 启动服务

以启动`user-api`服务为例

```powershell
cd app/user/api
run.bat
```

服务启动

### TODO

- [x] 用户服务
  - [x] 用户的RPC服务
  - [x] 用户的HTTP服务
- [x] 商品服务
  - [x] 商品的RPC服务
  - [x] 商品的HTTP服务
- [x] 库存服务
  - [x] 库存的RPC服务
- [x] 订单服务
  - [x] 订单的RPC服务
  - [x] 订单的HTTP服务
- [x] 使用consul实现服务注册
- [x] 使用consul实现服务发现
- [x] 使用consul实现配置中心
- [x] 实现同个服务多个实例的负载均衡
- [x] 使用jaeger实现链路追踪
- [x] 使用redis实现分布式锁
- [ ] 使用rocketMQ实现消息队列
- [ ] 完善管理员和普通用户的权限
- [ ] 完善用户只能给自己下单的限制
- [ ] 完善启动服务的文档

 
