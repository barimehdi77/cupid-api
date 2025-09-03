package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/barimehdi77/cupid-api/internal/database"
	"github.com/barimehdi77/cupid-api/internal/env"
	"github.com/barimehdi77/cupid-api/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		// Use standard log for this since logger isn't initialized yet
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	// Initialize logger
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Ensure logger is synced on exit
	defer logger.Sync()

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.Info("Shutting down gracefully...")
		logger.Sync()
		os.Exit(0)
	}()

	logger.Info("Starting application...")

	port := env.GetEnvInt("SERVER_PORT", 8080)

	// Set Gin mode based on environment
	env := os.Getenv("GO_ENV")
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin engine without default middleware
	r := gin.New()

	// Add custom Zap middleware
	r.Use(logger.GinMiddleware())
	r.Use(logger.GinRecoveryMiddleware())

	// Connect to database
	database.ConnectDatabase()

	// Example route
	r.GET("/ping", func(c *gin.Context) {
		logger.Info("Ping endpoint called",
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	logger.Info("Server starting",
		zap.Int("port", port),
		zap.String("environment", env),
	)

	// Start server
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
