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
		response.Success(c, "API Gateway is running", nil)
	})

	// TODO: Implement gateway logic to proxy requests to services

	r.Run(":8080")
}
