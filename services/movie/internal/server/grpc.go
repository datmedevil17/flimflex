package server

import (
	"fmt"
	"net"

	pb "github.com/datmedevil17/micro-flex/proto/movie"
	"github.com/datmedevil17/micro-flex/services/movie/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	movieService *service.MovieService
	port         string
}

func NewGRPCServer(movieService *service.MovieService, port string) *GRPCServer {
	return &GRPCServer{
		movieService: movieService,
		port:         port,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer()
	
	pb.RegisterMovieServiceServer(grpcServer, s.movieService)
	
	reflection.Register(grpcServer)

	fmt.Printf("Movie Service gRPC server listening on port %s\n", s.port)
	
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
