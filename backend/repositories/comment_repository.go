package repositories

import (
	"context"
	"project-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) GetByTaskID(ctx context.Context, taskID uuid.UUID) ([]models.Comment, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, task_id, user_id, content, created_at, updated_at
		 FROM comments
		 WHERE task_id = $1
		 ORDER BY created_at ASC`,
		taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.TaskID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *CommentRepository) GetByTaskIDWithUser(ctx context.Context, taskID uuid.UUID) ([]models.CommentWithUser, error) {
	rows, err := r.db.Query(ctx,
		`SELECT c.id, c.task_id, c.user_id, u.username, c.content, c.created_at, c.updated_at
		 FROM comments c
		 JOIN users u ON c.user_id = u.id
		 WHERE c.task_id = $1
		 ORDER BY c.created_at ASC`,
		taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.CommentWithUser
	for rows.Next() {
		var c models.CommentWithUser
		if err := rows.Scan(&c.ID, &c.TaskID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (r *CommentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Comment, error) {
	var c models.Comment
	err := r.db.QueryRow(ctx,
		`SELECT id, task_id, user_id, content, created_at, updated_at
		 FROM comments
		 WHERE id = $1`,
		id).
		Scan(&c.ID, &c.TaskID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CommentRepository) GetByIDWithUser(ctx context.Context, id uuid.UUID) (*models.CommentWithUser, error) {
	var c models.CommentWithUser
	err := r.db.QueryRow(ctx,
		`SELECT c.id, c.task_id, c.user_id, u.username, c.content, c.created_at, c.updated_at
		 FROM comments c
		 JOIN users u ON c.user_id = u.id
		 WHERE c.id = $1`,
		id).
		Scan(&c.ID, &c.TaskID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt, &c.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CommentRepository) Create(ctx context.Context, taskID uuid.UUID, userID uuid.UUID, req models.CreateCommentRequest) (*models.Comment, error) {
	id := uuid.New()
	var c models.Comment

	err := r.db.QueryRow(ctx,
		`INSERT INTO comments (id, task_id, user_id, content)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, task_id, user_id, content, created_at, updated_at`,
		id, taskID, userID, req.Content).
		Scan(&c.ID, &c.TaskID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CommentRepository) Update(ctx context.Context, id uuid.UUID, req models.UpdateCommentRequest) (*models.Comment, error) {
	var c models.Comment
	err := r.db.QueryRow(ctx,
		`UPDATE comments
		 SET content = $1, updated_at = CURRENT_TIMESTAMP
		 WHERE id = $2
		 RETURNING id, task_id, user_id, content, created_at, updated_at`,
		req.Content, id).
		Scan(&c.ID, &c.TaskID, &c.UserID, &c.Content, &c.CreatedAt, &c.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CommentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM comments WHERE id = $1", id)
	return err
}
