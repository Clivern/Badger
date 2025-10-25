// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// User represents a user in the database.
//
// Users can have different roles: admin, user, or readonly.
type User struct {
	ID          int64
	Email       string
	Password    string
	Role        string
	APIKey      *string
	IsActive    bool
	LastLoginAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserRepository handles database operations for users.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user into the database.
//
// Example:
//
//	user := &User{
//		Email:    "user@example.com",
//		Password: "hashed_password",
//		Role:     "user",
//		IsActive: true,
//	}
//	err := repo.Create(user)
func (r *UserRepository) Create(user *User) error {
	result, err := r.db.Exec(
		`INSERT INTO users (email, password, role, api_key, is_active, last_login_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		user.Email,
		user.Password,
		user.Role,
		user.APIKey,
		user.IsActive,
		user.LastLoginAt,
	)
	if err != nil {
		return err
	}

	user.ID, err = result.LastInsertId()
	return err
}

// GetByID retrieves a user by ID.
func (r *UserRepository) GetByID(id int64) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(
		`SELECT id, email, password, role, api_key, is_active, last_login_at, created_at, updated_at
		FROM users WHERE id = ?`,
		id,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.APIKey, &user.IsActive,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by email.
func (r *UserRepository) GetByEmail(email string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(
		`SELECT id, email, password, role, api_key, is_active, last_login_at, created_at, updated_at
		FROM users WHERE email = ?`,
		email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.APIKey, &user.IsActive,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByAPIKey retrieves a user by API key.
func (r *UserRepository) GetByAPIKey(apiKey string) (*User, error) {
	user := &User{}
	err := r.db.QueryRow(
		`SELECT id, email, password, role, api_key, is_active, last_login_at, created_at, updated_at
		FROM users WHERE api_key = ?`,
		apiKey,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.APIKey, &user.IsActive,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Update updates a user's information.
func (r *UserRepository) Update(user *User) error {
	_, err := r.db.Exec(
		`UPDATE users SET email = ?, password = ?, role = ?, api_key = ?, is_active = ?,
		last_login_at = ?, updated_at = ? WHERE id = ?`,
		user.Email,
		user.Password,
		user.Role,
		user.APIKey,
		user.IsActive,
		user.LastLoginAt,
		time.Now(),
		user.ID,
	)
	return err
}

// UpdateLastLogin updates the user's last login timestamp.
func (r *UserRepository) UpdateLastLogin(id int64) error {
	now := time.Now()
	_, err := r.db.Exec(
		`UPDATE users SET last_login_at = ?, updated_at = ? WHERE id = ?`,
		now,
		now,
		id,
	)
	return err
}

// Delete removes a user from the database.
func (r *UserRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// List retrieves all users with pagination.
func (r *UserRepository) List(limit, offset int) ([]*User, error) {
	rows, err := r.db.Query(
		`SELECT id, email, password, role, api_key, is_active, last_login_at, created_at, updated_at
		FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.APIKey,
			&user.IsActive, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

// Count returns the total number of users.
func (r *UserRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

// UserMeta represents metadata associated with a user.
type UserMeta struct {
	ID        int64
	Key       string
	Value     string
	UserID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserMetaRepository handles database operations for user metadata.
type UserMetaRepository struct {
	db *sql.DB
}

// NewUserMetaRepository creates a new user meta repository.
func NewUserMetaRepository(db *sql.DB) *UserMetaRepository {
	return &UserMetaRepository{db: db}
}

// Create inserts new metadata for a user.
func (r *UserMetaRepository) Create(userID int64, key, value string) error {
	_, err := r.db.Exec(
		"INSERT INTO users_meta (user_id, key, value) VALUES (?, ?, ?)",
		userID,
		key,
		value,
	)
	return err
}

// Get retrieves metadata for a user by key.
func (r *UserMetaRepository) Get(userID int64, key string) (*UserMeta, error) {
	meta := &UserMeta{}
	err := r.db.QueryRow(
		"SELECT id, key, value, user_id, created_at, updated_at FROM users_meta WHERE user_id = ? AND key = ?",
		userID,
		key,
	).Scan(&meta.ID, &meta.Key, &meta.Value, &meta.UserID, &meta.CreatedAt, &meta.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return meta, nil
}

// Update updates metadata for a user.
func (r *UserMetaRepository) Update(userID int64, key, value string) error {
	_, err := r.db.Exec(
		"UPDATE users_meta SET value = ?, updated_at = ? WHERE user_id = ? AND key = ?",
		value,
		time.Now(),
		userID,
		key,
	)
	return err
}

// Delete removes metadata for a user.
func (r *UserMetaRepository) Delete(userID int64, key string) error {
	_, err := r.db.Exec(
		"DELETE FROM users_meta WHERE user_id = ? AND key = ?",
		userID,
		key,
	)
	return err
}

// ListByUser retrieves all metadata for a user.
func (r *UserMetaRepository) ListByUser(userID int64) ([]*UserMeta, error) {
	rows, err := r.db.Query(
		"SELECT id, key, value, user_id, created_at, updated_at FROM users_meta WHERE user_id = ? ORDER BY key",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metadata []*UserMeta
	for rows.Next() {
		meta := &UserMeta{}
		if err := rows.Scan(&meta.ID, &meta.Key, &meta.Value, &meta.UserID, &meta.CreatedAt, &meta.UpdatedAt); err != nil {
			return nil, err
		}
		metadata = append(metadata, meta)
	}

	return metadata, rows.Err()
}

// Upsert inserts or updates metadata for a user.
func (r *UserMetaRepository) Upsert(userID int64, key, value string) error {
	existing, err := r.Get(userID, key)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(userID, key, value)
	}

	return r.Update(userID, key, value)
}
