package main

import (
	"context"
	"log"
	"os"

	"project-management/config"
	"project-management/repositories"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	if err := config.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer config.CloseDB()

	ctx := context.Background()

	// Initialize repositories
	sessionRepo := repositories.NewSessionRepository(config.DB)
	passwordResetRepo := repositories.NewPasswordResetRepository(config.DB)

	// Clean up expired sessions
	log.Println("Cleaning up expired sessions...")
	sessionCount, err := sessionRepo.DeleteExpired(ctx)
	if err != nil {
		log.Printf("Error cleaning up sessions: %v", err)
		os.Exit(1)
	}
	log.Printf("Deleted %d expired sessions", sessionCount)

	// Clean up expired password reset tokens
	log.Println("Cleaning up expired password reset tokens...")
	tokenCount, err := passwordResetRepo.DeleteExpired(ctx)
	if err != nil {
		log.Printf("Error cleaning up password reset tokens: %v", err)
		os.Exit(1)
	}
	log.Printf("Deleted %d expired password reset tokens", tokenCount)

	log.Println("Cleanup completed successfully")
}
