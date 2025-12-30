package services

import (
	"context"
	"project-management/models"
	"project-management/repositories"

	"github.com/google/uuid"
)

type TaskService struct {
	repo        *repositories.TaskRepository
	projectRepo *repositories.ProjectRepository
}

func NewTaskService(repo *repositories.TaskRepository, projectRepo *repositories.ProjectRepository) *TaskService {
	return &TaskService{repo: repo, projectRepo: projectRepo}
}

func (s *TaskService) GetTasksByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.Task, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil || project == nil {
		return nil, models.ErrNotFound
	}
	return s.repo.GetByProjectID(ctx, projectID)
}

func (s *TaskService) GetTaskByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) CreateTask(ctx context.Context, projectID uuid.UUID, req models.CreateTaskRequest) (*models.Task, error) {
	if req.Title == "" {
		return nil, models.ErrValidation
	}
	if req.Priority == "" {
		req.Priority = "Medium"
	}
	return s.repo.Create(ctx, projectID, req)
}

func (s *TaskService) UpdateTask(ctx context.Context, id uuid.UUID, req models.UpdateTaskRequest) (*models.Task, error) {
	if req.Title == "" {
		return nil, models.ErrValidation
	}
	return s.repo.Update(ctx, id, req)
}

func (s *TaskService) ToggleTaskCompletion(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil || task == nil {
		return nil, models.ErrNotFound
	}

	return s.repo.Update(ctx, id, models.UpdateTaskRequest{
		Title:     task.Title,
		Priority:  task.Priority,
		Completed: !task.Completed,
	})
}

func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
