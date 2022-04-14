package config

//
// Postgres
//  @Description: postgres 配置
//
type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Type     string `mapstructure:"type"`
}

//
// ServiceInfo
//  @Description: 当前服务的配置
//
type ServiceInfo struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

// ConsulInfo 这里是指注册服务和服务发现的Consul
type ConsulInfo struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

//
// JaegerConfig
//  @Description: jaeger 的配置
//
type JaegerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// ALLConfig 需要用的远程配置文件
type ALLConfig struct {
	Postgres    Postgres     `mapstructure:"postgres"`
	ServiceInfo ServiceInfo  `mapstructure:"service-info"`
	ConsulInfo  ConsulInfo   `mapstructure:"consul-info"`
	JaegerInfo  JaegerConfig `mapstructure:"jaeger-info"`
	ThirdServer ThirdServer  `mapstructure:"third-server"`
}

//
// GrpcServer
//  @Description: grpc服务的配置
//
type GrpcServer struct {
	Name string `mapstructure:"name"` // 服务的名称 服务的名称应该是唯一的
}

//
// ThirdServer
//  @Description: 第三方服务的配置
//
type ThirdServer struct {
	GoodsGrpcServer     GrpcServer `mapstructure:"goods-grpc-server"`
	InventoryGrpcServer GrpcServer `mapstructure:"inventory-grpc-server"`
}
