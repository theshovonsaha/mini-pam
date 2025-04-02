package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/theshovonaha/mini-pam/internal/database"
)

// CredentialRepository handles database operations related to credentials
type CredentialRepository struct {
	DB *database.Connection
}

// NewCredentialRepository creates a new credential repository
func NewCredentialRepository(db *database.Connection) *CredentialRepository {
	return &CredentialRepository{
		DB: db,
	}
}

// Create inserts a new credential into the database
func (r *CredentialRepository) Create(credential *Credential) error {
	query := `
		INSERT INTO credentials (name, description, type, username, secret, system, expires_at, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(
		ctx,
		query,
		credential.Name,
		credential.Description,
		credential.Type,
		credential.Username,
		credential.Secret,
		credential.System,
		credential.ExpiresAt,
		credential.CreatedBy,
	).Scan(&credential.ID, &credential.CreatedAt, &credential.UpdatedAt)

	return err
}

// GetByID retrieves a credential by its ID
func (r *CredentialRepository) GetByID(id int) (*Credential, error) {
	query := `
		SELECT id, name, description, type, username, secret, system, expires_at, created_at, updated_at, created_by
		FROM credentials
		WHERE id = $1`

	var credential Credential

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(ctx, query, id).Scan(
		&credential.ID,
		&credential.Name,
		&credential.Description,
		&credential.Type,
		&credential.Username,
		&credential.Secret,
		&credential.System,
		&credential.ExpiresAt,
		&credential.CreatedAt,
		&credential.UpdatedAt,
		&credential.CreatedBy,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &credential, nil
}

// Update updates an existing credential
func (r *CredentialRepository) Update(credential *Credential) error {
	query := `
		UPDATE credentials
		SET name = $1, description = $2, type = $3, username = $4, secret = $5, 
		    system = $6, expires_at = $7, updated_at = NOW()
		WHERE id = $8
		RETURNING updated_at`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(
		ctx,
		query,
		credential.Name,
		credential.Description,
		credential.Type,
		credential.Username,
		credential.Secret,
		credential.System,
		credential.ExpiresAt,
		credential.ID,
	).Scan(&credential.UpdatedAt)

	return err
}

// Delete deletes a credential by its ID
func (r *CredentialRepository) Delete(id int) error {
	query := `
		DELETE FROM credentials
		WHERE id = $1`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	result, err := r.DB.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// List returns a filtered and paginated list of credentials
func (r *CredentialRepository) List(system string, page, pageSize int) ([]*Credential, error) {
	// Ensure valid pagination parameters
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	// Calculate the offset
	offset := (page - 1) * pageSize

	// Base query with optional system filter
	var query string
	var args []interface{}

	if system == "" {
		query = `
			SELECT id, name, description, type, username, secret, system, expires_at, created_at, updated_at, created_by
			FROM credentials
			ORDER BY name
			LIMIT $1 OFFSET $2`
		args = []interface{}{pageSize, offset}
	} else {
		query = `
			SELECT id, name, description, type, username, secret, system, expires_at, created_at, updated_at, created_by
			FROM credentials
			WHERE system = $1
			ORDER BY name
			LIMIT $2 OFFSET $3`
		args = []interface{}{system, pageSize, offset}
	}

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	rows, err := r.DB.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	credentials := []*Credential{}
	for rows.Next() {
		var credential Credential
		err := rows.Scan(
			&credential.ID,
			&credential.Name,
			&credential.Description,
			&credential.Type,
			&credential.Username,
			&credential.Secret,
			&credential.System,
			&credential.ExpiresAt,
			&credential.CreatedAt,
			&credential.UpdatedAt,
			&credential.CreatedBy,
		)
		if err != nil {
			return nil, err
		}
		credentials = append(credentials, &credential)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return credentials, nil
}

// LogAccess logs an access to a credential
func (r *CredentialRepository) LogAccess(access *CredentialAccess) error {
	query := `
		INSERT INTO credential_access (user_id, credential_id, ip_address, user_agent, reason)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, accessed_at`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(
		ctx,
		query,
		access.UserID,
		access.CredentialID,
		access.IPAddress,
		access.UserAgent,
		access.Reason,
	).Scan(&access.ID, &access.AccessedAt)

	return err
}

// GetAccessHistory retrieves the access history for a credential
func (r *CredentialRepository) GetAccessHistory(credentialID int, limit int) ([]*CredentialAccess, error) {
	if limit <= 0 {
		limit = 50 // Default limit
	}

	query := `
		SELECT id, user_id, credential_id, accessed_at, ip_address, user_agent, reason
		FROM credential_access
		WHERE credential_id = $1
		ORDER BY accessed_at DESC
		LIMIT $2`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	rows, err := r.DB.DB.QueryContext(ctx, query, credentialID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	accesses := []*CredentialAccess{}
	for rows.Next() {
		var access CredentialAccess
		err := rows.Scan(
			&access.ID,
			&access.UserID,
			&access.CredentialID,
			&access.AccessedAt,
			&access.IPAddress,
			&access.UserAgent,
			&access.Reason,
		)
		if err != nil {
			return nil, err
		}
		accesses = append(accesses, &access)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accesses, nil
}
