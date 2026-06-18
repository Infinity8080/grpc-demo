package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/Infinity8080/grpc-demo/proto"
	"google.golang.org/grpc"
)

func runChat(ctx context.Context, conn *grpc.ClientConn) {
	client := pb.NewChatServiceClient(conn)
	stream, err := client.Chat(ctx)
	if err != nil {
		log.Fatalf("Client did not connect %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			rec, err := stream.Recv()
			if err == io.EOF {
				log.Print("End of receiving!")
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed in recieved in client side %v", err)
			}
			log.Print("Response from recieving:", rec.OutgoingMessage)
		}
	}()
	texts := []string{
		"Hello",
		"jatin",
		"Sharaya",
	}

	for _, val := range texts {
		log.Printf("Sending from client %v", val)
		err := stream.Send(&pb.ChatRequest{
			IncomingMessage: val,
		})
		if err != nil {
			log.Fatalf("Failing to send stream from client %v", err)

		}
		time.Sleep(time.Second)
	}
	stream.CloseSend()
	<-waitc
	log.Println("Finished streaming.")
}
