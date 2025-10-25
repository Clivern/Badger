// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Gateway represents a gateway in the database.
type Gateway struct {
	ID          int64
	Name        string
	Slug        string
	GatewayType string
	Config      *string
	IsEnabled   bool
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// GatewayRepository handles database operations for gateways.
type GatewayRepository struct {
	db *sql.DB
}

// NewGatewayRepository creates a new gateway repository.
func NewGatewayRepository(db *sql.DB) *GatewayRepository {
	return &GatewayRepository{db: db}
}

// Create inserts a new gateway into the database.
func (r *GatewayRepository) Create(gateway *Gateway) error {
	result, err := r.db.Exec(
		`INSERT INTO gateways (name, slug, gateway_type, config, is_enabled, description)
		VALUES (?, ?, ?, ?, ?, ?)`,
		gateway.Name,
		gateway.Slug,
		gateway.GatewayType,
		gateway.Config,
		gateway.IsEnabled,
		gateway.Description,
	)
	if err != nil {
		return err
	}

	gateway.ID, err = result.LastInsertId()
	return err
}

// GetByID retrieves a gateway by ID.
func (r *GatewayRepository) GetByID(id int64) (*Gateway, error) {
	gateway := &Gateway{}
	err := r.db.QueryRow(
		`SELECT id, name, slug, gateway_type, config, is_enabled, description, created_at, updated_at
		FROM gateways WHERE id = ?`,
		id,
	).Scan(&gateway.ID, &gateway.Name, &gateway.Slug, &gateway.GatewayType, &gateway.Config,
		&gateway.IsEnabled, &gateway.Description, &gateway.CreatedAt, &gateway.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return gateway, nil
}

// GetBySlug retrieves a gateway by slug.
func (r *GatewayRepository) GetBySlug(slug string) (*Gateway, error) {
	gateway := &Gateway{}
	err := r.db.QueryRow(
		`SELECT id, name, slug, gateway_type, config, is_enabled, description, created_at, updated_at
		FROM gateways WHERE slug = ?`,
		slug,
	).Scan(&gateway.ID, &gateway.Name, &gateway.Slug, &gateway.GatewayType, &gateway.Config,
		&gateway.IsEnabled, &gateway.Description, &gateway.CreatedAt, &gateway.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return gateway, nil
}

// Update updates a gateway's information.
func (r *GatewayRepository) Update(gateway *Gateway) error {
	_, err := r.db.Exec(
		`UPDATE gateways SET name = ?, slug = ?, gateway_type = ?, config = ?,
		is_enabled = ?, description = ?, updated_at = ? WHERE id = ?`,
		gateway.Name,
		gateway.Slug,
		gateway.GatewayType,
		gateway.Config,
		gateway.IsEnabled,
		gateway.Description,
		time.Now(),
		gateway.ID,
	)
	return err
}

// Delete removes a gateway from the database.
func (r *GatewayRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM gateways WHERE id = ?", id)
	return err
}

// List retrieves all gateways with pagination.
func (r *GatewayRepository) List(limit, offset int) ([]*Gateway, error) {
	rows, err := r.db.Query(
		`SELECT id, name, slug, gateway_type, config, is_enabled, description, created_at, updated_at
		FROM gateways ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gateways []*Gateway
	for rows.Next() {
		gateway := &Gateway{}
		if err := rows.Scan(&gateway.ID, &gateway.Name, &gateway.Slug, &gateway.GatewayType,
			&gateway.Config, &gateway.IsEnabled, &gateway.Description, &gateway.CreatedAt,
			&gateway.UpdatedAt); err != nil {
			return nil, err
		}
		gateways = append(gateways, gateway)
	}

	return gateways, rows.Err()
}

// ListByType retrieves all gateways of a specific type.
func (r *GatewayRepository) ListByType(gatewayType string, limit, offset int) ([]*Gateway, error) {
	rows, err := r.db.Query(
		`SELECT id, name, slug, gateway_type, config, is_enabled, description, created_at, updated_at
		FROM gateways WHERE gateway_type = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		gatewayType,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gateways []*Gateway
	for rows.Next() {
		gateway := &Gateway{}
		if err := rows.Scan(&gateway.ID, &gateway.Name, &gateway.Slug, &gateway.GatewayType,
			&gateway.Config, &gateway.IsEnabled, &gateway.Description, &gateway.CreatedAt,
			&gateway.UpdatedAt); err != nil {
			return nil, err
		}
		gateways = append(gateways, gateway)
	}

	return gateways, rows.Err()
}

// Count returns the total number of gateways.
func (r *GatewayRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM gateways").Scan(&count)
	return count, err
}

// GatewayMeta represents metadata associated with a gateway.
type GatewayMeta struct {
	ID        int64
	Key       string
	Value     string
	GatewayID int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GatewayMetaRepository handles database operations for gateway metadata.
type GatewayMetaRepository struct {
	db *sql.DB
}

// NewGatewayMetaRepository creates a new gateway meta repository.
func NewGatewayMetaRepository(db *sql.DB) *GatewayMetaRepository {
	return &GatewayMetaRepository{db: db}
}

// Create inserts new metadata for a gateway.
func (r *GatewayMetaRepository) Create(gatewayID int64, key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO gateways_meta (gateway_id, key, value) VALUES (?, ?, ?)",
		gatewayID,
		key,
		value,
	)
	return err
}

// Get retrieves metadata for a gateway by key.
func (r *GatewayMetaRepository) Get(gatewayID int64, key string) (*GatewayMeta, error) {
	meta := &GatewayMeta{}
	err := r.db.QueryRow(
		"SELECT id, key, value, gateway_id, created_at, updated_at FROM gateways_meta WHERE gateway_id = ? AND key = ?",
		gatewayID,
		key,
	).Scan(&meta.ID, &meta.Key, &meta.Value, &meta.GatewayID, &meta.CreatedAt, &meta.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return meta, nil
}

// Update updates metadata for a gateway.
func (r *GatewayMetaRepository) Update(gatewayID int64, key, value string) error {
	_, err := r.db.Exec(
		"UPDATE gateways_meta SET value = ?, updated_at = ? WHERE gateway_id = ? AND key = ?",
		value,
		time.Now(),
		gatewayID,
		key,
	)
	return err
}

// Delete removes metadata for a gateway.
func (r *GatewayMetaRepository) Delete(gatewayID int64, key string) error {
	_, err := r.db.Exec(
		"DELETE FROM gateways_meta WHERE gateway_id = ? AND key = ?",
		gatewayID,
		key,
	)
	return err
}

// ListByGateway retrieves all metadata for a gateway.
func (r *GatewayMetaRepository) ListByGateway(gatewayID int64) ([]*GatewayMeta, error) {
	rows, err := r.db.Query(
		"SELECT id, key, value, gateway_id, created_at, updated_at FROM gateways_meta WHERE gateway_id = ? ORDER BY key",
		gatewayID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []*GatewayMeta
	for rows.Next() {
		meta := &GatewayMeta{}
		if err := rows.Scan(&meta.ID, &meta.Key, &meta.Value, &meta.GatewayID, &meta.CreatedAt, &meta.UpdatedAt); err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}

	return metadata, rows.Err()
}

// Upsert inserts or updates metadata for a gateway.
func (r *GatewayMetaRepository) Upsert(gatewayID int64, key, value string) error {
	existing, err := r.Get(gatewayID, key)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(gatewayID, key, value)
	}

	return r.Update(gatewayID, key, value)
}
