package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"wrpc/proto/pb"
	srv2 "wrpc/stream_rpc/server/srv"
)

var (
	port = flag.Int("port", 8090, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", *port))
	if err != nil {
		panic(err)
	}
	var (
		// TODO: add some grpc options
		opts []grpc.ServerOption
	)
	srv := grpc.NewServer(opts...)
	pb.RegisterStreamChatServer(srv, &srv2.ChatServer{})
	log.Print("Start Server 8090 !!! ")
	err = srv.Serve(lis)
	if err != nil {
		panic(err)
	}
}
