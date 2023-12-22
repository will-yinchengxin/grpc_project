package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"wrpc/interceptor_grpc/client/client"
	"wrpc/proto/pb"
)

func main() {
	// conn, err := grpc.Dial(":8090", grpc.WithInsecure())
	conn, err := grpc.Dial(":8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client.Cli = pb.NewInterceptorServiceClient(conn)
	client.DoInterceptor()

}
