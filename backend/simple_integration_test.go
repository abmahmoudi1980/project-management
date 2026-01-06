package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"os"
	"project-management/config"
	"project-management/services"
	"testing"

	"github.com/google/uuid"
)

// TestAttachmentEndToEnd tests the complete attachment workflow
func TestAttachmentEndToEnd(t *testing.T) {
	fmt.Println("=== End-to-End Attachment Integration Test ===")

	// Set up test environment
	os.Setenv("UPLOAD_PATH", "./test_uploads_e2e")
	os.Setenv("MAX_FILE_SIZE", "10485760")    // 10MB
	os.Setenv("MAX_TOTAL_SIZE", "104857600")  // 100MB

	err := config.InitFileStorage()
	if err != nil {
		t.Fatalf("Failed to initialize file storage: %v", err)
	}

	testUploadPath := "./test_uploads_e2e"
	err = os.MkdirAll(testUploadPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testUploadPath)

	// Initialize services
	fileStorageService := services.NewFileStorageService()
	fileValidationService := services.NewFileValidationService()

	// Test 1: Complete upload-to-download workflow
	t.Run("Complete Upload-Download Workflow", func(t *testing.T) {
		fmt.Println("\n1. Testing complete upload-download workflow...")

		testFiles := []struct {
			name    string
			content []byte
		}{
			{"document.txt", []byte("This is a test document with some content.")},
			{"data.pdf", []byte("%PDF-1.4\nTest PDF content for integration testing.")},
			{"image.png", createTestPNGImageE2E()},
		}

		var uploadedFiles []UploadedFileInfo

		// Upload each file
		for _, testFile := range testFiles {
			fmt.Printf("  Uploading %s...\n", testFile.name)

			// Create multipart file header simulation
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("files", testFile.name)
			if err != nil {
				t.Fatalf("Failed to create form file: %v", err)
			}
			_, err = part.Write(testFile.content)
			if err != nil {
				t.Fatalf("Failed to write content: %v", err)
			}
			writer.Close()

			// Parse the multipart form to get file header
			req := &MockRequest{
				Body:        body,
				ContentType: writer.FormDataContentType(),
			}
			form, err := req.ParseMultipartForm()
			if err != nil {
				t.Fatalf("Failed to parse multipart form: %v", err)
			}

			fileHeader := form.File["files"][0]

			// Validate file
			validationResult, err := fileValidationService.ValidateFile(fileHeader)
			if err != nil {
				t.Fatalf("Validation failed: %v", err)
			}

			if !validationResult.IsValid {
				t.Fatalf("File validation failed: %s", validationResult.ErrorMessage)
			}

			// Generate secure file path
			filePath, storedFilename, err := fileStorageService.GenerateSecureFilePath(testFile.name)
			if err != nil {
				t.Fatalf("Failed to generate file path: %v", err)
			}

			// Store file
			file, err := fileHeader.Open()
			if err != nil {
				t.Fatalf("Failed to open file: %v", err)
			}

			err = fileStorageService.StoreFile(file, filePath)
			file.Close()
			if err != nil {
				t.Fatalf("Failed to store file: %v", err)
			}

			// Verify file was stored correctly
			if !fileStorageService.FileExists(filePath) {
				t.Fatalf("File was not stored: %s", filePath)
			}

			uploadedFiles = append(uploadedFiles, UploadedFileInfo{
				ID:               uuid.New(),
				OriginalFilename: testFile.name,
				StoredFilename:   storedFilename,
				FilePath:         filePath,
				FileSize:         int64(len(testFile.content)),
				MimeType:         validationResult.MimeType,
				OriginalContent:  testFile.content,
			})

			fmt.Printf("  ✓ %s uploaded successfully\n", testFile.name)
		}

		// Verify all files can be read back
		fmt.Println("  Verifying file downloads...")
		for _, uploadedFile := range uploadedFiles {
			// Read file back
			storedContent, err := os.ReadFile(uploadedFile.FilePath)
			if err != nil {
				t.Fatalf("Failed to read stored file %s: %v", uploadedFile.OriginalFilename, err)
			}

			// Verify content matches
			if !bytes.Equal(uploadedFile.OriginalContent, storedContent) {
				t.Fatalf("Content mismatch for file %s", uploadedFile.OriginalFilename)
			}

			fmt.Printf("  ✓ %s download verified\n", uploadedFile.OriginalFilename)
		}

		// Clean up files
		fmt.Println("  Cleaning up files...")
		for _, uploadedFile := range uploadedFiles {
			err := fileStorageService.DeleteFile(uploadedFile.FilePath)
			if err != nil {
				t.Fatalf("Failed to delete file %s: %v", uploadedFile.OriginalFilename, err)
			}

			// Verify file is deleted
			if fileStorageService.FileExists(uploadedFile.FilePath) {
				t.Fatalf("File was not deleted: %s", uploadedFile.FilePath)
			}

			fmt.Printf("  ✓ %s deleted successfully\n", uploadedFile.OriginalFilename)
		}

		fmt.Println("✓ Complete upload-download workflow test passed")
	})

	// Test 2: Error scenarios
	t.Run("Error Scenarios", func(t *testing.T) {
		fmt.Println("\n2. Testing error scenarios...")

		errorTestCases := []struct {
			name        string
			filename    string
			content     []byte
			expectError bool
			errorType   string
		}{
			{
				name:        "File Too Large",
				filename:    "huge.txt",
				content:     bytes.Repeat([]byte("A"), 15*1024*1024), // 15MB
				expectError: true,
				errorType:   "size",
			},
			{
				name:        "Blocked Extension",
				filename:    "virus.exe",
				content:     []byte("MZ\x90\x00"), // PE header
				expectError: true,
				errorType:   "extension",
			},
			{
				name:        "Invalid Extension",
				filename:    "unknown.xyz",
				content:     []byte("Unknown file type"),
				expectError: true,
				errorType:   "extension",
			},
			{
				name:        "Empty File",
				filename:    "empty.txt",
				content:     []byte{},
				expectError: true,
				errorType:   "size",
			},
		}

		for _, testCase := range errorTestCases {
			t.Run(testCase.name, func(t *testing.T) {
				fmt.Printf("  Testing %s...\n", testCase.name)

				// Create multipart form
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("files", testCase.filename)
				if err != nil {
					t.Fatalf("Failed to create form file: %v", err)
				}
				_, err = part.Write(testCase.content)
				if err != nil {
					t.Fatalf("Failed to write content: %v", err)
				}
				writer.Close()

				// Parse form
				req := &MockRequest{
					Body:        body,
					ContentType: writer.FormDataContentType(),
				}
				form, err := req.ParseMultipartForm()
				if err != nil {
					t.Fatalf("Failed to parse multipart form: %v", err)
				}

				fileHeader := form.File["files"][0]

				// Validate file
				validationResult, err := fileValidationService.ValidateFile(fileHeader)

				if testCase.expectError {
					if err == nil && validationResult.IsValid {
						t.Errorf("Expected validation to fail for %s, but it passed", testCase.name)
					} else {
						fmt.Printf("  ✓ %s correctly rejected\n", testCase.name)
					}
				} else {
					if err != nil || !validationResult.IsValid {
						t.Errorf("Expected validation to pass for %s, but it failed: %v", testCase.name, err)
					}
				}
			})
		}

		fmt.Println("✓ Error scenarios test passed")
	})

	// Test 3: Various file types and sizes
	t.Run("File Type and Size Variations", func(t *testing.T) {
		fmt.Println("\n3. Testing various file types and sizes...")

		variations := []struct {
			name     string
			filename string
			content  []byte
			valid    bool
		}{
			{
				name:     "Small Text",
				filename: "small.txt",
				content:  []byte("Small"),
				valid:    true,
			},
			{
				name:     "Medium Document",
				filename: "medium.doc",
				content:  bytes.Repeat([]byte("Medium content "), 1000),
				valid:    true,
			},
			{
				name:     "Large Valid File",
				filename: "large.txt",
				content:  bytes.Repeat([]byte("Large content "), 50000), // ~750KB
				valid:    true,
			},
			{
				name:     "ZIP Archive",
				filename: "archive.zip",
				content:  []byte("PK\x03\x04\x14\x00\x00\x00\x08\x00"),
				valid:    true,
			},
			{
				name:     "Excel File",
				filename: "spreadsheet.xlsx",
				content:  []byte("PK\x03\x04\x14\x00\x06\x00"), // XLSX header
				valid:    true,
			},
		}

		for _, variation := range variations {
			t.Run(variation.name, func(t *testing.T) {
				fmt.Printf("  Testing %s (%s)...\n", variation.name, variation.filename)

				// Create and validate file
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, err := writer.CreateFormFile("files", variation.filename)
				if err != nil {
					t.Fatalf("Failed to create form file: %v", err)
				}
				_, err = part.Write(variation.content)
				if err != nil {
					t.Fatalf("Failed to write content: %v", err)
				}
				writer.Close()

				req := &MockRequest{
					Body:        body,
					ContentType: writer.FormDataContentType(),
				}
				form, err := req.ParseMultipartForm()
				if err != nil {
					t.Fatalf("Failed to parse multipart form: %v", err)
				}

				fileHeader := form.File["files"][0]
				validationResult, err := fileValidationService.ValidateFile(fileHeader)

				if variation.valid {
					if err != nil || !validationResult.IsValid {
						t.Errorf("Expected %s to be valid, but validation failed", variation.name)
					} else {
						fmt.Printf("  ✓ %s validated successfully\n", variation.name)
					}
				} else {
					if err == nil && validationResult.IsValid {
						t.Errorf("Expected %s to be invalid, but validation passed", variation.name)
					}
				}
			})
		}

		fmt.Println("✓ File type and size variations test passed")
	})

	fmt.Println("\n=== All End-to-End Integration Tests Completed Successfully ===")
}

