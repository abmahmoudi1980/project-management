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
	// Get refresh token from cookie
	refreshToken := c.Cookies("refresh_token")

	// Revoke session in database if refresh token exists
	if refreshToken != "" {
		// Note: We don't fail logout if revocation fails (best effort)
		_ = h.authService.RevokeSession(c.Context(), refreshToken)
	}

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

// ForgotPassword handles password reset request
func (h *AuthHandler) ForgotPassword(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "درخواست نامعتبر است",
				"code":    "INVALID_REQUEST",
			},
		})
	}

	// Call service (always returns success to prevent email enumeration)
	err := h.authService.RequestPasswordReset(c.Context(), req.Email)
	if err != nil {
		// Log error but don't expose to user
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "خطا در ارسال ایمیل بازیابی",
				"code":    "SERVER_ERROR",
			},
		})
	}

	// Always return success (security best practice)
	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": "اگر ایمیل شما در سیستم ثبت باشد، لینک بازیابی ارسال می‌شود",
		},
	})
}

// ResetPassword handles password reset with token
func (h *AuthHandler) ResetPassword(c *fiber.Ctx) error {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "درخواست نامعتبر است",
				"code":    "INVALID_REQUEST",
			},
		})
	}

	// Validate token and reset password
	err := h.authService.ResetPassword(c.Context(), req.Token, req.NewPassword)
	if err != nil {
		statusCode := fiber.StatusBadRequest
		if err == services.ErrInvalidToken {
			statusCode = fiber.StatusUnauthorized
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": err.Error(),
				"code":    "PASSWORD_RESET_FAILED",
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"message": "رمز عبور شما با موفقیت تغییر یافت",
		},
	})
}
