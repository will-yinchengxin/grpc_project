package main

import (
	"fmt"
	"google.golang.org/grpc"
	_ "grpc_product/internal"
	"grpc_product/internal/consul"
	srv2 "grpc_product/product_srv/srv"
	"grpc_product/proto/pb"
	"grpc_product/util"
	"log"
	"net"
)

func main() {
	addr := util.GetAddr()
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterProductServiceServer(srv, &srv2.ProductServer{})

	consul.RegisterGRPCService(srv, addr)

	fmt.Println("Start Server: " + addr)
	err = srv.Serve(listen)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
