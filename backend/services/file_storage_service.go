package services

import (
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"project-management/config"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FileStorageService struct {
	config *config.FileStorageConfig
}

func NewFileStorageService() *FileStorageService {
	return &FileStorageService{
		config: config.FileStorage,
	}
}

// GenerateSecureFilePath creates a secure file path using UUID and date structure
// Returns: filePath, storedFilename, error
func (s *FileStorageService) GenerateSecureFilePath(originalFilename string) (string, string, error) {
	// Generate UUID for secure filename
	fileUUID := uuid.New()
	
	// Get file extension
	ext := filepath.Ext(originalFilename)
	
	// Create stored filename with UUID
	storedFilename := fmt.Sprintf("%s%s", fileUUID.String(), ext)
	
	// Create date-based directory structure
	now := time.Now()
	dateDir := filepath.Join(
		fmt.Sprintf("%04d", now.Year()),
		fmt.Sprintf("%02d", now.Month()),
		fmt.Sprintf("%02d", now.Day()),
	)
	
	// Create full path
	attachmentDir := filepath.Join(s.config.BaseUploadPath, "attachments", dateDir)
	filePath := filepath.Join(attachmentDir, storedFilename)
	
	return filePath, storedFilename, nil
}

// GenerateThumbnailPath creates a path for thumbnail storage
func (s *FileStorageService) GenerateThumbnailPath(originalFilePath string) (string, string, error) {
	// Get the directory structure from original path
	relPath, err := filepath.Rel(filepath.Join(s.config.BaseUploadPath, "attachments"), originalFilePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to get relative path: %w", err)
	}
	
	// Extract directory and filename
	dir := filepath.Dir(relPath)
	filename := filepath.Base(relPath)
	
	// Remove extension and add thumbnail suffix
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	thumbnailFilename := fmt.Sprintf("%s-thumb.jpg", nameWithoutExt)
	
	// Create thumbnail path
	thumbnailDir := filepath.Join(s.config.BaseUploadPath, "thumbnails", dir)
	thumbnailPath := filepath.Join(thumbnailDir, thumbnailFilename)
	
	return thumbnailPath, thumbnailFilename, nil
}

// EnsureDirectoryExists creates directory structure if it doesn't exist
func (s *FileStorageService) EnsureDirectoryExists(filePath string) error {
	dir := filepath.Dir(filePath)
	
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create directory with proper permissions
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	
	return nil
}

// StoreFile saves the file to the specified path
func (s *FileStorageService) StoreFile(src io.Reader, filePath string) error {
	// Ensure directory exists
	if err := s.EnsureDirectoryExists(filePath); err != nil {
		return err
	}
	
	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer dst.Close()
	
	// Copy data from source to destination
	_, err = io.Copy(dst, src)
	if err != nil {
		// Clean up the file if copy failed
		os.Remove(filePath)
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}
	
	return nil
}

// DeleteFile removes a file from storage
func (s *FileStorageService) DeleteFile(filePath string) error {
	if filePath == "" {
		return nil // Nothing to delete
	}
	
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File doesn't exist, consider it deleted
	}
	
	// Remove the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file %s: %w", filePath, err)
	}
	
	return nil
}

// DeleteFileWithThumbnail removes both the main file and its thumbnail
func (s *FileStorageService) DeleteFileWithThumbnail(filePath, thumbnailPath string) error {
	// Delete main file
	if err := s.DeleteFile(filePath); err != nil {
		return fmt.Errorf("failed to delete main file: %w", err)
	}
	
	// Delete thumbnail if it exists
	if thumbnailPath != "" {
		if err := s.DeleteFile(thumbnailPath); err != nil {
			// Log error but don't fail the operation
			// Thumbnail deletion is not critical
			fmt.Printf("Warning: failed to delete thumbnail %s: %v\n", thumbnailPath, err)
		}
	}
	
	return nil
}

// CleanupEmptyDirectories removes empty directories in the upload path
func (s *FileStorageService) CleanupEmptyDirectories(filePath string) error {
	dir := filepath.Dir(filePath)
	
	// Don't remove the base upload directories
	baseAttachments := filepath.Join(s.config.BaseUploadPath, "attachments")
	baseThumbnails := filepath.Join(s.config.BaseUploadPath, "thumbnails")
	
	if dir == baseAttachments || dir == baseThumbnails || dir == s.config.BaseUploadPath {
		return nil
	}
	
	// Check if directory is empty
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil // Directory might not exist, ignore
	}
	
	if len(entries) == 0 {
		// Directory is empty, remove it
		if err := os.Remove(dir); err != nil {
			return fmt.Errorf("failed to remove empty directory %s: %w", dir, err)
		}
		
		// Recursively check parent directory
		return s.CleanupEmptyDirectories(dir)
	}
	
	return nil
}

// GetFileSize returns the size of a file
func (s *FileStorageService) GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}
	
	return info.Size(), nil
}

// FileExists checks if a file exists at the given path
func (s *FileStorageService) FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// InitializeStorageDirectories creates the base directory structure
func (s *FileStorageService) InitializeStorageDirectories() error {
	directories := []string{
		filepath.Join(s.config.BaseUploadPath, "attachments"),
		filepath.Join(s.config.BaseUploadPath, "thumbnails"),
	}
	
	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}
	
	return nil
}

// GenerateRandomBytes generates cryptographically secure random bytes
func (s *FileStorageService) GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return bytes, nil
}