package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	ctx := context.Background()

	databaseURL := "postgres://postgres:1@localhost:5432/project_management?sslmode=disable"
	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	email := "admin@example.com"
	newPassword := "Admin123!"

	userID, passwordHash, err := getUserID(ctx, db, email)
	if err != nil {
		log.Fatalf("Failed to get user: %v", err)
	}

	if userID == nil {
		log.Fatalf("User %s not found", email)
	}

	log.Printf("Found user: %s (ID: %s)\n", email, *userID)

	err = updatePassword(ctx, db, *userID, passwordHash, newPassword)
	if err != nil {
		log.Fatalf("Failed to update password: %v", err)
	}

	log.Printf("Password successfully reset for %s\n", email)
	log.Printf("New password: %s\n", newPassword)
	log.Println("Please change this password after first login!")
}

func getUserID(ctx context.Context, db *pgxpool.Pool, email string) (*uuid.UUID, string, error) {
	var userID uuid.UUID
	var currentHash string

	query := `SELECT id, password_hash FROM users WHERE email = $1`
	err := db.QueryRow(ctx, query, email).Scan(&userID, &currentHash)
	if err != nil {
		return nil, "", err
	}

	return &userID, currentHash, nil
}

func updatePassword(ctx context.Context, db *pgxpool.Pool, userID uuid.UUID, currentHash, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `UPDATE users SET password_hash = $1, failed_login_attempts = 0, locked_until = NULL, updated_at = $2 WHERE id = $3`
	_, err = db.Exec(ctx, query, string(hashedPassword), time.Now(), userID)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}
