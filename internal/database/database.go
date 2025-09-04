// internal/database/database.go
package database

import (
	"database/sql"
	"fmt"

	"github.com/barimehdi77/cupid-api/internal/env"
	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB() (*DB, error) {
	// Get database configuration
	driver := env.GetEnvString("DB_DRIVER", "postgres")
	host := env.GetEnvString("DB_HOST", "localhost")
	port := env.GetEnvInt("DB_PORT", 5432)
	user := env.GetEnvString("DB_USER", "root")
	dbname := env.GetEnvString("DB_NAME", "cupid")
	password := env.GetEnvString("DB_PASSWORD", "")

	psqlSetup := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)

	db, err := sql.Open(driver, psqlSetup)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{DB: db}, nil
}

// Add helper methods if needed
func (db *DB) Close() error {
	return db.DB.Close()
}
