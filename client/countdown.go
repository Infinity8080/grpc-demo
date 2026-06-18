package main

import (
	"context"
	"io"
	"log"

	pb "github.com/Infinity8080/grpc-demo/proto"
	"google.golang.org/grpc"
)

func runCountdown(ctx context.Context, conn *grpc.ClientConn) {
	client := pb.NewCountdownServiceClient(conn)
	stream, err := client.CountdownNumber(ctx, &pb.CountdownRequest{
		StartingNumber: 10,
	})
	if err != nil {
		log.Printf("Error in stream client")
	}
	for {
		rec, err := stream.Recv()
		if err == io.EOF {
			log.Printf("End of streaming!")
			break
		}
		if err != nil {
			log.Fatalf("Error while recieving the stream %v", err)
			break
		}
		log.Print(rec.CurrentNumber)
	}
}
