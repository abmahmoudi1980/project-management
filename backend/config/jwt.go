package config

import (
	"os"
	"time"
)

// JWT Configuration
var (
	JWTSecret        = getEnv("JWT_SECRET", "")
	JWTAccessExpiry  = parseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	JWTRefreshExpiry = parseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))
)

// SMTP Configuration
var (
	SMTPHost     = getEnv("SMTP_HOST", "smtp.gmail.com")
	SMTPPort     = getEnv("SMTP_PORT", "587")
	SMTPUser     = getEnv("SMTP_USER", "")
	SMTPPassword = getEnv("SMTP_PASSWORD", "")
)

// Application URLs
var (
	AppURL = getEnv("APP_URL", "http://localhost:5173")
	APIURL = getEnv("API_URL", "http://localhost:3000")
)

// Helper function to get environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to parse duration
func parseDuration(duration string) time.Duration {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return 15 * time.Minute
	}
	return d
}
