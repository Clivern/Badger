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
			Version:     "20250101000016",
			Description: "Create users_meta table",
			Up:          createUsersMetaTable,
			Down:        dropUsersMetaTable,
		},
		{
			Version:     "20250101000005",
			Description: "Create servers table",
			Up:          createServersTable,
			Down:        dropServersTable,
		},
		{
			Version:     "20250101000018",
			Description: "Create servers_meta table",
			Up:          createServersMetaTable,
			Down:        dropServersMetaTable,
		},
		{
			Version:     "20250101000006",
			Description: "Create mcps table",
			Up:          createMcpsTable,
			Down:        dropMcpsTable,
		},
		{
			Version:     "20250101000017",
			Description: "Create mcps_meta table",
			Up:          createMcpsMetaTable,
			Down:        dropMcpsMetaTable,
		},
		{
			Version:     "20250101000007",
			Description: "Create gateways table",
			Up:          createGatewaysTable,
			Down:        dropGatewaysTable,
		},
		{
			Version:     "20250101000022",
			Description: "Create gateways_meta table",
			Up:          createGatewaysMetaTable,
			Down:        dropGatewaysMetaTable,
		},
		{
			Version:     "20250101000008",
			Description: "Create tools table",
			Up:          createToolsTable,
			Down:        dropToolsTable,
		},
		{
			Version:     "20250101000019",
			Description: "Create tools_meta table",
			Up:          createToolsMetaTable,
			Down:        dropToolsMetaTable,
		},
		{
			Version:     "20250101000009",
			Description: "Create resources table",
			Up:          createResourcesTable,
			Down:        dropResourcesTable,
		},
		{
			Version:     "20250101000020",
			Description: "Create resources_meta table",
			Up:          createResourcesMetaTable,
			Down:        dropResourcesMetaTable,
		},
		{
			Version:     "20250101000010",
			Description: "Create prompts table",
			Up:          createPromptsTable,
			Down:        dropPromptsTable,
		},
		{
			Version:     "20250101000021",
			Description: "Create prompts_meta table",
			Up:          createPromptsMetaTable,
			Down:        dropPromptsMetaTable,
		},
		{
			Version:     "20250101000011",
			Description: "Create server_tools table",
			Up:          createServerToolsTable,
			Down:        dropServerToolsTable,
		},
		{
			Version:     "20250101000012",
			Description: "Create server_resources table",
			Up:          createServerResourcesTable,
			Down:        dropServerResourcesTable,
		},
		{
			Version:     "20250101000013",
			Description: "Create server_prompts table",
			Up:          createServerPromptsTable,
			Down:        dropServerPromptsTable,
		},
		{
			Version:     "20250101000014",
			Description: "Create tool_metrics table",
			Up:          createToolMetricsTable,
			Down:        dropToolMetricsTable,
		},
		{
			Version:     "20250101000015",
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
		// role is admin, user or readonly
		query = `
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			api_key VARCHAR(255) UNIQUE,
			is_active BOOLEAN DEFAULT 1,
			last_login_at DATETIME NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		query = `
		CREATE TABLE users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			role VARCHAR(50) NOT NULL DEFAULT 'user',
			api_key VARCHAR(255) UNIQUE,
			is_active BOOLEAN DEFAULT true,
			last_login_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_email (email),
			INDEX idx_api_key (api_key),
			INDEX idx_role (role)
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

// createMcpsTable creates the mcps table (Backend MCP Server Connections)
func createMcpsTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE mcps (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			url VARCHAR(767) NOT NULL,
			transport VARCHAR(20) NOT NULL DEFAULT 'sse',
			auth_type VARCHAR(20) DEFAULT 'none',
			auth_token TEXT,
			timeout_ms INTEGER DEFAULT 30000,
			max_retries INTEGER DEFAULT 3,
			headers TEXT,
			status VARCHAR(20) DEFAULT 'active',
			health_check_url VARCHAR(767),
			last_health_check_at DATETIME NULL,
			health_status VARCHAR(20) DEFAULT 'unknown',
			capabilities TEXT,
			protocol_version VARCHAR(20),
			description TEXT,
			tags VARCHAR(500),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		query = `
		CREATE TABLE mcps (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			url VARCHAR(767) NOT NULL,
			transport VARCHAR(20) NOT NULL DEFAULT 'sse',
			auth_type VARCHAR(20) DEFAULT 'none',
			auth_token TEXT,
			timeout_ms INT DEFAULT 30000,
			max_retries INT DEFAULT 3,
			headers TEXT,
			status VARCHAR(20) DEFAULT 'active',
			health_check_url VARCHAR(767),
			last_health_check_at TIMESTAMP NULL,
			health_status VARCHAR(20) DEFAULT 'unknown',
			capabilities TEXT,
			protocol_version VARCHAR(20),
			description TEXT,
			tags VARCHAR(500),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_slug (slug),
			INDEX idx_status (status),
			INDEX idx_transport (transport)
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

// createServersTable creates the servers table (Virtual Servers)
func createServersTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE servers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			description TEXT,
			is_public BOOLEAN DEFAULT 0,
			allowed_user_ids TEXT,
			enable_tools BOOLEAN DEFAULT 1,
			enable_resources BOOLEAN DEFAULT 1,
			enable_prompts BOOLEAN DEFAULT 1,
			tags VARCHAR(500),
			created_by INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
		)`
	} else {
		query = `
		CREATE TABLE servers (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			description TEXT,
			is_public BOOLEAN DEFAULT false,
			allowed_user_ids TEXT,
			enable_tools BOOLEAN DEFAULT true,
			enable_resources BOOLEAN DEFAULT true,
			enable_prompts BOOLEAN DEFAULT true,
			tags VARCHAR(500),
			created_by INT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
			INDEX idx_slug (slug),
			INDEX idx_created_by (created_by)
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

// createToolsTable creates the tools table (Discovered from MCPs)
func createToolsTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE tools (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			original_name VARCHAR(255) NOT NULL,
			mcp_id INTEGER NOT NULL,
			description TEXT,
			input_schema TEXT NOT NULL,
			is_enabled BOOLEAN DEFAULT 1,
			timeout_ms INTEGER DEFAULT 30000,
			max_retries INTEGER DEFAULT 3,
			tags VARCHAR(500),
			category VARCHAR(100),
			call_count INTEGER DEFAULT 0,
			last_called_at DATETIME NULL,
			avg_response_time_ms INTEGER DEFAULT 0,
			error_count INTEGER DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (mcp_id) REFERENCES mcps(id) ON DELETE CASCADE,
			UNIQUE(mcp_id, original_name)
		)`
	} else {
		query = `
		CREATE TABLE tools (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			original_name VARCHAR(255) NOT NULL,
			mcp_id INT NOT NULL,
			description TEXT,
			input_schema TEXT NOT NULL,
			is_enabled BOOLEAN DEFAULT true,
			timeout_ms INT DEFAULT 30000,
			max_retries INT DEFAULT 3,
			tags VARCHAR(500),
			category VARCHAR(100),
			call_count INT DEFAULT 0,
			last_called_at TIMESTAMP NULL,
			avg_response_time_ms INT DEFAULT 0,
			error_count INT DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (mcp_id) REFERENCES mcps(id) ON DELETE CASCADE,
			UNIQUE KEY unique_tool_per_mcp (mcp_id, original_name),
			INDEX idx_name (name),
			INDEX idx_mcp_id (mcp_id),
			INDEX idx_enabled (is_enabled),
			INDEX idx_category (category)
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

// createResourcesTable creates the resources table (Discovered from MCPs)
func createResourcesTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE resources (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			original_name VARCHAR(255) NOT NULL,
			uri VARCHAR(767) NOT NULL,
			mcp_id INTEGER NOT NULL,
			description TEXT,
			mime_type VARCHAR(100),
			is_enabled BOOLEAN DEFAULT 1,
			tags VARCHAR(500),
			access_count INTEGER DEFAULT 0,
			last_accessed_at DATETIME NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (mcp_id) REFERENCES mcps(id) ON DELETE CASCADE,
			UNIQUE(mcp_id, uri)
		)`
	} else {
		query = `
		CREATE TABLE resources (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			original_name VARCHAR(255) NOT NULL,
			uri VARCHAR(767) NOT NULL,
			mcp_id INT NOT NULL,
			description TEXT,
			mime_type VARCHAR(100),
			is_enabled BOOLEAN DEFAULT true,
			tags VARCHAR(500),
			access_count INT DEFAULT 0,
			last_accessed_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (mcp_id) REFERENCES mcps(id) ON DELETE CASCADE,
			UNIQUE KEY unique_resource_per_mcp (mcp_id, uri),
			INDEX idx_name (name),
			INDEX idx_mcp_id (mcp_id),
			INDEX idx_uri (uri)
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

// createPromptsTable creates the prompts table (Discovered from MCPs)
func createPromptsTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE prompts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			original_name VARCHAR(255) NOT NULL,
			mcp_id INTEGER NOT NULL,
			description TEXT,
			template TEXT NOT NULL,
			arguments TEXT,
			is_enabled BOOLEAN DEFAULT 1,
			tags VARCHAR(500),
			use_count INTEGER DEFAULT 0,
			last_used_at DATETIME NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (mcp_id) REFERENCES mcps(id) ON DELETE CASCADE,
			UNIQUE(mcp_id, original_name)
		)`
	} else {
		query = `
		CREATE TABLE prompts (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			original_name VARCHAR(255) NOT NULL,
			mcp_id INT NOT NULL,
			description TEXT,
			template TEXT NOT NULL,
			arguments TEXT,
			is_enabled BOOLEAN DEFAULT true,
			tags VARCHAR(500),
			use_count INT DEFAULT 0,
			last_used_at TIMESTAMP NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (mcp_id) REFERENCES mcps(id) ON DELETE CASCADE,
			UNIQUE KEY unique_prompt_per_mcp (mcp_id, original_name),
			INDEX idx_name (name),
			INDEX idx_mcp_id (mcp_id)
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

// createGatewaysTable creates the gateways table
func createGatewaysTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE gateways (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			gateway_type VARCHAR(50) NOT NULL,
			config TEXT,
			is_enabled BOOLEAN DEFAULT 1,
			description TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`
	} else {
		query = `
		CREATE TABLE gateways (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			gateway_type VARCHAR(50) NOT NULL,
			config TEXT,
			is_enabled BOOLEAN DEFAULT true,
			description TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_slug (slug),
			INDEX idx_type (gateway_type)
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

// createServerToolsTable creates the server_tools table
func createServerToolsTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE server_tools (
			server_id INTEGER NOT NULL,
			tool_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (server_id, tool_id),
			FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE,
			FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE
		)`
	} else {
		query = `
		CREATE TABLE server_tools (
			server_id INT NOT NULL,
			tool_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (server_id, tool_id),
			FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE,
			FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE,
			INDEX idx_server (server_id),
			INDEX idx_tool (tool_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropServerToolsTable drops the server_tools table
func dropServerToolsTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS server_tools")
	return err
}

// createServerResourcesTable creates the server_resources table
func createServerResourcesTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE server_resources (
			server_id INTEGER NOT NULL,
			resource_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (server_id, resource_id),
			FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE,
			FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE
		)`
	} else {
		query = `
		CREATE TABLE server_resources (
			server_id INT NOT NULL,
			resource_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (server_id, resource_id),
			FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE,
			FOREIGN KEY (resource_id) REFERENCES resources(id) ON DELETE CASCADE,
			INDEX idx_server (server_id),
			INDEX idx_resource (resource_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropServerResourcesTable drops the server_resources table
func dropServerResourcesTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS server_resources")
	return err
}

// createServerPromptsTable creates the server_prompts table
func createServerPromptsTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE server_prompts (
			server_id INTEGER NOT NULL,
			prompt_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (server_id, prompt_id),
			FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE,
			FOREIGN KEY (prompt_id) REFERENCES prompts(id) ON DELETE CASCADE
		)`
	} else {
		query = `
		CREATE TABLE server_prompts (
			server_id INT NOT NULL,
			prompt_id INT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (server_id, prompt_id),
			FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE CASCADE,
			FOREIGN KEY (prompt_id) REFERENCES prompts(id) ON DELETE CASCADE,
			INDEX idx_server (server_id),
			INDEX idx_prompt (prompt_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	}

	_, err = db.Exec(query)
	return err
}

// dropServerPromptsTable drops the server_prompts table
func dropServerPromptsTable(db *sql.DB) error {
	_, err := db.Exec("DROP TABLE IF EXISTS server_prompts")
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
			user_id INTEGER,
			request_id VARCHAR(100),
			arguments TEXT,
			success BOOLEAN,
			response_time_ms INTEGER,
			error_message TEXT,
			server_id INTEGER,
			client_ip VARCHAR(45),
			user_agent VARCHAR(500),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
			FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE SET NULL
		)`
	} else {
		query = `
		CREATE TABLE tool_metrics (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			tool_id INT NOT NULL,
			user_id INT,
			request_id VARCHAR(100),
			arguments TEXT,
			success BOOLEAN,
			response_time_ms INT,
			error_message TEXT,
			server_id INT,
			client_ip VARCHAR(45),
			user_agent VARCHAR(500),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
			FOREIGN KEY (server_id) REFERENCES servers(id) ON DELETE SET NULL,
			INDEX idx_tool_id (tool_id),
			INDEX idx_user_id (user_id),
			INDEX idx_created_at (created_at),
			INDEX idx_success (success)
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

// createAuditsTable creates the audits table
func createAuditsTable(db *sql.DB) error {
	_, err := db.Exec("SELECT sqlite_version()")
	isSQLite := err == nil

	var query string

	if isSQLite {
		query = `
		CREATE TABLE audits (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			user_email VARCHAR(255),
			action VARCHAR(100) NOT NULL,
			entity_type VARCHAR(50) NOT NULL,
			entity_id INTEGER,
			entity_name VARCHAR(255),
			details TEXT,
			status VARCHAR(20),
			error_message TEXT,
			ip_address VARCHAR(45),
			user_agent VARCHAR(500),
			request_id VARCHAR(100),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
		)`
	} else {
		query = `
		CREATE TABLE audits (
			id BIGINT AUTO_INCREMENT PRIMARY KEY,
			user_id INT,
			user_email VARCHAR(255),
			action VARCHAR(100) NOT NULL,
			entity_type VARCHAR(50) NOT NULL,
			entity_id INT,
			entity_name VARCHAR(255),
			details TEXT,
			status VARCHAR(20),
			error_message TEXT,
			ip_address VARCHAR(45),
			user_agent VARCHAR(500),
			request_id VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
			INDEX idx_user_id (user_id),
			INDEX idx_action (action),
			INDEX idx_entity (entity_type, entity_id),
			INDEX idx_created_at (created_at)
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

// createUsersMetaTable creates the users_meta table
func createUsersMetaTable(db *sql.DB) error {
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
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
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
