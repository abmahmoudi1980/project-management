package repositories

import (
	"context"
	"project-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) GetAll(ctx context.Context) ([]models.Project, error) {
	rows, err := r.db.Query(ctx, "SELECT id, title, description, status, identifier, homepage, is_public, user_id, created_by, created_at, updated_at FROM projects ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		if err := rows.Scan(&p.ID, &p.Title, &p.Description, &p.Status, &p.Identifier, &p.Homepage, &p.IsPublic, &p.UserID, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var p models.Project
	err := r.db.QueryRow(ctx, "SELECT id, title, description, status, identifier, homepage, is_public, user_id, created_by, created_at, updated_at FROM projects WHERE id = $1", id).
		Scan(&p.ID, &p.Title, &p.Description, &p.Status, &p.Identifier, &p.Homepage, &p.IsPublic, &p.UserID, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProjectRepository) Create(ctx context.Context, req models.CreateProjectRequest, createdBy *uuid.UUID) (*models.Project, error) {
	id := uuid.New()
	var p models.Project

	err := r.db.QueryRow(ctx,
		"INSERT INTO projects (id, title, description, status, identifier, homepage, is_public, created_by) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, title, description, status, identifier, homepage, is_public, user_id, created_by, created_at, updated_at",
		id, req.Title, req.Description, req.Status, req.Identifier, req.Homepage, req.IsPublic, createdBy).
		Scan(&p.ID, &p.Title, &p.Description, &p.Status, &p.Identifier, &p.Homepage, &p.IsPublic, &p.UserID, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProjectRepository) Update(ctx context.Context, id uuid.UUID, req models.UpdateProjectRequest) (*models.Project, error) {
	var p models.Project

	err := r.db.QueryRow(ctx,
		"UPDATE projects SET title = $1, description = $2, status = $3, identifier = $4, homepage = $5, is_public = $6 WHERE id = $7 RETURNING id, title, description, status, identifier, homepage, is_public, created_at, updated_at",
		req.Title, req.Description, req.Status, req.Identifier, req.Homepage, req.IsPublic, id).
		Scan(&p.ID, &p.Title, &p.Description, &p.Status, &p.Identifier, &p.Homepage, &p.IsPublic, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM projects WHERE id = $1", id)
	return err
}
