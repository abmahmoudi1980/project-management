package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	Title     string    `json:"title"`
	Priority  string    `json:"priority"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title    string `json:"title"`
	Priority string `json:"priority"`
}

type UpdateTaskRequest struct {
	Title     string `json:"title"`
	Priority  string `json:"priority"`
	Completed bool   `json:"completed"`
}
