package server

import (
	"net/http"
	"runtime/debug"
	"time"
)

// loggingMiddleware logs all requests with their path, method, and duration
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Record the request details
		s.logger.Printf("Request: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		
		// Call the next handler
		next.ServeHTTP(w, r)
		
		// Log the request duration
		s.logger.Printf("Request completed: %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// recoverPanicMiddleware recovers from panics and logs the error
func (s *Server) recoverPanicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Recover from panic and log the error
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace
				s.logger.Printf("PANIC: %v\n%s", err, debug.Stack())
				
				// Return a 500 Internal Server Error response
				w.Header().Set("Connection", "close")
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		
		next.ServeHTTP(w, r)
	})
}