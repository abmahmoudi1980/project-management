package config

import (
	"os"
	"path/filepath"
	"strconv"
)

type FileStorageConfig struct {
	BaseUploadPath    string
	MaxFileSize       int64
	MaxTotalSize      int64
	AllowedTypes      []string
	AllowedExtensions []string
	BlockedExtensions []string
	ThumbnailSize     int
}

var FileStorage *FileStorageConfig

func InitFileStorage() error {
	baseUploadPath := os.Getenv("UPLOAD_PATH")
	if baseUploadPath == "" {
		baseUploadPath = "./uploads"
	}

	// Ensure the base upload path is absolute
	absPath, err := filepath.Abs(baseUploadPath)
	if err != nil {
		return err
	}

	maxFileSize := int64(10 * 1024 * 1024) // 10MB default
	if envSize := os.Getenv("MAX_FILE_SIZE"); envSize != "" {
		if size, err := strconv.ParseInt(envSize, 10, 64); err == nil {
			maxFileSize = size
		}
	}

	maxTotalSize := int64(100 * 1024 * 1024) // 100MB default
	if envSize := os.Getenv("MAX_TOTAL_SIZE"); envSize != "" {
		if size, err := strconv.ParseInt(envSize, 10, 64); err == nil {
			maxTotalSize = size
		}
	}

	thumbnailSize := 200 // 200px default
	if envSize := os.Getenv("THUMBNAIL_SIZE"); envSize != "" {
		if size, err := strconv.Atoi(envSize); err == nil {
			thumbnailSize = size
		}
	}

	FileStorage = &FileStorageConfig{
		BaseUploadPath: absPath,
		MaxFileSize:    maxFileSize,
		MaxTotalSize:   maxTotalSize,
		AllowedTypes: []string{
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.ms-excel",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"application/vnd.ms-powerpoint",
			"application/vnd.openxmlformats-officedocument.presentationml.presentation",
			"text/plain",
			"image/jpeg",
			"image/png",
			"image/gif",
			"application/zip",
		},
		AllowedExtensions: []string{
			".pdf", ".doc", ".docx", ".xls", ".xlsx",
			".ppt", ".pptx", ".txt", ".jpg", ".jpeg",
			".png", ".gif", ".zip",
		},
		BlockedExtensions: []string{
			".exe", ".bat", ".cmd", ".scr", ".pif",
			".com", ".js", ".vbs", ".jar",
		},
		ThumbnailSize: thumbnailSize,
	}

	return nil
}