// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package db

import (
	"database/sql"
	"time"
)

// Audit represents an audit log entry in the database.
type Audit struct {
	ID           int64
	UserID       *int64
	UserEmail    *string
	Action       string
	EntityType   string
	EntityID     *int64
	EntityName   *string
	Details      *string
	Status       *string
	ErrorMessage *string
	IPAddress    *string
	UserAgent    *string
	RequestID    *string
	CreatedAt    time.Time
}

// AuditRepository handles database operations for audit logs.
type AuditRepository struct {
	db *sql.DB
}

// NewAuditRepository creates a new audit repository.
func NewAuditRepository(db *sql.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

// Create inserts a new audit log entry into the database.
//
// Example:
//
//	audit := &Audit{
//		UserID:     &userID,
//		UserEmail:  &email,
//		Action:     "user.login",
//		EntityType: "user",
//		EntityID:   &userID,
//		Status:     &"success",
//		IPAddress:  &ip,
//	}
//	err := repo.Create(audit)
func (r *AuditRepository) Create(audit *Audit) error {
	result, err := r.db.Exec(
		`INSERT INTO audits (user_id, user_email, action, entity_type, entity_id, entity_name,
		details, status, error_message, ip_address, user_agent, request_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		audit.UserID,
		audit.UserEmail,
		audit.Action,
		audit.EntityType,
		audit.EntityID,
		audit.EntityName,
		audit.Details,
		audit.Status,
		audit.ErrorMessage,
		audit.IPAddress,
		audit.UserAgent,
		audit.RequestID,
	)
	if err != nil {
		return err
	}

	audit.ID, err = result.LastInsertId()
	return err
}

// GetByID retrieves an audit log entry by ID.
func (r *AuditRepository) GetByID(id int64) (*Audit, error) {
	audit := &Audit{}
	err := r.db.QueryRow(
		`SELECT id, user_id, user_email, action, entity_type, entity_id, entity_name,
		details, status, error_message, ip_address, user_agent, request_id, created_at
		FROM audits WHERE id = ?`,
		id,
	).Scan(&audit.ID, &audit.UserID, &audit.UserEmail, &audit.Action, &audit.EntityType,
		&audit.EntityID, &audit.EntityName, &audit.Details, &audit.Status, &audit.ErrorMessage,
		&audit.IPAddress, &audit.UserAgent, &audit.RequestID, &audit.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return audit, nil
}

// List retrieves all audit logs with pagination.
func (r *AuditRepository) List(limit, offset int) ([]*Audit, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, user_email, action, entity_type, entity_id, entity_name,
		details, status, error_message, ip_address, user_agent, request_id, created_at
		FROM audits ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAudits(rows)
}

// ListByUser retrieves all audit logs for a specific user.
func (r *AuditRepository) ListByUser(userID int64, limit, offset int) ([]*Audit, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, user_email, action, entity_type, entity_id, entity_name,
		details, status, error_message, ip_address, user_agent, request_id, created_at
		FROM audits WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAudits(rows)
}

// ListByAction retrieves all audit logs for a specific action.
func (r *AuditRepository) ListByAction(action string, limit, offset int) ([]*Audit, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, user_email, action, entity_type, entity_id, entity_name,
		details, status, error_message, ip_address, user_agent, request_id, created_at
		FROM audits WHERE action = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		action,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAudits(rows)
}

// ListByEntity retrieves all audit logs for a specific entity.
func (r *AuditRepository) ListByEntity(entityType string, entityID int64, limit, offset int) ([]*Audit, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, user_email, action, entity_type, entity_id, entity_name,
		details, status, error_message, ip_address, user_agent, request_id, created_at
		FROM audits WHERE entity_type = ? AND entity_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		entityType,
		entityID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAudits(rows)
}

// ListByDateRange retrieves audit logs within a date range.
func (r *AuditRepository) ListByDateRange(startDate, endDate time.Time, limit, offset int) ([]*Audit, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, user_email, action, entity_type, entity_id, entity_name,
		details, status, error_message, ip_address, user_agent, request_id, created_at
		FROM audits WHERE created_at >= ? AND created_at <= ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		startDate,
		endDate,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanAudits(rows)
}

// Count returns the total number of audit logs.
func (r *AuditRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM audits").Scan(&count)
	return count, err
}

// CountByUser returns the total number of audit logs for a specific user.
func (r *AuditRepository) CountByUser(userID int64) (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM audits WHERE user_id = ?", userID).Scan(&count)
	return count, err
}

// Delete removes audit logs older than a specific date (for cleanup).
func (r *AuditRepository) DeleteOlderThan(date time.Time) (int64, error) {
	result, err := r.db.Exec("DELETE FROM audits WHERE created_at < ?", date)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *AuditRepository) scanAudits(rows *sql.Rows) ([]*Audit, error) {
	var audits []*Audit
	for rows.Next() {
		audit := &Audit{}
		if err := rows.Scan(&audit.ID, &audit.UserID, &audit.UserEmail, &audit.Action,
			&audit.EntityType, &audit.EntityID, &audit.EntityName, &audit.Details,
			&audit.Status, &audit.ErrorMessage, &audit.IPAddress, &audit.UserAgent,
			&audit.RequestID, &audit.CreatedAt); err != nil {
			return nil, err
		}
		audits = append(audits, audit)
	}
	return audits, rows.Err()
}

// ToolMetric represents a tool execution metric in the database.
type ToolMetric struct {
	ID             int64
	ToolID         int64
	UserID         *int64
	RequestID      *string
	Arguments      *string
	Success        *bool
	ResponseTimeMs *int
	ErrorMessage   *string
	ServerID       *int64
	ClientIP       *string
	UserAgent      *string
	CreatedAt      time.Time
}

// ToolMetricRepository handles database operations for tool metrics.
type ToolMetricRepository struct {
	db *sql.DB
}

// NewToolMetricRepository creates a new tool metric repository.
func NewToolMetricRepository(db *sql.DB) *ToolMetricRepository {
	return &ToolMetricRepository{db: db}
}

// Create inserts a new tool metric into the database.
//
// Example:
//
//	metric := &ToolMetric{
//		ToolID:         toolID,
//		UserID:         &userID,
//		RequestID:      &reqID,
//		Success:        &true,
//		ResponseTimeMs: &150,
//		ServerID:       &serverID,
//	}
//	err := repo.Create(metric)
func (r *ToolMetricRepository) Create(metric *ToolMetric) error {
	result, err := r.db.Exec(
		`INSERT INTO tool_metrics (tool_id, user_id, request_id, arguments, success,
		response_time_ms, error_message, server_id, client_ip, user_agent)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		metric.ToolID,
		metric.UserID,
		metric.RequestID,
		metric.Arguments,
		metric.Success,
		metric.ResponseTimeMs,
		metric.ErrorMessage,
		metric.ServerID,
		metric.ClientIP,
		metric.UserAgent,
	)
	if err != nil {
		return err
	}

	metric.ID, err = result.LastInsertId()
	return err
}

// GetByID retrieves a tool metric by ID.
func (r *ToolMetricRepository) GetByID(id int64) (*ToolMetric, error) {
	metric := &ToolMetric{}
	err := r.db.QueryRow(
		`SELECT id, tool_id, user_id, request_id, arguments, success, response_time_ms,
		error_message, server_id, client_ip, user_agent, created_at
		FROM tool_metrics WHERE id = ?`,
		id,
	).Scan(&metric.ID, &metric.ToolID, &metric.UserID, &metric.RequestID, &metric.Arguments,
		&metric.Success, &metric.ResponseTimeMs, &metric.ErrorMessage, &metric.ServerID,
		&metric.ClientIP, &metric.UserAgent, &metric.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return metric, nil
}

// List retrieves all tool metrics with pagination.
func (r *ToolMetricRepository) List(limit, offset int) ([]*ToolMetric, error) {
	rows, err := r.db.Query(
		`SELECT id, tool_id, user_id, request_id, arguments, success, response_time_ms,
		error_message, server_id, client_ip, user_agent, created_at
		FROM tool_metrics ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMetrics(rows)
}

// ListByTool retrieves all metrics for a specific tool.
func (r *ToolMetricRepository) ListByTool(toolID int64, limit, offset int) ([]*ToolMetric, error) {
	rows, err := r.db.Query(
		`SELECT id, tool_id, user_id, request_id, arguments, success, response_time_ms,
		error_message, server_id, client_ip, user_agent, created_at
		FROM tool_metrics WHERE tool_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		toolID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMetrics(rows)
}

// ListByUser retrieves all metrics for a specific user.
func (r *ToolMetricRepository) ListByUser(userID int64, limit, offset int) ([]*ToolMetric, error) {
	rows, err := r.db.Query(
		`SELECT id, tool_id, user_id, request_id, arguments, success, response_time_ms,
		error_message, server_id, client_ip, user_agent, created_at
		FROM tool_metrics WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMetrics(rows)
}

// ListByDateRange retrieves metrics within a date range.
func (r *ToolMetricRepository) ListByDateRange(startDate, endDate time.Time, limit, offset int) ([]*ToolMetric, error) {
	rows, err := r.db.Query(
		`SELECT id, tool_id, user_id, request_id, arguments, success, response_time_ms,
		error_message, server_id, client_ip, user_agent, created_at
		FROM tool_metrics WHERE created_at >= ? AND created_at <= ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		startDate,
		endDate,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanMetrics(rows)
}

// GetAverageResponseTime calculates the average response time for a tool.
func (r *ToolMetricRepository) GetAverageResponseTime(toolID int64) (float64, error) {
	var avg sql.NullFloat64
	err := r.db.QueryRow(
		"SELECT AVG(response_time_ms) FROM tool_metrics WHERE tool_id = ? AND response_time_ms IS NOT NULL",
		toolID,
	).Scan(&avg)

	if err != nil {
		return 0, err
	}

	if !avg.Valid {
		return 0, nil
	}

	return avg.Float64, nil
}

// GetSuccessRate calculates the success rate for a tool.
func (r *ToolMetricRepository) GetSuccessRate(toolID int64) (float64, error) {
	var total, successful int64

	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM tool_metrics WHERE tool_id = ?",
		toolID,
	).Scan(&total)
	if err != nil {
		return 0, err
	}

	if total == 0 {
		return 0, nil
	}

	err = r.db.QueryRow(
		"SELECT COUNT(*) FROM tool_metrics WHERE tool_id = ? AND success = 1",
		toolID,
	).Scan(&successful)
	if err != nil {
		return 0, err
	}

	return float64(successful) / float64(total) * 100, nil
}

// Count returns the total number of tool metrics.
func (r *ToolMetricRepository) Count() (int64, error) {
	var count int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM tool_metrics").Scan(&count)
	return count, err
}

// Delete removes metrics older than a specific date (for cleanup).
func (r *ToolMetricRepository) DeleteOlderThan(date time.Time) (int64, error) {
	result, err := r.db.Exec("DELETE FROM tool_metrics WHERE created_at < ?", date)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (r *ToolMetricRepository) scanMetrics(rows *sql.Rows) ([]*ToolMetric, error) {
	var metrics []*ToolMetric
	for rows.Next() {
		metric := &ToolMetric{}
		if err := rows.Scan(&metric.ID, &metric.ToolID, &metric.UserID, &metric.RequestID,
			&metric.Arguments, &metric.Success, &metric.ResponseTimeMs, &metric.ErrorMessage,
			&metric.ServerID, &metric.ClientIP, &metric.UserAgent, &metric.CreatedAt); err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}
	return metrics, rows.Err()
}
