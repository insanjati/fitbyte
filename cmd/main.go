package main

import (
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	"github.com/insanjati/fitbyte/internal/database"
	"github.com/insanjati/fitbyte/internal/handler"
	"github.com/insanjati/fitbyte/internal/middleware"
	"github.com/insanjati/fitbyte/internal/repository"
	"github.com/insanjati/fitbyte/internal/service"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	HTTPPort    string `env:"HTTP_PORT" envDefault:"8080"`
	DatabaseURL string `env:"DATABASE_URL"`

	// JWT Configuration
	JWTSecret   string        `env:"JWT_SECRET" envDefault:"your-secret-key"`
	JWTDuration time.Duration `env:"JWT_DURATION" envDefault:"24h"`
	JWTIssuer   string        `env:"JWT_ISSUER" envDefault:"fitbyte-app"`

	// Redis Configuration
	RedisAddr     string `env:"REDIS_ADDR" envDefault:"redis:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" envDefault:""`
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`
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

	// Initialize Redis cache
	cacheRepo := repository.NewRedisCache(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)

	// Initialize JWT service
	jwtConfig := &service.SecurityConfig{
		Key:    cfg.JWTSecret,
		Durasi: cfg.JWTDuration,
		Issues: cfg.JWTIssuer,
	}
	jwtService := service.NewJwtService(jwtConfig)

	// Initialize users layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cacheRepo, jwtService)
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
		v1.GET("/health", healthCheckHandler(db, cacheRepo))
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

	log.Printf("Server starting on port %s", cfg.HTTPPort)
	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func healthCheckHandler(db *database.DB, cache repository.CacheRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "down",
				"database": "disconnected",
				"cache":    "unknown",
			})
			return
		}

		ctx := c.Request.Context()
		if err := cache.Set(ctx, "health_check", "ok", 10*time.Second); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "degraded",
				"database": "connected",
				"cache":    "disconnected",
			})
			return
		}

		_ = cache.Delete(ctx, "health_check")

		c.JSON(http.StatusOK, gin.H{
			"status":   "up",
			"database": "connected",
			"cache":    "connected",
		})
	}
}
