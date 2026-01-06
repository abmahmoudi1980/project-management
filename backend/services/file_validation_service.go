package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"project-management/config"
	"strings"
)

type FileValidationService struct {
	config *config.FileStorageConfig
}

type ValidationResult struct {
	IsValid      bool
	ErrorMessage string
	MimeType     string
	FileSize     int64
}

func NewFileValidationService() *FileValidationService {
	return &FileValidationService{
		config: config.FileStorage,
	}
}

// ValidateFile performs comprehensive validation on an uploaded file
func (s *FileValidationService) ValidateFile(fileHeader *multipart.FileHeader) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:  true,
		FileSize: fileHeader.Size,
	}

	// Open the file to read its content
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Validate file size
	if err := s.validateFileSize(fileHeader.Size, result); err != nil {
		return result, nil
	}

	// Validate file extension
	if err := s.validateFileExtension(fileHeader.Filename, result); err != nil {
		return result, nil
	}

	// Detect and validate MIME type
	if err := s.validateMimeType(file, result); err != nil {
		return result, nil
	}

	// Basic malicious content detection
	if err := s.detectMaliciousContent(file, fileHeader.Filename, result); err != nil {
		return result, nil
	}

	return result, nil
}

// ValidateTotalSize checks if adding a new file would exceed the total size limit
func (s *FileValidationService) ValidateTotalSize(currentTotalSize, newFileSize int64) error {
	if currentTotalSize+newFileSize > s.config.MaxTotalSize {
		return fmt.Errorf("total attachment size would exceed limit of %d bytes", s.config.MaxTotalSize)
	}
	return nil
}

// validateFileSize checks individual file size limits
func (s *FileValidationService) validateFileSize(fileSize int64, result *ValidationResult) error {
	if fileSize > s.config.MaxFileSize {
		result.IsValid = false
		result.ErrorMessage = fmt.Sprintf("file size %d bytes exceeds maximum allowed size of %d bytes", fileSize, s.config.MaxFileSize)
		return fmt.Errorf("file too large")
	}

	if fileSize == 0 {
		result.IsValid = false
		result.ErrorMessage = "file is empty"
		return fmt.Errorf("empty file")
	}

	return nil
}

// validateFileExtension checks if the file extension is allowed
func (s *FileValidationService) validateFileExtension(filename string, result *ValidationResult) error {
	ext := strings.ToLower(filepath.Ext(filename))

	// Check if extension is blocked
	for _, blockedExt := range s.config.BlockedExtensions {
		if ext == strings.ToLower(blockedExt) {
			result.IsValid = false
			result.ErrorMessage = fmt.Sprintf("file extension %s is not allowed for security reasons", ext)
			return fmt.Errorf("blocked extension")
		}
	}

	// Check if extension is in allowed list
	allowed := false
	for _, allowedExt := range s.config.AllowedExtensions {
		if ext == strings.ToLower(allowedExt) {
			allowed = true
			break
		}
	}

	if !allowed {
		result.IsValid = false
		result.ErrorMessage = fmt.Sprintf("file extension %s is not allowed", ext)
		return fmt.Errorf("extension not allowed")
	}

	return nil
}

// validateMimeType detects and validates the MIME type of the file
func (s *FileValidationService) validateMimeType(file multipart.File, result *ValidationResult) error {
	// Read first 512 bytes to detect MIME type
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		result.IsValid = false
		result.ErrorMessage = "failed to read file content for MIME type detection"
		return fmt.Errorf("failed to read file")
	}

	// Reset file pointer to beginning
	if _, err := file.Seek(0, 0); err != nil {
		result.IsValid = false
		result.ErrorMessage = "failed to reset file pointer"
		return fmt.Errorf("failed to reset file")
	}

	// Detect MIME type
	mimeType := http.DetectContentType(buffer[:n])
	result.MimeType = mimeType

	// Check if MIME type is allowed
	allowed := false
	for _, allowedType := range s.config.AllowedTypes {
		if mimeType == allowedType {
			allowed = true
			break
		}
	}

	// Special handling for some MIME types that might be detected differently
	if !allowed {
		// Handle cases where DetectContentType might return generic types
		if strings.HasPrefix(mimeType, "application/") {
			// For some document types, we might get generic application/octet-stream
			// In this case, we rely on extension validation
			if mimeType == "application/octet-stream" {
				// Extension validation already passed, so we allow it
				allowed = true
				// But we should try to get a more specific MIME type based on extension
				result.MimeType = s.getMimeTypeFromExtension(result.MimeType)
			}
		}
	}

	if !allowed {
		result.IsValid = false
		result.ErrorMessage = fmt.Sprintf("file type %s is not allowed", mimeType)
		return fmt.Errorf("MIME type not allowed")
	}

	return nil
}

