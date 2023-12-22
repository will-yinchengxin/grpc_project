package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "grpc_product/internal"
	"grpc_product/internal/config"
	"grpc_product/internal/consul"
	"grpc_product/product_web/router"
	"grpc_product/util"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	randomId string
)

func init() {
	randomId = util.GetUuid()
	err := consul.RegisterSrv(
		config.AppConf.ProductWebConfig.Host,
		config.AppConf.ProductWebConfig.SrvName,
		randomId,
		config.AppConf.ProductWebConfig.Port,
		config.AppConf.ProductWebConfig.Tags,
	)
	if err != nil {
		log.Fatal("Web-Service Registration Error:", err)
	}
}

func main() {
	addr := fmt.Sprintf("%s:%d", config.AppConf.ProductWebConfig.Host,
		config.AppConf.ProductWebConfig.Port)
	r := gin.Default()
	router.InitRouter(r)
	go func() {
		err := r.Run(addr)
		if err != nil {
			log.Fatal(addr + "启动失败" + err.Error())
		}
		log.Println("Web Server Start: ", addr)
	}()
	q := make(chan os.Signal)
	signal.Notify(q, syscall.SIGINT, syscall.SIGTERM)
	<-q
	err := consul.DeRegister(randomId)
	if err != nil {
		log.Fatal("注销失败" + randomId + ":" + err.Error())
	}
	log.Println("注销成功" + randomId)
}
