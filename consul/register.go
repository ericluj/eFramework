package consul

import (
	"fmt"
	"time"

	log "github.com/ericluj/elog"
	"github.com/hashicorp/consul/api"
)

const ConsulAddress = "consul:8500"

// const ConsulAddress = "127.0.0.1:8500"

type ConsulService struct {
	IP   string
	Port int
	Tag  []string
	Name string
}

func RegitserService(s *ConsulService) {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = ConsulAddress

	client, err := api.NewClient(consulConfig)
	if err != nil {
		log.Infof("NewClient error: %v", err)
		return
	}

	agent := client.Agent()
	interval := 10 * time.Second
	deregister := 1 * time.Minute

	reg := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%v-%v-%v", s.Name, s.IP, s.Port), // 服务节点的名称
		Name:    s.Name,                                        // 服务名称
		Tags:    s.Tag,                                         // tag，可以为空
		Port:    s.Port,                                        // 服务端口
		Address: s.IP,                                          // 服务 IP
		Check: &api.AgentServiceCheck{ // 健康检查
			Interval:                       interval.String(),                              // 健康检查间隔
			HTTP:                           fmt.Sprintf("http://%s:%d/health", s.IP, 8001), // 健康检查地址
			DeregisterCriticalServiceAfter: deregister.String(),                            // 注销时间，相当于过期时间
		},
	}

	log.Infof("service %v registing to %v\n", s.Name, ConsulAddress)
	if err := agent.ServiceRegister(reg); err != nil {
		log.Infof("ServiceRegister error: %v", err)
		return
	}

}
