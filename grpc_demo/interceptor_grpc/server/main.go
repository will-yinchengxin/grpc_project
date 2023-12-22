package main

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"wrpc/interceptor_grpc/server/server"
	"wrpc/proto/pb"
)

func main() {
	// 注入拦截器
	opt := grpc.UnaryInterceptor(server.MyInterceptor)
	srv := grpc.NewServer(opt)
	pb.RegisterInterceptorServiceServer(srv, &server.Interceptor{})
	listen, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	fmt.Println("Start Server 8090 !!! ")
	err = srv.Serve(listen)
	if err != nil {
		panic(err)
	}
}
