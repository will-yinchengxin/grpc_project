package main

import (
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"wrpc/proto/pb"
	"wrpc/stream_rpc/client/cli"
)

var (
	serverAddr = flag.String("addr", "0.0.0.0:8090", "The server address in the format of host:port")
)

func main() {
	flag.Parse()
	var (
		opts = make([]grpc.DialOption, 0)
	)
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	cli.Chat(pb.NewStreamChatClient(conn))
}