// Helper structures
type UploadedFileInfo struct {
	ID               uuid.UUID
	OriginalFilename string
	StoredFilename   string
	FilePath         string
	FileSize         int64
	MimeType         string
	OriginalContent  []byte
}

type MockRequest struct {
	Body        *bytes.Buffer
	ContentType string
}

func (r *MockRequest) ParseMultipartForm() (*multipart.Form, error) {
	// This is a simplified mock - in real implementation this would be more complex
	reader := multipart.NewReader(r.Body, extractBoundary(r.ContentType))
	return reader.ReadForm(32 << 20) // 32MB max memory
}

func extractBoundary(contentType string) string {
	// Extract boundary from content type
	parts := bytes.Split([]byte(contentType), []byte("boundary="))
	if len(parts) < 2 {
		return ""
	}
	return string(parts[1])
}

func createTestPNGImageE2E() []byte {
	// Minimal valid PNG file (1x1 pixel transparent)
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
		0x00, 0x00, 0x00, 0x0D, // IHDR chunk length
		0x49, 0x48, 0x44, 0x52, // IHDR
		0x00, 0x00, 0x00, 0x01, // Width: 1
		0x00, 0x00, 0x00, 0x01, // Height: 1
		0x08, 0x06, 0x00, 0x00, 0x00, // Bit depth, color type, compression, filter, interlace
		0x1F, 0x15, 0xC4, 0x89, // CRC
		0x00, 0x00, 0x00, 0x0A, // IDAT chunk length
		0x49, 0x44, 0x41, 0x54, // IDAT
		0x78, 0x9C, 0x62, 0x00, 0x00, 0x00, 0x02, 0x00, 0x01, // Compressed data
		0xE2, 0x21, 0xBC, 0x33, // CRC
		0x00, 0x00, 0x00, 0x00, // IEND chunk length
		0x49, 0x45, 0x4E, 0x44, // IEND
		0xAE, 0x42, 0x60, 0x82, // CRC
	}
}

