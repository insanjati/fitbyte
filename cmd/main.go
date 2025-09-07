package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/insanjati/fitbyte/internal/cache"
	"github.com/insanjati/fitbyte/internal/database"
	"github.com/insanjati/fitbyte/internal/handler"
	"github.com/insanjati/fitbyte/internal/middleware"
	"github.com/insanjati/fitbyte/internal/repository"
	"github.com/insanjati/fitbyte/internal/service"
	"github.com/insanjati/fitbyte/internal/storage"

	"github.com/caarlos0/env"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
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

	// MinIO Configuration
	MinIOEndpoint       string `env:"MINIO_ENDPOINT" envDefault:"minio:9000"`
	MinIOAccessKey      string `env:"MINIO_ACCESS_KEY" envDefault:"minioadmin"`
	MinIOSecretKey      string `env:"MINIO_SECRET_KEY" envDefault:"minioadmin"`
	MinIOBucket         string `env:"MINIO_BUCKET" envDefault:"fitbyte-uploads"`
	MinIOPublicEndpoint string `env:"MINIO_PUBLIC_ENDPOINT" envDefault:"http://localhost:9000"`
	MinIOUseSSL         bool   `env:"MINIO_USE_SSL" envDefault:"false"`
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
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		}
	}()

	cache, err := cache.NewRedis(cache.RedisConfig{DB: redisClient})
	if err != nil {
		log.Fatal("Failed to initialize Redis cache:", err)
	}

	// Initialize JWT service
	jwtConfig := &service.SecurityConfig{
		Key:    cfg.JWTSecret,
		Durasi: cfg.JWTDuration,
		Issues: cfg.JWTIssuer,
	}
	jwtService := service.NewJwtService(jwtConfig)

	// Initialize MinIO storage
	minioConfig := &storage.MinIOConfig{
		Endpoint:       cfg.MinIOEndpoint,
		AccessKey:      cfg.MinIOAccessKey,
		SecretKey:      cfg.MinIOSecretKey,
		BucketName:     cfg.MinIOBucket,
		PublicEndpoint: cfg.MinIOPublicEndpoint,
		UseSSL:         cfg.MinIOUseSSL,
	}
	minioStorage, err := storage.NewMinIOStorage(minioConfig)
	if err != nil {
		log.Fatal("Failed to initialize MinIO storage:", err)
	}

	// Initialize health handler
	healthHandler := handler.NewHealthHandler(db, cache)

	// Initialize users layers
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, cache, jwtService)
	userHandler := handler.NewUserHandler(userService)

	// Initialize activities layers
	activityRepo := repository.NewActivityRepository(db)
	activityService := service.NewActivityService(activityRepo, cache)
	activityHandler := handler.NewActivityHandler(activityService)

	// Initialize file handler
	fileHandler := handler.NewFileHandler(minioStorage)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	// Setup Gin router
	r := gin.Default()

	// Routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/healthz", healthHandler.Check)
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

		protected.POST("/file", fileHandler.UploadFile)
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give 30 seconds for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
