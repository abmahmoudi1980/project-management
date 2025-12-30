package services

import (
	"context"
	"project-management/models"
	"project-management/repositories"

	"github.com/google/uuid"
)

type TimeLogService struct {
	repo     *repositories.TimeLogRepository
	taskRepo *repositories.TaskRepository
}

func NewTimeLogService(repo *repositories.TimeLogRepository, taskRepo *repositories.TaskRepository) *TimeLogService {
	return &TimeLogService{repo: repo, taskRepo: taskRepo}
}

func (s *TimeLogService) GetTimeLogsByTaskID(ctx context.Context, taskID uuid.UUID) ([]models.TimeLog, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return nil, models.ErrNotFound
	}
	return s.repo.GetByTaskID(ctx, taskID)
}

func (s *TimeLogService) GetTimeLogByID(ctx context.Context, id uuid.UUID) (*models.TimeLog, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TimeLogService) CreateTimeLog(ctx context.Context, taskID uuid.UUID, req models.CreateTimeLogRequest) (*models.TimeLog, error) {
	if req.DurationMinutes <= 0 {
		return nil, models.ErrValidation
	}
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return nil, models.ErrNotFound
	}
	return s.repo.Create(ctx, taskID, req)
}

func (s *TimeLogService) DeleteTimeLog(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
