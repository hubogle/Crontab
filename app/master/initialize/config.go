package initialize

import (
	"github.com/fsnotify/fsnotify"
	"github.com/hubogle/Crontab/app/master/global"
	"github.com/spf13/viper"
)

// InitConfig 配置文件初始化
func InitConfig() {
	var (
		confFileName string
		v            *viper.Viper
		err          error
	)
	confFileName = "app/master/config/config-dev.json"

	v = viper.New()
	v.SetConfigFile(confFileName)
	if err = v.ReadInConfig(); err != nil {
		panic(err)
	}
	if err = v.Unmarshal(global.ServerConfig); err != nil {
		panic(err)
	}
	// viper的功能 - 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
	})
}
