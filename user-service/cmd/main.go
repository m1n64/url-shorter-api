package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
	"time"
	"user-service/internal/users/grpc/auth"
	"user-service/internal/users/handlers"
	"user-service/internal/users/workers"
	"user-service/pkg/di"
)

var dependencies *di.Dependencies

func init() {
	dependencies = di.InitDependencies()
}

func main() {
	fmt.Println("User service started!")

	port := os.Getenv("SERVICE_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	go workers.StartRemoveExpiredTokensWorker(dependencies.TokenRepo, dependencies.Logger, 24*time.Hour)

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(tokenAuthInterceptor),
	)

	authServiceServer := handlers.NewAuthServiceServer(dependencies.AuthService, dependencies.Logger)
	auth.RegisterAuthServiceServer(grpcServer, authServiceServer)

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
