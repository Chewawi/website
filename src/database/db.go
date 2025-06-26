package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "blog")

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")

	// Create tables if they don't exist
	createTables()
}

// createTables creates the necessary tables if they don't exist
func createTables() {
	// Create posts table
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			slug VARCHAR(255) NOT NULL UNIQUE,
			created TIMESTAMP NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create posts table: %v", err)
	}

	log.Println("Tables created successfully")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
