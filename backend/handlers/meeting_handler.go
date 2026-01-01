package handlers

import (
	"project-management/middleware"
	"project-management/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MeetingHandler struct {
	service *services.MeetingService
}

func NewMeetingHandler(service *services.MeetingService) *MeetingHandler {
	return &MeetingHandler{service: service}
}

func (h *MeetingHandler) GetNextMeeting(c *fiber.Ctx) error {
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	meeting, err := h.service.GetNextMeetingForUser(c.Context(), userContext.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch next meeting"})
	}

	if meeting == nil {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.JSON(meeting)
}

func (h *MeetingHandler) CreateMeeting(c *fiber.Ctx) error {
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var input services.CreateMeetingInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	meeting, err := h.service.CreateMeeting(c.Context(), userContext.UserID, input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(meeting)
}

func (h *MeetingHandler) ListMeetings(c *fiber.Ctx) error {
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")

	from := time.Now().AddDate(0, 0, -30)
	if fromStr != "" {
		if f, err := time.Parse(time.RFC3339, fromStr); err == nil {
			from = f
		}
	}

	to := time.Now().AddDate(0, 0, 30)
	if toStr != "" {
		if t, err := time.Parse(time.RFC3339, toStr); err == nil {
			to = t
		}
	}

	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	meetings, err := h.service.ListMeetings(c.Context(), userContext.UserID, from, to, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch meetings"})
	}

	return c.JSON(meetings)
}

func (h *MeetingHandler) GetMeeting(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid meeting id"})
	}

	meeting, err := h.service.GetMeetingByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to fetch meeting"})
	}

	if meeting == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "meeting not found"})
	}

	return c.JSON(meeting)
}
