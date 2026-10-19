package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	// For local testing, we can use environment variables or a default connection string
	// Format: postgres://username:password@localhost:5432/dbname?sslmode=disable
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Secure fallback with a placeholder instead of hardcoding sensitive credentials
		connStr = "postgres://postgres:your_password_here@localhost:5432/adaptive_budgeting?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database!")
	return db, nil
}
