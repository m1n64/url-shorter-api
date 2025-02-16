package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
	"user-service/internal/users/grpc/auth"
	"user-service/internal/users/handlers"
	"user-service/internal/users/workers"
	"user-service/pkg/di"
)

func main() {
	fmt.Println("User service started!")

	dependencies := di.InitDependencies()

	port := os.Getenv("SERVICE_PORT")
	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", port))
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	go workers.StartRemoveExpiredTokensWorker(dependencies.TokenRepo, dependencies.Logger, 24*time.Hour)

	grpcServer := grpc.NewServer()
	authServiceServer := handlers.NewAuthServiceServer(dependencies.AuthService, dependencies.Logger)
	auth.RegisterAuthServiceServer(grpcServer, authServiceServer)

	log.Printf("gRPC server is running on port %s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
