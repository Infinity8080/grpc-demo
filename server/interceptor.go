package main

import (
	"context"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	errorDenied = status.Error(codes.PermissionDenied, "Context is not ok coming from the metadata")
)

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	return token == "some-secret-token"
}
func AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errorDenied
	}
	if !valid(md["authorization"]) {
		return nil, errorDenied
	}

	log.Print(info.FullMethod)
	m, err := handler(ctx, req)
	if err != nil {
		return nil, status.Error(codes.Aborted, "RPC function call didn't work in the auth interceptor!")

	}
	return m, nil

}
