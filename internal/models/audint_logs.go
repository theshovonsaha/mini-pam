package models

import (
	"context"
	"time"

	"github.com/theshovonaha/mini-pam/internal/database"
)

// AuditLogRepository handles database operations related to audit logs
type AuditLogRepository struct {
	DB *database.Connection
}

// NewAuditLogRepository creates a new audit log repository
func NewAuditLogRepository(db *database.Connection) *AuditLogRepository {
	return &AuditLogRepository{
		DB: db,
	}
}

// Create inserts a new audit log entry into the database
func (r *AuditLogRepository) Create(log *AuditLog) error {
	query := `
		INSERT INTO audit_logs (user_id, action, resource, resource_id, ip_address, user_agent, details)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, timestamp`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(
		ctx,
		query,
		log.UserID,
		log.Action,
		log.Resource,
		log.ResourceID,
		log.IPAddress,
		log.UserAgent,
		log.Details,
	).Scan(&log.ID, &log.Timestamp)

	return err
}

// GetByResourceID retrieves audit logs for a specific resource
func (r *AuditLogRepository) GetByResourceID(resource string, resourceID int, limit int) ([]*AuditLog, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}

	query := `
		SELECT id, user_id, action, resource, resource_id, timestamp, ip_address, user_agent, details
		FROM audit_logs
		WHERE resource = $1 AND resource_id = $2
		ORDER BY timestamp DESC
		LIMIT $3`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	rows, err := r.DB.DB.QueryContext(ctx, query, resource, resourceID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	logs := []*AuditLog{}
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.Resource,
			&log.ResourceID,
			&log.Timestamp,
			&log.IPAddress,
			&log.UserAgent,
			&log.Details,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

// GetByUserID retrieves audit logs for a specific user
func (r *AuditLogRepository) GetByUserID(userID int, limit int) ([]*AuditLog, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}

	query := `
		SELECT id, user_id, action, resource, resource_id, timestamp, ip_address, user_agent, details
		FROM audit_logs
		WHERE user_id = $1
		ORDER BY timestamp DESC
		LIMIT $2`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	rows, err := r.DB.DB.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	logs := []*AuditLog{}
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.Resource,
			&log.ResourceID,
			&log.Timestamp,
			&log.IPAddress,
			&log.UserAgent,
			&log.Details,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

// List retrieves recent audit logs with pagination
func (r *AuditLogRepository) List(page, pageSize int) ([]*AuditLog, error) {
	// Ensure valid pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	// Calculate the offset
	offset := (page - 1) * pageSize

	query := `
		SELECT id, user_id, action, resource, resource_id, timestamp, ip_address, user_agent, details
		FROM audit_logs
		ORDER BY timestamp DESC
		LIMIT $1 OFFSET $2`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	rows, err := r.DB.DB.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	logs := []*AuditLog{}
	for rows.Next() {
		var log AuditLog
		err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.Resource,
			&log.ResourceID,
			&log.Timestamp,
			&log.IPAddress,
			&log.UserAgent,
			&log.Details,
		)
		if err != nil {
			return nil, err
		}
		logs = append(logs, &log)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return logs, nil
}
