package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"` // Never expose password hash
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Active         bool      `json:"active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// Role represents a role that can be assigned to users
type Role struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	UserID    int       `json:"user_id"`
	RoleID    int       `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// Credential represents a stored privileged credential
type Credential struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"` // e.g., "password", "ssh_key", "api_key"
	Username    string    `json:"username"`
	Secret      string    `json:"-"` // Encrypted secret, never exposed directly
	System      string    `json:"system"`
	ExpiresAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   int       `json:"created_by"` // User ID who created this credential
}

// CredentialAccess represents a record of credential access
type CredentialAccess struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	CredentialID int       `json:"credential_id"`
	AccessedAt   time.Time `json:"accessed_at"`
	IPAddress    string    `json:"ip_address"`
	UserAgent    string    `json:"user_agent"`
	Reason       string    `json:"reason"`
}

// AuditLog represents a system audit log entry
type AuditLog struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	ResourceID int       `json:"resource_id"`
	Timestamp  time.Time `json:"timestamp"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	Details    string    `json:"details"`
}
