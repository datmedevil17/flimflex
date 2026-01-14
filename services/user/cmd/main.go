package main

import (
	"log"

	"go.uber.org/zap"

	"github.com/datmedevil17/micro-flex/pkg/logger"
	"github.com/datmedevil17/micro-flex/services/user/config"
	"github.com/datmedevil17/micro-flex/services/user/db"
	"github.com/datmedevil17/micro-flex/services/user/internal/models"
	"github.com/datmedevil17/micro-flex/services/user/internal/repository"
	"github.com/datmedevil17/micro-flex/services/user/internal/server"
	"github.com/datmedevil17/micro-flex/services/user/internal/service"
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

	logger.Info("Starting User Service...")

	// Connect to database
	database, err := db.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto migrate
	if err := database.AutoMigrate(&models.Profile{}, &models.WatchHistory{}, &models.Watchlist{}); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	logger.Info("Database connected and migrated successfully")

	// Initialize repository
	repo := repository.NewUserRepository(database)

	// Initialize service
	userService := service.NewUserService(repo)

	// Start gRPC server
	grpcServer := server.NewGRPCServer(userService, cfg.Server.Port)
	if err := grpcServer.Start(); err != nil {
		logger.Fatal("Failed to start gRPC server", zap.Error(err))
	}
}
