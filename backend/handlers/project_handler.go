package handlers

import (
	"project-management/middleware"
	"project-management/models"
	"project-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	service *services.ProjectService
}

func NewProjectHandler(service *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

func (h *ProjectHandler) GetAllProjects(c *fiber.Ctx) error {
	// Get user from context (set by RequireAuth middleware)
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Get projects filtered by user role
	projects, err := h.service.GetProjectsByUser(c.Context(), userContext.UserID, userContext.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch projects"})
	}
	return c.JSON(projects)
}

func (h *ProjectHandler) CreateProject(c *fiber.Ctx) error {
	var req models.CreateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	// Get authenticated user to set created_by
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	project, err := h.service.CreateProject(c.Context(), req, &userContext.UserID)
	if err != nil {
		if err == models.ErrValidation {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		// Handle validation errors from service layer (identifier, URL, etc.)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(project)
}

func (h *ProjectHandler) GetProject(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid project id"})
	}

	project, err := h.service.GetProjectByID(c.Context(), id)
	if err != nil || project == nil {
		return c.Status(404).JSON(fiber.Map{"error": "project not found"})
	}

	return c.JSON(project)
}

func (h *ProjectHandler) UpdateProject(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid project id"})
	}

	var req models.UpdateProjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	project, err := h.service.UpdateProject(c.Context(), id, req)
	if err != nil {
		if err == models.ErrValidation {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		// Handle validation errors from service layer (identifier, URL, etc.)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(project)
}

func (h *ProjectHandler) DeleteProject(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid project id"})
	}

	if err := h.service.DeleteProject(c.Context(), id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete project"})
	}

	return c.Status(204).Send(nil)
}
