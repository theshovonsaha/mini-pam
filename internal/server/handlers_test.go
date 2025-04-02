package server

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/theshovonsaha/miniPam/internal/database"
)

func TestHealthHandler(t *testing.T) {
	// Create a new server for testing
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db := &database.Connection{} // Add a nil database connection for testing
	srv := NewServer("test", logger, db)

	// Create a request to the health endpoint
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Use the router to handle the request
	srv.Routes().ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	// Check that the status field is correct
	if status, exists := response["status"]; !exists || status != "available" {
		t.Errorf("handler returned wrong status: got %v want %v", status, "available")
	}

	// Check that the environment field is correct
	if env, exists := response["environment"]; !exists || env != "test" {
		t.Errorf("handler returned wrong environment: got %v want %v", env, "test")
	}
}

func TestVersionHandler(t *testing.T) {
	// Create a new server for testing
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	db := &database.Connection{} // Add a nil or mock database connection
	srv := NewServer("test", logger, db)

	// Create a request to the version endpoint
	req, err := http.NewRequest("GET", "/api/v1/version", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a response recorder to record the response
	rr := httptest.NewRecorder()

	// Use the router to handle the request
	srv.Routes().ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	var response map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	// Check that the version field exists
	if _, exists := response["version"]; !exists {
		t.Errorf("handler returned no version field")
	}

	// Check that the buildTime field exists
	if _, exists := response["buildTime"]; !exists {
		t.Errorf("handler returned no buildTime field")
	}
}
