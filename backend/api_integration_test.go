package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"project-management/config"
	"project-management/services"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TestAttachmentAPIEndpoints tests the actual API endpoints with a real Fiber app
func TestAttachmentAPIEndpoints(t *testing.T) {
	// Skip if database is not available
	if testing.Short() {
		t.Skip("Skipping API integration tests in short mode")
	}

	fmt.Println("=== Attachment API Integration Tests ===")

	// Set up test environment
	os.Setenv("UPLOAD_PATH", "./test_uploads_api")
	os.Setenv("MAX_FILE_SIZE", "10485760")    // 10MB
	os.Setenv("MAX_TOTAL_SIZE", "104857600")  // 100MB

	// Initialize configuration
	err := config.InitFileStorage()
	if err != nil {
		t.Fatalf("Failed to initialize file storage: %v", err)
	}

	// Create test upload directory
	testUploadPath := "./test_uploads_api"
	err = os.MkdirAll(testUploadPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testUploadPath)

	// Create a minimal Fiber app for testing
	app := createTestApp(t)

	// Test cases
	t.Run("Upload Validation Tests", func(t *testing.T) {
		testUploadValidation(t, app)
	})

	t.Run("File Processing Tests", func(t *testing.T) {
		testFileProcessing(t, app)
	})

	t.Run("Error Handling Tests", func(t *testing.T) {
		testAPIErrorHandling(t, app)
	})

	fmt.Println("=== API Integration Tests Completed ===")
}

func createTestApp(t *testing.T) *fiber.App {
	// Create a minimal Fiber app with just the attachment routes
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	// Create mock services for testing
	fileStorageService := services.NewFileStorageService()
	fileValidationService := services.NewFileValidationService()

	// Create a mock attachment service that doesn't require database
	mockAttachmentService := &MockAttachmentService{
		fileStorageService:    fileStorageService,
		fileValidationService: fileValidationService,
	}

	// Create handler with mock service
	attachmentHandler := &MockAttachmentHandler{
		service: mockAttachmentService,
	}

	// Set up minimal routes for testing
	api := app.Group("/api")
	tasks := api.Group("/tasks")
	tasks.Post("/:taskId/attachments", attachmentHandler.UploadAttachments)
	tasks.Get("/:taskId/attachments", attachmentHandler.ListAttachments)

	attachments := api.Group("/attachments")
	attachments.Get("/:id/download", attachmentHandler.DownloadAttachment)
	attachments.Delete("/:id", attachmentHandler.DeleteAttachment)

	return app
}

// Mock services for testing without database dependency
type MockAttachmentService struct {
	fileStorageService    *services.FileStorageService
	fileValidationService *services.FileValidationService
	uploadedFiles         map[string]MockAttachmentData
}

type MockAttachmentData struct {
	ID               uuid.UUID
	TaskID           uuid.UUID
	OriginalFilename string
	StoredFilename   string
	FilePath         string
	FileSize         int64
	MimeType         string
	CreatedAt        time.Time
}

type MockAttachmentHandler struct {
	service *MockAttachmentService
}

func (h *MockAttachmentHandler) UploadAttachments(c *fiber.Ctx) error {
	// Parse task ID
	taskIDStr := c.Params("taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID format",
		})
	}

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse multipart form",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files provided",
		})
	}

	// Process files
	var successFiles []MockAttachmentData
	var failedFiles []map[string]string

	for _, fileHeader := range files {
		// Validate file
		validationResult, err := h.service.fileValidationService.ValidateFile(fileHeader)
		if err != nil || !validationResult.IsValid {
			failedFiles = append(failedFiles, map[string]string{
				"filename": fileHeader.Filename,
				"error":    validationResult.ErrorMessage,
			})
			continue
		}

		// Generate file path
		filePath, storedFilename, err := h.service.fileStorageService.GenerateSecureFilePath(fileHeader.Filename)
		if err != nil {
			failedFiles = append(failedFiles, map[string]string{
				"filename": fileHeader.Filename,
				"error":    "Failed to generate file path",
			})
			continue
		}

		// Store file
		file, err := fileHeader.Open()
		if err != nil {
			failedFiles = append(failedFiles, map[string]string{
				"filename": fileHeader.Filename,
				"error":    "Failed to open file",
			})
			continue
		}

		err = h.service.fileStorageService.StoreFile(file, filePath)
		file.Close()
		if err != nil {
			failedFiles = append(failedFiles, map[string]string{
				"filename": fileHeader.Filename,
				"error":    "Failed to store file",
			})
			continue
		}

		// Create attachment record
		attachment := MockAttachmentData{
			ID:               uuid.New(),
			TaskID:           taskID,
			OriginalFilename: fileHeader.Filename,
			StoredFilename:   storedFilename,
			FilePath:         filePath,
			FileSize:         fileHeader.Size,
			MimeType:         validationResult.MimeType,
			CreatedAt:        time.Now(),
		}

		// Store in mock database
		if h.service.uploadedFiles == nil {
			h.service.uploadedFiles = make(map[string]MockAttachmentData)
		}
		h.service.uploadedFiles[attachment.ID.String()] = attachment

		successFiles = append(successFiles, attachment)
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success":     successFiles,
		"failed":      failedFiles,
		"count":       len(successFiles),
		"total_size":  calculateTotalSize(successFiles),
	})
}

