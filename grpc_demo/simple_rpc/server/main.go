package main

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"wrpc/proto/pb"
	"wrpc/simple_rpc/server/server"
)

func main() {
	listen, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer()
	pb.RegisterUserServiceServer(srv, &server.UserServer{})
	fmt.Println("Start Server 8090 !!! ")
	err = srv.Serve(listen)
	if err != nil {
		panic(err)
	}
}
