package models

import (
	"time"

	"github.com/google/uuid"
)

// ProjectRole represents a predefined role that can be assigned to project members
type ProjectRole struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProjectMember represents a user's membership in a project with an assigned role
type ProjectMember struct {
	ProjectID     uuid.UUID  `json:"project_id"`
	UserID        uuid.UUID  `json:"user_id"`
	ProjectRoleID uuid.UUID  `json:"project_role_id"`
	JoinedAt      time.Time  `json:"joined_at"`
	AddedBy       *uuid.UUID `json:"added_by,omitempty"`

	// Joined fields (populated when listing members)
	User        *User        `json:"user,omitempty"`
	ProjectRole *ProjectRole `json:"project_role,omitempty"`
}

// AddProjectMemberRequest represents the request to add a member to a project
type AddProjectMemberRequest struct {
	UserID        uuid.UUID `json:"user_id"`
	ProjectRoleID uuid.UUID `json:"project_role_id"`
}

// ProjectMemberResponse represents a member in the project member list
type ProjectMemberResponse struct {
	ProjectID     uuid.UUID  `json:"project_id"`
	UserID        uuid.UUID  `json:"user_id"`
	ProjectRoleID uuid.UUID  `json:"project_role_id"`
	JoinedAt      time.Time  `json:"joined_at"`
	AddedBy       *uuid.UUID `json:"added_by,omitempty"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	RoleName      string     `json:"role_name"`
	RoleDisplay   string     `json:"role_display"`
}

// EligibleUser represents a system user eligible to be added to a project
type EligibleUser struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

// ProjectRoleResponse represents a role in the API response
type ProjectRoleResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	IsActive    bool      `json:"is_active"`
}
