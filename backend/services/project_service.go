package services

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"project-management/models"
	"project-management/repositories"
	"regexp"

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

// GetProjectsByUser returns all projects for admins, or only projects created by the user for regular users
func (s *ProjectService) GetProjectsByUser(ctx context.Context, userID uuid.UUID, role string) ([]models.Project, error) {
	// Admins can see all projects
	if role == "admin" {
		return s.repo.GetAll(ctx)
	}

	// Regular users only see projects they created
	allProjects, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Filter projects by user_id or created_by
	var userProjects []models.Project
	for _, p := range allProjects {
		if (p.UserID != nil && *p.UserID == userID) || (p.CreatedBy != nil && *p.CreatedBy == userID) {
			userProjects = append(userProjects, p)
		}
	}

	return userProjects, nil
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

	// Validate identifier
	if err := s.ValidateProjectIdentifier(ctx, req.Identifier, nil); err != nil {
		return nil, err
	}

	// Validate homepage URL
	if err := s.ValidateHomepageURL(req.Homepage); err != nil {
		return nil, err
	}

	return s.repo.Create(ctx, req)
}

func (s *ProjectService) UpdateProject(ctx context.Context, id uuid.UUID, req models.UpdateProjectRequest) (*models.Project, error) {
	if req.Title == "" {
		return nil, models.ErrValidation
	}

	// Validate identifier with exclusion of current project
	if err := s.ValidateProjectIdentifier(ctx, req.Identifier, &id); err != nil {
		return nil, err
	}

	// Validate homepage URL
	if err := s.ValidateHomepageURL(req.Homepage); err != nil {
		return nil, err
	}

	return s.repo.Update(ctx, id, req)
}

func (s *ProjectService) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// ValidateProjectIdentifier validates the project identifier format and checks uniqueness
func (s *ProjectService) ValidateProjectIdentifier(ctx context.Context, identifier string, excludeID *uuid.UUID) error {
	if identifier == "" {
		return errors.New("identifier is required")
	}

	// Validate format: only alphanumeric, underscore, and hyphen allowed
	matched, err := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, identifier)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("identifier can only contain alphanumeric characters, underscores, and hyphens")
	}

	// Check uniqueness by querying all projects
	projects, err := s.repo.GetAll(ctx)
	if err != nil {
		return err
	}

	for _, p := range projects {
		// Skip the project being updated
		if excludeID != nil && p.ID == *excludeID {
			continue
		}
		if p.Identifier == identifier {
			return fmt.Errorf("identifier '%s' is already in use", identifier)
		}
	}

	return nil
}

// ValidateHomepageURL validates the homepage URL format if provided
func (s *ProjectService) ValidateHomepageURL(homepage *string) error {
	if homepage == nil || *homepage == "" {
		return nil // Homepage is optional
	}

	_, err := url.ParseRequestURI(*homepage)
	if err != nil {
		return errors.New("invalid homepage URL format")
	}

	return nil
}
