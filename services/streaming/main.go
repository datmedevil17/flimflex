package main

import (
	"github.com/gin-gonic/gin"

	"micro-flex/pkg/logger"
	"micro-flex/pkg/response"
)

func main() {
	logger.Init()
	defer logger.Log.Sync()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		response.Success(c, "Streaming Service is running", nil)
	})

	// Streaming logic typically doesn't need gRPC in the same way, but can if needed.
	// Just HTTP serving for now.

	r.Run(":8085")
}
