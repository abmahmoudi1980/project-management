package services

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"project-management/models"
	"project-management/repositories"
	"strings"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
)

type AttachmentService struct {
	attachmentRepo    *repositories.AttachmentRepository
	taskRepo          *repositories.TaskRepository
	projectRepo       *repositories.ProjectRepository
	fileStorageService *FileStorageService
	fileValidationService *FileValidationService
}

func NewAttachmentService(
	attachmentRepo *repositories.AttachmentRepository,
	taskRepo *repositories.TaskRepository,
	projectRepo *repositories.ProjectRepository,
	fileStorageService *FileStorageService,
	fileValidationService *FileValidationService,
) *AttachmentService {
	return &AttachmentService{
		attachmentRepo:        attachmentRepo,
		taskRepo:              taskRepo,
		projectRepo:           projectRepo,
		fileStorageService:    fileStorageService,
		fileValidationService: fileValidationService,
	}
}

// UploadAttachments processes batch upload with individual error handling
func (s *AttachmentService) UploadAttachments(ctx context.Context, taskID uuid.UUID, files []*multipart.FileHeader, userID *uuid.UUID) (*models.UploadResponse, error) {
	// Verify task exists and user has access
	if err := s.verifyTaskAccess(ctx, taskID, userID); err != nil {
		return nil, err
	}

	// Get current total size for the task
	currentTotalSize, err := s.attachmentRepo.GetTotalSizeByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current attachment size: %w", err)
	}

	response := &models.UploadResponse{
		Success: []models.TaskAttachment{},
		Failed:  []models.UploadError{},
	}

	// Process each file individually
	for _, fileHeader := range files {
		attachment, err := s.processFileUpload(ctx, taskID, fileHeader, userID, currentTotalSize)
		if err != nil {
			response.Failed = append(response.Failed, models.UploadError{
				Filename: fileHeader.Filename,
				Error:    err.Error(),
			})
			continue
		}

		response.Success = append(response.Success, *attachment)
		currentTotalSize += attachment.FileSize
	}

	// Update response totals
	response.Count = len(response.Success)
	response.TotalSize = currentTotalSize

	return response, nil
}

// processFileUpload handles individual file upload with validation and storage
func (s *AttachmentService) processFileUpload(ctx context.Context, taskID uuid.UUID, fileHeader *multipart.FileHeader, userID *uuid.UUID, currentTotalSize int64) (*models.TaskAttachment, error) {
	// Validate filename
	if err := s.fileValidationService.ValidateFilename(fileHeader.Filename); err != nil {
		return nil, fmt.Errorf("invalid filename: %w", err)
	}

	// Validate file
	validationResult, err := s.fileValidationService.ValidateFile(fileHeader)
	if err != nil {
		return nil, fmt.Errorf("file validation failed: %w", err)
	}

	if !validationResult.IsValid {
		return nil, fmt.Errorf("file validation failed: %s", validationResult.ErrorMessage)
	}

	// Check total size limit
	if err := s.fileValidationService.ValidateTotalSize(currentTotalSize, validationResult.FileSize); err != nil {
		return nil, err
	}

	// Generate secure file path
	filePath, storedFilename, err := s.fileStorageService.GenerateSecureFilePath(fileHeader.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to generate file path: %w", err)
	}

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Store the file
	if err := s.fileStorageService.StoreFile(file, filePath); err != nil {
		return nil, fmt.Errorf("failed to store file: %w", err)
	}

	// Prepare attachment record
	createReq := models.CreateAttachmentRequest{
		TaskID:           taskID,
		OriginalFilename: fileHeader.Filename,
		StoredFilename:   storedFilename,
		FilePath:         filePath,
		FileSize:         validationResult.FileSize,
		MimeType:         validationResult.MimeType,
		UploadedBy:       userID,
		HasThumbnail:     false,
		ThumbnailPath:    nil,
	}

	// Generate thumbnail for images
	if s.fileValidationService.IsImageFile(validationResult.MimeType) {
		thumbnailPath, err := s.generateThumbnail(filePath)
		if err != nil {
			// Log error but don't fail the upload
			fmt.Printf("Warning: failed to generate thumbnail for %s: %v\n", fileHeader.Filename, err)
		} else {
			createReq.HasThumbnail = true
			createReq.ThumbnailPath = &thumbnailPath
		}
	}

	// Create database record
	attachment, err := s.attachmentRepo.Create(ctx, createReq)
	if err != nil {
		// Clean up the stored file if database operation fails
		s.fileStorageService.DeleteFile(filePath)
		if createReq.ThumbnailPath != nil {
			s.fileStorageService.DeleteFile(*createReq.ThumbnailPath)
		}
		return nil, fmt.Errorf("failed to create attachment record: %w", err)
	}

	return attachment, nil
}

