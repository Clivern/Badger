// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package migration

import "database/sql"

// GetAll returns all registered migrations
func GetAll() []Migration {
	return []Migration{
		{
			Version:     "20250101000003",
			Description: "Create options table",
			Up:          createOptionsTable,
			Down:        dropOptionsTable,
		},
	}
}

// createOptionsTable creates the options table
func createOptionsTable(db *sql.DB) error {
	// Try to determine if it's SQLite
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE options (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key VARCHAR(255) NOT NULL UNIQUE,
			value TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		query = `
		CREATE TABLE options (
			id INT AUTO_INCREMENT PRIMARY KEY,
			key VARCHAR(255) NOT NULL UNIQUE,
			value TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_key (key)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropOptionsTable drops the options table
func dropOptionsTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS options")
	return err
}
