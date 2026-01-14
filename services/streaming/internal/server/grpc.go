package server

import (
	"fmt"
	"net"

	pb "github.com/datmedevil17/micro-flex/proto/streaming"
	"github.com/datmedevil17/micro-flex/services/streaming/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	streamingService *service.StreamingService
	port             string
}

func NewGRPCServer(streamingService *service.StreamingService, port string) *GRPCServer {
	return &GRPCServer{
		streamingService: streamingService,
		port:             port,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterStreamingServiceServer(grpcServer, s.streamingService)
	reflection.Register(grpcServer)

	fmt.Printf("Streaming Service gRPC server listening on port %s\n", s.port)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