// generateThumbnail creates a thumbnail for image files
func (s *AttachmentService) generateThumbnail(originalPath string) (string, error) {
	// Generate thumbnail path
	thumbnailPath, _, err := s.fileStorageService.GenerateThumbnailPath(originalPath)
	if err != nil {
		return "", fmt.Errorf("failed to generate thumbnail path: %w", err)
	}

	// Ensure thumbnail directory exists
	if err := s.fileStorageService.EnsureDirectoryExists(thumbnailPath); err != nil {
		return "", fmt.Errorf("failed to create thumbnail directory: %w", err)
	}

	// Open original image
	originalFile, err := os.Open(originalPath)
	if err != nil {
		return "", fmt.Errorf("failed to open original image: %w", err)
	}
	defer originalFile.Close()

	// Decode image based on file extension
	var img image.Image
	ext := strings.ToLower(filepath.Ext(originalPath))
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(originalFile)
	case ".png":
		img, err = png.Decode(originalFile)
	default:
		return "", fmt.Errorf("unsupported image format: %s", ext)
	}

	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize image to thumbnail size (200x200 with aspect ratio preserved)
	thumbnail := resize.Thumbnail(200, 200, img, resize.Lanczos3)

	// Create thumbnail file
	thumbnailFile, err := os.Create(thumbnailPath)
	if err != nil {
		return "", fmt.Errorf("failed to create thumbnail file: %w", err)
	}
	defer thumbnailFile.Close()

	// Encode thumbnail as JPEG
	if err := jpeg.Encode(thumbnailFile, thumbnail, &jpeg.Options{Quality: 85}); err != nil {
		// Clean up the file if encoding fails
		os.Remove(thumbnailPath)
		return "", fmt.Errorf("failed to encode thumbnail: %w", err)
	}

	return thumbnailPath, nil
}

// GetAttachmentsByTaskID retrieves all attachments for a task
func (s *AttachmentService) GetAttachmentsByTaskID(ctx context.Context, taskID uuid.UUID, userID *uuid.UUID) (*models.AttachmentResponse, error) {
	// Verify task access
	if err := s.verifyTaskAccess(ctx, taskID, userID); err != nil {
		return nil, err
	}

	// Get attachments
	attachments, err := s.attachmentRepo.GetByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachments: %w", err)
	}

	// Get totals
	totalSize, err := s.attachmentRepo.GetTotalSizeByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get total size: %w", err)
	}

	count, err := s.attachmentRepo.GetCountByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get count: %w", err)
	}

	return &models.AttachmentResponse{
		Attachments: attachments,
		TotalSize:   totalSize,
		Count:       count,
	}, nil
}

// GetAttachmentByID retrieves a single attachment with access control
func (s *AttachmentService) GetAttachmentByID(ctx context.Context, attachmentID uuid.UUID, userID *uuid.UUID) (*models.TaskAttachment, error) {
	// Get attachment with uploader info
	attachment, err := s.attachmentRepo.GetByIDWithUploader(ctx, attachmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get attachment: %w", err)
	}

	if attachment == nil {
		return nil, fmt.Errorf("attachment not found")
	}

	// Verify task access
	if err := s.verifyTaskAccess(ctx, attachment.TaskID, userID); err != nil {
		return nil, err
	}

	return attachment, nil
}

// DeleteAttachment removes an attachment with proper cleanup
func (s *AttachmentService) DeleteAttachment(ctx context.Context, attachmentID uuid.UUID, userID *uuid.UUID) error {
	// Check if user can delete this attachment
	canDelete, err := s.CanDeleteAttachment(ctx, attachmentID, userID)
	if err != nil {
		return err
	}

	if !canDelete {
		return fmt.Errorf("access denied: insufficient permissions to delete attachment")
	}

	// Get attachment for cleanup
	attachment, err := s.attachmentRepo.GetAttachmentToCleanup(ctx, attachmentID)
	if err != nil {
		return fmt.Errorf("failed to get attachment for cleanup: %w", err)
	}

	if attachment == nil {
		return fmt.Errorf("attachment not found")
	}

	// Delete from database first
	if err := s.attachmentRepo.Delete(ctx, attachmentID); err != nil {
		return fmt.Errorf("failed to delete attachment record: %w", err)
	}

	// Clean up files
	var thumbnailPath string
	if attachment.ThumbnailPath != nil {
		thumbnailPath = *attachment.ThumbnailPath
	}

	if err := s.fileStorageService.DeleteFileWithThumbnail(attachment.FilePath, thumbnailPath); err != nil {
		// Log error but don't fail the operation since database record is already deleted
		fmt.Printf("Warning: failed to delete attachment files: %v\n", err)
	}

	// Clean up empty directories
	if err := s.fileStorageService.CleanupEmptyDirectories(attachment.FilePath); err != nil {
		// Log error but don't fail the operation
		fmt.Printf("Warning: failed to cleanup empty directories: %v\n", err)
	}

	return nil
}