// Test concurrent operations
func TestConcurrentAttachmentOperations(t *testing.T) {
	fmt.Println("=== Concurrent Operations Test ===")

	// Set up test environment
	os.Setenv("UPLOAD_PATH", "./test_uploads_concurrent")
	err := config.InitFileStorage()
	if err != nil {
		t.Fatalf("Failed to initialize file storage: %v", err)
	}

	testUploadPath := "./test_uploads_concurrent"
	err = os.MkdirAll(testUploadPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testUploadPath)

	fileStorageService := services.NewFileStorageService()

	// Test concurrent file operations
	numConcurrent := 10
	results := make(chan bool, numConcurrent)

	for i := 0; i < numConcurrent; i++ {
		go func(index int) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Goroutine %d panicked: %v\n", index, r)
					results <- false
					return
				}
			}()

			// Create unique test file
			filename := fmt.Sprintf("concurrent_%d.txt", index)
			content := []byte(fmt.Sprintf("Concurrent test file %d content", index))

			// Generate file path
			filePath, _, err := fileStorageService.GenerateSecureFilePath(filename)
			if err != nil {
				fmt.Printf("Goroutine %d: Failed to generate path: %v\n", index, err)
				results <- false
				return
			}

			// Store file
			err = fileStorageService.StoreFile(bytes.NewReader(content), filePath)
			if err != nil {
				fmt.Printf("Goroutine %d: Failed to store file: %v\n", index, err)
				results <- false
				return
			}

			// Verify file exists
			if !fileStorageService.FileExists(filePath) {
				fmt.Printf("Goroutine %d: File not found after storage\n", index)
				results <- false
				return
			}

			// Read and verify content
			storedContent, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Goroutine %d: Failed to read file: %v\n", index, err)
				results <- false
				return
			}

			if !bytes.Equal(content, storedContent) {
				fmt.Printf("Goroutine %d: Content mismatch\n", index)
				results <- false
				return
			}

			// Clean up
			err = fileStorageService.DeleteFile(filePath)
			if err != nil {
				fmt.Printf("Goroutine %d: Failed to delete file: %v\n", index, err)
				results <- false
				return
			}

			results <- true
		}(i)
	}

	// Wait for all operations to complete
	successCount := 0
	for i := 0; i < numConcurrent; i++ {
		if <-results {
			successCount++
		}
	}

	if successCount != numConcurrent {
		t.Errorf("Expected %d successful operations, got %d", numConcurrent, successCount)
	} else {
		fmt.Printf("✓ All %d concurrent operations completed successfully\n", numConcurrent)
	}

	fmt.Println("=== Concurrent Operations Test Completed ===")
}