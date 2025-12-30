package models

import (
	"time"

	"github.com/google/uuid"
)

type TimeLog struct {
	ID              uuid.UUID `json:"id"`
	TaskID          uuid.UUID `json:"task_id"`
	Date            time.Time `json:"date"`
	DurationMinutes int       `json:"duration_minutes"`
	Note            *string   `json:"note,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type CreateTimeLogRequest struct {
	Date            time.Time `json:"date"`
	DurationMinutes int       `json:"duration_minutes"`
	Note            *string   `json:"note,omitempty"`
}
