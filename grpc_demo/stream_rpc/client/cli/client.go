package cli

import (
	"context"
	"io"
	"log"
	"time"
	"wrpc/proto/pb"
)

func Chat(client pb.StreamChatClient) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	stream, err := client.Chat(ctx)
	if err != nil {
		log.Fatalf("client.RouteChat failed: %v", err)
		return
	}

	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("client.RouteChat failed: %v", err)
				return
			}
			log.Printf("Got Server message [%s] at point (%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()
	notes := []*pb.RouteNote{
		{Location: &pb.Point{Latitude: 0, Longitude: 1}, Message: "First msg"},
		{Location: &pb.Point{Latitude: 0, Longitude: 2}, Message: "Second msg"},
		{Location: &pb.Point{Latitude: 0, Longitude: 3}, Message: "Third msg"},
		{Location: &pb.Point{Latitude: 0, Longitude: 4}, Message: "Fourth msg"},
		{Location: &pb.Point{Latitude: 0, Longitude: 5}, Message: "Fifth msg"},
		{Location: &pb.Point{Latitude: 0, Longitude: 6}, Message: "Sixth msg"},
	}
	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			log.Fatalf("client.RouteChat: stream.Send(%v) failed: %v", note, err)
		}
	}
	stream.CloseSend()
	<-waitc
}
