package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	databaseURL := "postgres://postgres:1@localhost:5432/project_management?sslmode=disable"

	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	migrationSQL, err := os.ReadFile("migration/001_enhance_entities.sql")
	if err != nil {
		log.Fatalf("Unable to read migration file: %v\n", err)
	}

	_, err = conn.Exec(context.Background(), string(migrationSQL))
	if err != nil {
		log.Fatalf("Migration failed: %v\n", err)
	}

	fmt.Println("Migration completed successfully!")
}
