package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:1@localhost:5432/project_management?sslmode=disable"
	}

	var err error
	poolConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return fmt.Errorf("unable to parse database config: %w", err)
	}

	// Optional SQL logging for debugging (very verbose).
	// Enable by setting DB_LOG_SQL=true
	if strings.EqualFold(os.Getenv("DB_LOG_SQL"), "true") {
		poolConfig.ConnConfig.Tracer = &SQLTracer{}
		log.Println("SQL logging enabled (DB_LOG_SQL=true)")
	}

	DB, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := DB.Ping(context.Background()); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Database connection established")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
