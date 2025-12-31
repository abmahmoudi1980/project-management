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
		"SELECT id, project_id, title, description, priority, completed, assignee_id, author_id, category, start_date, due_date, estimated_hours, done_ratio, created_at, updated_at FROM tasks WHERE project_id = $1 ORDER BY created_at DESC",
		projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Priority, &t.Completed, &t.AssigneeID, &t.AuthorID, &t.Category, &t.StartDate, &t.DueDate, &t.EstimatedHours, &t.DoneRatio, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	var t models.Task
	err := r.db.QueryRow(ctx,
		"SELECT id, project_id, title, description, priority, completed, assignee_id, author_id, category, start_date, due_date, estimated_hours, done_ratio, created_at, updated_at FROM tasks WHERE id = $1", id).
		Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Priority, &t.Completed, &t.AssigneeID, &t.AuthorID, &t.Category, &t.StartDate, &t.DueDate, &t.EstimatedHours, &t.DoneRatio, &t.CreatedAt, &t.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TaskRepository) GetByIDWithUsers(ctx context.Context, id uuid.UUID) (*models.TaskWithUsers, error) {
	var t models.TaskWithUsers
	err := r.db.QueryRow(ctx,
		`SELECT t.id, t.project_id, t.title, t.description, t.priority, t.completed,
		        t.assignee_id, t.author_id, t.category, t.start_date, t.due_date,
		        t.estimated_hours, t.done_ratio, t.created_at, t.updated_at,
		        assignee.username as assignee_name,
		        author.username as author_name
		 FROM tasks t
		 LEFT JOIN users assignee ON t.assignee_id = assignee.id
		 LEFT JOIN users author ON t.author_id = author.id
		 WHERE t.id = $1`, id).
		Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Priority, &t.Completed,
			&t.AssigneeID, &t.AuthorID, &t.Category, &t.StartDate, &t.DueDate,
			&t.EstimatedHours, &t.DoneRatio, &t.CreatedAt, &t.UpdatedAt,
			&t.AssigneeName, &t.AuthorName)

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
		"INSERT INTO tasks (id, project_id, title, description, priority, assignee_id, author_id, category, start_date, due_date, estimated_hours, done_ratio) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, project_id, title, description, priority, completed, assignee_id, author_id, category, start_date, due_date, estimated_hours, done_ratio, created_at, updated_at",
		id, projectID, req.Title, req.Description, req.Priority, req.AssigneeID, req.AuthorID, req.Category, req.StartDate, req.DueDate, req.EstimatedHours, req.DoneRatio).
		Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Priority, &t.Completed, &t.AssigneeID, &t.AuthorID, &t.Category, &t.StartDate, &t.DueDate, &t.EstimatedHours, &t.DoneRatio, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TaskRepository) Update(ctx context.Context, id uuid.UUID, req models.UpdateTaskRequest) (*models.Task, error) {
	var t models.Task

	err := r.db.QueryRow(ctx,
		"UPDATE tasks SET title = $1, description = $2, priority = $3, completed = $4, assignee_id = $5, author_id = $6, category = $7, start_date = $8, due_date = $9, estimated_hours = $10, done_ratio = $11 WHERE id = $12 RETURNING id, project_id, title, description, priority, completed, assignee_id, author_id, category, start_date, due_date, estimated_hours, done_ratio, created_at, updated_at",
		req.Title, req.Description, req.Priority, req.Completed, req.AssigneeID, req.AuthorID, req.Category, req.StartDate, req.DueDate, req.EstimatedHours, req.DoneRatio, id).
		Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Priority, &t.Completed, &t.AssigneeID, &t.AuthorID, &t.Category, &t.StartDate, &t.DueDate, &t.EstimatedHours, &t.DoneRatio, &t.CreatedAt, &t.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	return err
}
