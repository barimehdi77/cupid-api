// cmd/api/api.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/barimehdi77/cupid-api/docs"
	"github.com/barimehdi77/cupid-api/internal/store"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type application struct {
	config  config
	logger  *zap.Logger
	storage *store.Storage
}

type config struct {
	port int
	env  string
}

// mount configures all routes, middleware, and handlers
func (app *application) mount() *gin.Engine {
	// Set Gin mode based on environment
	if app.config.env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin engine without default middleware
	r := gin.New()

	// Add custom middleware
	r.Use(app.ginLogger())   // Custom logger middleware
	r.Use(app.ginRecovery()) // Custom recovery middleware

	// Initialize Swagger docs
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Health check routes
	r.GET("/health", app.healthHandler)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", app.pingHandler)
		// Add more routes here as your API grows
	}

	// Swagger endpoint
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

// run starts the server and handles graceful shutdown
func (app *application) run() error {
	// Mount routes
	router := app.mount()

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	// Channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		app.logger.Info("Server starting",
			zap.Int("port", app.config.port),
			zap.String("environment", app.config.env),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-shutdown
	app.logger.Info("Shutting down server gracefully...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		app.logger.Error("Server forced to shutdown", zap.Error(err))
		return err
	}

	app.logger.Info("Server stopped")
	return nil
}

// Handler methods (receiver functions on app)
func (app *application) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"info": gin.H{
			"environment": app.config.env,
		},
	})
}

func (app *application) pingHandler(c *gin.Context) {
	app.logger.Info("Ping endpoint called",
		zap.String("ip", c.ClientIP()),
		zap.String("user_agent", c.Request.UserAgent()),
	)
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// Custom middleware methods
func (app *application) ginLogger() gin.HandlerFunc {
	return gin.LoggerWithWriter(gin.DefaultWriter)
	// You can customize this to use your zap logger
}

func (app *application) ginRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		app.logger.Error("Panic recovered",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.Any("panic", recovered),
		)
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
