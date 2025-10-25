// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Prompt represents a prompt discovered from MCP servers.
type Prompt struct {
	ID           int64
	Name         string
	OriginalName string
	MCPID        int64
	Description  *string
	Template     string
	Arguments    *string
	IsEnabled    bool
	Tags         *string
	UseCount     int
	LastUsedAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// PromptRepository handles database operations for prompts.
type PromptRepository struct {
	db *sql.DB
}

// NewPromptRepository creates a new prompt repository.
func NewPromptRepository(db *sql.DB) *PromptRepository {
	return &PromptRepository{db: db}
}

// Create inserts a new prompt into the database.
func (r *PromptRepository) Create(prompt *Prompt) error {
	result, err := r.db.Exec(
		`INSERT INTO prompts (name, original_name, mcp_id, description, template, arguments,
		is_enabled, tags, use_count, last_used_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		prompt.Name,
		prompt.OriginalName,
		prompt.MCPID,
		prompt.Description,
		prompt.Template,
		prompt.Arguments,
		prompt.IsEnabled,
		prompt.Tags,
		prompt.UseCount,
		prompt.LastUsedAt,
	)
	if err != nil {
		return err
	}

	prompt.ID, err = result.LastInsertId()
	return err
}

// GetByID retrieves a prompt by ID.
func (r *PromptRepository) GetByID(id int64) (*Prompt, error) {
	prompt := &Prompt{}
	err := r.db.QueryRow(
		`SELECT id, name, original_name, mcp_id, description, template, arguments,
		is_enabled, tags, use_count, last_used_at, created_at, updated_at
		FROM prompts WHERE id = ?`,
		id,
	).Scan(&prompt.ID, &prompt.Name, &prompt.OriginalName, &prompt.MCPID, &prompt.Description,
		&prompt.Template, &prompt.Arguments, &prompt.IsEnabled, &prompt.Tags, &prompt.UseCount,
		&prompt.LastUsedAt, &prompt.CreatedAt, &prompt.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return prompt, nil
}

// GetByMCPAndName retrieves a prompt by MCP ID and original name.
func (r *PromptRepository) GetByMCPAndName(mcpID int64, originalName string) (*Prompt, error) {
	prompt := &Prompt{}
	err := r.db.QueryRow(
		`SELECT id, name, original_name, mcp_id, description, template, arguments,
		is_enabled, tags, use_count, last_used_at, created_at, updated_at
		FROM prompts WHERE mcp_id = ? AND original_name = ?`,
		mcpID,
		originalName,
	).Scan(&prompt.ID, &prompt.Name, &prompt.OriginalName, &prompt.MCPID, &prompt.Description,
		&prompt.Template, &prompt.Arguments, &prompt.IsEnabled, &prompt.Tags, &prompt.UseCount,
		&prompt.LastUsedAt, &prompt.CreatedAt, &prompt.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return prompt, nil
}

// Update updates a prompt's information.
func (r *PromptRepository) Update(prompt *Prompt) error {
	_, err := r.db.Exec(
		`UPDATE prompts SET name = ?, original_name = ?, mcp_id = ?, description = ?,
		template = ?, arguments = ?, is_enabled = ?, tags = ?, use_count = ?,
		last_used_at = ?, updated_at = ? WHERE id = ?`,
		prompt.Name,
		prompt.OriginalName,
		prompt.MCPID,
		prompt.Description,
		prompt.Template,
		prompt.Arguments,
		prompt.IsEnabled,
		prompt.Tags,
		prompt.UseCount,
		prompt.LastUsedAt,
		time.Now(),
		prompt.ID,
	)
	return err
}

// UpdateUseMetrics updates the usage metrics for a prompt.
func (r *PromptRepository) UpdateUseMetrics(id int64) error {
	now := time.Now()
	_, err := r.db.Exec(
		`UPDATE prompts SET use_count = use_count + 1, last_used_at = ?,
		updated_at = ? WHERE id = ?`,
		now,
		now,
		id,
	)
	return err
}

// Delete removes a prompt from the database.
func (r *PromptRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM prompts WHERE id = ?", id)
	return err
}

// List retrieves all prompts with pagination.
func (r *PromptRepository) List(limit, offset int) ([]*Prompt, error) {
	rows, err := r.db.Query(
		`SELECT id, name, original_name, mcp_id, description, template, arguments,
		is_enabled, tags, use_count, last_used_at, created_at, updated_at
		FROM prompts ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPrompts(rows)
}

// ListByMCP retrieves all prompts for a specific MCP connection.
func (r *PromptRepository) ListByMCP(mcpID int64, limit, offset int) ([]*Prompt, error) {
	rows, err := r.db.Query(
		`SELECT id, name, original_name, mcp_id, description, template, arguments,
		is_enabled, tags, use_count, last_used_at, created_at, updated_at
		FROM prompts WHERE mcp_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		mcpID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanPrompts(rows)
}

// Count returns the total number of prompts.
func (r *PromptRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM prompts").Scan(&count)
	return count, err
}

func (r *PromptRepository) scanPrompts(rows *sql.Rows) ([]*Prompt, error) {
	var prompts []*Prompt
	for rows.Next() {
		prompt := &Prompt{}
		if err := rows.Scan(&prompt.ID, &prompt.Name, &prompt.OriginalName, &prompt.MCPID,
			&prompt.Description, &prompt.Template, &prompt.Arguments, &prompt.IsEnabled,
			&prompt.Tags, &prompt.UseCount, &prompt.LastUsedAt,
			&prompt.CreatedAt, &prompt.UpdatedAt); err != nil {
			return nil, err
		}
		prompts = append(prompts, prompt)
	}
	return prompts, rows.Err()
}

// PromptMeta represents metadata associated with a prompt.
type PromptMeta struct {
	ID        int64
	Key       string
	Value     string
	PromptID  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PromptMetaRepository handles database operations for prompt metadata.
type PromptMetaRepository struct {
	db *sql.DB
}

// NewPromptMetaRepository creates a new prompt meta repository.
func NewPromptMetaRepository(db *sql.DB) *PromptMetaRepository {
	return &PromptMetaRepository{db: db}
}

// Create inserts new metadata for a prompt.
func (r *PromptMetaRepository) Create(promptID int64, key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO prompts_meta (prompt_id, key, value) VALUES (?, ?, ?)",
		promptID,
		key,
		value,
	)
	return err
}

// Get retrieves metadata for a prompt by key.
func (r *PromptMetaRepository) Get(promptID int64, key string) (*PromptMeta, error) {
	meta := &PromptMeta{}
	err := r.db.QueryRow(
		"SELECT id, key, value, prompt_id, created_at, updated_at FROM prompts_meta WHERE prompt_id = ? AND key = ?",
		promptID,
		key,
	).Scan(&meta.ID, &meta.Key, &meta.Value, &meta.PromptID, &meta.CreatedAt, &meta.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return meta, nil
}

// Update updates metadata for a prompt.
func (r *PromptMetaRepository) Update(promptID int64, key, value string) error {
	_, err := r.db.Exec(
		"UPDATE prompts_meta SET value = ?, updated_at = ? WHERE prompt_id = ? AND key = ?",
		value,
		time.Now(),
		promptID,
		key,
	)
	return err
}

// Delete removes metadata for a prompt.
func (r *PromptMetaRepository) Delete(promptID int64, key string) error {
	_, err := r.db.Exec(
		"DELETE FROM prompts_meta WHERE prompt_id = ? AND key = ?",
		promptID,
		key,
	)
	return err
}

// ListByPrompt retrieves all metadata for a prompt.
func (r *PromptMetaRepository) ListByPrompt(promptID int64) ([]*PromptMeta, error) {
	rows, err := r.db.Query(
		"SELECT id, key, value, prompt_id, created_at, updated_at FROM prompts_meta WHERE prompt_id = ? ORDER BY key",
		promptID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []*PromptMeta
	for rows.Next() {
		meta := &PromptMeta{}
		if err := rows.Scan(&meta.ID, &meta.Key, &meta.Value, &meta.PromptID, &meta.CreatedAt, &meta.UpdatedAt); err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}

	return metadata, rows.Err()
}

// Upsert inserts or updates metadata for a prompt.
func (r *PromptMetaRepository) Upsert(promptID int64, key, value string) error {
	existing, err := r.Get(promptID, key)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(promptID, key, value)
	}

	return r.Update(promptID, key, value)
}
