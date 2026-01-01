package repositories

import (
	"context"
	"errors"
	"project-management/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MeetingRepository struct {
	db *pgxpool.Pool
}

func NewMeetingRepository(db *pgxpool.Pool) *MeetingRepository {
	return &MeetingRepository{db: db}
}

// CreateMeeting inserts a new meeting
func (r *MeetingRepository) CreateMeeting(ctx context.Context, meeting *models.Meeting) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO meetings (id, title, description, meeting_date, duration_minutes, project_id, created_by, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		meeting.ID, meeting.Title, meeting.Description, meeting.MeetingDate, meeting.DurationMinutes,
		meeting.ProjectID, meeting.CreatedBy, meeting.CreatedAt, meeting.UpdatedAt)
	return err
}

// GetMeetingByID retrieves a meeting by ID with attendees
func (r *MeetingRepository) GetMeetingByID(ctx context.Context, meetingID uuid.UUID) (*models.MeetingWithAttendees, error) {
	var m models.MeetingWithAttendees
	err := r.db.QueryRow(ctx,
		`SELECT id, title, description, meeting_date, duration_minutes, project_id, created_by, created_at, updated_at
		 FROM meetings WHERE id = $1`,
		meetingID).
		Scan(&m.ID, &m.Title, &m.Description, &m.MeetingDate, &m.DurationMinutes, &m.ProjectID, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get attendees
	rows, err := r.db.Query(ctx,
		`SELECT u.id, u.full_name FROM users u
		 JOIN meeting_attendees ma ON u.id = ma.user_id
		 WHERE ma.meeting_id = $1 LIMIT 3`,
		meetingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var attendee models.UserInfo
		if err := rows.Scan(&attendee.ID, &attendee.FullName); err != nil {
			return nil, err
		}
		m.Attendees = append(m.Attendees, attendee)
	}

	// Get total attendees count
	var count int
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM meeting_attendees WHERE meeting_id = $1", meetingID).Scan(&count)
	if err != nil {
		return nil, err
	}
	m.TotalAttendees = count

	return &m, nil
}

// GetNextMeetingForUser retrieves the next scheduled meeting for a user
func (r *MeetingRepository) GetNextMeetingForUser(ctx context.Context, userID uuid.UUID) (*models.MeetingWithAttendees, error) {
	var m models.MeetingWithAttendees
	err := r.db.QueryRow(ctx,
		`SELECT id, title, description, meeting_date, duration_minutes, project_id, created_by, created_at, updated_at
		 FROM meetings
		 WHERE id IN (
		   SELECT meeting_id FROM meeting_attendees WHERE user_id = $1
		 )
		 AND meeting_date > NOW()
		 ORDER BY meeting_date ASC
		 LIMIT 1`,
		userID).
		Scan(&m.ID, &m.Title, &m.Description, &m.MeetingDate, &m.DurationMinutes, &m.ProjectID, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get attendees
	rows, err := r.db.Query(ctx,
		`SELECT u.id, u.full_name FROM users u
		 JOIN meeting_attendees ma ON u.id = ma.user_id
		 WHERE ma.meeting_id = $1 LIMIT 3`,
		m.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var attendee models.UserInfo
		if err := rows.Scan(&attendee.ID, &attendee.FullName); err != nil {
			return nil, err
		}
		m.Attendees = append(m.Attendees, attendee)
	}

	// Get total attendees count
	var count int
	err = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM meeting_attendees WHERE meeting_id = $1", m.ID).Scan(&count)
	if err != nil {
		return nil, err
	}
	m.TotalAttendees = count

	return &m, nil
}

// AddAttendees adds users to a meeting
func (r *MeetingRepository) AddAttendees(ctx context.Context, meetingID uuid.UUID, userIDs []uuid.UUID) error {
	if len(userIDs) == 0 {
		return errors.New("at least one attendee required")
	}

	for _, userID := range userIDs {
		_, err := r.db.Exec(ctx,
			`INSERT INTO meeting_attendees (meeting_id, user_id, response_status)
			 VALUES ($1, $2, 'pending')
			 ON CONFLICT DO NOTHING`,
			meetingID, userID)
		if err != nil {
			return err
		}
	}

	return nil
}

// ListMeetings retrieves meetings for a user within a date range
func (r *MeetingRepository) ListMeetings(ctx context.Context, userID uuid.UUID, from, to time.Time, limit, offset int) ([]models.Meeting, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, title, description, meeting_date, duration_minutes, project_id, created_by, created_at, updated_at
		 FROM meetings
		 WHERE id IN (SELECT meeting_id FROM meeting_attendees WHERE user_id = $1)
		 AND meeting_date BETWEEN $2 AND $3
		 ORDER BY meeting_date ASC
		 LIMIT $4 OFFSET $5`,
		userID, from, to, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []models.Meeting
	for rows.Next() {
		var m models.Meeting
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.MeetingDate, &m.DurationMinutes, &m.ProjectID, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		meetings = append(meetings, m)
	}

	return meetings, nil
}
