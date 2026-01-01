package services

import (
	"context"
	"errors"
	"project-management/models"
	"project-management/repositories"
	"time"

	"github.com/google/uuid"
)

type CreateMeetingInput struct {
	Title           string      `json:"title"`
	Description     *string     `json:"description"`
	MeetingDate     time.Time   `json:"meeting_date"`
	DurationMinutes int         `json:"duration_minutes"`
	ProjectID       *uuid.UUID  `json:"project_id"`
	AttendeeIDs     []uuid.UUID `json:"attendee_ids"`
}

type MeetingService struct {
	meetingRepo *repositories.MeetingRepository
	userRepo    repositories.UserRepository
}

func NewMeetingService(meetingRepo *repositories.MeetingRepository, userRepo repositories.UserRepository) *MeetingService {
	return &MeetingService{
		meetingRepo: meetingRepo,
		userRepo:    userRepo,
	}
}

func (s *MeetingService) GetNextMeetingForUser(ctx context.Context, userID uuid.UUID) (*models.MeetingWithAttendees, error) {
	return s.meetingRepo.GetNextMeetingForUser(ctx, userID)
}

func (s *MeetingService) CreateMeeting(ctx context.Context, userID uuid.UUID, input CreateMeetingInput) (*models.MeetingWithAttendees, error) {
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
			return nil, errors.New("one or more attendees not found")
		}
	}

	meeting := &models.Meeting{
		Title:           input.Title,
		Description:     input.Description,
		MeetingDate:     input.MeetingDate,
		DurationMinutes: input.DurationMinutes,
		ProjectID:       input.ProjectID,
		CreatedBy:       userID,
	}

	err := s.meetingRepo.CreateMeeting(ctx, meeting)
	if err != nil {
		return nil, err
	}

	// Add attendees (including creator)
	attendees := input.AttendeeIDs
	creatorIncluded := false
	for _, id := range attendees {
		if id == userID {
			creatorIncluded = true
			break
		}
	}
	if !creatorIncluded {
		attendees = append(attendees, userID)
	}

	err = s.meetingRepo.AddAttendees(ctx, meeting.ID, attendees)
	if err != nil {
		return nil, err
	}

	return s.meetingRepo.GetMeetingByID(ctx, meeting.ID)
}

func (s *MeetingService) validateMeetingInput(input CreateMeetingInput) error {
	if input.Title == "" || len(input.Title) > 200 {
		return errors.New("title must be between 1 and 200 characters")
	}
	if input.Description != nil && len(*input.Description) > 5000 {
		return errors.New("description must be less than 5000 characters")
	}
	if input.MeetingDate.Before(time.Now()) {
		return errors.New("meeting date must be in the future")
	}
	if input.DurationMinutes <= 0 || input.DurationMinutes > 1440 {
		return errors.New("duration must be between 1 and 1440 minutes")
	}
	if len(input.AttendeeIDs) == 0 {
		return errors.New("at least one attendee is required")
	}
	return nil
}

func (s *MeetingService) ListMeetings(ctx context.Context, userID uuid.UUID, from, to time.Time, limit, offset int) ([]models.Meeting, error) {
	return s.meetingRepo.ListMeetings(ctx, userID, from, to, limit, offset)
}

func (s *MeetingService) GetMeetingByID(ctx context.Context, meetingID uuid.UUID) (*models.MeetingWithAttendees, error) {
	return s.meetingRepo.GetMeetingByID(ctx, meetingID)
}
