package services

import (
	"context"
	"errors"
	"project-management/models"
	"project-management/repositories"

	"github.com/google/uuid"
)

var (
	ErrMemberAlreadyExists = errors.New("user is already a member of this project")
	ErrInvalidRole         = errors.New("invalid or inactive project role")
	ErrUserNotEligible     = errors.New("user is not eligible for this project")
	ErrProjectNotFound     = errors.New("project not found")
	ErrMemberNotFound      = errors.New("member not found in this project")
)

type ProjectMemberService struct {
	memberRepo  *repositories.ProjectMemberRepository
	projectRepo *repositories.ProjectRepository
	userRepo    repositories.UserRepository
}

func NewProjectMemberService(
	memberRepo *repositories.ProjectMemberRepository,
	projectRepo *repositories.ProjectRepository,
	userRepo repositories.UserRepository,
) *ProjectMemberService {
	return &ProjectMemberService{
		memberRepo:  memberRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// AddMember adds a user to a project with a specific role
func (s *ProjectMemberService) AddMember(ctx context.Context, projectID uuid.UUID, req models.AddProjectMemberRequest, addedBy uuid.UUID) error {
	// Check if project exists
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return ErrProjectNotFound
	}

	// Check if role exists and is active
	role, err := s.memberRepo.GetRoleByID(ctx, req.ProjectRoleID)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrInvalidRole
	}

	// Check if user exists and is active
	user, err := s.userRepo.GetByID(ctx, req.UserID)
	if err != nil {
		return err
	}
	if user == nil || !user.IsActive {
		return ErrUserNotEligible
	}

	// Check if user is already a member
	exists, err := s.memberRepo.IsUserProjectMember(ctx, projectID, req.UserID)
	if err != nil {
		return err
	}
	if exists {
		return ErrMemberAlreadyExists
	}

	// Add the member
	return s.memberRepo.AddMember(ctx, projectID, req.UserID, req.ProjectRoleID, addedBy)
}

// ListMembers returns all members of a project
func (s *ProjectMemberService) ListMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMemberResponse, error) {
	return s.memberRepo.ListMembers(ctx, projectID)
}

// GetEligibleUsers returns users who can be added to a project
func (s *ProjectMemberService) GetEligibleUsers(ctx context.Context, projectID uuid.UUID) ([]models.EligibleUser, error) {
	return s.memberRepo.ListEligibleUsers(ctx, projectID)
}

// GetActiveRoles returns all active project roles
func (s *ProjectMemberService) GetActiveRoles(ctx context.Context) ([]models.ProjectRole, error) {
	return s.memberRepo.GetActiveRoles(ctx)
}

// GetMember returns a specific membership by project and user
func (s *ProjectMemberService) GetMember(ctx context.Context, projectID, userID uuid.UUID) (*models.ProjectMember, error) {
	return s.memberRepo.GetMemberByProjectAndUser(ctx, projectID, userID)
}

// RemoveMember removes a user from a project
func (s *ProjectMemberService) RemoveMember(ctx context.Context, projectID, userID uuid.UUID) error {
	// Check if project exists
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return ErrProjectNotFound
	}

	// Check if member exists
	member, err := s.memberRepo.GetMemberByProjectAndUser(ctx, projectID, userID)
	if err != nil {
		return err
	}
	if member == nil {
		return ErrMemberNotFound
	}

	// Remove the member
	return s.memberRepo.RemoveMember(ctx, projectID, userID)
}
