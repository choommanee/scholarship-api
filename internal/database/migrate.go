package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Migration represents a database migration
type Migration struct {
	Version   string
	Name      string
	UpSQL     string
	DownSQL   string
	AppliedAt *time.Time
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db            *sql.DB
	migrationsDir string
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sql.DB, migrationsDir string) *MigrationManager {
	return &MigrationManager{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// InitMigrationTable creates the migration tracking table if it doesn't exist
func (m *MigrationManager) InitMigrationTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	log.Println("Migration table initialized")
	return nil
}

// GetAppliedMigrations returns a list of applied migrations
func (m *MigrationManager) GetAppliedMigrations() (map[string]*Migration, error) {
	query := `SELECT version, name, applied_at FROM schema_migrations ORDER BY version`
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]*Migration)
	for rows.Next() {
		var migration Migration
		err := rows.Scan(&migration.Version, &migration.Name, &migration.AppliedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan migration: %w", err)
		}
		applied[migration.Version] = &migration
	}

	return applied, nil
}

// LoadMigrationFiles reads all migration files from the migrations directory
func (m *MigrationManager) LoadMigrationFiles() ([]*Migration, error) {
	var migrations []*Migration

	err := filepath.WalkDir(m.migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Skip non-SQL files
		if !strings.HasSuffix(path, ".sql") {
			return nil
		}

		// Parse migration file name
		fileName := d.Name()

		// Skip down migration files for now (we'll process them separately)
		if strings.Contains(fileName, ".down.") {
			return nil
		}

		// Extract version and name from file name
		// Expected format: 001_initial_schema.sql or 000004_create_news_table.up.sql
		baseName := strings.TrimSuffix(fileName, ".sql")
		baseName = strings.TrimSuffix(baseName, ".up")

		parts := strings.SplitN(baseName, "_", 2)
		if len(parts) < 2 {
			// Handle files like 001_initial_schema.sql
			if len(parts) == 1 && strings.HasPrefix(parts[0], "00") {
				parts = []string{parts[0], "migration"}
			} else {
				log.Printf("Skipping file with unexpected format: %s", fileName)
				return nil
			}
		}

		version := parts[0]
		name := parts[1]

		// Read the migration file content
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", path, err)
		}

		// Look for corresponding down migration
		downContent := ""
		downPath := strings.Replace(path, ".up.sql", ".down.sql", 1)
		if !strings.HasSuffix(path, ".up.sql") {
			// For files like 001_initial_schema.sql, look for down migration
			downPath = filepath.Join(m.migrationsDir, version+"_"+name+".down.sql")
		}

		if _, err := os.Stat(downPath); err == nil {
			downBytes, err := os.ReadFile(downPath)
			if err == nil {
				downContent = string(downBytes)
			}
		}

		migration := &Migration{
			Version: version,
			Name:    name,
			UpSQL:   string(content),
			DownSQL: downContent,
		}

		migrations = append(migrations, migration)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load migration files: %w", err)
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// ApplyMigration applies a single migration
func (m *MigrationManager) ApplyMigration(migration *Migration) error {
	log.Printf("Applying migration %s: %s", migration.Version, migration.Name)

	// Start transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Execute migration SQL
	_, err = tx.Exec(migration.UpSQL)
	if err != nil {
		return fmt.Errorf("failed to execute migration %s: %w", migration.Version, err)
	}

	// Record the migration as applied
	_, err = tx.Exec(
		`INSERT INTO schema_migrations (version, name) VALUES ($1, $2)`,
		migration.Version, migration.Name,
	)
	if err != nil {
		return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration %s: %w", migration.Version, err)
	}

	log.Printf("Successfully applied migration %s", migration.Version)
	return nil
}

// Migrate runs all pending migrations
func (m *MigrationManager) Migrate() error {
	log.Println("Starting database migration...")

	// Initialize migration table
	if err := m.InitMigrationTable(); err != nil {
		return err
	}

	// Get applied migrations
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	// Load migration files
	migrations, err := m.LoadMigrationFiles()
	if err != nil {
		return err
	}

	// Apply pending migrations
	pendingCount := 0
	for _, migration := range migrations {
		if _, exists := applied[migration.Version]; !exists {
			if err := m.ApplyMigration(migration); err != nil {
				return fmt.Errorf("migration failed at version %s: %w", migration.Version, err)
			}
			pendingCount++
		}
	}

	if pendingCount == 0 {
		log.Println("No pending migrations found")
	} else {
		log.Printf("Successfully applied %d migrations", pendingCount)
	}

	return nil
}

// GetStatus returns the current migration status
func (m *MigrationManager) GetStatus() error {
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}

	migrations, err := m.LoadMigrationFiles()
	if err != nil {
		return err
	}

	log.Println("Migration Status:")
	log.Println("================")

	for _, migration := range migrations {
		status := "[ ]"
		appliedAt := ""
		if appliedMigration, exists := applied[migration.Version]; exists {
			status = "[✓]"
			if appliedMigration.AppliedAt != nil {
				appliedAt = appliedMigration.AppliedAt.Format("2006-01-02 15:04:05")
			}
		}

		log.Printf("%s %s %s %s", status, migration.Version, migration.Name, appliedAt)
	}

	pendingCount := len(migrations) - len(applied)
	log.Printf("\nTotal migrations: %d", len(migrations))
	log.Printf("Applied: %d", len(applied))
	log.Printf("Pending: %d", pendingCount)

	return nil
}

// RunMigrations is a convenience function to run migrations
func RunMigrations(db *sql.DB, migrationsDir string) error {
	manager := NewMigrationManager(db, migrationsDir)
	return manager.Migrate()
}

// GetMigrationStatus is a convenience function to get migration status
func GetMigrationStatus(db *sql.DB, migrationsDir string) error {
	manager := NewMigrationManager(db, migrationsDir)
	return manager.GetStatus()
}
