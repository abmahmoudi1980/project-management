package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID             uuid.UUID  `json:"id"`
	ProjectID      uuid.UUID  `json:"project_id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Priority       string     `json:"priority"`
	Completed      bool       `json:"completed"`
	AssigneeID     *uuid.UUID `json:"assignee_id,omitempty"`
	AuthorID       *uuid.UUID `json:"author_id,omitempty"`
	Category       *string    `json:"category,omitempty"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	EstimatedHours *float64   `json:"estimated_hours,omitempty"`
	DoneRatio      int        `json:"done_ratio"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type CreateTaskRequest struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Priority       string     `json:"priority"`
	AssigneeID     *uuid.UUID `json:"assignee_id,omitempty"`
	AuthorID       *uuid.UUID `json:"author_id,omitempty"`
	Category       *string    `json:"category,omitempty"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	EstimatedHours *float64   `json:"estimated_hours,omitempty"`
	DoneRatio      int        `json:"done_ratio"`
}

type UpdateTaskRequest struct {
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Priority       string     `json:"priority"`
	Completed      bool       `json:"completed"`
	AssigneeID     *uuid.UUID `json:"assignee_id,omitempty"`
	AuthorID       *uuid.UUID `json:"author_id,omitempty"`
	Category       *string    `json:"category,omitempty"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	EstimatedHours *float64   `json:"estimated_hours,omitempty"`
	DoneRatio      int        `json:"done_ratio"`
}

type TaskWithUsers struct {
	ID             uuid.UUID  `json:"id"`
	ProjectID      uuid.UUID  `json:"project_id"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	Priority       string     `json:"priority"`
	Completed      bool       `json:"completed"`
	AssigneeID     *uuid.UUID `json:"assignee_id,omitempty"`
	AssigneeName   *string    `json:"assignee_name,omitempty"`
	AuthorID       *uuid.UUID `json:"author_id,omitempty"`
	AuthorName     *string    `json:"author_name,omitempty"`
	Category       *string    `json:"category,omitempty"`
	StartDate      *time.Time `json:"start_date,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	EstimatedHours *float64   `json:"estimated_hours,omitempty"`
	DoneRatio      int        `json:"done_ratio"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
