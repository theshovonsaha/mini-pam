package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/theshovonaha/mini-pam/internal/database"
	"github.com/theshovonaha/mini-pam/internal/models"
)

// Server is our API server
type Server struct {
	environment string
	logger      *log.Logger
	router      *mux.Router
	db          *database.Connection
	models      Models
}

// Models holds all the repository instances
type Models struct {
	Users       *models.UserRepository
	Roles       *models.RoleRepository
	Credentials *models.CredentialRepository
	AuditLogs   *models.AuditLogRepository
}

// NewServer creates a new server instance
func NewServer(environment string, logger *log.Logger, db *database.Connection) *Server {
	s := &Server{
		environment: environment,
		logger:      logger,
		router:      mux.NewRouter(),
		db:          db,
	}

	// Initialize repositories
	s.models = Models{
		Users:       models.NewUserRepository(db),
		Roles:       models.NewRoleRepository(db),
		Credentials: models.NewCredentialRepository(db),
		AuditLogs:   models.NewAuditLogRepository(db),
	}

	return s
}

// Routes sets up all the routes for our application
func (s *Server) Routes() http.Handler {
	// API version prefix
	v1 := s.router.PathPrefix("/api/v1").Subrouter()

	// Health check endpoint
	v1.HandleFunc("/health", s.handleHealth()).Methods("GET")

	// Version info endpoint
	v1.HandleFunc("/version", s.handleVersion()).Methods("GET")

	// User routes
	v1.HandleFunc("/users", s.handleListUsers()).Methods("GET")
	v1.HandleFunc("/users", s.handleCreateUser()).Methods("POST")
	v1.HandleFunc("/users/{id:[0-9]+}", s.handleGetUser()).Methods("GET")
	v1.HandleFunc("/users/{id:[0-9]+}", s.handleUpdateUser()).Methods("PUT")
	v1.HandleFunc("/users/{id:[0-9]+}", s.handleDeleteUser()).Methods("DELETE")

	// // Role routes
	// v1.HandleFunc("/roles", s.handleListRoles()).Methods("GET")
	// v1.HandleFunc("/roles", s.handleCreateRole()).Methods("POST")
	// v1.HandleFunc("/roles/{id:[0-9]+}", s.handleGetRole()).Methods("GET")
	// v1.HandleFunc("/roles/{id:[0-9]+}", s.handleUpdateRole()).Methods("PUT")
	// v1.HandleFunc("/roles/{id:[0-9]+}", s.handleDeleteRole()).Methods("DELETE")
	// v1.HandleFunc("/users/{id:[0-9]+}/roles", s.handleGetUserRoles()).Methods("GET")
	// v1.HandleFunc("/users/{id:[0-9]+}/roles/{roleId:[0-9]+}", s.handleAssignRole()).Methods("POST")
	// v1.HandleFunc("/users/{id:[0-9]+}/roles/{roleId:[0-9]+}", s.handleRemoveRole()).Methods("DELETE")

	// // Credential routes
	// v1.HandleFunc("/credentials", s.handleListCredentials()).Methods("GET")
	// v1.HandleFunc("/credentials", s.handleCreateCredential()).Methods("POST")
	// v1.HandleFunc("/credentials/{id:[0-9]+}", s.handleGetCredential()).Methods("GET")
	// v1.HandleFunc("/credentials/{id:[0-9]+}", s.handleUpdateCredential()).Methods("PUT")
	// v1.HandleFunc("/credentials/{id:[0-9]+}", s.handleDeleteCredential()).Methods("DELETE")
	// v1.HandleFunc("/credentials/{id:[0-9]+}/access", s.handleGetCredentialAccessHistory()).Methods("GET")
	// v1.HandleFunc("/credentials/{id:[0-9]+}/access", s.handleLogCredentialAccess()).Methods("POST")

	// // Audit log routes
	// v1.HandleFunc("/audit-logs", s.handleListAuditLogs()).Methods("GET")
	// v1.HandleFunc("/audit-logs/users/{id:[0-9]+}", s.handleGetUserAuditLogs()).Methods("GET")
	// v1.HandleFunc("/audit-logs/resources/{resource}/{id:[0-9]+}", s.handleGetResourceAuditLogs()).Methods("GET")

	// Add middleware (order matters)
	s.router.Use(s.corsMiddleware)
	s.router.Use(s.rateLimitMiddleware)
	s.router.Use(s.authMiddleware)
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.recoverPanicMiddleware)

	return s.router
}
