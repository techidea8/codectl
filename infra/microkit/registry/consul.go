package registry

import (
	"fmt"
	"time"

	consulapi "github.com/hashicorp/consul/api"
)

type ConsulRegistry struct {
	config *consulapi.Config
	client *consulapi.Client
}
type ConsulRegistryOption func(*ConsulRegistry)

func SetAddress(addr string) ConsulRegistryOption {
	return func(cr *ConsulRegistry) {
		cr.config.Address = addr
	}
}

func NewConsulRegistry(config *consulapi.Config, opts ...ConsulRegistryOption) *ConsulRegistry {
	config.Address = consulapi.DefaultConfig().Address
	ret := &ConsulRegistry{
		config: config,
	}
	for _, v := range opts {
		v(ret)
	}
	return ret
}
func (r *ConsulRegistry) Register(service *Service) (err error) {
	// 创建连接consul服务配置
	config := r.config
	r.client, err = consulapi.NewClient(config)
	if err != nil {
		return
	}

	// 创建注册到consul的服务到
	registration := new(consulapi.AgentServiceRegistration)
	registration.Name = service.Name
	// registration.ID = registration.Name
	registration.Port = service.Port
	registration.Tags = service.Tags
	registration.Address = service.Host

	// 增加consul健康检查回调函数
	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d%s", registration.Address, registration.Port, service.Check.HealthPath)
	if service.Check.Timeout > 0 {
		service.Check.Timeout = time.Second * 10
	}
	check.Timeout = service.Check.Timeout.String()
	if service.Check.Interval == 0 {
		service.Check.Interval = time.Second * 5
	}
	check.Interval = service.Check.Interval.String()
	if service.Check.DeregisterCriticalServiceAfter == 0 {
		service.Check.DeregisterCriticalServiceAfter = time.Second * 15
	}
	check.DeregisterCriticalServiceAfter = service.Check.DeregisterCriticalServiceAfter.String()
	registration.Check = check
	// 注册服务到consul
	err = r.client.Agent().ServiceRegister(registration)
	return
}
func (r *ConsulRegistry) DeRegister(serviceId string) (err error) {
	return r.client.Agent().ServiceDeregister(serviceId)
}

func (r *ConsulRegistry) Run() {
	for {
		time.Sleep(1 * time.Second)
	}
}
