package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	pb "github.com/Infinity8080/grpc-demo/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedGreetingServiceServer
	pb.UnimplementedCountdownServiceServer
	pb.UnimplementedSumServiceServer
	pb.UnimplementedChatServiceServer
}

func (s *Server) Chat(stream grpc.BidiStreamingServer[pb.ChatRequest, pb.ChatResponse]) error {

	for {
		rec, err := stream.Recv()
		if err == io.EOF {
			log.Print("End of streaming server side")
			break
		}
		if err != nil {
			return status.Error(codes.Aborted, "Server is not able to recieve the stream")
		}
		log.Printf("Recieved: %v", rec.IncomingMessage)
		err = stream.Send(&pb.ChatResponse{
			OutgoingMessage: rec.IncomingMessage,
		})
		if err != nil {
			return status.Error(codes.Aborted, "Server is not able to send the stream")
		}
	}
	return nil
}

func (s *Server) Sum(stream grpc.ClientStreamingServer[pb.SumRequest, pb.SumResponse]) error {
	sum := pb.SumResponse{
		Sum: 0,
	}
	for {
		st, err := stream.Recv()
		if err == io.EOF {
			log.Printf("End of streaming from client!")
			return stream.SendAndClose(&sum)
		}

		if err != nil {
			return err
		}
		log.Print("Got: ", st.Number)

		sum.Sum += st.Number
	}
}

func (s *Server) CountdownNumber(req *pb.CountdownRequest, stream grpc.ServerStreamingServer[pb.CountdownResponse]) error {
	countdownNumber := req.StartingNumber
	for i := countdownNumber; i > 0; i-- {
		err := stream.Send(
			&pb.CountdownResponse{
				CurrentNumber: i,
			},
		)
		if err != nil {
			log.Printf("Error while sending stream")
			return status.Error(codes.Aborted, "Error while sending stream!")

		}
		time.Sleep(time.Second)
	}
	return nil
}

func (s *Server) GreetUser(ctx context.Context, req *pb.GreetingRequest) (*pb.GreetingResponse, error) {
	log.Print("Recived the request")
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Name should not be empty!")
	}
	log.Print("Returning response")
	return &pb.GreetingResponse{
		GreetingMessage: fmt.Sprintf("Hello %v", req.Name),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50001")
	if err != nil {
		log.Fatalf("Error running TCP server %v", err)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(AuthInterceptor))
	pb.RegisterGreetingServiceServer(grpcServer, &Server{})
	pb.RegisterCountdownServiceServer(grpcServer, &Server{})
	pb.RegisterSumServiceServer(grpcServer, &Server{})
	pb.RegisterChatServiceServer(grpcServer, &Server{})
	log.Println("Service starting on localhost:50001")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to server %v", err)
	}
}
