// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

/*
Package migration provides database migration management for Badger.

# Overview

The migration system allows you to version control your database schema
and apply changes incrementally. It supports both MySQL and SQLite databases
with automatic driver detection.

# Basic Usage

Create a migration manager and run migrations:

	import (
		"github.com/clivern/badger/db"
		"github.com/clivern/badger/migration"
	)

	// Connect to database
	config := db.Config{
		Driver:     "sqlite",
		DataSource: "./cache/badger.db",
	}
	conn, err := db.NewConnection(config)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Create migration manager
	mgr := migration.NewManager(conn.DB, conn.Driver)

	// Register all migrations
	for _, m := range migration.GetAll() {
		mgr.Register(m)
	}

	// Run pending migrations
	if err := mgr.Up(); err != nil {
		log.Fatal(err)
	}

# CLI Usage

The easiest way to run migrations is through the CLI:

	# Run all pending migrations
	badger migrate up -c config.dist.yml

	# Roll back the last migration
	badger migrate down -c config.dist.yml

	# Check migration status
	badger migrate status -c config.dist.yml

# Migration Structure

Each migration has a version, description, and up/down functions:

	Migration{
		Version:     "20250101000001",
		Description: "Create users table",
		Up:          createUsersTable,
		Down:        dropUsersTable,
	}

# Creating New Migrations

Add your migration to registry.go:

	func GetAll() []Migration {
		return []Migration{
			{
				Version:     "20250101000001",
				Description: "Create users table",
				Up:          createUsersTable,
				Down:        dropUsersTable,
			},
			{
				Version:     "20250102000001",
				Description: "Create your new table",
				Up:          createYourTable,
				Down:        dropYourTable,
			},
		}
	}

	// Implement your migration functions
	func createYourTable(db *sql.DB) error {
		query := `CREATE TABLE your_table (...)`
		_, err := db.Exec(query)
		return err
	}

	func dropYourTable(db *sql.DB) error {
		_, err := db.Exec("DROP TABLE IF EXISTS your_table")
		return err
	}

# Version Format

Migration versions should follow the format: YYYYMMDDHHMMSS
Example: 20250101000001 (January 1, 2025, 00:00:01)

This ensures migrations are applied in chronological order.

# Database Support

The migration system automatically detects the database driver and
uses appropriate SQL syntax for:
  - Auto-increment fields (AUTO_INCREMENT vs AUTOINCREMENT)
  - Boolean defaults (TRUE vs 1)
  - Date/time functions
  - Table engines (InnoDB for MySQL)
  - Character sets (utf8mb4 for MySQL)

# Transaction Safety

All migrations run within transactions. If a migration fails,
changes are rolled back automatically.

# Migration Tracking

A migrations table is automatically created to track which
migrations have been applied. This ensures each migration
runs only once.

# Best Practices

  - Always provide both Up and Down functions
  - Test migrations on a copy of production data
  - Keep migrations small and focused
  - Never modify existing migrations after they've been deployed
  - Use descriptive names and comments
  - Back up your database before running migrations in production
*/
package migration
