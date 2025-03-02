package main

import (
	"analytics-service/pkg/di"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
)

var dependencies *di.Dependencies

func init() {
	dependencies = di.InitDependencies()
}

func main() {
	fmt.Println("Analytics Service started!")

	port := os.Getenv("SERVICE_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	go di.InitializeQueuesConsumer(dependencies)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(tokenAuthInterceptor),
	)

	log.Printf("gRPC server is running on port %s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

func tokenAuthInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	authHeader, exists := md["authorization"]
	if !exists || len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	expectedToken := os.Getenv("GRPC_TOKEN")
	if authHeader[0] != expectedToken {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}
