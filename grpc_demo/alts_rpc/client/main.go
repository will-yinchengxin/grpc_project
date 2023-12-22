package main

import (
	"flag"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/alts"
	"wrpc/alts_rpc/client/cli"
	"wrpc/proto/pb"
)

var (
	serverAddr = flag.String("serverInfo", "127.0.0.1:8090", "server info")
)

func main() {
	flag.Parse()
	// 创建ALTS凭证
	altsTC := alts.NewClientCreds(alts.DefaultClientOptions())

	conn, err := grpc.Dial(*serverAddr, grpc.WithTransportCredentials(altsTC))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewGreeterServiceClient(conn)

	cli.Greeter(client)
}
