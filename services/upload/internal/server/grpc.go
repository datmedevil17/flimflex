package server

import (
	"fmt"
	"net"

	pb "github.com/datmedevil17/micro-flex/proto/upload"
	"github.com/datmedevil17/micro-flex/services/upload/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	uploadService *service.UploadService
	port          string
}

func NewGRPCServer(uploadService *service.UploadService, port string) *GRPCServer {
	return &GRPCServer{
		uploadService: uploadService,
		port:          port,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.MaxRecvMsgSize(50 * 1024 * 1024), // 50MB for video uploads
		grpc.MaxSendMsgSize(50 * 1024 * 1024),
	)
	
	pb.RegisterUploadServiceServer(grpcServer, s.uploadService)
	reflection.Register(grpcServer)

	fmt.Printf("Upload Service gRPC server listening on port %s\n", s.port)
	
	if err := grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
