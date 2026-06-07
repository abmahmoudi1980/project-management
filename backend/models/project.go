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
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
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
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
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
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

var (
	ErrValidation = &Error{Message: "validation error", Code: 400}
	ErrNotFound   = &Error{Message: "resource not found", Code: 404}
)

// ProjectTreeNode is a project in the tree response. Children is left empty in
// the wire payload — the client folds the flat list into a tree by parent_id.
type ProjectTreeNode struct {
	Project
	Children []ProjectTreeNode `json:"children,omitempty"`
}

// ProjectTree is the response shape for GET /api/projects/tree.
type ProjectTree struct {
	Nodes []ProjectTreeNode `json:"nodes"`
	Total int               `json:"total"`
}

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (e *Error) Error() string {
	return e.Message
}
