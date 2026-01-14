package main

import (
	"fmt"
	"log"
	"time"

	"github.com/datmedevil17/micro-flex/api-gateway/config"
	"github.com/datmedevil17/micro-flex/api-gateway/internal/clients"
	"github.com/datmedevil17/micro-flex/api-gateway/internal/router"
	"github.com/datmedevil17/micro-flex/pkg/jwt"
	"github.com/datmedevil17/micro-flex/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	if err := logger.Init(cfg.Server.Env); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Log.Sync()

	logger.Info("Starting API Gateway...")

	// Initialize JWT manager
	jwtManager := jwt.NewJWTManager(
		cfg.JWT.SecretKey,
		15*time.Minute, // access token expiry
		7*24*time.Hour, // refresh token expiry
	)

	// Initialize gRPC clients
	grpcClients, err := clients.NewGRPCClients(cfg.Services)
	if err != nil {
		logger.Fatal("Failed to initialize gRPC clients", zap.Error(err))
	}

	logger.Info("gRPC clients initialized successfully")

	// Setup router
	r := router.SetupRouter(grpcClients, jwtManager)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Info(fmt.Sprintf("API Gateway listening on %s", addr))

	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
