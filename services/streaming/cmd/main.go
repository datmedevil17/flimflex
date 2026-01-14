package main

import (
	"log"

	"github.com/datmedevil17/micro-flex/pkg/logger"
	"github.com/datmedevil17/micro-flex/services/streaming/config"
	"github.com/datmedevil17/micro-flex/services/streaming/db"
	"github.com/datmedevil17/micro-flex/services/streaming/internal/models"
	"github.com/datmedevil17/micro-flex/services/streaming/internal/repository"
	"github.com/datmedevil17/micro-flex/services/streaming/internal/server"
	"github.com/datmedevil17/micro-flex/services/streaming/internal/service"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := logger.Init(cfg.Server.Env); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Log.Sync()

	logger.Info("Starting Streaming Service...")

	database, err := db.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := database.AutoMigrate(&models.StreamSession{}, &models.StreamAnalytics{}); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	logger.Info("Database connected and migrated successfully")

	repo := repository.NewStreamingRepository(database)
	streamingService := service.NewStreamingService(repo, cfg.Streaming.SecretKey)

	grpcServer := server.NewGRPCServer(streamingService, cfg.Server.Port)
	if err := grpcServer.Start(); err != nil {
		logger.Fatal("Failed to start gRPC server", zap.Error(err))
	}
}
