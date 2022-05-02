package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hubogle/Crontab/app/master/config"
)

// InitConsul 初始化链接 Consul 服务
func InitConsul() (*api.Client, error) {
	var (
		client  *api.Client
		err     error
		address string
	)
	cfg := config.GetConfig().Consul
	address = fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	apiConfig := &api.Config{
		Address: address,
		Scheme:  "http",
	}
	client, err = api.NewClient(apiConfig)
	return client, err
}
