package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Config holds the database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Connection wraps a SQL database connection with additional functionality
type Connection struct {
	DB     *sql.DB
	Logger *log.Logger
}

// NewConnection creates a new database connection
func NewConnection(cfg Config, logger *log.Logger) (*Connection, error) {
	// Create the connection string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	// Open a connection to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Check the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Println("Database connection established")

	return &Connection{
		DB:     db,
		Logger: logger,
	}, nil
}

// Close closes the database connection
func (c *Connection) Close() error {
	c.Logger.Println("Closing database connection")
	return c.DB.Close()
}
