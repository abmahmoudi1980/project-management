package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Printf("Warning: .env file not found, using defaults: %v\n", err)
	}

	// Build database URL from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "1")
	dbName := getEnv("DB_NAME", "project_management")

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	// Get all migration files
	migrationFiles, err := filepath.Glob("migration/*.sql")
	if err != nil {
		log.Fatalf("Unable to read migration files: %v\n", err)
	}

	// Sort migrations to ensure they run in order
	sort.Strings(migrationFiles)

	// Run each migration
	for _, file := range migrationFiles {
		fmt.Printf("Running migration: %s\n", file)

		migrationSQL, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("Unable to read migration file %s: %v\n", file, err)
		}

		_, err = conn.Exec(context.Background(), string(migrationSQL))
		if err != nil {
			log.Fatalf("Migration %s failed: %v\n", file, err)
		}

		fmt.Printf("✓ Migration %s completed successfully!\n", file)
	}

	fmt.Println("\n✓ All migrations completed successfully!")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
