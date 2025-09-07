package main

import (
	"context"
	"fmt"
	"os"

	"github.com/barimehdi77/cupid-api/internal/cupid"
	"github.com/barimehdi77/cupid-api/internal/database"
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

	logger.LogStartup("Cupid API Data Fetcher")

	// Create context
	ctx := context.Background()

	// Initialize database
	db, err := database.NewDB()
	if err != nil {
		logger.LogError("Failed to connect to database", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create storage
	storage := store.NewStorage(db)

	// Create service
	service := cupid.NewService()

	// Fetch all properties
	properties, err := service.FetchAllProperties(ctx)
	if err != nil {
		logger.LogError("Failed to fetch properties", err)
		os.Exit(1)
	}

	logger.LogSuccess("Data fetching completed",
		zap.Int("total_properties", len(properties)),
	)

	// Store properties in database
	successCount := 0
	errorCount := 0

	for i, propertyData := range properties {
		logger.LogProgress("Storing property",
			zap.Int("current", i+1),
			zap.Int("total", len(properties)),
			zap.Int64("property_id", propertyData.Property.HotelID),
		)

		if err := storage.StoreProperty(ctx, propertyData); err != nil {
			logger.LogError("Failed to store property", err,
				zap.Int64("property_id", propertyData.Property.HotelID),
			)
			errorCount++
		} else {
			successCount++
		}
	}

	logger.LogSuccess("Data storage completed",
		zap.Int("successful", successCount),
		zap.Int("failed", errorCount),
		zap.Int("total", len(properties)),
	)

	// Test fetching a single property
	logger.Info("Testing property retrieval...")
	testProperty, err := storage.GetProperty(ctx, 1018946)
	if err != nil {
		logger.LogError("Failed to retrieve test property", err)
	} else {
		logger.LogSuccess("Test property retrieved successfully",
			zap.Int64("property_id", testProperty.Property.HotelID),
			zap.String("hotel_name", testProperty.Property.HotelName),
			zap.Int("review_count", len(testProperty.Reviews)),
			zap.Int("translation_count", len(testProperty.Translations)),
		)
	}
}
