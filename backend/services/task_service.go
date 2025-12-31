package services

import (
	"context"
	"errors"
	"project-management/models"
	"project-management/repositories"
	"time"

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

// GetTasksByUser returns tasks filtered by user role
// Admins can see all tasks, regular users only see tasks from projects they own
func (s *TaskService) GetTasksByUser(ctx context.Context, userID uuid.UUID, role string, projectID uuid.UUID) ([]models.Task, error) {
	// Verify project exists
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil || project == nil {
		return nil, models.ErrNotFound
	}

	// Admins can see all tasks
	if role == "admin" {
		return s.repo.GetByProjectID(ctx, projectID)
	}

	// Regular users can only see tasks from projects they own
	if (project.UserID != nil && *project.UserID == userID) || (project.CreatedBy != nil && *project.CreatedBy == userID) {
		return s.repo.GetByProjectID(ctx, projectID)
	}

	// User doesn't own this project, return empty list
	return []models.Task{}, nil
}

func (s *TaskService) GetTasksByUserPaginated(ctx context.Context, userID uuid.UUID, role string, projectID uuid.UUID, page int, pageSize int) (*models.PaginatedTasksResponse, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil || project == nil {
		return nil, models.ErrNotFound
	}

	offset := (page - 1) * pageSize

	total, err := s.repo.GetTotalTasksByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var tasks []models.Task
	if role == "admin" {
		tasks, err = s.repo.GetByProjectIDPaginated(ctx, projectID, pageSize, offset)
	} else if (project.UserID != nil && *project.UserID == userID) || (project.CreatedBy != nil && *project.CreatedBy == userID) {
		tasks, err = s.repo.GetByProjectIDPaginated(ctx, projectID, pageSize, offset)
	} else {
		tasks = []models.Task{}
	}

	if err != nil {
		return nil, err
	}

	hasMore := (page * pageSize) < total

	return &models.PaginatedTasksResponse{
		Tasks:    tasks,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		HasMore:  hasMore,
	}, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TaskService) GetTaskByIDWithUsers(ctx context.Context, id uuid.UUID) (*models.TaskWithUsers, error) {
	return s.repo.GetByIDWithUsers(ctx, id)
}

func (s *TaskService) CreateTask(ctx context.Context, projectID uuid.UUID, req models.CreateTaskRequest) (*models.Task, error) {
	if req.Title == "" {
		return nil, models.ErrValidation
	}
	if req.Priority == "" {
		req.Priority = "Medium"
	}

	// Validate dates
	if err := s.ValidateTaskDates(req.StartDate, req.DueDate); err != nil {
		return nil, err
	}

	// Validate done_ratio
	if err := s.ValidateDoneRatio(req.DoneRatio); err != nil {
		return nil, err
	}

	// Validate estimated_hours
	if err := s.ValidateEstimatedHours(req.EstimatedHours); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, projectID, req)
}

func (s *TaskService) UpdateTask(ctx context.Context, id uuid.UUID, req models.UpdateTaskRequest) (*models.Task, error) {
	if req.Title == "" {
		return nil, models.ErrValidation
	}

	// Validate dates
	if err := s.ValidateTaskDates(req.StartDate, req.DueDate); err != nil {
		return nil, err
	}

	// Validate done_ratio
	if err := s.ValidateDoneRatio(req.DoneRatio); err != nil {
		return nil, err
	}

	// Validate estimated_hours
	if err := s.ValidateEstimatedHours(req.EstimatedHours); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, id, req)
}

func (s *TaskService) ToggleTaskCompletion(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil || task == nil {
		return nil, models.ErrNotFound
	}

	return s.repo.Update(ctx, id, models.UpdateTaskRequest{
		Title:          task.Title,
		Description:    task.Description,
		Priority:       task.Priority,
		Completed:      !task.Completed,
		AssigneeID:     task.AssigneeID,
		AuthorID:       task.AuthorID,
		Category:       task.Category,
		StartDate:      task.StartDate,
		DueDate:        task.DueDate,
		EstimatedHours: task.EstimatedHours,
		DoneRatio:      task.DoneRatio,
	})
}

func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// ValidateTaskDates ensures due_date >= start_date when both are provided
func (s *TaskService) ValidateTaskDates(startDate, dueDate *time.Time) error {
	if startDate != nil && dueDate != nil {
		if dueDate.Before(*startDate) {
			return errors.New("due date must be equal to or after start date")
		}
	}
	return nil
}

// ValidateDoneRatio ensures done_ratio is between 0 and 100
func (s *TaskService) ValidateDoneRatio(doneRatio int) error {
	if doneRatio < 0 || doneRatio > 100 {
		return errors.New("done ratio must be between 0 and 100")
	}
	return nil
}

// ValidateEstimatedHours ensures estimated_hours >= 0 if provided
func (s *TaskService) ValidateEstimatedHours(estimatedHours *float64) error {
	if estimatedHours != nil && *estimatedHours < 0 {
		return errors.New("estimated hours must be greater than or equal to 0")
	}
	return nil
}
