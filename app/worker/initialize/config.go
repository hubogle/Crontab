package initialize

import (
	"encoding/json"
	"github.com/hubogle/Crontab/app/worker/config"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
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
	// v.WatchConfig()
	// v.OnConfigChange(func(e fsnotify.Event) {
	// 	_ = v.ReadInConfig()
	// 	_ = v.Unmarshal(Config)
	// })
	clientConfig := constant.ClientConfig{
		NamespaceId:         Config.NacosConfig.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp/nacos/log",
		CacheDir:            "tmp/nacos/cache",
		LogLevel:            "debug",
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: Config.NacosConfig.Host,
			Port:   Config.NacosConfig.Port,
		},
	}
	configClient, _ := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	content, _ := configClient.GetConfig(vo.ConfigParam{
		DataId: Config.NacosConfig.DataId,
		Group:  Config.NacosConfig.Group})
	err = json.Unmarshal([]byte(content), Config)
	if err != nil {
		zap.S().Fatalf("读取 Nacos 配置失败 %s", zap.Error(err))
	}
}
