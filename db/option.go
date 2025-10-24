// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Option represents a key-value option in the database.
//
// Options are stored in the options table and provide a flexible
// way to store application settings and configuration values.
type Option struct {
	ID        int64
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// OptionRepository handles database operations for options.
//
// It provides CRUD operations for the options table using the repository pattern.
type OptionRepository struct {
	db *sql.DB
}

// NewOptionRepository creates a new option repository.
//
// Example:
//
//	repo := db.NewOptionRepository(conn.DB)
func NewOptionRepository(db *sql.DB) *OptionRepository {
	return &OptionRepository{db: db}
}

// Create inserts a new option into the database.
//
// The key must be unique. If a key already exists, this will return an error.
//
// Example:
//
//	err := repo.Create("app_name", "Yun")
//	if err != nil {
//		log.Fatal(err)
//	}
func (r *OptionRepository) Create(key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO options (key, value) VALUES (?, ?)",
		key,
		value,
	)
	return err
}

// Get retrieves an option by key.
//
// Returns nil if the option doesn't exist.
//
// Example:
//
//	opt, err := repo.Get("app_name")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if opt != nil {
//		fmt.Printf("Value: %s\n", opt.Value)
//	}
func (r *OptionRepository) Get(key string) (*Option, error) {
	opt := &Option{}
	err := r.db.QueryRow(
		"SELECT id, key, value, created_at, updated_at FROM options WHERE key = ?",
		key,
	).Scan(&opt.ID, &opt.Key, &opt.Value, &opt.CreatedAt, &opt.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return opt, nil
}

// Update updates an option value.
//
// The updated_at timestamp is automatically set to the current time.
//
// Example:
//
//	err := repo.Update("app_name", "Yun v2")
//	if err != nil {
//		log.Fatal(err)
//	}
func (r *OptionRepository) Update(key, value string) error {
	_, err := r.db.Exec(
		"UPDATE options SET value = ?, updated_at = ? WHERE key = ?",
		value,
		time.Now(),
		key,
	)
	return err
}

// Delete removes an option from the database.
//
// Example:
//
//	err := repo.Delete("app_name")
//	if err != nil {
//		log.Fatal(err)
//	}
func (r *OptionRepository) Delete(key string) error {
	_, err := r.db.Exec("DELETE FROM options WHERE key = ?", key)
	return err
}

// List retrieves all options from the database.
//
// Results are ordered by key alphabetically.
//
// Example:
//
//	options, err := repo.List()
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, opt := range options {
//		fmt.Printf("%s = %s\n", opt.Key, opt.Value)
//	}
func (r *OptionRepository) List() ([]*Option, error) {
	rows, err := r.db.Query("SELECT id, key, value, created_at, updated_at FROM options ORDER BY key")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []*Option
	for rows.Next() {
		opt := &Option{}
		if err := rows.Scan(&opt.ID, &opt.Key, &opt.Value, &opt.CreatedAt, &opt.UpdatedAt); err != nil {
			return nil, err
		}
		options = append(options, opt)
	}

	return options, rows.Err()
}
