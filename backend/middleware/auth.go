package middleware

import (
	"strings"

	"project-management/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// User context key
type contextKey string

const (
	UserContextKey contextKey = "user"
)

// UserContext represents the authenticated user in the request context
type UserContext struct {
	UserID uuid.UUID
	Role   string
}

// RequireAuth middleware validates JWT token and sets user context
func RequireAuth(c *fiber.Ctx) error {
	// Get token from cookie
	tokenString := c.Cookies("access_token")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "احراز هویت نشده است",
				"code":    "UNAUTHORIZED",
			},
		})
	}

	// Parse and validate token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "توکن نامعتبر یا منقضی شده است",
				"code":    "INVALID_TOKEN",
			},
		})
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "توکن نامعتبر است",
				"code":    "INVALID_TOKEN",
			},
		})
	}

	// Check token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "access" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "نوع توکن نامعتبر است",
				"code":    "INVALID_TOKEN_TYPE",
			},
		})
	}

	// Extract user ID and role
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "توکن نامعتبر است",
				"code":    "INVALID_TOKEN",
			},
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "توکن نامعتبر است",
				"code":    "INVALID_TOKEN",
			},
		})
	}

	role, ok := claims["role"].(string)
	if !ok {
		role = "user" // Default role
	}

	// Set user context
	userContext := &UserContext{
		UserID: userID,
		Role:   role,
	}
	c.Locals(string(UserContextKey), userContext)

	return c.Next()
}

// RequireRole middleware checks if user has required role
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user context
		userContext, ok := c.Locals(string(UserContextKey)).(*UserContext)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"message": "احراز هویت نشده است",
					"code":    "UNAUTHORIZED",
				},
			})
		}

		// Check if user has required role
		hasRole := false
		for _, role := range allowedRoles {
			if strings.EqualFold(userContext.Role, role) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"message": "دسترسی غیرمجاز",
					"code":    "FORBIDDEN",
				},
			})
		}

		return c.Next()
	}
}

// GetUserFromContext extracts user info from Fiber context
func GetUserFromContext(c *fiber.Ctx) (*UserContext, error) {
	userContext, ok := c.Locals(string(UserContextKey)).(*UserContext)
	if !ok {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "احراز هویت نشده است")
	}
	return userContext, nil
}
