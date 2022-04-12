package model

//
// ConfigCenterInfo
//  @Description: 配置中心的位置信息
//
type ConfigCenterInfo struct {
	Host        string `mapstructure:"host"`         // consul所在的host
	Port        int    `mapstructure:"port"`         // consul的port
	ReleasePath string `mapstructure:"release-path"` // 线上版本的配置文件路径
	DebugPath   string `mapstructure:"debug-path"`   //开发环境的配置文件路径
	Type        string `mapstructure:"type"`         //配置的类型 consul etcd 这个是consul
	FileType    string `mapstructure:"file-type"`    // 配置文件的类型 yaml toml
	EnvName     string `mapstructure:"env-name"`     // 开发环境的环境变量的值
}
