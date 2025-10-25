// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// MCP represents a backend MCP server connection in the database.
type MCP struct {
	ID                int64
	Name              string
	Slug              string
	URL               string
	Transport         string
	AuthType          string
	AuthToken         *string
	TimeoutMs         int
	MaxRetries        int
	Headers           *string
	Status            string
	HealthCheckURL    *string
	LastHealthCheckAt *time.Time
	HealthStatus      string
	Capabilities      *string
	ProtocolVersion   *string
	Description       *string
	Tags              *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// MCPRepository handles database operations for MCP connections.
type MCPRepository struct {
	db *sql.DB
}

// NewMCPRepository creates a new MCP repository.
func NewMCPRepository(db *sql.DB) *MCPRepository {
	return &MCPRepository{db: db}
}

// Create inserts a new MCP connection into the database.
func (r *MCPRepository) Create(mcp *MCP) error {
	result, err := r.db.Exec(
		`INSERT INTO mcps (name, slug, url, transport, auth_type, auth_token, timeout_ms,
		max_retries, headers, status, health_check_url, last_health_check_at, health_status,
		capabilities, protocol_version, description, tags)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		mcp.Name,
		mcp.Slug,
		mcp.URL,
		mcp.Transport,
		mcp.AuthType,
		mcp.AuthToken,
		mcp.TimeoutMs,
		mcp.MaxRetries,
		mcp.Headers,
		mcp.Status,
		mcp.HealthCheckURL,
		mcp.LastHealthCheckAt,
		mcp.HealthStatus,
		mcp.Capabilities,
		mcp.ProtocolVersion,
		mcp.Description,
		mcp.Tags,
	)
	if err != nil {
		return err
	}

	mcp.ID, err = result.LastInsertId()
	return err
}

// GetByID retrieves an MCP connection by ID.
func (r *MCPRepository) GetByID(id int64) (*MCP, error) {
	mcp := &MCP{}
	err := r.db.QueryRow(
		`SELECT id, name, slug, url, transport, auth_type, auth_token, timeout_ms, max_retries,
		headers, status, health_check_url, last_health_check_at, health_status, capabilities,
		protocol_version, description, tags, created_at, updated_at
		FROM mcps WHERE id = ?`,
		id,
	).Scan(&mcp.ID, &mcp.Name, &mcp.Slug, &mcp.URL, &mcp.Transport, &mcp.AuthType, &mcp.AuthToken,
		&mcp.TimeoutMs, &mcp.MaxRetries, &mcp.Headers, &mcp.Status, &mcp.HealthCheckURL,
		&mcp.LastHealthCheckAt, &mcp.HealthStatus, &mcp.Capabilities, &mcp.ProtocolVersion,
		&mcp.Description, &mcp.Tags, &mcp.CreatedAt, &mcp.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mcp, nil
}

// GetBySlug retrieves an MCP connection by slug.
func (r *MCPRepository) GetBySlug(slug string) (*MCP, error) {
	mcp := &MCP{}
	err := r.db.QueryRow(
		`SELECT id, name, slug, url, transport, auth_type, auth_token, timeout_ms, max_retries,
		headers, status, health_check_url, last_health_check_at, health_status, capabilities,
		protocol_version, description, tags, created_at, updated_at
		FROM mcps WHERE slug = ?`,
		slug,
	).Scan(&mcp.ID, &mcp.Name, &mcp.Slug, &mcp.URL, &mcp.Transport, &mcp.AuthType, &mcp.AuthToken,
		&mcp.TimeoutMs, &mcp.MaxRetries, &mcp.Headers, &mcp.Status, &mcp.HealthCheckURL,
		&mcp.LastHealthCheckAt, &mcp.HealthStatus, &mcp.Capabilities, &mcp.ProtocolVersion,
		&mcp.Description, &mcp.Tags, &mcp.CreatedAt, &mcp.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return mcp, nil
}

// Update updates an MCP connection's information.
func (r *MCPRepository) Update(mcp *MCP) error {
	_, err := r.db.Exec(
		`UPDATE mcps SET name = ?, slug = ?, url = ?, transport = ?, auth_type = ?,
		auth_token = ?, timeout_ms = ?, max_retries = ?, headers = ?, status = ?,
		health_check_url = ?, last_health_check_at = ?, health_status = ?, capabilities = ?,
		protocol_version = ?, description = ?, tags = ?, updated_at = ? WHERE id = ?`,
		mcp.Name,
		mcp.Slug,
		mcp.URL,
		mcp.Transport,
		mcp.AuthType,
		mcp.AuthToken,
		mcp.TimeoutMs,
		mcp.MaxRetries,
		mcp.Headers,
		mcp.Status,
		mcp.HealthCheckURL,
		mcp.LastHealthCheckAt,
		mcp.HealthStatus,
		mcp.Capabilities,
		mcp.ProtocolVersion,
		mcp.Description,
		mcp.Tags,
		time.Now(),
		mcp.ID,
	)
	return err
}

// UpdateHealthStatus updates the health status of an MCP connection.
func (r *MCPRepository) UpdateHealthStatus(id int64, healthStatus string) error {
	now := time.Now()
	_, err := r.db.Exec(
		`UPDATE mcps SET health_status = ?, last_health_check_at = ?, updated_at = ? WHERE id = ?`,
		healthStatus,
		now,
		now,
		id,
	)
	return err
}

// Delete removes an MCP connection from the database.
func (r *MCPRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM mcps WHERE id = ?", id)
	return err
}

// List retrieves all MCP connections with pagination.
func (r *MCPRepository) List(limit, offset int) ([]*MCP, error) {
	rows, err := r.db.Query(
		`SELECT id, name, slug, url, transport, auth_type, auth_token, timeout_ms, max_retries,
		headers, status, health_check_url, last_health_check_at, health_status, capabilities,
		protocol_version, description, tags, created_at, updated_at
		FROM mcps ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mcps []*MCP
	for rows.Next() {
		mcp := &MCP{}
		if err := rows.Scan(&mcp.ID, &mcp.Name, &mcp.Slug, &mcp.URL, &mcp.Transport, &mcp.AuthType,
			&mcp.AuthToken, &mcp.TimeoutMs, &mcp.MaxRetries, &mcp.Headers, &mcp.Status,
			&mcp.HealthCheckURL, &mcp.LastHealthCheckAt, &mcp.HealthStatus, &mcp.Capabilities,
			&mcp.ProtocolVersion, &mcp.Description, &mcp.Tags, &mcp.CreatedAt, &mcp.UpdatedAt); err != nil {
			return nil, err
		}
		mcps = append(mcps, mcp)
	}

	return mcps, rows.Err()
}

// ListByStatus retrieves all MCP connections with a specific status.
func (r *MCPRepository) ListByStatus(status string, limit, offset int) ([]*MCP, error) {
	rows, err := r.db.Query(
		`SELECT id, name, slug, url, transport, auth_type, auth_token, timeout_ms, max_retries,
		headers, status, health_check_url, last_health_check_at, health_status, capabilities,
		protocol_version, description, tags, created_at, updated_at
		FROM mcps WHERE status = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		status,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mcps []*MCP
	for rows.Next() {
		mcp := &MCP{}
		if err := rows.Scan(&mcp.ID, &mcp.Name, &mcp.Slug, &mcp.URL, &mcp.Transport, &mcp.AuthType,
			&mcp.AuthToken, &mcp.TimeoutMs, &mcp.MaxRetries, &mcp.Headers, &mcp.Status,
			&mcp.HealthCheckURL, &mcp.LastHealthCheckAt, &mcp.HealthStatus, &mcp.Capabilities,
			&mcp.ProtocolVersion, &mcp.Description, &mcp.Tags, &mcp.CreatedAt, &mcp.UpdatedAt); err != nil {
			return nil, err
		}
		mcps = append(mcps, mcp)
	}

	return mcps, rows.Err()
}

// Count returns the total number of MCP connections.
func (r *MCPRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM mcps").Scan(&count)
	return count, err
}

// MCPMeta represents metadata associated with an MCP connection.
type MCPMeta struct {
	ID        int64
	Key       string
	Value     string
	MCPID     int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MCPMetaRepository handles database operations for MCP metadata.
type MCPMetaRepository struct {
	db *sql.DB
}

// NewMCPMetaRepository creates a new MCP meta repository.
func NewMCPMetaRepository(db *sql.DB) *MCPMetaRepository {
	return &MCPMetaRepository{db: db}
}

// Create inserts new metadata for an MCP connection.
func (r *MCPMetaRepository) Create(mcpID int64, key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO mcps_meta (mcp_id, key, value) VALUES (?, ?, ?)",
		mcpID,
		key,
		value,
	)
	return err
}

// Get retrieves metadata for an MCP connection by key.
func (r *MCPMetaRepository) Get(mcpID int64, key string) (*MCPMeta, error) {
	meta := &MCPMeta{}
	err := r.db.QueryRow(
		"SELECT id, key, value, mcp_id, created_at, updated_at FROM mcps_meta WHERE mcp_id = ? AND key = ?",
		mcpID,
		key,
	).Scan(&meta.ID, &meta.Key, &meta.Value, &meta.MCPID, &meta.CreatedAt, &meta.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return meta, nil
}

// Update updates metadata for an MCP connection.
func (r *MCPMetaRepository) Update(mcpID int64, key, value string) error {
	_, err := r.db.Exec(
		"UPDATE mcps_meta SET value = ?, updated_at = ? WHERE mcp_id = ? AND key = ?",
		value,
		time.Now(),
		mcpID,
		key,
	)
	return err
}

// Delete removes metadata for an MCP connection.
func (r *MCPMetaRepository) Delete(mcpID int64, key string) error {
	_, err := r.db.Exec(
		"DELETE FROM mcps_meta WHERE mcp_id = ? AND key = ?",
		mcpID,
		key,
	)
	return err
}

// ListByMCP retrieves all metadata for an MCP connection.
func (r *MCPMetaRepository) ListByMCP(mcpID int64) ([]*MCPMeta, error) {
	rows, err := r.db.Query(
		"SELECT id, key, value, mcp_id, created_at, updated_at FROM mcps_meta WHERE mcp_id = ? ORDER BY key",
		mcpID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []*MCPMeta
	for rows.Next() {
		meta := &MCPMeta{}
		if err := rows.Scan(&meta.ID, &meta.Key, &meta.Value, &meta.MCPID, &meta.CreatedAt, &meta.UpdatedAt); err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}

	return metadata, rows.Err()
}

// Upsert inserts or updates metadata for an MCP connection.
func (r *MCPMetaRepository) Upsert(mcpID int64, key, value string) error {
	existing, err := r.Get(mcpID, key)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(mcpID, key, value)
	}

	return r.Update(mcpID, key, value)
}
