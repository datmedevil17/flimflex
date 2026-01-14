package main

import (
	"context"
	"log"

	"github.com/datmedevil17/micro-flex/pkg/logger"
	"github.com/datmedevil17/micro-flex/services/recommendation/config"
	"github.com/datmedevil17/micro-flex/services/recommendation/db"
	"github.com/datmedevil17/micro-flex/services/recommendation/internal/models"
	"github.com/datmedevil17/micro-flex/services/recommendation/internal/repository"
	"github.com/datmedevil17/micro-flex/services/recommendation/internal/server"
	"github.com/datmedevil17/micro-flex/services/recommendation/internal/service"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load .env file from root
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("No .env file found or failed to load, using default env variables")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := logger.Init(cfg.Server.Env); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Log.Sync()

	logger.Info("Starting Recommendation Service...")

	database, err := db.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := database.AutoMigrate(&models.UserPreference{}, &models.TrendingCache{}); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	logger.Info("Database connected and migrated successfully")

	repo := repository.NewRecommendationRepository(database)
	recommendationService, err := service.NewRecommendationService(context.Background(), repo, cfg.Gemini.APIKey)
	if err != nil {
		logger.Fatal("Failed to initialize recommendation service", zap.Error(err))
	}

	grpcServer := server.NewGRPCServer(recommendationService, cfg.Server.Port)
	if err := grpcServer.Start(); err != nil {
		logger.Fatal("Failed to start gRPC server", zap.Error(err))
	}
}
