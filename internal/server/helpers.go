package server

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// respondError is a helper for sending error responses
func (s *Server) respondError(w http.ResponseWriter, status int, message string) {
	s.respondJSON(w, status, ErrorResponse{Error: message})
}

// readJSON is a helper function for reading JSON from a request
func (s *Server) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	// Set the maximum size of the request body (adjust as needed)
	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB

	// Initialize the JSON decoder
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	// Decode the request body into the destination
	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	// Check for additional JSON values
	err = dec.Decode(&struct{}{})
	if err != nil && err.Error() != "EOF" {
		return err
	}

	return nil
}
