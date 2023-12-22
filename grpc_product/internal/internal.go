package internal

import (
	"grpc_product/internal/config"
	"grpc_product/internal/mysql"
)

func init() {
	config.InitConfig()
	mysql.InitMysql()
}
