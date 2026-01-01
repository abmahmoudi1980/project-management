package services

import (
	"context"
	"errors"
	"project-management/models"
	"project-management/repositories"
	"strings"
	"time"

	"github.com/google/uuid"
)

type MeetingService struct {
	repo     *repositories.MeetingRepository
	userRepo repositories.UserRepository
}

func NewMeetingService(repo *repositories.MeetingRepository, userRepo repositories.UserRepository) *MeetingService {
	return &MeetingService{
		repo:     repo,
		userRepo: userRepo,
	}
}

// GetNextMeetingForUser returns the next upcoming meeting for a user
func (s *MeetingService) GetNextMeetingForUser(ctx context.Context, userID uuid.UUID) (*models.MeetingWithAttendees, error) {
	return s.repo.GetNextMeetingForUser(ctx, userID)
}

// CreateMeeting validates and creates a new meeting
func (s *MeetingService) CreateMeeting(ctx context.Context, userID uuid.UUID, input *models.CreateMeetingRequest) (*models.MeetingWithAttendees, error) {
	// Validate input
	if err := s.validateMeetingInput(input); err != nil {
		return nil, err
	}

	// Verify attendees exist
	for _, attendeeID := range input.AttendeeIDs {
		user, err := s.userRepo.GetByID(ctx, attendeeID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("attendee not found: " + attendeeID.String())
		}
	}

	// Create meeting
	meeting := &models.Meeting{
		ID:              uuid.New(),
		Title:           input.Title,
		Description:     input.Description,
		MeetingDate:     input.MeetingDate,
		DurationMinutes: input.DurationMinutes,
		ProjectID:       input.ProjectID,
		CreatedBy:       userID,
		CreatedAt:       time.Now().UTC(),
		UpdatedAt:       time.Now().UTC(),
	}

	if err := s.repo.CreateMeeting(ctx, meeting); err != nil {
		return nil, err
	}

	// Add attendees (including creator)
	attendeeIDs := input.AttendeeIDs
	if !containsUUID(attendeeIDs, userID) {
		attendeeIDs = append(attendeeIDs, userID)
	}

	if err := s.repo.AddAttendees(ctx, meeting.ID, attendeeIDs); err != nil {
		return nil, err
	}

	// Retrieve the created meeting with attendees
	return s.repo.GetMeetingByID(ctx, meeting.ID)
}

// validateMeetingInput validates meeting creation input
func (s *MeetingService) validateMeetingInput(input *models.CreateMeetingRequest) error {
	// Title validation
	title := strings.TrimSpace(input.Title)
	if len(title) < 1 || len(title) > 200 {
		return errors.New("title must be between 1 and 200 characters")
	}

	// Description validation
	if input.Description != nil {
		desc := strings.TrimSpace(*input.Description)
		if len(desc) > 5000 {
			return errors.New("description must not exceed 5000 characters")
		}
	}

	// Duration validation
	if input.DurationMinutes < 1 || input.DurationMinutes > 1440 {
		return errors.New("duration must be between 1 and 1440 minutes")
	}

	// Meeting date validation - must be in future
	if input.MeetingDate.Before(time.Now().UTC()) {
		return errors.New("meeting date must be in the future")
	}

	// Attendees validation
	if len(input.AttendeeIDs) == 0 {
		return errors.New("at least one attendee required")
	}

	return nil
}

// Helper function to check if UUID exists in slice
func containsUUID(slice []uuid.UUID, id uuid.UUID) bool {
	for _, item := range slice {
		if item == id {
			return true
		}
	}
	return false
}
