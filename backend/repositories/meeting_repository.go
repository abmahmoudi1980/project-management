package repositories

import (
	"context"
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

func (r *MeetingRepository) CreateMeeting(ctx context.Context, m *models.Meeting) error {
	err := r.db.QueryRow(ctx, `
		INSERT INTO meetings (title, description, meeting_date, duration_minutes, project_id, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`, m.Title, m.Description, m.MeetingDate, m.DurationMinutes, m.ProjectID, m.CreatedBy).
		Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
	return err
}

func (r *MeetingRepository) AddAttendees(ctx context.Context, meetingID uuid.UUID, userIDs []uuid.UUID) error {
	batch := &pgx.Batch{}
	for _, userID := range userIDs {
		batch.Queue(`
			INSERT INTO meeting_attendees (meeting_id, user_id)
			VALUES ($1, $2)
			ON CONFLICT (meeting_id, user_id) DO NOTHING
		`, meetingID, userID)
	}
	br := r.db.SendBatch(ctx, batch)
	defer br.Close()

	for i := 0; i < len(userIDs); i++ {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *MeetingRepository) GetNextMeetingForUser(ctx context.Context, userID uuid.UUID) (*models.MeetingWithAttendees, error) {
	var m models.MeetingWithAttendees
	err := r.db.QueryRow(ctx, `
		SELECT 
			m.id, m.title, m.description, m.meeting_date, m.duration_minutes, m.project_id, m.created_by, m.created_at, m.updated_at
		FROM meetings m
		JOIN meeting_attendees ma ON ma.meeting_id = m.id
		WHERE ma.user_id = $1 AND m.meeting_date >= NOW()
		ORDER BY m.meeting_date ASC
		LIMIT 1
	`, userID).Scan(&m.ID, &m.Title, &m.Description, &m.MeetingDate, &m.DurationMinutes, &m.ProjectID, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get attendees
	m.Attendees, m.TotalAttendees, err = r.getMeetingAttendees(ctx, m.ID)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *MeetingRepository) GetMeetingByID(ctx context.Context, meetingID uuid.UUID) (*models.MeetingWithAttendees, error) {
	var m models.MeetingWithAttendees
	err := r.db.QueryRow(ctx, `
		SELECT 
			id, title, description, meeting_date, duration_minutes, project_id, created_by, created_at, updated_at
		FROM meetings
		WHERE id = $1
	`, meetingID).Scan(&m.ID, &m.Title, &m.Description, &m.MeetingDate, &m.DurationMinutes, &m.ProjectID, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Get attendees
	m.Attendees, m.TotalAttendees, err = r.getMeetingAttendees(ctx, m.ID)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *MeetingRepository) getMeetingAttendees(ctx context.Context, meetingID uuid.UUID) ([]models.User, int, error) {
	// Total count
	var total int
	err := r.db.QueryRow(ctx, `
		SELECT COUNT(*) FROM meeting_attendees WHERE meeting_id = $1
	`, meetingID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Attendees (max 3 for dashboard, but maybe more for detail? Let's just get all for now or limit to 10)
	rows, err := r.db.Query(ctx, `
		SELECT u.id, u.username, u.email, u.role, u.is_active, u.created_at, u.updated_at
		FROM users u
		JOIN meeting_attendees ma ON ma.user_id = u.id
		WHERE ma.meeting_id = $1
		LIMIT 10
	`, meetingID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	attendees := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		attendees = append(attendees, u)
	}

	return attendees, total, nil
}

func (r *MeetingRepository) ListMeetings(ctx context.Context, userID uuid.UUID, from, to time.Time, limit, offset int) ([]models.Meeting, error) {
	rows, err := r.db.Query(ctx, `
		SELECT 
			m.id, m.title, m.description, m.meeting_date, m.duration_minutes, m.project_id, m.created_by, m.created_at, m.updated_at
		FROM meetings m
		JOIN meeting_attendees ma ON ma.meeting_id = m.id
		WHERE ma.user_id = $1 AND m.meeting_date BETWEEN $2 AND $3
		ORDER BY m.meeting_date ASC
		LIMIT $4 OFFSET $5
	`, userID, from, to, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meetings := []models.Meeting{}
	for rows.Next() {
		var m models.Meeting
		if err := rows.Scan(&m.ID, &m.Title, &m.Description, &m.MeetingDate, &m.DurationMinutes, &m.ProjectID, &m.CreatedBy, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		meetings = append(meetings, m)
	}

	return meetings, nil
}
