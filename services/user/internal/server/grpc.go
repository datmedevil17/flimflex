package server

import (
	"fmt"
	"net"

	pb "github.com/datmedevil17/micro-flex/proto/user"
	"github.com/datmedevil17/micro-flex/services/user/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	userService *service.UserService
	port        string
}

func NewGRPCServer(userService *service.UserService, port string) *GRPCServer {
	return &GRPCServer{
		userService: userService,
		port:        port,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	
	pb.RegisterUserServiceServer(grpcServer, s.userService)
	
	// Register reflection service for development
	reflection.Register(grpcServer)

	fmt.Printf("User Service gRPC server listening on port %s\n", s.port)
	
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}