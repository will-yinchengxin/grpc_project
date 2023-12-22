package main

import (
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/alts"
	"log"
	"net"
	srv2 "wrpc/alts_rpc/server/srv"
	"wrpc/proto/pb"
)

var (
	port = flag.String("port", "0.0.0.0:8090", "port to listen on")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", *port)
	if err != nil {
		panic(err)
	}
	var (
		opts []grpc.ServerOption
	)
	// 创建ALTS凭证
	altsTC := alts.NewServerCreds(alts.DefaultServerOptions())
	opts = append(opts, grpc.Creds(altsTC))

	srv := grpc.NewServer(opts...)
	pb.RegisterGreeterServiceServer(srv, &srv2.Greeter{})
	log.Print("Start Server 8090 !!! ")
	err = srv.Serve(lis)
	if err != nil {
		panic(err)
	}
}
