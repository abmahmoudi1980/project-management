package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:1@localhost:5432/project_management?sslmode=disable"
	}

	var err error
	DB, err = pgxpool.New(context.Background(), databaseURL)
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
