package handlers

import (
	"strconv"

	"project-management/middleware"
	"project-management/models"
	"project-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUsers returns paginated list of users (admin only)
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// Validate and cap limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if page < 1 {
		page = 1
	}

	// Parse filter parameters
	role := c.Query("role", "")
	isActiveStr := c.Query("is_active", "")
	var isActive *bool
	if isActiveStr != "" {
		val := isActiveStr == "true"
		isActive = &val
	}

	// Get users with pagination
	users, total, err := h.userService.GetUsers(c.Context(), page, limit, role, isActive)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "خطا در دریافت لیست کاربران",
				"code":    "SERVER_ERROR",
			},
		})
	}

	totalPages := (total + limit - 1) / limit

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"users": users,
			"pagination": fiber.Map{
				"page":        page,
				"limit":       limit,
				"total":       total,
				"total_pages": totalPages,
			},
		},
	})
}

// GetUserByID returns user details by ID (admin only)
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "شناسه کاربر نامعتبر است",
				"code":    "INVALID_USER_ID",
			},
		})
	}

	user, err := h.userService.GetUserByID(c.Context(), id)
	if err != nil || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "کاربر یافت نشد",
				"code":    "USER_NOT_FOUND",
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": user,
		},
	})
}

// UpdateUserRole changes user role (admin only)
func (h *UserHandler) UpdateUserRole(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "شناسه کاربر نامعتبر است",
				"code":    "INVALID_USER_ID",
			},
		})
	}

	var req models.UpdateUserRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "درخواست نامعتبر است",
				"code":    "INVALID_REQUEST",
			},
		})
	}

	// Validate role
	if req.Role != "admin" && req.Role != "user" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "نقش نامعتبر است",
				"code":    "INVALID_ROLE",
			},
		})
	}

	// Update user role
	user, err := h.userService.UpdateUserRole(c.Context(), id, req.Role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": err.Error(),
				"code":    "UPDATE_ROLE_FAILED",
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": user,
		},
	})
}

// UpdateUserActivation changes user activation status (admin only)
func (h *UserHandler) UpdateUserActivation(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "شناسه کاربر نامعتبر است",
				"code":    "INVALID_USER_ID",
			},
		})
	}

	var req models.UpdateUserActivationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "درخواست نامعتبر است",
				"code":    "INVALID_REQUEST",
			},
		})
	}

	// Get current user from context
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "احراز هویت نشده است",
				"code":    "UNAUTHORIZED",
			},
		})
	}

	// Prevent admin from deactivating themselves
	if id == userContext.UserID && !req.IsActive {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "نمی‌توانید حساب کاربری خود را غیرفعال کنید",
				"code":    "CANNOT_DEACTIVATE_SELF",
			},
		})
	}

	// Update activation status
	user, err := h.userService.UpdateUserActivation(c.Context(), id, req.IsActive)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": err.Error(),
				"code":    "UPDATE_ACTIVATION_FAILED",
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": user,
		},
	})
}
