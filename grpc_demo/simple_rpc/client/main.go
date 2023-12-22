package main

import (
	"google.golang.org/grpc"
	"wrpc/proto/pb"
	"wrpc/simple_rpc/client/client"
)

func main() {
	// conn, err := grpc.Dial(":8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(":8090", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client.Cli = pb.NewUserServiceClient(conn)
	client.GetUserInfo()
}
