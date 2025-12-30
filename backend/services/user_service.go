package services

import (
	"context"
	"errors"

	"project-management/models"
	"project-management/repositories"

	"github.com/google/uuid"
)

var (
	ErrCannotDeactivateLastAdmin = errors.New("نمی‌توانید آخرین ادمین را غیرفعال کنید")
)

type UserService interface {
	GetUsers(ctx context.Context, page, limit int, role string, isActive *bool) ([]models.User, int, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	UpdateUserRole(ctx context.Context, userID uuid.UUID, role string) (*models.User, error)
	UpdateUserActivation(ctx context.Context, userID uuid.UUID, isActive bool) (*models.User, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUsers(ctx context.Context, page, limit int, role string, isActive *bool) ([]models.User, int, error) {
	offset := (page - 1) * limit
	return s.userRepo.ListPaginated(ctx, limit, offset, role, isActive)
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) UpdateUserRole(ctx context.Context, userID uuid.UUID, role string) (*models.User, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("کاربر یافت نشد")
	}

	// Update role
	user.Role = role
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUserActivation(ctx context.Context, userID uuid.UUID, isActive bool) (*models.User, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return nil, errors.New("کاربر یافت نشد")
	}

	// If deactivating an admin, check if they're the last admin
	if !isActive && user.Role == "admin" {
		// Count active admins
		activeAdminCount, err := s.userRepo.CountActiveAdmins(ctx)
		if err != nil {
			return nil, err
		}

		// If this is the last active admin, prevent deactivation
		if activeAdminCount <= 1 {
			return nil, ErrCannotDeactivateLastAdmin
		}
	}

	// Update activation status
	user.IsActive = isActive
	err = s.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
