package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Identifier  string     `json:"identifier"`
	Homepage    *string    `json:"homepage,omitempty"`
	IsPublic    bool       `json:"is_public"`
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateProjectRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Identifier  string     `json:"identifier"`
	Homepage    *string    `json:"homepage,omitempty"`
	IsPublic    bool       `json:"is_public"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type UpdateProjectRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Identifier  string     `json:"identifier"`
	Homepage    *string    `json:"homepage,omitempty"`
	IsPublic    bool       `json:"is_public"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

var (
	ErrValidation = &Error{Message: "validation error", Code: 400}
	ErrNotFound   = &Error{Message: "resource not found", Code: 404}
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *Error) Error() string {
	return e.Message
}
