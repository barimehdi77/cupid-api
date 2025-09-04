// cmd/api/main.go
package main

import (
	"fmt"
	"os"

	"github.com/barimehdi77/cupid-api/internal/database"
	"github.com/barimehdi77/cupid-api/internal/env"
	"github.com/barimehdi77/cupid-api/internal/logger"
	"github.com/barimehdi77/cupid-api/internal/store"
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

	// Create application instance with dependencies
	app := &application{
		config: config{
			port: env.GetEnvInt("SERVER_PORT", 8080),
			env:  env.GetEnvString("GO_ENV", "development"),
		},
		logger:  logger.Logger,
		storage: storage,
	}

	// Start the server
	if err := app.run(); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
