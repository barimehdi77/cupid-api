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
	"github.com/barimehdi77/cupid-api/internal/api"
	"github.com/barimehdi77/cupid-api/internal/logger"
	"github.com/barimehdi77/cupid-api/internal/store"
	"github.com/barimehdi77/cupid-api/internal/sync"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

type application struct {
	config      config
	logger      *zap.Logger
	storage     store.Storage
	handlers    *api.Handlers
	syncService *sync.SyncService
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

	// Add enhanced logging middleware
	r.Use(logger.GinMiddleware())         // Enhanced HTTP request logging
	r.Use(logger.GinRecoveryMiddleware()) // Enhanced panic recovery logging

	// Initialize Swagger docs
	docs.SwaggerInfo.BasePath = "/api/v1"

	// Create handlers
	app.handlers = api.NewHandlers(app.storage)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Health check routes
		v1.GET("/health", app.handlers.HealthCheckHandler)

		// Property routes
		v1.GET("/properties", app.handlers.ListPropertiesHandler)
		v1.GET("/properties/:id", app.handlers.GetPropertyHandler)
		v1.GET("/properties/:id/reviews", app.handlers.GetPropertyReviewsHandler)
		v1.GET("/properties/:id/translations", app.handlers.GetPropertyTranslationsHandler)
		v1.GET("/properties/location", app.handlers.GetPropertiesByLocationHandler)
		v1.GET("/properties/rating", app.handlers.GetPropertiesByRatingHandler)

		// Search routes
		v1.GET("/search", app.handlers.SearchPropertiesHandler)

		// Admin sync routes (only if sync service is available)
		if app.syncService != nil {
			syncHandlers := api.NewSyncHandlers(app.syncService)
			admin := v1.Group("/admin")
			{
				admin.POST("/sync", syncHandlers.TriggerSyncHandler)
				admin.GET("/sync/status", syncHandlers.GetSyncStatusHandler)
				admin.POST("/sync/start", syncHandlers.StartSyncHandler)
				admin.POST("/sync/stop", syncHandlers.StopSyncHandler)
				admin.GET("/sync/logs", syncHandlers.GetSyncLogsHandler)
				admin.GET("/sync/settings", syncHandlers.GetSyncSettingsHandler)
				admin.PUT("/sync/settings", syncHandlers.UpdateSyncSettingsHandler)
				admin.GET("/sync/health", syncHandlers.GetSyncHealthHandler)
			}
		}
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
		logger.LogStartup("HTTP Server",
			zap.Int("port", app.config.port),
			zap.String("environment", app.config.env),
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.LogError("Server startup", err)
		}
	}()

	// Wait for interrupt signal
	<-shutdown
	logger.LogShutdown("HTTP Server", zap.String("reason", "interrupt signal received"))

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		logger.LogError("Graceful shutdown", err)
		return err
	}

	logger.LogSuccess("Server shutdown")
	return nil
}
