// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Tool represents a tool discovered from MCP servers.
type Tool struct {
	ID                int64
	Name              string
	OriginalName      string
	MCPID             int64
	Description       *string
	InputSchema       string
	IsEnabled         bool
	TimeoutMs         int
	MaxRetries        int
	Tags              *string
	Category          *string
	CallCount         int
	LastCalledAt      *time.Time
	AvgResponseTimeMs int
	ErrorCount        int
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// ToolRepository handles database operations for tools.
type ToolRepository struct {
	db *sql.DB
}

// NewToolRepository creates a new tool repository.
func NewToolRepository(db *sql.DB) *ToolRepository {
	return &ToolRepository{db: db}
}

// Create inserts a new tool into the database.
func (r *ToolRepository) Create(tool *Tool) error {
	result, err := r.db.Exec(
		`INSERT INTO tools (name, original_name, mcp_id, description, input_schema, is_enabled,
		timeout_ms, max_retries, tags, category, call_count, last_called_at, avg_response_time_ms,
		error_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		tool.Name,
		tool.OriginalName,
		tool.MCPID,
		tool.Description,
		tool.InputSchema,
		tool.IsEnabled,
		tool.TimeoutMs,
		tool.MaxRetries,
		tool.Tags,
		tool.Category,
		tool.CallCount,
		tool.LastCalledAt,
		tool.AvgResponseTimeMs,
		tool.ErrorCount,
	)
	if err != nil {
		return err
	}

	tool.ID, err = result.LastInsertId()
	return err
}

// GetByID retrieves a tool by ID.
func (r *ToolRepository) GetByID(id int64) (*Tool, error) {
	tool := &Tool{}
	err := r.db.QueryRow(
		`SELECT id, name, original_name, mcp_id, description, input_schema, is_enabled,
		timeout_ms, max_retries, tags, category, call_count, last_called_at, avg_response_time_ms,
		error_count, created_at, updated_at FROM tools WHERE id = ?`,
		id,
	).Scan(&tool.ID, &tool.Name, &tool.OriginalName, &tool.MCPID, &tool.Description,
		&tool.InputSchema, &tool.IsEnabled, &tool.TimeoutMs, &tool.MaxRetries, &tool.Tags,
		&tool.Category, &tool.CallCount, &tool.LastCalledAt, &tool.AvgResponseTimeMs,
		&tool.ErrorCount, &tool.CreatedAt, &tool.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return tool, nil
}

// GetByMCPAndName retrieves a tool by MCP ID and original name.
func (r *ToolRepository) GetByMCPAndName(mcpID int64, originalName string) (*Tool, error) {
	tool := &Tool{}
	err := r.db.QueryRow(
		`SELECT id, name, original_name, mcp_id, description, input_schema, is_enabled,
		timeout_ms, max_retries, tags, category, call_count, last_called_at, avg_response_time_ms,
		error_count, created_at, updated_at FROM tools WHERE mcp_id = ? AND original_name = ?`,
		mcpID,
		originalName,
	).Scan(&tool.ID, &tool.Name, &tool.OriginalName, &tool.MCPID, &tool.Description,
		&tool.InputSchema, &tool.IsEnabled, &tool.TimeoutMs, &tool.MaxRetries, &tool.Tags,
		&tool.Category, &tool.CallCount, &tool.LastCalledAt, &tool.AvgResponseTimeMs,
		&tool.ErrorCount, &tool.CreatedAt, &tool.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return tool, nil
}

// Update updates a tool's information.
func (r *ToolRepository) Update(tool *Tool) error {
	_, err := r.db.Exec(
		`UPDATE tools SET name = ?, original_name = ?, mcp_id = ?, description = ?,
		input_schema = ?, is_enabled = ?, timeout_ms = ?, max_retries = ?, tags = ?,
		category = ?, call_count = ?, last_called_at = ?, avg_response_time_ms = ?,
		error_count = ?, updated_at = ? WHERE id = ?`,
		tool.Name,
		tool.OriginalName,
		tool.MCPID,
		tool.Description,
		tool.InputSchema,
		tool.IsEnabled,
		tool.TimeoutMs,
		tool.MaxRetries,
		tool.Tags,
		tool.Category,
		tool.CallCount,
		tool.LastCalledAt,
		tool.AvgResponseTimeMs,
		tool.ErrorCount,
		time.Now(),
		tool.ID,
	)
	return err
}

// UpdateCallMetrics updates the call metrics for a tool.
func (r *ToolRepository) UpdateCallMetrics(id int64, responseTimeMs int, success bool) error {
	// Get current values
	tool, err := r.GetByID(id)
	if err != nil || tool == nil {
		return err
	}

	// Calculate new average
	newCallCount := tool.CallCount + 1
	newAvg := ((tool.AvgResponseTimeMs * tool.CallCount) + responseTimeMs) / newCallCount
	newErrorCount := tool.ErrorCount
	if !success {
		newErrorCount++
	}

	now := time.Now()
	_, err = r.db.Exec(
		`UPDATE tools SET call_count = ?, last_called_at = ?, avg_response_time_ms = ?,
		error_count = ?, updated_at = ? WHERE id = ?`,
		newCallCount,
		now,
		newAvg,
		newErrorCount,
		now,
		id,
	)
	return err
}

// Delete removes a tool from the database.
func (r *ToolRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM tools WHERE id = ?", id)
	return err
}

// List retrieves all tools with pagination.
func (r *ToolRepository) List(limit, offset int) ([]*Tool, error) {
	rows, err := r.db.Query(
		`SELECT id, name, original_name, mcp_id, description, input_schema, is_enabled,
		timeout_ms, max_retries, tags, category, call_count, last_called_at, avg_response_time_ms,
		error_count, created_at, updated_at FROM tools ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTools(rows)
}

// ListByMCP retrieves all tools for a specific MCP connection.
func (r *ToolRepository) ListByMCP(mcpID int64, limit, offset int) ([]*Tool, error) {
	rows, err := r.db.Query(
		`SELECT id, name, original_name, mcp_id, description, input_schema, is_enabled,
		timeout_ms, max_retries, tags, category, call_count, last_called_at, avg_response_time_ms,
		error_count, created_at, updated_at FROM tools WHERE mcp_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		mcpID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTools(rows)
}

// ListByCategory retrieves all tools in a specific category.
func (r *ToolRepository) ListByCategory(category string, limit, offset int) ([]*Tool, error) {
	rows, err := r.db.Query(
		`SELECT id, name, original_name, mcp_id, description, input_schema, is_enabled,
		timeout_ms, max_retries, tags, category, call_count, last_called_at, avg_response_time_ms,
		error_count, created_at, updated_at FROM tools WHERE category = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		category,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTools(rows)
}

// Count returns the total number of tools.
func (r *ToolRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM tools").Scan(&count)
	return count, err
}

func (r *ToolRepository) scanTools(rows *sql.Rows) ([]*Tool, error) {
	var tools []*Tool
	for rows.Next() {
		tool := &Tool{}
		if err := rows.Scan(&tool.ID, &tool.Name, &tool.OriginalName, &tool.MCPID,
			&tool.Description, &tool.InputSchema, &tool.IsEnabled, &tool.TimeoutMs,
			&tool.MaxRetries, &tool.Tags, &tool.Category, &tool.CallCount,
			&tool.LastCalledAt, &tool.AvgResponseTimeMs, &tool.ErrorCount,
			&tool.CreatedAt, &tool.UpdatedAt); err != nil {
			return nil, err
		}
		tools = append(tools, tool)
	}
	return tools, rows.Err()
}

// ToolMeta represents metadata associated with a tool.
type ToolMeta struct {
	ID        int64
	Key       string
	Value     string
	ToolID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ToolMetaRepository handles database operations for tool metadata.
type ToolMetaRepository struct {
	db *sql.DB
}

// NewToolMetaRepository creates a new tool meta repository.
func NewToolMetaRepository(db *sql.DB) *ToolMetaRepository {
	return &ToolMetaRepository{db: db}
}

// Create inserts new metadata for a tool.
func (r *ToolMetaRepository) Create(toolID int64, key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO tools_meta (tool_id, key, value) VALUES (?, ?, ?)",
		toolID,
		key,
		value,
	)
	return err
}

// Get retrieves metadata for a tool by key.
func (r *ToolMetaRepository) Get(toolID int64, key string) (*ToolMeta, error) {
	meta := &ToolMeta{}
	err := r.db.QueryRow(
		"SELECT id, key, value, tool_id, created_at, updated_at FROM tools_meta WHERE tool_id = ? AND key = ?",
		toolID,
		key,
	).Scan(&meta.ID, &meta.Key, &meta.Value, &meta.ToolID, &meta.CreatedAt, &meta.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return meta, nil
}

// Update updates metadata for a tool.
func (r *ToolMetaRepository) Update(toolID int64, key, value string) error {
	_, err := r.db.Exec(
		"UPDATE tools_meta SET value = ?, updated_at = ? WHERE tool_id = ? AND key = ?",
		value,
		time.Now(),
		toolID,
		key,
	)
	return err
}

// Delete removes metadata for a tool.
func (r *ToolMetaRepository) Delete(toolID int64, key string) error {
	_, err := r.db.Exec(
		"DELETE FROM tools_meta WHERE tool_id = ? AND key = ?",
		toolID,
		key,
	)
	return err
}

// ListByTool retrieves all metadata for a tool.
func (r *ToolMetaRepository) ListByTool(toolID int64) ([]*ToolMeta, error) {
	rows, err := r.db.Query(
		"SELECT id, key, value, tool_id, created_at, updated_at FROM tools_meta WHERE tool_id = ? ORDER BY key",
		toolID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []*ToolMeta
	for rows.Next() {
		meta := &ToolMeta{}
		if err := rows.Scan(&meta.ID, &meta.Key, &meta.Value, &meta.ToolID, &meta.CreatedAt, &meta.UpdatedAt); err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}

	return metadata, rows.Err()
}

// Upsert inserts or updates metadata for a tool.
func (r *ToolMetaRepository) Upsert(toolID int64, key, value string) error {
	existing, err := r.Get(toolID, key)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(toolID, key, value)
	}

	return r.Update(toolID, key, value)
}
