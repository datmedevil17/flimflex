package server

import (
	"fmt"
	"net"

	pb "github.com/datmedevil17/micro-flex/proto/auth"
	"github.com/datmedevil17/micro-flex/services/auth/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	authService *service.AuthService
	port        string
}

func NewGRPCServer(authService *service.AuthService, port string) *GRPCServer {
	return &GRPCServer{
		authService: authService,
		port:        port,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	
	pb.RegisterAuthServiceServer(grpcServer, s.authService)
	
	// Register reflection service for development
	reflection.Register(grpcServer)

	fmt.Printf("Auth Service gRPC server listening on port %s\n", s.port)
	
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}