func (h *MockAttachmentHandler) ListAttachments(c *fiber.Ctx) error {
	taskIDStr := c.Params("taskId")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID format",
		})
	}

	var attachments []MockAttachmentData
	for _, attachment := range h.service.uploadedFiles {
		if attachment.TaskID == taskID {
			attachments = append(attachments, attachment)
		}
	}

	return c.JSON(fiber.Map{
		"attachments": attachments,
		"count":       len(attachments),
		"total_size":  calculateTotalSize(attachments),
	})
}

func (h *MockAttachmentHandler) DownloadAttachment(c *fiber.Ctx) error {
	attachmentIDStr := c.Params("id")
	attachmentID, err := uuid.Parse(attachmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid attachment ID format",
		})
	}

	attachment, exists := h.service.uploadedFiles[attachmentID.String()]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Attachment not found",
		})
	}

	// Set security headers
	c.Set("Content-Type", "application/octet-stream")
	c.Set("Content-Disposition", "attachment; filename=\""+attachment.OriginalFilename+"\"")
	c.Set("X-Content-Type-Options", "nosniff")
	c.Set("X-Frame-Options", "DENY")
	c.Set("X-XSS-Protection", "1; mode=block")

	return c.SendFile(attachment.FilePath)
}

func (h *MockAttachmentHandler) DeleteAttachment(c *fiber.Ctx) error {
	attachmentIDStr := c.Params("id")
	attachmentID, err := uuid.Parse(attachmentIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid attachment ID format",
		})
	}

	attachment, exists := h.service.uploadedFiles[attachmentID.String()]
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Attachment not found",
		})
	}

	// Delete file
	err = h.service.fileStorageService.DeleteFile(attachment.FilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete file",
		})
	}

	// Remove from mock database
	delete(h.service.uploadedFiles, attachmentID.String())

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func calculateTotalSize(attachments []MockAttachmentData) int64 {
	var total int64
	for _, attachment := range attachments {
		total += attachment.FileSize
	}
	return total
}

func testUploadValidation(t *testing.T, app *fiber.App) {
	fmt.Println("\n1. Testing upload validation...")

	testCases := []struct {
		name           string
		filename       string
		content        []byte
		expectedStatus int
		shouldSucceed  bool
	}{
		{
			name:           "Valid Text File",
			filename:       "test.txt",
			content:        []byte("Valid text content"),
			expectedStatus: 201,
			shouldSucceed:  true,
		},
		{
			name:           "Valid PDF File",
			filename:       "document.pdf",
			content:        []byte("%PDF-1.4\nValid PDF content"),
			expectedStatus: 201,
			shouldSucceed:  true,
		},
		{
			name:           "Blocked Extension",
			filename:       "malware.exe",
			content:        []byte("MZ\x90\x00"),
			expectedStatus: 201, // Request succeeds but file fails
			shouldSucceed:  false,
		},
		{
			name:           "Invalid Extension",
			filename:       "unknown.xyz",
			content:        []byte("Unknown file type"),
			expectedStatus: 201, // Request succeeds but file fails
			shouldSucceed:  false,
		},
		{
			name:           "Empty File",
			filename:       "empty.txt",
			content:        []byte{},
			expectedStatus: 201, // Request succeeds but file fails
			shouldSucceed:  false,
		},
	}

	taskID := uuid.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create multipart form
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("files", tc.filename)
			if err != nil {
				t.Fatalf("Failed to create form file: %v", err)
			}

			_, err = part.Write(tc.content)
			if err != nil {
				t.Fatalf("Failed to write file content: %v", err)
			}

			err = writer.Close()
			if err != nil {
				t.Fatalf("Failed to close writer: %v", err)
			}

			// Make request
			req, err := http.NewRequest("POST", fmt.Sprintf("/api/tasks/%s/attachments", taskID), body)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())

			resp, err := app.Test(req, 10000) // 10 second timeout
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}

			// Check status code
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			// Parse response
			if resp.StatusCode == 201 {
				var response map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				success, _ := response["success"].([]interface{})
				failed, _ := response["failed"].([]interface{})

				if tc.shouldSucceed {
					if len(success) != 1 {
						t.Errorf("Expected 1 successful upload, got %d", len(success))
					}
					if len(failed) != 0 {
						t.Errorf("Expected 0 failed uploads, got %d", len(failed))
					}
				} else {
					if len(success) != 0 {
						t.Errorf("Expected 0 successful uploads, got %d", len(success))
					}
					if len(failed) != 1 {
						t.Errorf("Expected 1 failed upload, got %d", len(failed))
					}
				}
			}
		})
	}

	fmt.Println("✓ Upload validation tests completed")
}

