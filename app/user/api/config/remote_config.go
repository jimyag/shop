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

//
// Email
//  @Description: 邮件信息的配置
//
type Email struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"` // host
	Form     string `mapstructure:"form" json:"form" yaml:"form"` //
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	IsSsl    bool   `mapstructure:"is-ssl" json:"isSsl" yaml:"is-ssl"`
	Secret   string `mapstructure:"secret" json:"secret" yaml:"secret"`
	Nickname string `mapstructure:"nickname" json:"nickname" yaml:"nickname"`
}

//
// Timeout
//  @Description: 所有超时的配置
//
type Timeout struct {
	CreateUserEmail int `mapstructure:"create-user-email"` // 多少分钟
}

type GrpcServer struct {
	Name string `mapstructure:"name"` // 服务的名称 服务的名称应该是唯一的
}

//
// ServerInfo
//  @Description: 整个服务的配置
//
type ServerInfo struct {
	ServiceInfo    ServiceInfo  `mapstructure:"server-info"`      // 本地服务的信息
	JaegerInfo     JaegerConfig `mapstructure:"jaeger-info"`      // jaeger 的 配置
	Secret         Secret       `mapstructure:"secret"`           // token 密钥
	Redis          RedisInfo    `mapstructure:"redis"`            // redis 的配置
	Email          Email        `mapstructure:"email"`            // 邮件的信息
	Timeout        Timeout      `mapstructure:"timeout"`          // 各种超时配置
	UserGrpcServer GrpcServer   `mapstructure:"user-grpc-server"` // user grpc server 的配置
}
