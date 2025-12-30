package handlers

import (
	"project-management/models"
	"project-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type TimeLogHandler struct {
	service *services.TimeLogService
}

func NewTimeLogHandler(service *services.TimeLogService) *TimeLogHandler {
	return &TimeLogHandler{service: service}
}

func (h *TimeLogHandler) GetTimeLogsByTask(c *fiber.Ctx) error {
	taskID, err := uuid.Parse(c.Params("taskId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid task id"})
	}

	timeLogs, err := h.service.GetTimeLogsByTaskID(c.Context(), taskID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "task not found"})
	}

	return c.JSON(timeLogs)
}

func (h *TimeLogHandler) CreateTimeLog(c *fiber.Ctx) error {
	taskID, err := uuid.Parse(c.Params("taskId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid task id"})
	}

	var req models.CreateTimeLogRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request body"})
	}

	timeLog, err := h.service.CreateTimeLog(c.Context(), taskID, req)
	if err != nil {
		if err == models.ErrValidation || err == models.ErrNotFound {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to create time log"})
	}

	return c.Status(201).JSON(timeLog)
}

func (h *TimeLogHandler) GetTimeLog(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid time log id"})
	}

	timeLog, err := h.service.GetTimeLogByID(c.Context(), id)
	if err != nil || timeLog == nil {
		return c.Status(404).JSON(fiber.Map{"error": "time log not found"})
	}

	return c.JSON(timeLog)
}

func (h *TimeLogHandler) DeleteTimeLog(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid time log id"})
	}

	if err := h.service.DeleteTimeLog(c.Context(), id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete time log"})
	}

	return c.Status(204).Send(nil)
}
