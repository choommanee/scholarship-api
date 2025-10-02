package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"scholarship-system/internal/config"
	"scholarship-system/internal/database"
)

func main() {
	var (
		action        = flag.String("action", "migrate", "Action to perform: migrate, status, rollback")
		migrationsDir = flag.String("migrations", "./migrations", "Directory containing migration files")
	)
	flag.Parse()

	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// Ensure migrations directory exists and is absolute
	absDir, err := filepath.Abs(*migrationsDir)
	if err != nil {
		log.Fatal("Failed to get absolute path for migrations directory:", err)
	}

	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		log.Fatalf("Migrations directory does not exist: %s", absDir)
	}

	log.Printf("Using migrations directory: %s", absDir)

	// Create migration manager
	manager := database.NewMigrationManager(database.DB, absDir)

	// Perform the requested action
	switch *action {
	case "migrate", "up":
		if err := manager.Migrate(); err != nil {
			log.Fatal("Migration failed:", err)
		}
		log.Println("Migration completed successfully")

	case "status":
		if err := manager.GetStatus(); err != nil {
			log.Fatal("Failed to get migration status:", err)
		}

	case "init":
		if err := manager.InitMigrationTable(); err != nil {
			log.Fatal("Failed to initialize migration table:", err)
		}
		log.Println("Migration table initialized successfully")

	default:
		log.Fatalf("Unknown action: %s. Available actions: migrate, status, init", *action)
	}
}
