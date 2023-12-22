package server

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
	"wrpc/proto/pb"
)

type Interceptor struct {
	pb.UnimplementedInterceptorServiceServer
}

func (i *Interceptor) Interceptor(ctx context.Context, req *pb.InterceptorRequest) (*pb.InterceptorResponse, error) {
	// 测试客户端响应超时
	//time.Sleep(time.Second * 5)

	// 测试服务端返回错误码
	//return nil, status.Error(codes.Internal, "server error")
	return nil, status.Error(codes.Unknown, "未知错误")

	log.Printf("Get Client %s Msg", req.GetName())
	return &pb.InterceptorResponse{
		Message: "Hello " + req.Name,
	}, nil
}

func MyInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	now := time.Now()
	resp, err = handler(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	time := time.Now().Sub(now)
	log.Printf("interceptor time: %v", time)
	return
}
