package server

import (
	"context"
	"fmt"
	"wrpc/proto/pb"
)

type UserServer struct {
	pb.UnsafeUserServiceServer
}

func (u *UserServer) GetById(ctx context.Context, req *pb.GetByIdReq) (*pb.GetByIdResp, error) {
	fmt.Printf("get client id %d", req.Id)
	user := &pb.User{
		Status: 100,
		Id:     100,
	}
	return &pb.GetByIdResp{
		User: user,
	}, nil
}
