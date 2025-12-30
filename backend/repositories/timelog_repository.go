package repositories

import (
	"context"
	"project-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TimeLogRepository struct {
	db *pgxpool.Pool
}

func NewTimeLogRepository(db *pgxpool.Pool) *TimeLogRepository {
	return &TimeLogRepository{db: db}
}

func (r *TimeLogRepository) GetByTaskID(ctx context.Context, taskID uuid.UUID) ([]models.TimeLog, error) {
	rows, err := r.db.Query(ctx,
		"SELECT id, task_id, date, duration_minutes, note, created_at FROM time_logs WHERE task_id = $1 ORDER BY date DESC",
		taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var timeLogs []models.TimeLog
	for rows.Next() {
		var tl models.TimeLog
		if err := rows.Scan(&tl.ID, &tl.TaskID, &tl.Date, &tl.DurationMinutes, &tl.Note, &tl.CreatedAt); err != nil {
			return nil, err
		}
		timeLogs = append(timeLogs, tl)
	}

	return timeLogs, nil
}

func (r *TimeLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TimeLog, error) {
	var tl models.TimeLog
	err := r.db.QueryRow(ctx,
		"SELECT id, task_id, date, duration_minutes, note, created_at FROM time_logs WHERE id = $1", id).
		Scan(&tl.ID, &tl.TaskID, &tl.Date, &tl.DurationMinutes, &tl.Note, &tl.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &tl, nil
}

func (r *TimeLogRepository) Create(ctx context.Context, taskID uuid.UUID, req models.CreateTimeLogRequest) (*models.TimeLog, error) {
	id := uuid.New()
	var tl models.TimeLog

	err := r.db.QueryRow(ctx,
		"INSERT INTO time_logs (id, task_id, date, duration_minutes, note) VALUES ($1, $2, $3, $4, $5) RETURNING id, task_id, date, duration_minutes, note, created_at",
		id, taskID, req.Date, req.DurationMinutes, req.Note).
		Scan(&tl.ID, &tl.TaskID, &tl.Date, &tl.DurationMinutes, &tl.Note, &tl.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &tl, nil
}

func (r *TimeLogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM time_logs WHERE id = $1", id)
	return err
}
