package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server - API server
type Server struct {
	environment string
	logger      *log.Logger
	router      *mux.Router
}

// NewServer - Create a new server instance
func NewServer(environment string, logger *log.Logger) *Server {
	return &Server{
		environment: environment,
		logger:      logger,
		router:      mux.NewRouter(),
	}
}

// Routes - Define the API routes
func (s *Server) Routes() http.Handler {
	// API version prefix
	v1 := s.router.PathPrefix("/api/v1").Subrouter()

	// Health check endpoint
	v1.HandleFunc("/health", s.handleHealth()).Methods("GET")

	// Version info endpoint
	v1.HandleFunc("/version", s.handleVersion()).Methods("GET")

	// Add middleware
	s.router.Use(s.loggingMiddleware)
	s.router.Use(s.recoverPanicMiddleware)

	return s.router
}