// getMimeTypeFromExtension returns a more specific MIME type based on file extension
func (s *FileValidationService) getMimeTypeFromExtension(detectedType string) string {
	// This is a fallback when DetectContentType returns generic types
	// We keep the detected type as is, but this could be enhanced
	// to map extensions to specific MIME types if needed
	return detectedType
}

// detectMaliciousContent performs basic checks for potentially malicious content
func (s *FileValidationService) detectMaliciousContent(file multipart.File, filename string, result *ValidationResult) error {
	// Reset file pointer to beginning
	if _, err := file.Seek(0, 0); err != nil {
		result.IsValid = false
		result.ErrorMessage = "failed to reset file pointer for content scanning"
		return fmt.Errorf("failed to reset file")
	}

	// Read first 1KB for basic content analysis
	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		result.IsValid = false
		result.ErrorMessage = "failed to read file content for security scanning"
		return fmt.Errorf("failed to read file")
	}

	content := string(buffer[:n])

	// Check for suspicious patterns
	suspiciousPatterns := []string{
		"<script",
		"javascript:",
		"vbscript:",
		"onload=",
		"onerror=",
		"eval(",
		"document.write",
		"<iframe",
		"<object",
		"<embed",
	}

	contentLower := strings.ToLower(content)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(contentLower, pattern) {
			result.IsValid = false
			result.ErrorMessage = "file contains potentially malicious content"
			return fmt.Errorf("malicious content detected")
		}
	}

	// Check for executable file signatures (magic numbers)
	if len(buffer) >= 2 {
		// Check for PE executable (Windows .exe)
		if buffer[0] == 0x4D && buffer[1] == 0x5A { // "MZ"
			result.IsValid = false
			result.ErrorMessage = "executable files are not allowed"
			return fmt.Errorf("executable file detected")
		}
	}

	if len(buffer) >= 4 {
		// Check for ELF executable (Linux)
		if buffer[0] == 0x7F && buffer[1] == 0x45 && buffer[2] == 0x4C && buffer[3] == 0x46 { // "\x7FELF"
			result.IsValid = false
			result.ErrorMessage = "executable files are not allowed"
			return fmt.Errorf("executable file detected")
		}
	}

	// Reset file pointer to beginning for subsequent operations
	if _, err := file.Seek(0, 0); err != nil {
		result.IsValid = false
		result.ErrorMessage = "failed to reset file pointer after scanning"
		return fmt.Errorf("failed to reset file")
	}

	return nil
}

// IsImageFile checks if the file is an image based on MIME type
func (s *FileValidationService) IsImageFile(mimeType string) bool {
	imageTypes := []string{
		"image/jpeg",
		"image/png",
		"image/gif",
	}

	for _, imageType := range imageTypes {
		if mimeType == imageType {
			return true
		}
	}

	return false
}

// GetHumanReadableSize converts bytes to human-readable format
func (s *FileValidationService) GetHumanReadableSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// ValidateFilename checks if filename contains invalid characters
func (s *FileValidationService) ValidateFilename(filename string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	// Check for invalid characters
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "\x00"}
	for _, char := range invalidChars {
		if strings.Contains(filename, char) {
			return fmt.Errorf("filename contains invalid character: %s", char)
		}
	}

	// Check for reserved names (Windows)
	reservedNames := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}

	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	for _, reserved := range reservedNames {
		if strings.EqualFold(nameWithoutExt, reserved) {
			return fmt.Errorf("filename uses reserved name: %s", reserved)
		}
	}

	return nil
}