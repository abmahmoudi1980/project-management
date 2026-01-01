package handlers

import (
	"project-management/middleware"
	"project-management/services"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler(service *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetDashboard retrieves all dashboard data
func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
	// Get user from context (set by RequireAuth middleware)
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	// Get dashboard data
	dashboard, err := h.service.GetDashboardData(c.Context(), userContext.UserID, userContext.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch dashboard data", "details": err.Error()})
	}

	return c.JSON(dashboard)
}
