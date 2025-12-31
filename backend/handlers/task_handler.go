package handlers

import (
	"project-management/middleware"
	"project-management/models"
	"project-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) GetTasksByProject(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid project id"})
	}

	// Get user from context (set by RequireAuth middleware)
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Get tasks filtered by user role and project ownership
	tasks, err := h.service.GetTasksByUser(c.Context(), userContext.UserID, userContext.Role, projectID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "project not found"})
	}

	return c.JSON(tasks)
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx) error {
	projectID, err := uuid.Parse(c.Params("projectId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid project id"})
	}

	var req models.CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	task, err := h.service.CreateTask(c.Context(), projectID, req)
	if err != nil {
		if err == models.ErrValidation || err == models.ErrNotFound {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		// Handle validation errors from service layer (dates, done_ratio, etc.)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(task)
}

func (h *TaskHandler) GetTask(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid task id"})
	}

	task, err := h.service.GetTaskByIDWithUsers(c.Context(), id)
	if err != nil || task == nil {
		return c.Status(404).JSON(fiber.Map{"error": "task not found"})
	}

	return c.JSON(task)
}

func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid task id"})
	}

	var req models.UpdateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	task, err := h.service.UpdateTask(c.Context(), id, req)
	if err != nil {
		if err == models.ErrValidation {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		// Handle validation errors from service layer (dates, done_ratio, etc.)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(task)
}

func (h *TaskHandler) ToggleTaskCompletion(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid task id"})
	}

	task, err := h.service.ToggleTaskCompletion(c.Context(), id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update task"})
	}

	return c.JSON(task)
}

func (h *TaskHandler) DeleteTask(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid task id"})
	}

	if err := h.service.DeleteTask(c.Context(), id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete task"})
	}

	return c.Status(204).Send(nil)
}