// verifyTaskAccess checks if user has access to the task
func (s *AttachmentService) verifyTaskAccess(ctx context.Context, taskID uuid.UUID, userID *uuid.UUID) error {
	hasAccess, err := s.HasTaskAccess(ctx, taskID, userID)
	if err != nil {
		return err
	}

	if !hasAccess {
		return fmt.Errorf("access denied: insufficient permissions")
	}

	return nil
}

// HasTaskAccess checks if user has access to the task (public method for external use)
func (s *AttachmentService) HasTaskAccess(ctx context.Context, taskID uuid.UUID, userID *uuid.UUID) (bool, error) {
	// Get task to verify it exists and get project ID
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return false, fmt.Errorf("failed to get task: %w", err)
	}

	if task == nil {
		return false, fmt.Errorf("task not found")
	}

	return s.HasProjectAccess(ctx, task.ProjectID, userID, task)
}

// HasProjectAccess checks if user has access to the project
func (s *AttachmentService) HasProjectAccess(ctx context.Context, projectID uuid.UUID, userID *uuid.UUID, task *models.Task) (bool, error) {
	// Get project to check access
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return false, fmt.Errorf("failed to get project: %w", err)
	}

	if project == nil {
		return false, fmt.Errorf("project not found")
	}

	// Check access permissions
	// We allow access if:
	// 1. Project is public, OR
	// 2. User is the project creator, OR
	// 3. User is assigned to the task (if task provided), OR
	// 4. User is the task author (if task provided)
	if project.IsPublic {
		return true, nil
	}

	if userID == nil {
		return false, nil // No user ID provided for private project
	}

	// Check if user is project creator
	if project.CreatedBy != nil && *project.CreatedBy == *userID {
		return true, nil
	}

	// Check task-specific permissions if task is provided
	if task != nil {
		// Check if user is task assignee
		if task.AssigneeID != nil && *task.AssigneeID == *userID {
			return true, nil
		}

		// Check if user is task author
		if task.AuthorID != nil && *task.AuthorID == *userID {
			return true, nil
		}
	}

	return false, nil
}

// CanDeleteAttachment checks if user can delete a specific attachment
func (s *AttachmentService) CanDeleteAttachment(ctx context.Context, attachmentID uuid.UUID, userID *uuid.UUID) (bool, error) {
	// Get attachment
	attachment, err := s.attachmentRepo.GetByID(ctx, attachmentID)
	if err != nil {
		return false, fmt.Errorf("failed to get attachment: %w", err)
	}

	if attachment == nil {
		return false, fmt.Errorf("attachment not found")
	}

	// Check basic task access first
	hasAccess, err := s.HasTaskAccess(ctx, attachment.TaskID, userID)
	if err != nil {
		return false, err
	}

	if !hasAccess {
		return false, nil
	}

	// Additional deletion permissions:
	// 1. User uploaded the attachment, OR
	// 2. User has project-level access (creator)
	if userID != nil {
		// Check if user uploaded the attachment
		if attachment.UploadedBy != nil && *attachment.UploadedBy == *userID {
			return true, nil
		}

		// Check if user has project-level access
		task, err := s.taskRepo.GetByID(ctx, attachment.TaskID)
		if err != nil {
			return false, fmt.Errorf("failed to get task: %w", err)
		}

		if task != nil {
			project, err := s.projectRepo.GetByID(ctx, task.ProjectID)
			if err != nil {
				return false, fmt.Errorf("failed to get project: %w", err)
			}

			if project != nil && project.CreatedBy != nil && *project.CreatedBy == *userID {
				return true, nil
			}
		}
	}

	return false, nil
}

// VerifyAttachmentAccess verifies user can access a specific attachment
func (s *AttachmentService) VerifyAttachmentAccess(ctx context.Context, attachmentID uuid.UUID, userID *uuid.UUID) error {
	// Get attachment
	attachment, err := s.attachmentRepo.GetByID(ctx, attachmentID)
	if err != nil {
		return fmt.Errorf("failed to get attachment: %w", err)
	}

	if attachment == nil {
		return fmt.Errorf("attachment not found")
	}

	// Verify task access
	return s.verifyTaskAccess(ctx, attachment.TaskID, userID)
}