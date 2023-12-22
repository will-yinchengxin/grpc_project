package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc_product/internal/config"
	"grpc_product/internal/consul"
	"grpc_product/product_web/request"
	"grpc_product/proto/pb"
	"log"
	"net/http"
	"sync"
)

var (
	once    sync.Once
	gClient pb.ProductServiceClient
)

func init() {
	once.Do(func() {
		instances := consul.GetTargetService(config.AppConf.AppInfo.SrvName)
		// TODO 这里获取了服务, 怎么实现负载均衡等
		conn, err := grpc.Dial(instances[0],
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robbin"}`),
		)
		if err != nil {
			log.Fatal(err)
		}
		gClient = pb.NewProductServiceClient(conn)
	})
}

func AdvertiseList(c *gin.Context) {
	list, err := gClient.AdvertiseList(c, &empty.Empty{})
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusOK, gin.H{
			"code": 10001,
			"msg":  "获取广告列表失败",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": list,
	})
}

func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg": "ok",
	})
}

func convertAdvertiseReqToPb(productReq request.AdvertiseReq) *pb.AdvertiseReq {
	return &pb.AdvertiseReq{
		Id:     productReq.Id,
		Index:  productReq.Index,
		Images: productReq.Images,
		Url:    productReq.Url,
	}
}
