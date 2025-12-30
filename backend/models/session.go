package models

import (
	"time"

	"github.com/google/uuid"
)

// Session represents an authentication session with refresh token
type Session struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	RefreshTokenHash string    `json:"-"` // Never send token hash to client
	UserAgent        string    `json:"user_agent,omitempty"`
	IPAddress        string    `json:"ip_address,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	ExpiresAt        time.Time `json:"expires_at"`
	Revoked          bool      `json:"revoked"`
}
