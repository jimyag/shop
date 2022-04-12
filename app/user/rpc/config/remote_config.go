package config

//
// Postgres
//  @Description: postgres数据库的配置
//
type Postgres struct {
	Host     string `mapstructure:"host"`     // host
	Port     int    `mapstructure:"port"`     // port
	Database string `mapstructure:"database"` // database name
	User     string `mapstructure:"user"`     // user
	Password string `mapstructure:"password"` // password
	Type     string `mapstructure:"type"`     // 数据库的类型 postgres
}

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
//  @Description: jaeger的配置文件
//
type JaegerConfig struct {
	Host string `mapstructure:"host"` //host
	Port int    `mapstructure:"port"` //port
}

//
// ALLConfig
//  @Description: 需要用的远程配置文件
//
type ALLConfig struct {
	Postgres    Postgres     `mapstructure:"postgres"`     // postgres 的配置
	ServiceInfo ServiceInfo  `mapstructure:"service-info"` // 服务的配置
	JaegerInfo  JaegerConfig `mapstructure:"jaeger-info"`  // jaeger的配置文件
}
