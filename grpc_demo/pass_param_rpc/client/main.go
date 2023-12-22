package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"wrpc/pass_param_rpc/client/client"
	"wrpc/proto/pb"
)

func main() {
	//conn, err := grpc.Dial(":8090", grpc.WithInsecure())
	conn, err := grpc.Dial(":8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client.Cli = pb.NewPassParamClient(conn)
	client.SendMsg()
}
