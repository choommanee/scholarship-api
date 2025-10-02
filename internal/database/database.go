package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"scholarship-system/internal/config"
)

var DB *sql.DB

func Connect(cfg *config.Config) error {
	var err error

	DB, err = sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)

	log.Println("Database connected successfully")
	return nil
}

func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}