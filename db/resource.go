// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Resource represents a resource discovered from MCP servers.
type Resource struct {
	ID             int64
	Name           string
	OriginalName   string
	URI            string
	MCPID          int64
	Description    *string
	MimeType       *string
	IsEnabled      bool
	Tags           *string
	AccessCount    int
	LastAccessedAt *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ResourceRepository handles database operations for resources.
type ResourceRepository struct {
	db *sql.DB
}

// NewResourceRepository creates a new resource repository.
func NewResourceRepository(db *sql.DB) *ResourceRepository {
	return &ResourceRepository{db: db}
}

// Create inserts a new resource into the database.
func (r *ResourceRepository) Create(resource *Resource) error {
	result, err := r.db.Exec(
		`INSERT INTO resources (name, original_name, uri, mcp_id, description, mime_type,
		is_enabled, tags, access_count, last_accessed_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		resource.Name,
		resource.OriginalName,
		resource.URI,
		resource.MCPID,
		resource.Description,
		resource.MimeType,
		resource.IsEnabled,
		resource.Tags,
		resource.AccessCount,
		resource.LastAccessedAt,
	)
	if err != nil {
		return err
	}

	resource.ID, err = result.LastInsertId()
	return err
}

// GetByID retrieves a resource by ID.
func (r *ResourceRepository) GetByID(id int64) (*Resource, error) {
	resource := &Resource{}
	err := r.db.QueryRow(
		`SELECT id, name, original_name, uri, mcp_id, description, mime_type, is_enabled,
		tags, access_count, last_accessed_at, created_at, updated_at
		FROM resources WHERE id = ?`,
		id,
	).Scan(&resource.ID, &resource.Name, &resource.OriginalName, &resource.URI, &resource.MCPID,
		&resource.Description, &resource.MimeType, &resource.IsEnabled, &resource.Tags,
		&resource.AccessCount, &resource.LastAccessedAt, &resource.CreatedAt, &resource.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// GetByMCPAndURI retrieves a resource by MCP ID and URI.
func (r *ResourceRepository) GetByMCPAndURI(mcpID int64, uri string) (*Resource, error) {
	resource := &Resource{}
	err := r.db.QueryRow(
		`SELECT id, name, original_name, uri, mcp_id, description, mime_type, is_enabled,
		tags, access_count, last_accessed_at, created_at, updated_at
		FROM resources WHERE mcp_id = ? AND uri = ?`,
		mcpID,
		uri,
	).Scan(&resource.ID, &resource.Name, &resource.OriginalName, &resource.URI, &resource.MCPID,
		&resource.Description, &resource.MimeType, &resource.IsEnabled, &resource.Tags,
		&resource.AccessCount, &resource.LastAccessedAt, &resource.CreatedAt, &resource.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// Update updates a resource's information.
func (r *ResourceRepository) Update(resource *Resource) error {
	_, err := r.db.Exec(
		`UPDATE resources SET name = ?, original_name = ?, uri = ?, mcp_id = ?,
		description = ?, mime_type = ?, is_enabled = ?, tags = ?, access_count = ?,
		last_accessed_at = ?, updated_at = ? WHERE id = ?`,
		resource.Name,
		resource.OriginalName,
		resource.URI,
		resource.MCPID,
		resource.Description,
		resource.MimeType,
		resource.IsEnabled,
		resource.Tags,
		resource.AccessCount,
		resource.LastAccessedAt,
		time.Now(),
		resource.ID,
	)
	return err
}

// UpdateAccessMetrics updates the access metrics for a resource.
func (r *ResourceRepository) UpdateAccessMetrics(id int64) error {
	now := time.Now()
	_, err := r.db.Exec(
		`UPDATE resources SET access_count = access_count + 1, last_accessed_at = ?,
		updated_at = ? WHERE id = ?`,
		now,
		now,
		id,
	)
	return err
}

// Delete removes a resource from the database.
func (r *ResourceRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM resources WHERE id = ?", id)
	return err
}

// List retrieves all resources with pagination.
func (r *ResourceRepository) List(limit, offset int) ([]*Resource, error) {
	rows, err := r.db.Query(
		`SELECT id, name, original_name, uri, mcp_id, description, mime_type, is_enabled,
		tags, access_count, last_accessed_at, created_at, updated_at
		FROM resources ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanResources(rows)
}

// ListByMCP retrieves all resources for a specific MCP connection.
func (r *ResourceRepository) ListByMCP(mcpID int64, limit, offset int) ([]*Resource, error) {
	rows, err := r.db.Query(
		`SELECT id, name, original_name, uri, mcp_id, description, mime_type, is_enabled,
		tags, access_count, last_accessed_at, created_at, updated_at
		FROM resources WHERE mcp_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		mcpID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanResources(rows)
}

// Count returns the total number of resources.
func (r *ResourceRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM resources").Scan(&count)
	return count, err
}

func (r *ResourceRepository) scanResources(rows *sql.Rows) ([]*Resource, error) {
	var resources []*Resource
	for rows.Next() {
		resource := &Resource{}
		if err := rows.Scan(&resource.ID, &resource.Name, &resource.OriginalName, &resource.URI,
			&resource.MCPID, &resource.Description, &resource.MimeType, &resource.IsEnabled,
			&resource.Tags, &resource.AccessCount, &resource.LastAccessedAt,
			&resource.CreatedAt, &resource.UpdatedAt); err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	return resources, rows.Err()
}

// ResourceMeta represents metadata associated with a resource.
type ResourceMeta struct {
	ID         int64
	Key        string
	Value      string
	ResourceID int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ResourceMetaRepository handles database operations for resource metadata.
type ResourceMetaRepository struct {
	db *sql.DB
}

// NewResourceMetaRepository creates a new resource meta repository.
func NewResourceMetaRepository(db *sql.DB) *ResourceMetaRepository {
	return &ResourceMetaRepository{db: db}
}

// Create inserts new metadata for a resource.
func (r *ResourceMetaRepository) Create(resourceID int64, key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO resources_meta (resource_id, key, value) VALUES (?, ?, ?)",
		resourceID,
		key,
		value,
	)
	return err
}

// Get retrieves metadata for a resource by key.
func (r *ResourceMetaRepository) Get(resourceID int64, key string) (*ResourceMeta, error) {
	meta := &ResourceMeta{}
	err := r.db.QueryRow(
		"SELECT id, key, value, resource_id, created_at, updated_at FROM resources_meta WHERE resource_id = ? AND key = ?",
		resourceID,
		key,
	).Scan(&meta.ID, &meta.Key, &meta.Value, &meta.ResourceID, &meta.CreatedAt, &meta.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return meta, nil
}

// Update updates metadata for a resource.
func (r *ResourceMetaRepository) Update(resourceID int64, key, value string) error {
	_, err := r.db.Exec(
		"UPDATE resources_meta SET value = ?, updated_at = ? WHERE resource_id = ? AND key = ?",
		value,
		time.Now(),
		resourceID,
		key,
	)
	return err
}

// Delete removes metadata for a resource.
func (r *ResourceMetaRepository) Delete(resourceID int64, key string) error {
	_, err := r.db.Exec(
		"DELETE FROM resources_meta WHERE resource_id = ? AND key = ?",
		resourceID,
		key,
	)
	return err
}

// ListByResource retrieves all metadata for a resource.
func (r *ResourceMetaRepository) ListByResource(resourceID int64) ([]*ResourceMeta, error) {
	rows, err := r.db.Query(
		"SELECT id, key, value, resource_id, created_at, updated_at FROM resources_meta WHERE resource_id = ? ORDER BY key",
		resourceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []*ResourceMeta
	for rows.Next() {
		meta := &ResourceMeta{}
		if err := rows.Scan(&meta.ID, &meta.Key, &meta.Value, &meta.ResourceID, &meta.CreatedAt, &meta.UpdatedAt); err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}

	return metadata, rows.Err()
}

// Upsert inserts or updates metadata for a resource.
func (r *ResourceMetaRepository) Upsert(resourceID int64, key, value string) error {
	existing, err := r.Get(resourceID, key)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(resourceID, key, value)
	}

	return r.Update(resourceID, key, value)
}
