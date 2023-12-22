package srv

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	"grpc_product/internal/mysql/entity"
	"grpc_product/model/dao"
	"grpc_product/proto/pb"
	"log"
)

type ProductServer struct {
	pb.UnimplementedProductServiceServer
}

func (p ProductServer) AdvertiseList(ctx context.Context, empty *emptypb.Empty) (*pb.AdvertisesRes, error) {
	var adItemList []*pb.AdvertiseItemRes
	var advertiseRes pb.AdvertisesRes
	aModel, err := dao.AdvertiseWithContext(context.Background())
	if err != nil {
		return nil, err
	}
	adList, count, err := aModel.List()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	for _, item := range adList {
		adItemList = append(adItemList, ConvertAdModel2Pb(item))
	}
	advertiseRes.Total = int32(count)
	advertiseRes.ItemList = adItemList
	return &advertiseRes, nil
}

func (p ProductServer) CreateAdvertise(ctx context.Context, req *pb.AdvertiseReq) (*pb.AdvertiseItemRes, error) {
	panic("implicate me")
}

func (p ProductServer) DeleteAdvertise(ctx context.Context, req *pb.AdvertiseReq) (*emptypb.Empty, error) {
	panic("implicate me")
}

func (p ProductServer) UpdateAdvertise(ctx context.Context, req *pb.AdvertiseReq) (*emptypb.Empty, error) {
	panic("implicate me")
}

func ConvertAdModel2Pb(item entity.Advertise) *pb.AdvertiseItemRes {
	ad := &pb.AdvertiseItemRes{
		Index:  item.Index,
		Images: item.Image,
		Url:    item.Url,
	}
	if item.ID > 0 {
		ad.Id = int32(item.ID)
	}
	return ad
}
