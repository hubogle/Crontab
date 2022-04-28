package config

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"` // 项目端口
	Name string `mapstructure:"name"` // 项目名称
}
