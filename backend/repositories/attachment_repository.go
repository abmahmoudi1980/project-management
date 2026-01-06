package repositories

import (
	"context"
	"project-management/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AttachmentRepository struct {
	db *pgxpool.Pool
}

func NewAttachmentRepository(db *pgxpool.Pool) *AttachmentRepository {
	return &AttachmentRepository{db: db}
}

// Create creates a new task attachment record
func (r *AttachmentRepository) Create(ctx context.Context, req models.CreateAttachmentRequest) (*models.TaskAttachment, error) {
	id := uuid.New()
	var attachment models.TaskAttachment

	err := r.db.QueryRow(ctx,
		`INSERT INTO task_attachments (id, task_id, original_filename, stored_filename, file_path, file_size, mime_type, uploaded_by, has_thumbnail, thumbnail_path)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 RETURNING id, task_id, original_filename, stored_filename, file_path, file_size, mime_type, uploaded_by, has_thumbnail, thumbnail_path, created_at, updated_at`,
		id, req.TaskID, req.OriginalFilename, req.StoredFilename, req.FilePath, req.FileSize, req.MimeType, req.UploadedBy, req.HasThumbnail, req.ThumbnailPath).
		Scan(&attachment.ID, &attachment.TaskID, &attachment.OriginalFilename, &attachment.StoredFilename, &attachment.FilePath, &attachment.FileSize, &attachment.MimeType, &attachment.UploadedBy, &attachment.HasThumbnail, &attachment.ThumbnailPath, &attachment.CreatedAt, &attachment.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &attachment, nil
}

// GetByID retrieves a single attachment by its ID
func (r *AttachmentRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.TaskAttachment, error) {
	var attachment models.TaskAttachment
	err := r.db.QueryRow(ctx,
		`SELECT id, task_id, original_filename, stored_filename, file_path, file_size, mime_type, uploaded_by, has_thumbnail, thumbnail_path, created_at, updated_at
		 FROM task_attachments WHERE id = $1`, id).
		Scan(&attachment.ID, &attachment.TaskID, &attachment.OriginalFilename, &attachment.StoredFilename, &attachment.FilePath, &attachment.FileSize, &attachment.MimeType, &attachment.UploadedBy, &attachment.HasThumbnail, &attachment.ThumbnailPath, &attachment.CreatedAt, &attachment.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &attachment, nil
}

// GetByIDWithUploader retrieves a single attachment by its ID with uploader information
func (r *AttachmentRepository) GetByIDWithUploader(ctx context.Context, id uuid.UUID) (*models.TaskAttachment, error) {
	var attachment models.TaskAttachment
	err := r.db.QueryRow(ctx,
		`SELECT ta.id, ta.task_id, ta.original_filename, ta.stored_filename, ta.file_path, ta.file_size, ta.mime_type, 
		        ta.uploaded_by, ta.has_thumbnail, ta.thumbnail_path, ta.created_at, ta.updated_at,
		        u.username as uploader_name
		 FROM task_attachments ta
		 LEFT JOIN users u ON ta.uploaded_by = u.id
		 WHERE ta.id = $1`, id).
		Scan(&attachment.ID, &attachment.TaskID, &attachment.OriginalFilename, &attachment.StoredFilename, &attachment.FilePath, &attachment.FileSize, &attachment.MimeType, &attachment.UploadedBy, &attachment.HasThumbnail, &attachment.ThumbnailPath, &attachment.CreatedAt, &attachment.UpdatedAt, &attachment.UploaderName)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &attachment, nil
}

// GetByTaskID retrieves all attachments for a specific task
func (r *AttachmentRepository) GetByTaskID(ctx context.Context, taskID uuid.UUID) ([]models.TaskAttachment, error) {
	rows, err := r.db.Query(ctx,
		`SELECT ta.id, ta.task_id, ta.original_filename, ta.stored_filename, ta.file_path, ta.file_size, ta.mime_type,
		        ta.uploaded_by, ta.has_thumbnail, ta.thumbnail_path, ta.created_at, ta.updated_at,
		        u.username as uploader_name
		 FROM task_attachments ta
		 LEFT JOIN users u ON ta.uploaded_by = u.id
		 WHERE ta.task_id = $1
		 ORDER BY ta.created_at DESC`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []models.TaskAttachment
	for rows.Next() {
		var attachment models.TaskAttachment
		if err := rows.Scan(&attachment.ID, &attachment.TaskID, &attachment.OriginalFilename, &attachment.StoredFilename, &attachment.FilePath, &attachment.FileSize, &attachment.MimeType, &attachment.UploadedBy, &attachment.HasThumbnail, &attachment.ThumbnailPath, &attachment.CreatedAt, &attachment.UpdatedAt, &attachment.UploaderName); err != nil {
			return nil, err
		}
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

// GetTotalSizeByTaskID calculates the total size of all attachments for a task
func (r *AttachmentRepository) GetTotalSizeByTaskID(ctx context.Context, taskID uuid.UUID) (int64, error) {
	var totalSize int64
	err := r.db.QueryRow(ctx,
		`SELECT COALESCE(SUM(file_size), 0) FROM task_attachments WHERE task_id = $1`, taskID).
		Scan(&totalSize)
	if err != nil {
		return 0, err
	}
	return totalSize, nil
}

// GetCountByTaskID returns the number of attachments for a task
func (r *AttachmentRepository) GetCountByTaskID(ctx context.Context, taskID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM task_attachments WHERE task_id = $1`, taskID).
		Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Delete removes an attachment record from the database
func (r *AttachmentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM task_attachments WHERE id = $1", id)
	return err
}

// DeleteByTaskID removes all attachments for a specific task (cascade deletion support)
func (r *AttachmentRepository) DeleteByTaskID(ctx context.Context, taskID uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM task_attachments WHERE task_id = $1", taskID)
	return err
}

// GetAttachmentsToCleanup retrieves attachment file paths for cleanup before deletion
func (r *AttachmentRepository) GetAttachmentsToCleanup(ctx context.Context, taskID uuid.UUID) ([]models.TaskAttachment, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, task_id, original_filename, stored_filename, file_path, file_size, mime_type,
		        uploaded_by, has_thumbnail, thumbnail_path, created_at, updated_at
		 FROM task_attachments
		 WHERE task_id = $1`, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var attachments []models.TaskAttachment
	for rows.Next() {
		var attachment models.TaskAttachment
		if err := rows.Scan(&attachment.ID, &attachment.TaskID, &attachment.OriginalFilename, &attachment.StoredFilename, &attachment.FilePath, &attachment.FileSize, &attachment.MimeType, &attachment.UploadedBy, &attachment.HasThumbnail, &attachment.ThumbnailPath, &attachment.CreatedAt, &attachment.UpdatedAt); err != nil {
			return nil, err
		}
		attachments = append(attachments, attachment)
	}

	return attachments, nil
}

// GetAttachmentToCleanup retrieves a single attachment's file paths for cleanup before deletion
func (r *AttachmentRepository) GetAttachmentToCleanup(ctx context.Context, id uuid.UUID) (*models.TaskAttachment, error) {
	var attachment models.TaskAttachment
	err := r.db.QueryRow(ctx,
		`SELECT id, task_id, original_filename, stored_filename, file_path, file_size, mime_type,
		        uploaded_by, has_thumbnail, thumbnail_path, created_at, updated_at
		 FROM task_attachments WHERE id = $1`, id).
		Scan(&attachment.ID, &attachment.TaskID, &attachment.OriginalFilename, &attachment.StoredFilename, &attachment.FilePath, &attachment.FileSize, &attachment.MimeType, &attachment.UploadedBy, &attachment.HasThumbnail, &attachment.ThumbnailPath, &attachment.CreatedAt, &attachment.UpdatedAt)

	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &attachment, nil
}

// Update updates an attachment record (mainly for thumbnail path updates)
func (r *AttachmentRepository) Update(ctx context.Context, id uuid.UUID, hasThumb bool, thumbPath *string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE task_attachments SET has_thumbnail = $1, thumbnail_path = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`,
		hasThumb, thumbPath, id)
	return err
}