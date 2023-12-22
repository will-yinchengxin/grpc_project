package config

import (
	"encoding/json"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

var (
	// nacos: http://172.16.27.95:8848/nacos
	configFileName = "dev-config.yaml"
	NacosConf      NacosConfig
	AppConf        AppConfig
)

func InitConfig() {
	v := viper.New()
	v.SetConfigFile(configFileName)
	v.ReadInConfig()
	err := v.Unmarshal(&NacosConf)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	initFromNacos()
}

func initFromNacos() {
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: NacosConf.Host,
			Port:   NacosConf.Port,
		},
	}
	clientConfig := constant.ClientConfig{
		NamespaceId:         NacosConf.NameSpace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "nacos/log",
		CacheDir:            "nacos/cache",
		LogLevel:            "debug",
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: NacosConf.DataId,
		Group:  NacosConf.Group,
	})
	if err != nil {
		panic(err)
	}
	if content == "" {
		panic("配置文件为空")
	}
	json.Unmarshal([]byte(content), &AppConf)
}
