package services

import (
	"context"
	"project-management/models"
	"project-management/repositories"

	"github.com/google/uuid"
)

type ProjectService struct {
	repo *repositories.ProjectRepository
}

func NewProjectService(repo *repositories.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) GetAllProjects(ctx context.Context) ([]models.Project, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProjectService) CreateProject(ctx context.Context, req models.CreateProjectRequest) (*models.Project, error) {
	if req.Title == "" {
		return nil, models.ErrValidation
	}
	if req.Status == "" {
		req.Status = "active"
	}
	return s.repo.Create(ctx, req)
}

func (s *ProjectService) UpdateProject(ctx context.Context, id uuid.UUID, req models.UpdateProjectRequest) (*models.Project, error) {
	if req.Title == "" {
		return nil, models.ErrValidation
	}
	return s.repo.Update(ctx, id, req)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
