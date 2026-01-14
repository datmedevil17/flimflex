package server

import (
	"fmt"
	"net"

	pb "github.com/datmedevil17/micro-flex/proto/recommendation"
	"github.com/datmedevil17/micro-flex/services/recommendation/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	recommendationService *service.RecommendationService
	port                  string
}

func NewGRPCServer(recommendationService *service.RecommendationService, port string) *GRPCServer {
	return &GRPCServer{
		recommendationService: recommendationService,
		port:                  port,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterRecommendationServiceServer(grpcServer, s.recommendationService)
	reflection.Register(grpcServer)

	fmt.Printf("Recommendation Service gRPC server listening on port %s\n", s.port)

	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
