package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"grpc_product/internal/config"
	"grpc_product/util"
	"log"
	"strconv"
	"strings"
)

func RegisterGRPCService(srv *grpc.Server, addr string) {
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())
	defaultConfig := api.DefaultConfig()
	defaultConfig.Address = fmt.Sprintf("%s:%d", config.AppConf.ConsulConfig.Host, config.AppConf.ConsulConfig.Port)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		log.Fatal("NewClient error: ", err)
	}

	// 注册grpc服务
	uuid := util.GetUuid()
	reg := createServiceRegistration(addr, uuid)

	err = client.Agent().ServiceRegister(reg)
	if err != nil {
		log.Fatal("Grpc-Service Registration Error: ", err)
	}

	fmt.Println(fmt.Sprintf("%s启动在%#v", uuid, addr))
}

func createServiceRegistration(addr, randUUID string) *api.AgentServiceRegistration {
	port, _ := strconv.Atoi(strings.Split(addr, ":")[1])

	check := api.AgentServiceCheck{
		GRPC:                           addr,
		Timeout:                        "3s",
		Interval:                       "15s",
		DeregisterCriticalServiceAfter: "10s",
	}

	return &api.AgentServiceRegistration{
		Name:    config.AppConf.AppInfo.SrvName,
		Address: addr,
		ID:      randUUID,
		Port:    port,
		Tags:    config.AppConf.AppInfo.SrvTag,
		Check:   &check,
	}
}
