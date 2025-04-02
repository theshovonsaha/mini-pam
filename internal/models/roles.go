package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/theshovonsaha/mini-pam/internal/database"
)

// RoleRepository handles database operations related to roles
type RoleRepository struct {
	DB *database.Connection
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *database.Connection) *RoleRepository {
	return &RoleRepository{
		DB: db,
	}
}

// Create inserts a new role into the database
func (r *RoleRepository) Create(role *Role) error {
	query := `
		INSERT INTO roles (name, description)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(ctx, query, role.Name, role.Description).Scan(
		&role.ID, &role.CreatedAt, &role.UpdatedAt,
	)

	return err
}

// GetByID retrieves a role by its ID
func (r *RoleRepository) GetByID(id int) (*Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = $1`

	var role Role

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &role, nil
}

// GetByName retrieves a role by its name
func (r *RoleRepository) GetByName(name string) (*Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE name = $1`

	var role Role

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &role, nil
}

// Update updates an existing role
func (r *RoleRepository) Update(role *Role) error {
	query := `
		UPDATE roles
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(ctx, query, role.Name, role.Description, role.ID).Scan(&role.UpdatedAt)

	return err
}

// Delete deletes a role by its ID
func (r *RoleRepository) Delete(id int) error {
	query := `
		DELETE FROM roles
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

// List returns a list of all roles
func (r *RoleRepository) List() ([]*Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		ORDER BY name`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	rows, err := r.DB.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	roles := []*Role{}
	for rows.Next() {
		var role Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

// AssignRoleToUser assigns a role to a user
func (r *RoleRepository) AssignRoleToUser(userID, roleID int) error {
	query := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, role_id) DO NOTHING`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	_, err := r.DB.DB.ExecContext(ctx, query, userID, roleID)
	return err
}

// RemoveRoleFromUser removes a role from a user
func (r *RoleRepository) RemoveRoleFromUser(userID, roleID int) error {
	query := `
		DELETE FROM user_roles
		WHERE user_id = $1 AND role_id = $2`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	result, err := r.DB.DB.ExecContext(ctx, query, userID, roleID)
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

// GetUserRoles returns all roles assigned to a user
func (r *RoleRepository) GetUserRoles(userID int) ([]*Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
		ORDER BY r.name`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	rows, err := r.DB.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows
	roles := []*Role{}
	for rows.Next() {
		var role Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}
