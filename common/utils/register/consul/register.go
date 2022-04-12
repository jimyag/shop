package consul

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

// RegistryGrpc consul注册grpc服务
type RegistryGrpc struct {
	Host   string // consul 的 host
	Port   int    // consul 的 port
	client *api.Client
}

type RegisterClient interface {
	// Register host port 本地要注册服务的host port ，服务的名称 标签，id
	Register(host string, port int, name string, tags []string, id string) error // 注册的接口
	DeRegister(serviceID string) error                                           // 注销的接口
}

func NewRegistryRpcClient(host string, port int) RegisterClient {
	return &RegistryGrpc{
		Host: host,
		Port: port,
	}
}

func (r *RegistryGrpc) Register(host string, port int, name string, tags []string, id string) error {
	// consul服务注册 consul 的IP的port
	apiCfg := api.DefaultConfig()
	apiCfg.Address = fmt.Sprintf("%s:%d",
		r.Host,
		r.Port,
	)
	var err error
	r.client, err = api.NewClient(apiCfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// 服务的ip
	check := api.AgentServiceCheck{
		GRPC: fmt.Sprintf("%s:%d",
			host,
			port,
		),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "100s",
	}

	// consul 做健康检查的ip
	serviceRegistration := api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Port:    port,
		Address: host,
		Tags:    tags,
	}

	serviceRegistration.Check = &check
	err = r.client.Agent().ServiceRegister(&serviceRegistration)
	if err != nil {
		log.Fatalln("注册失败", err)
	}
	log.Printf("启动 grpc 健康检查 name: %s ,id: %s", name, id)
	return nil
}

func (r *RegistryGrpc) DeRegister(serviceID string) error {
	return r.client.Agent().ServiceDeregister(serviceID)
}

// RegistryHttp  consul注册http服务
type RegistryHttp struct {
	Host   string // consul 的 host
	Port   int    // consul 的 port
	client *api.Client
}

func NewRegistryHttpClient(host string, port int) RegisterClient {
	return &RegistryHttp{
		Host: host,
		Port: port,
	}
}

func (r *RegistryHttp) Register(host string, port int, name string, tags []string, id string) error {
	// consul服务注册 consul 的IP的port
	apiCfg := api.DefaultConfig()
	apiCfg.Address = fmt.Sprintf("%s:%d",
		r.Host,
		r.Port,
	)
	var err error
	r.client, err = api.NewClient(apiCfg)
	if err != nil {
		log.Fatalln(err.Error())
	}
	// 服务的ip
	check := api.AgentServiceCheck{
		HTTP: fmt.Sprintf("http://%s:%d/health",
			host,
			port,
		),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "100s",
	}

	// consul 做健康检查的ip
	serviceRegistration := api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Port:    port,
		Address: host,
		Tags:    tags,
	}

	serviceRegistration.Check = &check
	err = r.client.Agent().ServiceRegister(&serviceRegistration)
	if err != nil {
		log.Fatalln("注册失败", err)
	}
	log.Printf("启动 grpc 健康检查 name: %s ,id: %s", name, id)
	return nil
}

func (r *RegistryHttp) DeRegister(serviceID string) error {
	return r.client.Agent().ServiceDeregister(serviceID)
}
