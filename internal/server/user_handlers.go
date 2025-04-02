package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/theshovonaha/mini-pam/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// UserRequest represents the request body for creating or updating a user
type UserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password,omitempty"` // Only used for creation, not returned
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Active    bool   `json:"active"`
}

// handleListUsers returns a handler for listing users
func (s *Server) handleListUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters for pagination
		page := 1
		pageSize := 20

		if pageStr := r.URL.Query().Get("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				page = p
			}
		}

		if pageSizeStr := r.URL.Query().Get("page_size"); pageSizeStr != "" {
			if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
				pageSize = ps
			}
		}

		// Get users from the database
		users, err := s.models.Users.List(page, pageSize)
		if err != nil {
			s.logger.Printf("Error listing users: %v", err)
			s.respondError(w, http.StatusInternalServerError, "Failed to list users")
			return
		}

		// Return the users
		s.respondJSON(w, http.StatusOK, users)
	}
}

// handleCreateUser returns a handler for creating a new user
func (s *Server) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var req UserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Validate required fields
		if req.Username == "" || req.Email == "" || req.Password == "" {
			s.respondError(w, http.StatusBadRequest, "Username, email, and password are required")
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
		if err != nil {
			s.logger.Printf("Error hashing password: %v", err)
			s.respondError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		// Create the user model
		user := &models.User{
			Username:       req.Username,
			Email:          req.Email,
			HashedPassword: string(hashedPassword),
			FirstName:      req.FirstName,
			LastName:       req.LastName,
			Active:         req.Active,
		}

		// Save the user to the database
		err = s.models.Users.Create(user)
		if err != nil {
			s.logger.Printf("Error creating user: %v", err)
			s.respondError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		// Create an audit log entry
		auditLog := &models.AuditLog{
			UserID:     user.ID,
			Action:     "create",
			Resource:   "user",
			ResourceID: user.ID,
			IPAddress:  r.RemoteAddr,
			UserAgent:  r.UserAgent(),
			Details:    "User created",
		}
		if err := s.models.AuditLogs.Create(auditLog); err != nil {
			s.logger.Printf("Error creating audit log: %v", err)
		}

		// Return the created user
		s.respondJSON(w, http.StatusCreated, user)
	}
}

// handleGetUser returns a handler for getting a user by ID
func (s *Server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the user ID from the URL
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		// Get the user from the database
		user, err := s.models.Users.GetByID(id)
		if err != nil {
			if errors.Is(err, models.ErrRecordNotFound) {
				s.respondError(w, http.StatusNotFound, "User not found")
			} else {
				s.logger.Printf("Error getting user: %v", err)
				s.respondError(w, http.StatusInternalServerError, "Failed to get user")
			}
			return
		}

		// Return the user
		s.respondJSON(w, http.StatusOK, user)
	}
}

// handleUpdateUser returns a handler for updating a user
func (s *Server) handleUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the user ID from the URL
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		// Get the user from the database
		user, err := s.models.Users.GetByID(id)
		if err != nil {
			if errors.Is(err, models.ErrRecordNotFound) {
				s.respondError(w, http.StatusNotFound, "User not found")
			} else {
				s.logger.Printf("Error getting user: %v", err)
				s.respondError(w, http.StatusInternalServerError, "Failed to update user")
			}
			return
		}

		// Parse the request body
		var req UserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.respondError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Update the user fields
		user.Username = req.Username
		user.Email = req.Email
		user.FirstName = req.FirstName
		user.LastName = req.LastName
		user.Active = req.Active

		// Update the user in the database
		err = s.models.Users.Update(user)
		if err != nil {
			s.logger.Printf("Error updating user: %v", err)
			s.respondError(w, http.StatusInternalServerError, "Failed to update user")
			return
		}

		// If a password was provided, update it
		if req.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
			if err != nil {
				s.logger.Printf("Error hashing password: %v", err)
				s.respondError(w, http.StatusInternalServerError, "Failed to update password")
				return
			}

			err = s.models.Users.UpdatePassword(user.ID, string(hashedPassword))
			if err != nil {
				s.logger.Printf("Error updating password: %v", err)
				s.respondError(w, http.StatusInternalServerError, "Failed to update password")
				return
			}
		}

		// Create an audit log entry
		auditLog := &models.AuditLog{
			UserID:     user.ID,
			Action:     "update",
			Resource:   "user",
			ResourceID: user.ID,
			IPAddress:  r.RemoteAddr,
			UserAgent:  r.UserAgent(),
			Details:    "User updated",
		}
		if err := s.models.AuditLogs.Create(auditLog); err != nil {
			s.logger.Printf("Error creating audit log: %v", err)
		}

		// Return the updated user
		s.respondJSON(w, http.StatusOK, user)
	}
}

// handleDeleteUser returns a handler for deleting a user
func (s *Server) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the user ID from the URL
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.respondError(w, http.StatusBadRequest, "Invalid user ID")
			return
		}

		// Delete the user from the database
		err = s.models.Users.Delete(id)
		if err != nil {
			if errors.Is(err, models.ErrRecordNotFound) {
				s.respondError(w, http.StatusNotFound, "User not found")
			} else {
				s.logger.Printf("Error deleting user: %v", err)
				s.respondError(w, http.StatusInternalServerError, "Failed to delete user")
			}
			return
		}

		// Create an audit log entry
		auditLog := &models.AuditLog{
			Action:     "delete",
			Resource:   "user",
			ResourceID: id,
			IPAddress:  r.RemoteAddr,
			UserAgent:  r.UserAgent(),
			Details:    "User deleted",
		}
		if err := s.models.AuditLogs.Create(auditLog); err != nil {
			s.logger.Printf("Error creating audit log: %v", err)
		}

		// Return a success message
		s.respondJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
	}
}
