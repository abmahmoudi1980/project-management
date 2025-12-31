package handlers

import (
	"project-management/middleware"
	"project-management/models"
	"project-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CommentHandler struct {
	service *services.CommentService
}

func NewCommentHandler(service *services.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) GetCommentsByTask(c *fiber.Ctx) error {
	taskID, err := uuid.Parse(c.Params("taskId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "آیدی تسک نامعتبر است",
				"code":    "INVALID_TASK_ID",
			},
		})
	}

	comments, err := h.service.GetCommentsByTaskIDWithUser(c.Context(), taskID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "تسک یافت نشد",
				"code":    "TASK_NOT_FOUND",
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"comments": comments,
		},
	})
}

func (h *CommentHandler) CreateComment(c *fiber.Ctx) error {
	taskID, err := uuid.Parse(c.Params("taskId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "آیدی تسک نامعتبر است",
				"code":    "INVALID_TASK_ID",
			},
		})
	}

	var req models.CreateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "درخواست نامعتبر است",
				"code":    "INVALID_REQUEST",
			},
		})
	}

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

	comment, err := h.service.CreateComment(c.Context(), taskID, userContext.UserID, req)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err == services.ErrCommentNotFound {
			statusCode = fiber.StatusNotFound
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": err.Error(),
				"code":    "CREATE_COMMENT_FAILED",
			},
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"comment": comment,
		},
	})
}

func (h *CommentHandler) UpdateComment(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "آیدی کامنت نامعتبر است",
				"code":    "INVALID_COMMENT_ID",
			},
		})
	}

	var req models.UpdateCommentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "درخواست نامعتبر است",
				"code":    "INVALID_REQUEST",
			},
		})
	}

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

	comment, err := h.service.UpdateComment(c.Context(), id, userContext.UserID, req)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err == services.ErrCommentNotFound {
			statusCode = fiber.StatusNotFound
		} else if err == services.ErrCommentUnauthorized {
			statusCode = fiber.StatusForbidden
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": err.Error(),
				"code":    "UPDATE_COMMENT_FAILED",
			},
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"comment": comment,
		},
	})
}

func (h *CommentHandler) DeleteComment(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": "آیدی کامنت نامعتبر است",
				"code":    "INVALID_COMMENT_ID",
			},
		})
	}

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

	err = h.service.DeleteComment(c.Context(), id, userContext.UserID)
	if err != nil {
		statusCode := fiber.StatusInternalServerError
		if err == services.ErrCommentNotFound {
			statusCode = fiber.StatusNotFound
		} else if err == services.ErrCommentUnauthorized {
			statusCode = fiber.StatusForbidden
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"success": false,
			"error": fiber.Map{
				"message": err.Error(),
				"code":    "DELETE_COMMENT_FAILED",
			},
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
