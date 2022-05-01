package initialize

import (
	"github.com/fsnotify/fsnotify"
	"github.com/hubogle/Crontab/app/worker/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// InitConfig 配置文件初始化
func InitConfig() {
	var (
		confFileName string
		v            *viper.Viper
		err          error
	)
	confFileName = "app/worker/config/config-dev.toml"
	Config := config.SetConfig()
	v = viper.New()
	v.SetConfigFile(confFileName)
	if err = v.ReadInConfig(); err != nil {
		zap.S().Fatalf("读取配置失败 %s", err.Error())
	}
	if err = v.Unmarshal(Config); err != nil {
		zap.S().Fatalf("读取配置失败 %s", err.Error())
	}
	// viper的功能 - 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		_ = v.ReadInConfig()
		_ = v.Unmarshal(Config)
	})
}
