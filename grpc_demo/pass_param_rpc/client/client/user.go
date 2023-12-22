package client

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"wrpc/proto/pb"
)

var Cli pb.PassParamClient

func SendMsg() {
	// Add metadata to header
	md := metadata.Pairs("token", "123456")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	req := pb.PassParamRequest{
		Name: "will",
	}
	msg, err := Cli.SendParam(ctx, &req)
	if err != nil {
		fmt.Printf("Get server err return %v", err)
		return
	}
	fmt.Printf("Get server return %v", msg)
}
