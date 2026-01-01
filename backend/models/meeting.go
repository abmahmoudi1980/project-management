package models

import (
	"time"

	"github.com/google/uuid"
)

// Meeting represents a scheduled team meeting
type Meeting struct {
	ID              uuid.UUID  `json:"id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description,omitempty"`
	MeetingDate     time.Time  `json:"meeting_date"`
	DurationMinutes int        `json:"duration_minutes"`
	ProjectID       *uuid.UUID `json:"project_id,omitempty"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// MeetingAttendee represents a user attending a meeting
type MeetingAttendee struct {
	MeetingID      uuid.UUID `json:"meeting_id"`
	UserID         uuid.UUID `json:"user_id"`
	ResponseStatus string    `json:"response_status"` // pending, accepted, declined, maybe
	AddedAt        time.Time `json:"added_at"`
}

// MeetingWithAttendees includes meeting details with attendee information
type MeetingWithAttendees struct {
	ID              uuid.UUID  `json:"id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description,omitempty"`
	MeetingDate     time.Time  `json:"meeting_date"`
	DurationMinutes int        `json:"duration_minutes"`
	ProjectID       *uuid.UUID `json:"project_id,omitempty"`
	CreatedBy       uuid.UUID  `json:"created_by"`
	Attendees       []UserInfo `json:"attendees"`
	TotalAttendees  int        `json:"total_attendees"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CreateMeetingRequest is the request body for creating a meeting
type CreateMeetingRequest struct {
	Title           string      `json:"title"`
	Description     *string     `json:"description,omitempty"`
	MeetingDate     time.Time   `json:"meeting_date"`
	DurationMinutes int         `json:"duration_minutes"`
	ProjectID       *uuid.UUID  `json:"project_id,omitempty"`
	AttendeeIDs     []uuid.UUID `json:"attendee_ids"`
}

// UpdateMeetingRequest is the request body for updating a meeting
type UpdateMeetingRequest struct {
	Title           string      `json:"title"`
	Description     *string     `json:"description,omitempty"`
	MeetingDate     time.Time   `json:"meeting_date"`
	DurationMinutes int         `json:"duration_minutes"`
	ProjectID       *uuid.UUID  `json:"project_id,omitempty"`
	AttendeeIDs     []uuid.UUID `json:"attendee_ids"`
}

// UserInfo is a simplified user representation for attendees
type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
}
