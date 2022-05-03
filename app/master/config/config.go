package config

import "time"

var config = new(Config)

type Config struct {
	App struct {
		Name     string `mapstructure:"name"` // 项目名称
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"` // 项目端口
		RunMode  string `mapstructure:"runMode"`
		Language string `mapstructure:"language"`
	} `mapstructure:"app"`
	Grpc struct { // 提供 gRPC 服务的端口
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"` // grpc 提供的端口
	} `mapstructure:"grpc"`
	MySQL struct {
		Host            string        `mapstructure:"host"`
		Port            int           `mapstructure:"port"`
		User            string        `mapstructure:"user"`
		Pass            string        `mapstructure:"pass"`
		DBName          string        `mapstructure:"dbname"`
		MaxOpenConn     int           `mapstructure:"maxOpenConn"`
		MaxIdleConn     int           `mapstructure:"maxIdleConn"`
		MaxLifeTimeConn time.Duration `mapstructure:"maxLifeTimeConn"`
	} `mapstructure:"mysql"`
	Consul struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"consul"`
}

func GetConfig() Config {
	return *config
}

func SetConfig() *Config {
	return config
}
