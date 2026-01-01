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
	ResponseStatus string    `json:"response_status"`
	AddedAt        time.Time `json:"added_at"`
}

// MeetingWithAttendees extends Meeting with attendee information
type MeetingWithAttendees struct {
	Meeting
	Attendees      []User `json:"attendees"`
	TotalAttendees int    `json:"total_attendees"`
}
