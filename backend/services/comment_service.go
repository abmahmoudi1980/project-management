package services

import (
	"context"
	"errors"
	"project-management/models"
	"project-management/repositories"

	"github.com/google/uuid"
)

var (
	ErrCommentNotFound     = errors.New("کامنت یافت نشد")
	ErrCommentUnauthorized = errors.New("شما مجوز ویرایش یا حذف این کامنت را ندارید")
)

type CommentService struct {
	repo     *repositories.CommentRepository
	taskRepo *repositories.TaskRepository
}

func NewCommentService(repo *repositories.CommentRepository, taskRepo *repositories.TaskRepository) *CommentService {
	return &CommentService{repo: repo, taskRepo: taskRepo}
}

func (s *CommentService) GetCommentsByTaskID(ctx context.Context, taskID uuid.UUID) ([]models.Comment, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return nil, ErrCommentNotFound
	}
	return s.repo.GetByTaskID(ctx, taskID)
}

func (s *CommentService) GetCommentsByTaskIDWithUser(ctx context.Context, taskID uuid.UUID) ([]models.CommentWithUser, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return nil, ErrCommentNotFound
	}
	return s.repo.GetByTaskIDWithUser(ctx, taskID)
}

func (s *CommentService) GetCommentByID(ctx context.Context, id uuid.UUID) (*models.Comment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CommentService) CreateComment(ctx context.Context, taskID uuid.UUID, userID uuid.UUID, req models.CreateCommentRequest) (*models.Comment, error) {
	if req.Content == "" {
		return nil, errors.New("متن کامنت نمی‌تواند خالی باشد")
	}
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return nil, ErrCommentNotFound
	}
	return s.repo.Create(ctx, taskID, userID, req)
}

func (s *CommentService) UpdateComment(ctx context.Context, id uuid.UUID, userID uuid.UUID, req models.UpdateCommentRequest) (*models.Comment, error) {
	if req.Content == "" {
		return nil, errors.New("متن کامنت نمی‌تواند خالی باشد")
	}
	comment, err := s.repo.GetByID(ctx, id)
	if err != nil || comment == nil {
		return nil, ErrCommentNotFound
	}
	if comment.UserID != userID {
		return nil, ErrCommentUnauthorized
	}
	return s.repo.Update(ctx, id, req)
}

func (s *CommentService) DeleteComment(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	comment, err := s.repo.GetByID(ctx, id)
	if err != nil || comment == nil {
		return ErrCommentNotFound
	}
	if comment.UserID != userID {
		return ErrCommentUnauthorized
	}
	return s.repo.Delete(ctx, id)
}
