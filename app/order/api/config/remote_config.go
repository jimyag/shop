package config

//
// ServiceInfo
//  @Description:  服务的信息
//
type ServiceInfo struct {
	Host string `mapstructure:"host"` // 服务所在的ip
	Port int    `mapstructure:"port"` // 服务所在的port
	Name string `mapstructure:"name"` // 服务的名称 服务的名称应该是唯一的
}

//
// JaegerConfig
//  @Description: jaeger配置
//
type JaegerConfig struct {
	Host string `mapstructure:"host"` // jaeger 的 host
	Port int    `mapstructure:"port"` // jaeger 的 port
}

//
// Secret
//  @Description: token的配置
//
type Secret struct {
	PrivateKey string `mapstructure:"private-key"` // token的 私钥
	PublicKey  string `mapstructure:"public-key"`  // token的 公钥
	Duration   int    `mapstructure:"duration"`    // token的过期时间 单位 小时
}

//
// RedisInfo
//  @Description: redis的配置
//
type RedisInfo struct {
	Host string `mapstructure:"host"` // redis的host
	Port int    `mapstructure:"port"` // redis的port
}

type GrpcServer struct {
	Name string `mapstructure:"name"` // 服务的名称 服务的名称应该是唯一的
}

//
// ServerInfo
//  @Description: 整个服务的配置
//
type ServerInfo struct {
	ServiceInfo     ServiceInfo  `mapstructure:"server-info"`       // 本地服务的信息
	JaegerInfo      JaegerConfig `mapstructure:"jaeger-info"`       // jaeger 的 配置
	Secret          Secret       `mapstructure:"secret"`            // token 密钥
	Redis           RedisInfo    `mapstructure:"redis"`             // redis 的配置
	OrderGrpcServer GrpcServer   `mapstructure:"order-grpc-server"` // user grpc server 的配置
}
