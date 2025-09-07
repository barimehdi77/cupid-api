// @title           Cupid Hotel API
// @version         1.0
// @description     A comprehensive hotel property API that fetches and serves hotel data from Cupid API with reviews, translations, and search capabilities
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/barimehdi77/cupid-api/internal/cupid"
	"github.com/barimehdi77/cupid-api/internal/database"
	"github.com/barimehdi77/cupid-api/internal/env"
	"github.com/barimehdi77/cupid-api/internal/logger"
	"github.com/barimehdi77/cupid-api/internal/store"
	"github.com/barimehdi77/cupid-api/internal/sync"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Printf("Warning: Could not load .env file: %v\n", err)
	}

	// Initialize logger
	if err := logger.InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Initialize database
	db, err := database.NewDB()
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize storage
	storage := store.NewStorage(db)

	// Create sync service
	cupidService := cupid.NewService()
	syncConfig := sync.DefaultConfig()
	syncService := sync.NewSyncService(cupidService, storage, syncConfig)

	// Create application instance with dependencies
	app := &application{
		config: config{
			port: env.GetEnvInt("SERVER_PORT", 8080),
			env:  env.GetEnvString("GO_ENV", "development"),
		},
		logger:      logger.Logger,
		storage:     storage,
		syncService: syncService,
	}

	// Start the sync service
	ctx := context.Background()
	if err := app.syncService.Start(ctx); err != nil {
		logger.LogError("Failed to start sync service", err)
		// Don't exit, just log the error and continue
	}

	// Start the server
	if err := app.run(); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
