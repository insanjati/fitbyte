package main

import (
	"log"
	"net/http"
	"time"

	"github.com/insanjati/fitbyte/internal/database"
	"github.com/insanjati/fitbyte/internal/handler"
	"github.com/insanjati/fitbyte/internal/middleware"
	"github.com/insanjati/fitbyte/internal/repository"
	"github.com/insanjati/fitbyte/internal/service"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	HTTPPort    string `env:"HTTP_PORT" envDefault:"8080"`
	DatabaseURL string `env:"DATABASE_URL"`

	// JWT Configuration
	JWTSecret   string        `env:"JWT_SECRET" envDefault:"your-secret-key"`
	JWTDuration time.Duration `env:"JWT_DURATION" envDefault:"24h"`
	JWTIssuer   string        `env:"JWT_ISSUER" envDefault:"fitbyte-app"`
}

func main() {
	// Load config
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize JWT service
	jwtConfig := &service.SecurityConfig{
		Key:    cfg.JWTSecret,
		Durasi: cfg.JWTDuration,
		Issues: cfg.JWTIssuer,
	}
	jwtService := service.NewJwtService(jwtConfig)

	// Initialize users layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, jwtService)
	userHandler := handler.NewUserHandler(userService)

	// Initialize activities layers
	activityRepo := repository.NewActivityRepository(db)
	activityService := service.NewActivityService(activityRepo)
	activityHandler := handler.NewActivityHandler(activityService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Setup Gin router
	r := gin.Default()

	// Routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", healthCheckHandler(db))

		v1.POST("/register", userHandler.CreateNewUser)
		v1.POST("/login", userHandler.Login)
	}

	protected := v1.Group("/")
	protected.Use(authMiddleware.CheckToken())
	{
		protected.PATCH("/users", userHandler.UpdateUser)
		protected.GET("/users", userHandler.GetUsers)

		protected.POST("/activity", activityHandler.CreateActivity)
		protected.GET("/activity", activityHandler.GetUserActivities)
		protected.PATCH("/activity/:activityId", activityHandler.UpdateActivity)
		protected.DELETE("/activity/:activityId", activityHandler.DeleteActivity)
	}

	// Start server
	log.Printf("Server starting on port %s", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func healthCheckHandler(db *database.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check database connectivity
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "down",
				"database": "disconnected",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":   "up",
			"database": "connected",
		})
	}
}
