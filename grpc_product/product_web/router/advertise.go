package router

import (
	"github.com/gin-gonic/gin"
	"grpc_product/product_web/controller"
)

func InitRouter(r *gin.Engine) {

	productGroup := r.Group("/v1/product/advertise")
	{
		// curl -X GET http://localhost:9988/v1/product/advertise/list
		productGroup.GET("/list", controller.AdvertiseList)
	}
	r.GET("/health", controller.HealthHandler)
}
