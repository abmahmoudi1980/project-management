package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user account in the system
type User struct {
	ID                  uuid.UUID  `json:"id"`
	Username            string     `json:"username"`
	Email               string     `json:"email"`
	PasswordHash        string     `json:"-"` // Never send password hash to client
	Role                string     `json:"role"`
	IsActive            bool       `json:"is_active"`
	FailedLoginAttempts int        `json:"-"` // Internal use only
	LockedUntil         *time.Time `json:"-"` // Internal use only
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	LastLoginAt         *time.Time `json:"last_login_at,omitempty"`
}

// CreateUserRequest represents the registration request
type CreateUserRequest struct {
	Username             string `json:"username"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the login response with tokens
type LoginResponse struct {
	User         User   `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// UpdateUserRequest represents user profile update request
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// UpdateUserRoleRequest represents role change request (admin only)
type UpdateUserRoleRequest struct {
	Role string `json:"role"`
}

// UpdateUserActivationRequest represents activation status change (admin only)
type UpdateUserActivationRequest struct {
	IsActive bool `json:"is_active"`
}
