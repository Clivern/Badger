// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

// Connection represents a database connection
type Connection struct {
	DB     *sql.DB
	Driver string
}

// Config holds database configuration
type Config struct {
	Driver          string
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime int
	DataSource      string
}

// NewConnection creates a new database connection based on the driver
func NewConnection(config Config) (*Connection, error) {
	var dsn string
	var err error
	var db *sql.DB

	switch config.Driver {
	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Username,
			config.Password,
			config.Host,
			config.Port,
			config.Database,
		)
		db, err = sql.Open("mysql", dsn)
	case "sqlite":
		// Use DataSource if provided, otherwise default to database name
		dsn = config.DataSource
		if dsn == "" {
			dsn = config.Database
		}
		db, err = sql.Open("sqlite3", dsn)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	if config.MaxOpenConns > 0 {
		db.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.MaxIdleConns > 0 {
		db.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	}

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info().
		Str("driver", config.Driver).
		Str("host", config.Host).
		Int("port", config.Port).
		Str("database", config.Database).
		Msg("Database connection established")

	return &Connection{
		DB:     db,
		Driver: config.Driver,
	}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	if c.DB != nil {
		log.Info().Msg("Closing database connection")
		return c.DB.Close()
	}
	return nil
}

// Ping checks if the database connection is alive
func (c *Connection) Ping() error {
	return c.DB.Ping()
}
