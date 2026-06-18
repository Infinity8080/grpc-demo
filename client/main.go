package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func main() {
	creds := insecure.NewCredentials()
	ctx := context.Background()
	conn, err := grpc.NewClient(":50001", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Error generating new client : %v", err)
	}
	defer conn.Close()
	md := metadata.New(map[string]string{
		"authorization": "Bearer somsecret-token",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)
	runGreeting(ctx, conn)
	// runSum(ctx, conn)
	// runChat(ctx, conn)
}
