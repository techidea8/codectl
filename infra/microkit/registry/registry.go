package registry

import "time"

type Registry interface {
	Register(service *Service) (err error)
	DeRegister(serviceId string) (err error)
	Run()
}

type Check struct {
	HealthPath                     string        `json:"healthpath"`
	Timeout                        time.Duration `json:"timeout"`
	Interval                       time.Duration `json:"interval"`
	DeregisterCriticalServiceAfter time.Duration `json:"deregisterCriticalServiceAfter"`
}
type Service struct {
	Name  string   `json:"name"`
	Host  string   `json:"host"`
	Port  int      `json:"port"`
	Tags  []string `json:"tags"`
	Check Check    `json:"check"`
}
