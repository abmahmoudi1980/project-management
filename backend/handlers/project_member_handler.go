package handlers

import (
	"errors"
	"project-management/middleware"
	"project-management/models"
	"project-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProjectMemberHandler struct {
	service *services.ProjectMemberService
}

func NewProjectMemberHandler(service *services.ProjectMemberService) *ProjectMemberHandler {
	return &ProjectMemberHandler{service: service}
}

// AddMember handles POST /api/projects/:projectId/members
func (h *ProjectMemberHandler) AddMember(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid project ID"})
	}

	var req models.AddProjectMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Get authenticated user
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	err = h.service.AddMember(c.Context(), projectID, req, userContext.UserID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrProjectNotFound):
			return c.Status(404).JSON(fiber.Map{"error": "project not found"})
		case errors.Is(err, services.ErrInvalidRole):
			return c.Status(400).JSON(fiber.Map{"error": "invalid or inactive project role"})
		case errors.Is(err, services.ErrUserNotEligible):
			return c.Status(400).JSON(fiber.Map{"error": "user is not eligible for this project"})
		case errors.Is(err, services.ErrMemberAlreadyExists):
			return c.Status(409).JSON(fiber.Map{"error": "user is already a member of this project"})
		default:
			return c.Status(500).JSON(fiber.Map{"error": "failed to add member"})
		}
	}

	return c.Status(201).JSON(fiber.Map{"message": "member added successfully"})
}

// ListMembers handles GET /api/projects/:projectId/members
func (h *ProjectMemberHandler) ListMembers(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid project ID"})
	}

	members, err := h.service.ListMembers(c.Context(), projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch members"})
	}

	return c.JSON(members)
}

// GetEligibleUsers handles GET /api/projects/:projectId/members/eligible-users
func (h *ProjectMemberHandler) GetEligibleUsers(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid project ID"})
	}

	users, err := h.service.GetEligibleUsers(c.Context(), projectID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch eligible users"})
	}

	return c.JSON(users)
}

// GetProjectRoles handles GET /api/project-roles
func (h *ProjectMemberHandler) GetProjectRoles(c *fiber.Ctx) error {
	roles, err := h.service.GetActiveRoles(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch project roles"})
	}

	return c.JSON(roles)
}
