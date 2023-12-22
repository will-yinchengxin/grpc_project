package client

import (
	"context"
	"fmt"
	"time"
	"wrpc/proto/pb"
)

var Cli pb.InterceptorServiceClient

func DoInterceptor() {
	// 超时控制
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	req := pb.InterceptorRequest{
		Name: "will",
	}
	userInfo, err := Cli.Interceptor(ctx, &req)
	if err != nil {
		fmt.Printf("Get Server Err Return: %s", err.Error())
		return
	}
	fmt.Printf("Get server return %v", userInfo)
}
