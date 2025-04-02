package server

import (
    "encoding/json"
    "net/http"
    "time"
)

//version and build time
var (
	version = "1.0.0"
	buildTime = time.Now().Format(time.RFC3339)
)

func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := map[string]string{
			"status":      "available",
			"environment": s.environment,
			"timestamp":   time.Now().Format(time.RFC3339),
		}

		s.respondJSON(w, http.StatusOK, status)
	}
}

// handleVersion returns a handler for version info requests
func (s *Server) handleVersion() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := map[string]string{
			"version":   version,
			"buildTime": buildTime,
		}

		s.respondJSON(w, http.StatusOK, info)
	}
}

// respondJSON is a helper for sending JSON responses
func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	// Set Content-Type header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Encode the data to JSON and send it
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.logger.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}