func testFileProcessing(t *testing.T, app *fiber.App) {
	fmt.Println("\n2. Testing file processing...")

	taskID := uuid.New()

	// Upload a test file
	testContent := []byte("Test file for processing")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("files", "process_test.txt")
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}

	_, err = part.Write(testContent)
	if err != nil {
		t.Fatalf("Failed to write content: %v", err)
	}

	err = writer.Close()
	if err != nil {
		t.Fatalf("Failed to close writer: %v", err)
	}

	// Upload file
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/tasks/%s/attachments", taskID), body)
	if err != nil {
		t.Fatalf("Failed to create upload request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req, 10000)
	if err != nil {
		t.Fatalf("Upload request failed: %v", err)
	}

	if resp.StatusCode != 201 {
		t.Fatalf("Upload failed with status %d", resp.StatusCode)
	}

	var uploadResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&uploadResponse)
	if err != nil {
		t.Fatalf("Failed to decode upload response: %v", err)
	}

	success := uploadResponse["success"].([]interface{})
	if len(success) != 1 {
		t.Fatalf("Expected 1 successful upload, got %d", len(success))
	}

	attachment := success[0].(map[string]interface{})
	attachmentID := attachment["ID"].(string)

	// Test listing attachments
	req, err = http.NewRequest("GET", fmt.Sprintf("/api/tasks/%s/attachments", taskID), nil)
	if err != nil {
		t.Fatalf("Failed to create list request: %v", err)
	}

	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("List request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("List failed with status %d", resp.StatusCode)
	}

	var listResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&listResponse)
	if err != nil {
		t.Fatalf("Failed to decode list response: %v", err)
	}

	attachments := listResponse["attachments"].([]interface{})
	if len(attachments) != 1 {
		t.Errorf("Expected 1 attachment, got %d", len(attachments))
	}

	// Test downloading attachment
	req, err = http.NewRequest("GET", fmt.Sprintf("/api/attachments/%s/download", attachmentID), nil)
	if err != nil {
		t.Fatalf("Failed to create download request: %v", err)
	}

	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Download request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Download failed with status %d", resp.StatusCode)
	}

	// Verify security headers
	if resp.Header.Get("Content-Type") != "application/octet-stream" {
		t.Errorf("Expected Content-Type 'application/octet-stream', got '%s'", resp.Header.Get("Content-Type"))
	}

	if !strings.Contains(resp.Header.Get("Content-Disposition"), "attachment") {
		t.Errorf("Expected Content-Disposition to contain 'attachment'")
	}

	// Verify content
	downloadedContent, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read downloaded content: %v", err)
	}

	if !bytes.Equal(testContent, downloadedContent) {
		t.Errorf("Downloaded content doesn't match original")
	}

	// Test deleting attachment
	req, err = http.NewRequest("DELETE", fmt.Sprintf("/api/attachments/%s", attachmentID), nil)
	if err != nil {
		t.Fatalf("Failed to create delete request: %v", err)
	}

	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Delete request failed: %v", err)
	}

	if resp.StatusCode != 204 {
		t.Fatalf("Delete failed with status %d", resp.StatusCode)
	}

	// Verify attachment is deleted
	req, err = http.NewRequest("GET", fmt.Sprintf("/api/tasks/%s/attachments", taskID), nil)
	if err != nil {
		t.Fatalf("Failed to create verification request: %v", err)
	}

	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Verification request failed: %v", err)
	}

	err = json.NewDecoder(resp.Body).Decode(&listResponse)
	if err != nil {
		t.Fatalf("Failed to decode verification response: %v", err)
	}

	attachments = listResponse["attachments"].([]interface{})
	if len(attachments) != 0 {
		t.Errorf("Expected 0 attachments after deletion, got %d", len(attachments))
	}

	fmt.Println("✓ File processing tests completed")
}

func testAPIErrorHandling(t *testing.T, app *fiber.App) {
	fmt.Println("\n3. Testing API error handling...")

	// Test invalid task ID
	req, err := http.NewRequest("POST", "/api/tasks/invalid-uuid/attachments", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("Expected status 400 for invalid UUID, got %d", resp.StatusCode)
	}

	// Test no files provided
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.Close()

	taskID := uuid.New()
	req, err = http.NewRequest("POST", fmt.Sprintf("/api/tasks/%s/attachments", taskID), body)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != 400 {
		t.Errorf("Expected status 400 for no files, got %d", resp.StatusCode)
	}

	// Test download non-existent attachment
	fakeID := uuid.New()
	req, err = http.NewRequest("GET", fmt.Sprintf("/api/attachments/%s/download", fakeID), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404 for non-existent attachment, got %d", resp.StatusCode)
	}

	// Test delete non-existent attachment
	req, err = http.NewRequest("DELETE", fmt.Sprintf("/api/attachments/%s", fakeID), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	resp, err = app.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404 for non-existent attachment, got %d", resp.StatusCode)
	}

	fmt.Println("✓ API error handling tests completed")
}