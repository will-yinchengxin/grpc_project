package consul

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"grpc_product/internal/config"
	"log"
	"strconv"
)

func RegisterSrv(host, name, id string, port int, tags []string) error {
	client, err := getConsulAddr()
	if err != nil {
		return err
	}

	agentServiceRegistration := new(api.AgentServiceRegistration)
	agentServiceRegistration.Address = host + ":" + strconv.Itoa(port)
	agentServiceRegistration.Port = port
	agentServiceRegistration.ID = id
	agentServiceRegistration.Name = name
	agentServiceRegistration.Tags = tags
	serverAddr := fmt.Sprintf("http://%s:%d/health", host, port)
	check := api.AgentServiceCheck{
		HTTP:                           serverAddr,
		Timeout:                        "3s",
		Interval:                       "15s",
		DeregisterCriticalServiceAfter: "3s",
	}
	agentServiceRegistration.Check = &check
	return client.Agent().ServiceRegister(agentServiceRegistration)
}

func DeRegister(serviceId string) error {
	client, err := getConsulAddr()
	if err != nil {
		zap.S().Error(err)
		return err
	}
	return client.Agent().ServiceDeregister(serviceId)
}

func GetTargetService(srv string) []string {
	client, err := getConsulAddr()
	if err != nil {
		log.Fatal("Get CClient Err: " + err.Error())
		return nil
	}

	instances, _, err := client.Catalog().Service(srv, "", nil)
	if err != nil {
		fmt.Println("获取服务实例列表失败：", err)
		return nil
	}

	info := make([]string, len(instances))
	for key, instance := range instances {
		fmt.Printf("服务名：%s，地址：%s，端口：%d\n", instance.ServiceName, instance.Address, instance.ServicePort)
		info[key] = fmt.Sprintf("%s:%d", instance.Address, instance.ServicePort)
	}
	return info
}

func GetServiceList() error {
	client, err := getConsulAddr()
	if err != nil {
		log.Fatal("Get CClient Err: " + err.Error())
		return err
	}
	services, err := client.Agent().Services()
	if err != nil {
		return err
	}
	for k, v := range services {
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println("----------------------")
	}
	return nil
}

func FilterService() error {
	client, err := getConsulAddr()
	if err != nil {
		return err
	}

	services, err := client.Agent().ServicesWithFilter("Service==account_web")
	if err != nil {
		return err
	}
	for k, v := range services {
		fmt.Println(k)
		fmt.Println(v)
		fmt.Println("----------------------")
	}
	return nil
}

func getConsulAddr() (*api.Client, error) {
	defaultConfig := api.DefaultConfig()
	h := config.AppConf.ConsulConfig.Host
	p := config.AppConf.ConsulConfig.Port
	defaultConfig.Address = fmt.Sprintf("%s:%d", h, p)
	client, err := api.NewClient(defaultConfig)
	if err != nil {
		return nil, err
	}
	return client, nil
}
