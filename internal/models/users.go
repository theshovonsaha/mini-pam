package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/theshovonsaha/mini-pam/internal/database"
)

// Common errors
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrDuplicateKey   = errors.New("duplicate key value violates unique constraint")
)

// UserRepository handles database operations related to users
type UserRepository struct {
	DB *database.Connection
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.Connection) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(user *User) error {
	query := `
		INSERT INTO users (username, email, hashed_password, first_name, last_name, active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.HashedPassword,
		user.FirstName,
		user.LastName,
		user.Active,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(id int) (*User, error) {
	query := `
		SELECT id, username, email, hashed_password, first_name, last_name, active, created_at, updated_at
		FROM users
		WHERE id = $1`

	var user User

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.HashedPassword,
		&user.FirstName,
		&user.LastName,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by their email address
func (r *UserRepository) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, username, email, hashed_password, first_name, last_name, active, created_at, updated_at
		FROM users
		WHERE email = $1`

	var user User

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.HashedPassword,
		&user.FirstName,
		&user.LastName,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetByUsername retrieves a user by their username
func (r *UserRepository) GetByUsername(username string) (*User, error) {
	query := `
		SELECT id, username, email, hashed_password, first_name, last_name, active, created_at, updated_at
		FROM users
		WHERE username = $1`

	var user User

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.HashedPassword,
		&user.FirstName,
		&user.LastName,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *User) error {
	query := `
		UPDATE users
		SET username = $1, email = $2, first_name = $3, last_name = $4, active = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	err := r.DB.DB.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Active,
		user.ID,
	).Scan(&user.UpdatedAt)

	return err
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(id int, hashedPassword string) error {
	query := `
		UPDATE users
		SET hashed_password = $1, updated_at = NOW()
		WHERE id = $2`

	// Set a timeout for the query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Execute the query
	_, err := r.DB.DB.ExecContext(ctx, query, hashedPassword, id)
	return err
}

// Delete deletes a user by their ID
func (r *UserRepository) Delete(id int) error {
	query := `
		DELETE FROM users
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

// List returns a paginated list of users
func (r *UserRepository) List(page, pageSize int) ([]*User, error) {
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
		SELECT id, username, email, hashed_password, first_name, last_name, active, created_at, updated_at
		FROM users
		ORDER BY username
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
	users := []*User{}
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.HashedPassword,
			&user.FirstName,
			&user.LastName,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	// Check for errors after iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
