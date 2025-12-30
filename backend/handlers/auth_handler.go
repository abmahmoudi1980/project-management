package handlers

import (
	"time"

	"project-management/middleware"
	"project-management/models"
	"project-management/services"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "درخواست نامعتبر است",
				"code":    "INVALID_REQUEST",
			},
		})
	}

	// Register user
	user, accessToken, refreshToken, err := h.authService.Register(c.Context(), req)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		if err == services.ErrEmailExists {
			statusCode = fiber.StatusConflict
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": err.Error(),
				"code":    "REGISTRATION_FAILED",
			},
		})
	}

	// Set httpOnly cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
		Expires:  time.Now().Add(15 * time.Minute),
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Path:     "/",
	})

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": user,
		},
	})
}

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "درخواست نامعتبر است",
				"code":    "INVALID_REQUEST",
			},
		})
	}

	// Get user agent and IP address
	userAgent := c.Get("User-Agent")
	ipAddress := c.IP()

	// Login user
	user, accessToken, refreshToken, err := h.authService.Login(c.Context(), req, userAgent, ipAddress)
	if err != nil {
		statusCode := fiber.StatusUnauthorized
		if err == services.ErrAccountLocked {
			statusCode = fiber.StatusForbidden
		} else if err == services.ErrAccountDeactivated {
			statusCode = fiber.StatusForbidden
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": err.Error(),
				"code":    "LOGIN_FAILED",
			},
		})
	}

	// Set httpOnly cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
		Expires:  time.Now().Add(15 * time.Minute),
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HTTPOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: "Lax",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user": user,
		},
	})
}

// GetCurrentUser returns the current authenticated user
func (h *AuthHandler) GetCurrentUser(c *fiber.Ctx) error {
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

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user_id": userContext.UserID,
			"role":    userContext.Role,
		},
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// Clear cookies
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": "با موفقیت خارج شدید",
		},
	})
}
