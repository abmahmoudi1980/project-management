package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskAttachment struct {
	ID               uuid.UUID  `json:"id"`
	TaskID           uuid.UUID  `json:"task_id"`
	OriginalFilename string     `json:"original_filename"`
	StoredFilename   string     `json:"stored_filename"`
	FilePath         string     `json:"file_path"`
	FileSize         int64      `json:"file_size"`
	MimeType         string     `json:"mime_type"`
	UploadedBy       *uuid.UUID `json:"uploaded_by,omitempty"`
	UploaderName     *string    `json:"uploader_name,omitempty"`
	HasThumbnail     bool       `json:"has_thumbnail"`
	ThumbnailPath    *string    `json:"thumbnail_path,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type CreateAttachmentRequest struct {
	TaskID           uuid.UUID `json:"task_id"`
	OriginalFilename string    `json:"original_filename"`
	StoredFilename   string    `json:"stored_filename"`
	FilePath         string    `json:"file_path"`
	FileSize         int64     `json:"file_size"`
	MimeType         string    `json:"mime_type"`
	UploadedBy       *uuid.UUID `json:"uploaded_by,omitempty"`
	HasThumbnail     bool      `json:"has_thumbnail"`
	ThumbnailPath    *string   `json:"thumbnail_path,omitempty"`
}

type AttachmentResponse struct {
	Attachments []TaskAttachment `json:"attachments"`
	TotalSize   int64            `json:"total_size"`
	Count       int              `json:"count"`
}

type UploadResponse struct {
	Success     []TaskAttachment `json:"success"`
	Failed      []UploadError    `json:"failed"`
	TotalSize   int64            `json:"total_size"`
	Count       int              `json:"count"`
}

type UploadError struct {
	Filename string `json:"filename"`
	Error    string `json:"error"`
}

type AttachmentMetadata struct {
	ID               uuid.UUID `json:"id"`
	OriginalFilename string    `json:"original_filename"`
	FileSize         int64     `json:"file_size"`
	MimeType         string    `json:"mime_type"`
	UploaderName     *string   `json:"uploader_name,omitempty"`
	HasThumbnail     bool      `json:"has_thumbnail"`
	CreatedAt        time.Time `json:"created_at"`
}