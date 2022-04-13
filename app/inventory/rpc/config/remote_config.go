package config

//
// Postgres
//  @Description: 数据库的配置
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
//  @Description:  服务的配置
//
type ServiceInfo struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

//
// JaegerConfig
//  @Description: jaeger的配置
//
type JaegerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// ALLConfig 需要用的远程配置文件
type ALLConfig struct {
	Postgres    Postgres     `mapstructure:"postgres"`
	ServiceInfo ServiceInfo  `mapstructure:"service-info"`
	JaegerInfo  JaegerConfig `mapstructure:"jaeger-info"`
}
