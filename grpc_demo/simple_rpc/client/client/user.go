package client

import (
	"context"
	"fmt"
	"wrpc/proto/pb"
)

var Cli pb.InterceptorServiceClient

func GetUserInfo() {
	req := pb.InterceptorRequest{
		Name: "will",
	}
	userInfo, err := Cli.Interceptor(context.Background(), &req)
	if err != nil {
		fmt.Printf("Get server err return %v", err)
		return
	}
	fmt.Printf("Get server return %v", userInfo)
}
