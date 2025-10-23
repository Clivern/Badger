// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package migration

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/rs/zerolog/log"
)

// Migration represents a database migration
type Migration struct {
	Version     string
	Description string
	Up          func(*sql.DB) error
	Down        func(*sql.DB) error
}

// Manager handles database migrations
type Manager struct {
	db         *sql.DB
	driver     string
	migrations []Migration
}

// NewManager creates a new migration manager
func NewManager(db *sql.DB, driver string) *Manager {
	return &Manager{
		db:         db,
		driver:     driver,
		migrations: []Migration{},
	}
}

// Register adds a migration to the manager
func (m *Manager) Register(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// createMigrationsTable creates the migrations tracking table if it doesn't exist
func (m *Manager) createMigrationsTable() error {
	var query string

	if m.driver == "sqlite" {
		query = `
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			version VARCHAR(255) NOT NULL UNIQUE,
			description TEXT,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		query = `
		CREATE TABLE IF NOT EXISTS migrations (
			id INT AUTO_INCREMENT PRIMARY KEY,
			version VARCHAR(255) NOT NULL UNIQUE,
			description TEXT,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_version (version)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	return nil
}

// isApplied checks if a migration version has been applied
func (m *Manager) isApplied(version string) (bool, error) {
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM migrations WHERE version = ?", version).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check migration status: %w", err)
	}
	return count > 0, nil
}

// recordMigration records a migration as applied
func (m *Manager) recordMigration(version, description string) error {
	_, err := m.db.Exec(
		"INSERT INTO migrations (version, description, applied_at) VALUES (?, ?, ?)",
		version,
		description,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}
	return nil
}

// removeMigration removes a migration record
func (m *Manager) removeMigration(version string) error {
	_, err := m.db.Exec("DELETE FROM migrations WHERE version = ?", version)
	if err != nil {
		return fmt.Errorf("failed to remove migration record: %w", err)
	}
	return nil
}

// Up runs all pending migrations
func (m *Manager) Up() error {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return err
	}

	// Sort migrations by version
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	appliedCount := 0

	for _, migration := range m.migrations {
		applied, err := m.isApplied(migration.Version)
		if err != nil {
			return err
		}

		if applied {
			log.Debug().
				Str("version", migration.Version).
				Msg("Migration already applied, skipping")
			continue
		}

		log.Info().
			Str("version", migration.Version).
			Str("description", migration.Description).
			Msg("Running migration")

		// Start transaction
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		// Run the migration
		if err := migration.Up(m.db); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %s failed: %w", migration.Version, err)
		}

		// Record the migration
		if err := m.recordMigration(migration.Version, migration.Description); err != nil {
			tx.Rollback()
			return err
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		log.Info().
			Str("version", migration.Version).
			Msg("Migration applied successfully")

		appliedCount++
	}

	if appliedCount == 0 {
		log.Info().Msg("No pending migrations to apply")
	} else {
		log.Info().
			Int("count", appliedCount).
			Msg("Migrations applied successfully")
	}

	return nil
}

// Down rolls back the last migration
func (m *Manager) Down() error {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return err
	}

	// Get the last applied migration
	var version, description string
	err := m.db.QueryRow(`
		SELECT version, description
		FROM migrations
		ORDER BY applied_at DESC, id DESC
		LIMIT 1
	`).Scan(&version, &description)

	if err == sql.ErrNoRows {
		log.Info().Msg("No migrations to roll back")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	// Find the migration
	var migration *Migration
	for i := range m.migrations {
		if m.migrations[i].Version == version {
			migration = &m.migrations[i]
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration %s not found in registered migrations", version)
	}

	if migration.Down == nil {
		return fmt.Errorf("migration %s has no down function", version)
	}

	log.Info().
		Str("version", version).
		Str("description", description).
		Msg("Rolling back migration")

	// Start transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	// Run the rollback
	if err := migration.Down(m.db); err != nil {
		tx.Rollback()
		return fmt.Errorf("rollback of migration %s failed: %w", version, err)
	}

	// Remove the migration record
	if err := m.removeMigration(version); err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Info().
		Str("version", version).
		Msg("Migration rolled back successfully")

	return nil
}

// Status shows the status of all migrations
func (m *Manager) Status() error {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return err
	}

	// Sort migrations by version
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	log.Info().Msg("Migration Status:")
	log.Info().Msg("==================")

	for _, migration := range m.migrations {
		applied, err := m.isApplied(migration.Version)
		if err != nil {
			return err
		}

		status := "Pending"
		if applied {
			status = "Applied"
		}

		log.Info().
			Str("version", migration.Version).
			Str("description", migration.Description).
			Str("status", status).
			Msg("")
	}

	return nil
}
