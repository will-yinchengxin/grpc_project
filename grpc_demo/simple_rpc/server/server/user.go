package server

import (
	"context"
	"google.golang.org/grpc/metadata"
	"log"
	"wrpc/proto/pb"
)

type MetaData struct {
	pb.UnimplementedPassParamServer
}

func (u *MetaData) SendParam(ctx context.Context, req *pb.PassParamRequest) (*pb.PassParamResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Fatal("Must Send MetaData")
	}
	for key, val := range md {
		log.Print("--> ", key, val, " <--")
	}

	log.Printf("Get Client %s Msg", req.GetName())
	return &pb.PassParamResponse{
		Message: "Hello " + req.Name,
	}, nil
}
