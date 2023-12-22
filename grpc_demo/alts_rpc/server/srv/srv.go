package srv

import (
	"context"
	"log"
	"wrpc/proto/pb"
)

type Greeter struct {
	pb.UnsafeGreeterServiceServer
}

func (g Greeter) GetById(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Get : [%s] Msg", request.GetName())
	return &pb.HelloResponse{
		Message: "Hello, " + request.Name,
	}, nil
}
