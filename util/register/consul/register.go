package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

type RegisterClient interface {
	Register(address string, port int, name string, tags []string, id string) error
	DeRegister(serviceId string) error
}

type Registry struct {
	Host string
	Port int
}

func NewRegistryClient(host string, port int) RegisterClient {
	return &Registry{
		Host: host,
		Port: port,
	}
}

// Register 服务注册，使用前提必须提供 /health 接口
func (r *Registry) Register(address string, port int, name string, tags []string, id string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	// 生成对应的检查对象
	check := &api.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d/health", address, port),
		Timeout:                        "5s",  // 超时时间
		Interval:                       "5s",  // 健康检查间隔
		DeregisterCriticalServiceAfter: "30s", // 故障检查失败30s后 consul自动将注册服务删除
	}

	// 生成注册对象
	registration := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Address: address,
		Port:    port,
		Tags:    tags,
		Check:   check,
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	return nil
}

// DeRegister 服务注销
func (r *Registry) DeRegister(serviceId string) error {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", r.Host, r.Port)

	client, err := api.NewClient(cfg)
	if err != nil {
		return err
	}
	err = client.Agent().ServiceDeregister(serviceId)
	return err
}
