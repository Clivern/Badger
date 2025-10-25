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
		{
			Version:     "20250101000004",
			Description: "Create users table",
			Up:          createUsersTable,
			Down:        dropUsersTable,
		},
		{
			Version:     "20250101000005",
			Description: "Create users_meta table",
			Up:          createUsersMetaTable,
			Down:        dropUsersMetaTable,
		},
		{
			Version:     "20250101000006",
			Description: "Create servers table",
			Up:          createServersTable,
			Down:        dropServersTable,
		},
		{
			Version:     "20250101000007",
			Description: "Create servers_meta table",
			Up:          createServersMetaTable,
			Down:        dropServersMetaTable,
		},
		{
			Version:     "20250101000008",
			Description: "Create mcps table",
			Up:          createMcpsTable,
			Down:        dropMcpsTable,
		},
		{
			Version:     "20250101000009",
			Description: "Create mcps_meta table",
			Up:          createMcpsMetaTable,
			Down:        dropMcpsMetaTable,
		},
		{
			Version:     "20250101000010",
			Description: "Create gateways table",
			Up:          createGatewaysTable,
			Down:        dropGatewaysTable,
		},
		{
			Version:     "20250101000011",
			Description: "Create gateways_meta table",
			Up:          createGatewaysMetaTable,
			Down:        dropGatewaysMetaTable,
		},
		{
			Version:     "20250101000012",
			Description: "Create tools table",
			Up:          createToolsTable,
			Down:        dropToolsTable,
		},
		{
			Version:     "20250101000013",
			Description: "Create tools_meta table",
			Up:          createToolsMetaTable,
			Down:        dropToolsMetaTable,
		},
		{
			Version:     "20250101000014",
			Description: "Create tool_metrics table",
			Up:          createToolMetricsTable,
			Down:        dropToolMetricsTable,
		},
		{
			Version:     "20250101000015",
			Description: "Create resources table",
			Up:          createResourcesTable,
			Down:        dropResourcesTable,
		},
		{
			Version:     "20250101000016",
			Description: "Create resources_meta table",
			Up:          createResourcesMetaTable,
			Down:        dropResourcesMetaTable,
		},
		{
			Version:     "20250101000017",
			Description: "Create prompts table",
			Up:          createPromptsTable,
			Down:        dropPromptsTable,
		},
		{
			Version:     "20250101000018",
			Description: "Create prompts_meta table",
			Up:          createPromptsMetaTable,
			Down:        dropPromptsMetaTable,
		},
		{
			Version:     "20250101000019",
			Description: "Create audits table",
			Up:          createAuditsTable,
			Down:        dropAuditsTable,
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

// createUsersTable creates the users table
func createUsersTable(db *sql.DB) error {
	// Try to determine if it's SQLite
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email VARCHAR(255) NOT NULL UNIQUE,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			password VARCHAR(255) NOT NULL,
			api_key VARCHAR(255) UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		query = `
		CREATE TABLE users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			password VARCHAR(255) NOT NULL,
			api_key VARCHAR(255) UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_email (email),
			INDEX idx_api_key (api_key)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropUsersTable drops the users table
func dropUsersTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users")
	return err
}

// createUsersMetaTable creates the users_meta table
func createUsersMetaTable(db *sql.DB) error {
	// Try to determine if it's SQLite
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE users_meta (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			user_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(user_id, key)
		)`
	} else {
		query = `
		CREATE TABLE users_meta (
			id INT AUTO_INCREMENT PRIMARY KEY,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			user_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE KEY unique_user_key (user_id, key),
			INDEX idx_user_id (user_id),
			INDEX idx_key (key)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropUsersMetaTable drops the users_meta table
func dropUsersMetaTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS users_meta")
	return err
}

// createServersTable creates the servers table
func createServersTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE servers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		query = `
		CREATE TABLE servers (
			id INT AUTO_INCREMENT PRIMARY KEY,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropServersTable drops the servers table
func dropServersTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS servers")
	return err
}

// createServersMetaTable creates the servers_meta table
func createServersMetaTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE servers_meta (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			server_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE,
			UNIQUE(server_id, key)
		)`
	} else {
		query = `
		CREATE TABLE servers_meta (
			id INT AUTO_INCREMENT PRIMARY KEY,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			server_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE,
			UNIQUE KEY unique_server_key (server_id, key),
			INDEX idx_server_id (server_id),
			INDEX idx_key (key)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropServersMetaTable drops the servers_meta table
func dropServersMetaTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS servers_meta")
	return err
}

// createMcpsTable creates the mcps table
func createMcpsTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE mcps (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		query = `
		CREATE TABLE mcps (
			id INT AUTO_INCREMENT PRIMARY KEY,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropMcpsTable drops the mcps table
func dropMcpsTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS mcps")
	return err
}

// createMcpsMetaTable creates the mcps_meta table
func createMcpsMetaTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE mcps_meta (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			mcp_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (mcp_id) REFERENCES mcps(id) ON DELETE CASCADE,
			UNIQUE(mcp_id, key)
		)`
	} else {
		query = `
		CREATE TABLE mcps_meta (
			id INT AUTO_INCREMENT PRIMARY KEY,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			mcp_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (mcp_id) REFERENCES mcps(id) ON DELETE CASCADE,
			UNIQUE KEY unique_mcp_key (mcp_id, key),
			INDEX idx_mcp_id (mcp_id),
			INDEX idx_key (key)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropMcpsMetaTable drops the mcps_meta table
func dropMcpsMetaTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS mcps_meta")
	return err
}

// createGatewaysTable creates the gateways table
func createGatewaysTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE gateways (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		query = `
		CREATE TABLE gateways (
			id INT AUTO_INCREMENT PRIMARY KEY,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropGatewaysTable drops the gateways table
func dropGatewaysTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS gateways")
	return err
}

// createGatewaysMetaTable creates the gateways_meta table
func createGatewaysMetaTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE gateways_meta (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			gateway_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (gateway_id) REFERENCES gateways(id) ON DELETE CASCADE,
			UNIQUE(gateway_id, key)
		)`
	} else {
		query = `
		CREATE TABLE gateways_meta (
			id INT AUTO_INCREMENT PRIMARY KEY,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			gateway_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (gateway_id) REFERENCES gateways(id) ON DELETE CASCADE,
			UNIQUE KEY unique_gateway_key (gateway_id, key),
			INDEX idx_gateway_id (gateway_id),
			INDEX idx_key (key)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropGatewaysMetaTable drops the gateways_meta table
func dropGatewaysMetaTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS gateways_meta")
	return err
}

// createToolsTable creates the tools table
func createToolsTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE tools (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		query = `
		CREATE TABLE tools (
			id INT AUTO_INCREMENT PRIMARY KEY,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropToolsTable drops the tools table
func dropToolsTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS tools")
	return err
}

// createToolsMetaTable creates the tools_meta table
func createToolsMetaTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE tools_meta (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			tool_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE,
			UNIQUE(tool_id, key)
		)`
	} else {
		query = `
		CREATE TABLE tools_meta (
			id INT AUTO_INCREMENT PRIMARY KEY,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			tool_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE,
			UNIQUE KEY unique_tool_key (tool_id, key),
			INDEX idx_tool_id (tool_id),
			INDEX idx_key (key)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropToolsMetaTable drops the tools_meta table
func dropToolsMetaTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS tools_meta")
	return err
}

// createToolMetricsTable creates the tool_metrics table
func createToolMetricsTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE tool_metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			tool_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE
		)`
	} else {
		query = `
		CREATE TABLE tool_metrics (
			id INT AUTO_INCREMENT PRIMARY KEY,
			tool_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE,
			INDEX idx_tool_id (tool_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropToolMetricsTable drops the tool_metrics table
func dropToolMetricsTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS tool_metrics")
	return err
}

// createResourcesTable creates the resources table
func createResourcesTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE resources (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		query = `
		CREATE TABLE resources (
			id INT AUTO_INCREMENT PRIMARY KEY,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropResourcesTable drops the resources table
func dropResourcesTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS resources")
	return err
}

// createResourcesMetaTable creates the resources_meta table
func createResourcesMetaTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE resources_meta (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			resource_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
			UNIQUE(resource_id, key)
		)`
	} else {
		query = `
		CREATE TABLE resources_meta (
			id INT AUTO_INCREMENT PRIMARY KEY,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			resource_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
			UNIQUE KEY unique_resource_key (resource_id, key),
			INDEX idx_resource_id (resource_id),
			INDEX idx_key (key)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropResourcesMetaTable drops the resources_meta table
func dropResourcesMetaTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS resources_meta")
	return err
}

// createPromptsTable creates the prompts table
func createPromptsTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE prompts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		query = `
		CREATE TABLE prompts (
			id INT AUTO_INCREMENT PRIMARY KEY,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropPromptsTable drops the prompts table
func dropPromptsTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS prompts")
	return err
}

// createPromptsMetaTable creates the prompts_meta table
func createPromptsMetaTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE prompts_meta (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			prompt_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (prompt_id) REFERENCES prompts(id) ON DELETE CASCADE,
			UNIQUE(prompt_id, key)
		)`
	} else {
		query = `
		CREATE TABLE prompts_meta (
			id INT AUTO_INCREMENT PRIMARY KEY,
			key VARCHAR(255) NOT NULL,
			value TEXT,
			prompt_id INT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (prompt_id) REFERENCES prompts(id) ON DELETE CASCADE,
			UNIQUE KEY unique_prompt_key (prompt_id, key),
			INDEX idx_prompt_id (prompt_id),
			INDEX idx_key (key)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropPromptsMetaTable drops the prompts_meta table
func dropPromptsMetaTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS prompts_meta")
	return err
}

// createAuditsTable creates the audits table
func createAuditsTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE audits (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		query = `
		CREATE TABLE audits (
			id INT AUTO_INCREMENT PRIMARY KEY,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropAuditsTable drops the audits table
func dropAuditsTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS audits")
	return err
}
