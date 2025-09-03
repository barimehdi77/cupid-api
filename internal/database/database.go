package database

import (
	"database/sql"
	"fmt"

	"github.com/barimehdi77/cupid-api/internal/env"
	"github.com/barimehdi77/cupid-api/internal/logger"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

var Db *sql.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		logger.Warn("Could not load .env file", zap.Error(err))
	}

	// Use the env package for consistency
	host := env.GetEnvString("DB_HOST", "localhost")
	port := env.GetEnvInt("DB_PORT", 5432)
	user := env.GetEnvString("DB_USER", "root")
	dbname := env.GetEnvString("DB_NAME", "cupid")
	password := env.GetEnvString("DB_PASSWORD", "")

	// Set up postgres connection string
	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)

	logger.Debug("Connecting to database",
		zap.String("host", host),
		zap.Int("port", port),
		zap.String("user", user),
		zap.String("dbname", dbname),
	)

	db, err := sql.Open("postgres", psqlSetup)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}

	Db = db
	logger.Info("Successfully connected to database!")
}
