package main

import (
	"log"

	"github.com/datmedevil17/micro-flex/pkg/logger"
	"github.com/datmedevil17/micro-flex/services/movie/config"
	"github.com/datmedevil17/micro-flex/services/movie/db"
	"github.com/datmedevil17/micro-flex/services/movie/internal/models"
	"github.com/datmedevil17/micro-flex/services/movie/internal/repository"
	"github.com/datmedevil17/micro-flex/services/movie/internal/server"
	"github.com/datmedevil17/micro-flex/services/movie/internal/service"
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

	logger.Info("Starting Movie Service...")

	database, err := db.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	if err := database.AutoMigrate(&models.Movie{}); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	logger.Info("Database connected and migrated successfully")

	repo := repository.NewMovieRepository(database)
	movieService := service.NewMovieService(repo)

	grpcServer := server.NewGRPCServer(movieService, cfg.Server.Port)
	if err := grpcServer.Start(); err != nil {
		logger.Fatal("Failed to start gRPC server", zap.Error(err))
	}
}
