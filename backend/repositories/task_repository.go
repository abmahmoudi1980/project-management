package repositories

import (
	"context"
	"project-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.Task, error) {
	rows, err := r.db.Query(ctx,
		"SELECT id, project_id, title, priority, completed, created_at, updated_at FROM tasks WHERE project_id = $1 ORDER BY created_at DESC",
		projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Priority, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	var t models.Task
	err := r.db.QueryRow(ctx,
		"SELECT id, project_id, title, priority, completed, created_at, updated_at FROM tasks WHERE id = $1", id).
		Scan(&t.ID, &t.ProjectID, &t.Title, &t.Priority, &t.Completed, &t.CreatedAt, &t.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TaskRepository) Create(ctx context.Context, projectID uuid.UUID, req models.CreateTaskRequest) (*models.Task, error) {
	id := uuid.New()
	var t models.Task

	err := r.db.QueryRow(ctx,
		"INSERT INTO tasks (id, project_id, title, priority) VALUES ($1, $2, $3, $4) RETURNING id, project_id, title, priority, completed, created_at, updated_at",
		id, projectID, req.Title, req.Priority).
		Scan(&t.ID, &t.ProjectID, &t.Title, &t.Priority, &t.Completed, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TaskRepository) Update(ctx context.Context, id uuid.UUID, req models.UpdateTaskRequest) (*models.Task, error) {
	var t models.Task

	err := r.db.QueryRow(ctx,
		"UPDATE tasks SET title = $1, priority = $2, completed = $3 WHERE id = $4 RETURNING id, project_id, title, priority, completed, created_at, updated_at",
		req.Title, req.Priority, req.Completed, id).
		Scan(&t.ID, &t.ProjectID, &t.Title, &t.Priority, &t.Completed, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	return err
}
