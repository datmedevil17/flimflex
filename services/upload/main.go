package main

import (
	"log"
	"net"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"micro-flex/pkg/logger"
	"micro-flex/pkg/response"
)

func main() {
	logger.Init()
	defer logger.Log.Sync()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		response.Success(c, "Upload Service is running", nil)
	})

	go func() {
		lis, err := net.Listen("tcp", ":50054")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		// proto.RegisterUploadServiceServer(s, &server{})
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	r.Run(":8084")
}
