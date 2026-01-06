package handlers

import (
	"project-management/middleware"
	"project-management/services"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AttachmentHandler struct {
	service *services.AttachmentService
}

func NewAttachmentHandler(service *services.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{service: service}
}

// UploadAttachments handles multipart file uploads for a task
func (h *AttachmentHandler) UploadAttachments(c *fiber.Ctx) error {
	// Parse task ID from URL parameters
	taskID, err := uuid.Parse(c.Params("taskId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID format",
		})
	}

	// Get user context from middleware
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: user context not found",
		})
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse multipart form",
		})
	}

	// Get files from form
	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files provided",
		})
	}

	// Process upload through service
	response, err := h.service.UploadAttachments(c.Context(), taskID, files, &userContext.UserID)
	if err != nil {
		// Check for specific error types
		if err.Error() == "access denied: insufficient permissions" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: insufficient permissions to upload attachments to this task",
			})
		}
		if err.Error() == "task not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Task not found",
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process file uploads",
			"details": err.Error(),
		})
	}

	// Return successful response with upload results
	return c.Status(fiber.StatusCreated).JSON(response)
}

// ListAttachments retrieves all attachments for a task
func (h *AttachmentHandler) ListAttachments(c *fiber.Ctx) error {
	// Parse task ID from URL parameters
	taskID, err := uuid.Parse(c.Params("taskId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID format",
		})
	}

	// Get user context from middleware
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: user context not found",
		})
	}

	// Get attachments through service
	response, err := h.service.GetAttachmentsByTaskID(c.Context(), taskID, &userContext.UserID)
	if err != nil {
		// Check for specific error types
		if err.Error() == "access denied: insufficient permissions" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: insufficient permissions to view attachments for this task",
			})
		}
		if err.Error() == "task not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Task not found",
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve attachments",
			"details": err.Error(),
		})
	}

	return c.JSON(response)
}

// DownloadAttachment serves a file for download with proper security headers
func (h *AttachmentHandler) DownloadAttachment(c *fiber.Ctx) error {
	// Parse attachment ID from URL parameters
	attachmentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid attachment ID format",
		})
	}

	// Get user context from middleware
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: user context not found",
		})
	}

	// Verify access and get attachment metadata
	attachment, err := h.service.GetAttachmentByID(c.Context(), attachmentID, &userContext.UserID)
	if err != nil {
		// Check for specific error types
		if err.Error() == "access denied: insufficient permissions" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: insufficient permissions to download this attachment",
			})
		}
		if err.Error() == "attachment not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Attachment not found",
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve attachment",
			"details": err.Error(),
		})
	}

	// Set security headers to prevent XSS attacks
	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", "attachment; filename=\""+attachment.OriginalFilename+"\"")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-XSS-Protection", "1; mode=block")

	// Serve the file
	return c.SendFile(attachment.FilePath)
}

// GetThumbnail serves thumbnail images for image attachments
func (h *AttachmentHandler) GetThumbnail(c *fiber.Ctx) error {
	// Parse attachment ID from URL parameters
	attachmentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid attachment ID format",
		})
	}

	// Get user context from middleware
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: user context not found",
		})
	}

	// Verify access and get attachment metadata
	attachment, err := h.service.GetAttachmentByID(c.Context(), attachmentID, &userContext.UserID)
	if err != nil {
		// Check for specific error types
		if err.Error() == "access denied: insufficient permissions" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: insufficient permissions to view this attachment",
			})
		}
		if err.Error() == "attachment not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Attachment not found",
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve attachment",
			"details": err.Error(),
		})
	}

	// Check if attachment has a thumbnail
	if !attachment.HasThumbnail || attachment.ThumbnailPath == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Thumbnail not available for this attachment",
		})
	}

	// Set appropriate headers for image serving
	c.Set("Content-Type", "image/jpeg")
	c.Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	c.Set("X-Content-Type-Options", "nosniff")

	// Serve the thumbnail file
	return c.SendFile(*attachment.ThumbnailPath)
}

// DeleteAttachment removes an attachment with proper cleanup
func (h *AttachmentHandler) DeleteAttachment(c *fiber.Ctx) error {
	// Parse attachment ID from URL parameters
	attachmentID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid attachment ID format",
		})
	}

	// Get user context from middleware
	userContext, err := middleware.GetUserFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized: user context not found",
		})
	}

	// Delete attachment through service
	err = h.service.DeleteAttachment(c.Context(), attachmentID, &userContext.UserID)
	if err != nil {
		// Check for specific error types
		if err.Error() == "access denied: insufficient permissions to delete attachment" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied: insufficient permissions to delete this attachment",
			})
		}
		if err.Error() == "attachment not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Attachment not found",
			})
		}
		
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete attachment",
			"details": err.Error(),
		})
	}

	// Return success response with no content
	return c.Status(fiber.StatusNoContent).Send(nil)
}