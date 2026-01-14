package main

import (
	"log"

	"github.com/datmedevil17/micro-flex/pkg/logger"
	"github.com/datmedevil17/micro-flex/services/upload/config"
	"github.com/datmedevil17/micro-flex/services/upload/db"
	"github.com/datmedevil17/micro-flex/services/upload/internal/models"
	"github.com/datmedevil17/micro-flex/services/upload/internal/repository"
	"github.com/datmedevil17/micro-flex/services/upload/internal/server"
	"github.com/datmedevil17/micro-flex/services/upload/internal/service"
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

	logger.Info("Starting Upload Service...")

	database, err := db.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := database.AutoMigrate(&models.File{}); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	logger.Info("Database connected and migrated successfully")

	repo := repository.NewUploadRepository(database)

	uploadService, err := service.NewUploadService(cfg.Cloudinary, repo)
	if err != nil {
		logger.Fatal("Failed to initialize upload service", zap.Error(err))
	}

	logger.Info("Cloudinary initialized successfully")

	grpcServer := server.NewGRPCServer(uploadService, cfg.Server.Port)
	if err := grpcServer.Start(); err != nil {
		logger.Fatal("Failed to start gRPC server", zap.Error(err))
	}
}
