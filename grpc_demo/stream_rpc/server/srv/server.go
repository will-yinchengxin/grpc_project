package srv

import (
	"io"
	"log"
	"time"
	"wrpc/proto/pb"
)

type ChatServer struct {
	pb.UnsafeStreamChatServer
}

func (c ChatServer) Chat(stream pb.StreamChat_ChatServer) error {
	var (
		startTime = time.Now()
	)
	for {
		recive, err := stream.Recv()
		if err == io.EOF {
			log.Printf("All Chat Use %d Second", time.Now().Sub(startTime))
			return nil
		}
		if err != nil {
			log.Fatalf("Recive Get Err  %v", err)
			return err
		}

		log.Printf("Got Client message [%s] at point (%d, %d)", recive.Message, recive.Location.Latitude, recive.Location.Longitude)

		stream.Send(&pb.RouteNote{
			Message:  "Server " + recive.GetMessage(),
			Location: &pb.Point{Latitude: int32(time.Now().Sub(startTime)), Longitude: int32(time.Now().Sub(startTime))},
		})
	}
	return nil
}
