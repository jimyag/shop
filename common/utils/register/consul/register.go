package consul

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

//
// RegistryGrpc
//  @Description:  在consul中注册grpc服务
//
type RegistryGrpc struct {
	Host   string // consul 的 host
	Port   int    // consul 的 port
	client *api.Client
}

//
// NewRegistryRpcClient
//  @Description:
//  @param host
//  @param port
//  @return RegisterClient
//
func NewRegistryRpcClient(host string, port int) RegisterClient {
	return &RegistryGrpc{
		Host: host,
		Port: port,
	}
}

//
// RegisterClient
//  @Description: 注册服务时的接口
//
type RegisterClient interface {
	// Register host port 本地要注册服务的host port ，服务的名称 标签，id
	Register(host string, port int, name string, tags []string, id string) error // 注册的接口
	DeRegister(serviceID string) error                                           // 注销的接口
}

//
// Register
//  @Description: 注册grpc服务
//  @receiver r
//  @param host 服务所在的host
//  @param port 服务所在的port
//  @param name 服务的名称
//  @param tags 服务的标签
//  @param id 服务的id 唯一
//  @return error
//
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

//
// DeRegister
//  @Description: 注销服务
//  @receiver r
//  @param serviceID 服务的id
//  @return error
//
func (r *RegistryGrpc) DeRegister(serviceID string) error {
	return r.client.Agent().ServiceDeregister(serviceID)
}

//
// RegistryHttp
//  @Description: 在consul中注册http服务
//
type RegistryHttp struct {
	Host   string // consul 的 host
	Port   int    // consul 的 port
	client *api.Client
}

//
// NewRegistryHttpClient
//  @Description:
//  @param host
//  @param port
//  @return RegisterClient
//
func NewRegistryHttpClient(host string, port int) RegisterClient {
	return &RegistryHttp{
		Host: host,
		Port: port,
	}
}

//
// Register
//  @Description: 注册http的服务，默认注册时使用http://%s:%d/health的url进行健康检查
//  @receiver r
//  @param host 服务所在的host
//  @param port 服务所在的port
//  @param name 服务的名称
//  @param tags 服务的标签
//  @param id 服务的id 唯一
//  @return error
//
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

//
// DeRegister
//  @Description: 注销http服务
//  @receiver r
//  @param serviceID
//  @return error
//
func (r *RegistryHttp) DeRegister(serviceID string) error {
	return r.client.Agent().ServiceDeregister(serviceID)
}
