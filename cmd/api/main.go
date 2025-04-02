package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/theshovonaha/mini-pam/internal/database"
	"github.com/theshovonaha/mini-pam/internal/server"
)

func main() {
	// Define command-line flags for application configuration
	var (
		port        = flag.Int("port", 8080, "API server port")
		environment = flag.String("env", "development", "Environment (development|staging|production)")
		dbHost      = flag.String("db-host", "localhost", "Database host")
		dbPort      = flag.Int("db-port", 5432, "Database port")
		dbUser      = flag.String("db-user", "postgres", "Database user")
		dbPassword  = flag.String("db-password", "postgres", "Database password")
		dbName      = flag.String("db-name", "securevault", "Database name")
		dbSSLMode   = flag.String("db-sslmode", "disable", "Database SSL mode")
	)
	flag.Parse()

	// Initialize the logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Initialize the database
	db, err := database.NewConnection(database.Config{
		Host:     *dbHost,
		Port:     *dbPort,
		User:     *dbUser,
		Password: *dbPassword,
		DBName:   *dbName,
		SSLMode:  *dbSSLMode,
	}, logger)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a new server instance
	srv := server.NewServer(*environment, logger, db)

	// Start the HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      srv.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the server in a goroutine so that it doesn't block
	serverErrors := make(chan error, 1)
	go func() {
		logger.Printf("Starting server on port %d in %s mode", *port, *environment)
		serverErrors <- httpServer.ListenAndServe()
	}()

	// Create a channel to listen for OS signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive a signal or an error from the server
	select {
	case err := <-serverErrors:
		logger.Fatalf("Error starting server: %v", err)
	case <-shutdown:
		logger.Println("Shutting down server...")

		// Create a context with timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown the server
		if err := httpServer.Shutdown(ctx); err != nil {
			logger.Printf("Error during server shutdown: %v", err)
			httpServer.Close()
		}

		logger.Println("Server stopped")
	}
}
