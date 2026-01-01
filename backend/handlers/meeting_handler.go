package handlers

import (
	"project-management/middleware"
	"project-management/models"
	"project-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MeetingHandler struct {
	service *services.MeetingService
}

func NewMeetingHandler(service *services.MeetingService) *MeetingHandler {
	return &MeetingHandler{service: service}
}

// GetNextMeeting retrieves the next upcoming meeting for the authenticated user
func (h *MeetingHandler) GetNextMeeting(c *fiber.Ctx) error {
	// Get user from context (set by RequireAuth middleware)
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	meeting, err := h.service.GetNextMeetingForUser(c.Context(), userContext.UserID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch meeting", "details": err.Error()})
	}

	if meeting == nil {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.JSON(meeting)
}

// CreateMeeting creates a new meeting
func (h *MeetingHandler) CreateMeeting(c *fiber.Ctx) error {
	// Get user from context (set by RequireAuth middleware)
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req models.CreateMeetingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	meeting, err := h.service.CreateMeeting(c.Context(), userContext.UserID, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(meeting)
}

// ListMeetings retrieves meetings for the authenticated user
func (h *MeetingHandler) ListMeetings(c *fiber.Ctx) error {
	// Get user from context (set by RequireAuth middleware)
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// For now, return next meeting via GetNextMeeting endpoint
	// Full list implementation can be added later
	meeting, err := h.service.GetNextMeetingForUser(c.Context(), userContext.UserID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch meetings"})
	}

	if meeting == nil {
		return c.JSON([]interface{}{})
	}

	return c.JSON([]interface{}{meeting})
}

// GetMeeting retrieves a specific meeting by ID
func (h *MeetingHandler) GetMeeting(c *fiber.Ctx) error {
	// Get user from context (set by RequireAuth middleware)
	_, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	_, err = uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid meeting id"})
	}

	// Note: In production, you would have a proper method on the repository to expose this
	// For now, we can use GetNextMeetingForUser or implement a separate method
	// This is a temporary solution

	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{"error": "not implemented"})
}
