package main

import (
	"log"

	"go.uber.org/zap"

	"github.com/datmedevil17/micro-flex/pkg/jwt"
	"github.com/datmedevil17/micro-flex/pkg/logger"
	"github.com/datmedevil17/micro-flex/services/auth/config"
	"github.com/datmedevil17/micro-flex/services/auth/db"
	"github.com/datmedevil17/micro-flex/services/auth/internal/models"
	"github.com/datmedevil17/micro-flex/services/auth/internal/repository"
	"github.com/datmedevil17/micro-flex/services/auth/internal/server"
	"github.com/datmedevil17/micro-flex/services/auth/internal/service"
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

	logger.Info("Starting Auth Service...")

	// Connect to database
	database, err := db.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Auto migrate
	if err := database.AutoMigrate(&models.User{}, &models.RefreshToken{}); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	logger.Info("Database connected and migrated successfully")

	// Initialize JWT manager
	jwtManager := jwt.NewJWTManager(
		cfg.JWT.SecretKey,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)

	// Initialize repository
	repo := repository.NewAuthRepository(database)

	// Initialize service
	authService := service.NewAuthService(repo, jwtManager)

	// Start gRPC server
	grpcServer := server.NewGRPCServer(authService, cfg.Server.Port)
	if err := grpcServer.Start(); err != nil {
		logger.Fatal("Failed to start gRPC server", zap.Error(err))
	}
}
