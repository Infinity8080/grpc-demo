package main

import (
	"context"
	"log"

	pb "github.com/Infinity8080/grpc-demo/proto"
	"google.golang.org/grpc"
)

func runSum(ctx context.Context, conn *grpc.ClientConn) {
	client := pb.NewSumServiceClient(conn)
	stream, err := client.Sum(ctx)
	if err != nil {
		log.Printf("Error in  client")
	}
	for i := range 5 {
		err := stream.Send(&pb.SumRequest{
			Number: int32(i),
		})
		if err != nil {
			log.Print("Error duing client stream: ", err)
			break
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Println("Error during closeing: ", err)
	}
	log.Print(res)
}
