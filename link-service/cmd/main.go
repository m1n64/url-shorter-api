package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	links "link-service/internal/links/grpc"
	"link-service/internal/links/handlers"
	"link-service/internal/links/workers"
	"link-service/pkg/di"
	"log"
	"net"
	"os"
)

var dependencies *di.Dependencies

func init() {
	dependencies = di.InitDependencies()

	go workers.RevalidateLinksInCache(dependencies.LinkService)
}

func main() {
	fmt.Println("Link service started!")

	port := os.Getenv("SERVICE_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(tokenAuthInterceptor),
	)

	linkServer := handlers.NewLinkServiceServer(dependencies.LinkService, dependencies.SlugService, dependencies.Validator)
	links.RegisterLinkServiceServer(grpcServer, linkServer)

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
