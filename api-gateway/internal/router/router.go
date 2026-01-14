package router

import (
	"github.com/datmedevil17/micro-flex/api-gateway/internal/clients"
	"github.com/datmedevil17/micro-flex/api-gateway/internal/handlers"
	"github.com/datmedevil17/micro-flex/api-gateway/internal/middleware"
	"github.com/datmedevil17/micro-flex/pkg/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(clients *clients.GRPCClients, jwtManager *jwt.JWTManager) *gin.Engine {
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(clients.AuthClient)
	userHandler := handlers.NewUserHandler(clients.UserClient)
	movieHandler := handlers.NewMovieHandler(clients.MovieClient)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Health check
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.GET("/validate", authHandler.ValidateToken)
		}

		// User routes (protected)
		users := v1.Group("/users")
		users.Use(authMiddleware.RequireAuth())
		{
			users.GET("/profile", userHandler.GetProfile)
			users.PUT("/profile", userHandler.UpdateProfile)
			users.GET("/watchlist", userHandler.GetWatchlist)
			users.POST("/watchlist", userHandler.AddToWatchlist)
			users.DELETE("/watchlist/:id", userHandler.RemoveFromWatchlist)
		}

		// Movie routes
		movies := v1.Group("/movies")
		{
			// Public routes
			movies.GET("", movieHandler.ListMovies)
			movies.GET("/:id", movieHandler.GetMovie)
			movies.GET("/search", movieHandler.SearchMovies)
			movies.GET("/genre/:genre", movieHandler.GetMoviesByGenre)

			// Protected routes (Admin only ideally, but using auth for now)
			movies.Use(authMiddleware.RequireAuth())
			movies.POST("", movieHandler.CreateMovie)
			movies.PUT("/:id", movieHandler.UpdateMovie)
			movies.DELETE("/:id", movieHandler.DeleteMovie)
		}

		// Streaming routes (protected)
		streaming := v1.Group("/streaming")
		streaming.Use(authMiddleware.RequireAuth())
		{
			// Will be implemented
		}

		// Recommendations routes
		recommendationHandler := handlers.NewRecommendationHandler(clients.RecommendationClient)
		recommendations := v1.Group("/recommendations")
		recommendations.Use(authMiddleware.RequireAuth())
		{
			recommendations.GET("", recommendationHandler.GetRecommendations)
		}
	}

	return r
}
