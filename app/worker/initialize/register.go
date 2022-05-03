package initialize

import (
	"github.com/hubogle/Crontab/app/worker/config"
	"github.com/hubogle/Crontab/util/register/consul"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
)

var serviceId string

func InitRegister() {
	cfg := config.GetConfig()
	registerClient := consul.NewRegistryClient(cfg.Consul.Host, cfg.Consul.Port)
	serviceId = uuid.NewV4().String() // 随机生成 UUID
	err := registerClient.Register(cfg.App.Host, cfg.App.Port, cfg.App.Name, cfg.App.Tag, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
}

func UnRegister() {
	cfg := config.GetConfig()
	registerClient := consul.NewRegistryClient(cfg.Consul.Host, cfg.Consul.Port)
	err := registerClient.DeRegister(serviceId)
	if err != nil {
		zap.S().Panic("服务注销失败:", err.Error())
	}
}
