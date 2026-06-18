package main

import (
	"context"
	"log"

	pb "github.com/Infinity8080/grpc-demo/proto"
	"google.golang.org/grpc"
)

func runGreeting(ctx context.Context, conn *grpc.ClientConn) {
	client := pb.NewGreetingServiceClient(conn)

	resp, err := client.GreetUser(ctx, &pb.GreetingRequest{
		Name: "Apappasdf",
	})
	if err != nil {
		log.Printf("Error getting response %v", err)
		return
	}

	log.Print(resp.GreetingMessage)
}
