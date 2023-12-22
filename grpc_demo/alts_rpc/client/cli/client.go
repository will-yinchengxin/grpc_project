package cli

import (
	"context"
	"log"
	"wrpc/proto/pb"
)

func Greeter(conn pb.GreeterServiceClient) error {
	msg, err := conn.GetById(context.Background(), &pb.HelloRequest{Name: "will"})
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("Get Server Msg: %s", msg.GetMessage())
	return nil
}
