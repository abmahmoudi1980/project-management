# Design Document: Task Attachments

## Overview

The task attachments feature extends the existing project management system to support file uploads, storage, and management for tasks. This design integrates seamlessly with the current Go backend architecture (handlers → services → repositories) and Svelte 5 frontend, providing secure file handling with proper access controls and an intuitive user interface.

The system will support common file types (documents, images, archives) with security-first design principles, including file validation, secure storage, and access control based on task permissions.

## Architecture

### Backend Architecture

The attachment system follows the established three-layer architecture:

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────────┐
│   Handlers      │    │    Services      │    │   Repositories      │
│                 │    │                  │    │                     │
│ attachment_     │───▶│ attachment_      │───▶│ attachment_         │
│ handler.go      │    │ service.go       │    │ repository.go       │
│                 │    │                  │    │                     │
│ - Upload        │    │ - Validation     │    │ - Database ops      │
│ - Download      │    │ - File storage   │    │ - Metadata CRUD     │
│ - List/Delete   │    │ - Security       │    │ - Cleanup           │
└─────────────────┘    └──────────────────┘    └─────────────────────┘
```

### File Storage Strategy

```
uploads/
├── attachments/
│   ├── 2024/
│   │   ├── 01/
│   │   │   ├── 15/
│   │   │   │   ├── {uuid}-original-filename.pdf
│   │   │   │   └── {uuid}-document.docx
│   │   │   └── 16/
│   │   └── 02/
│   └── thumbnails/
│       ├── 2024/
│       │   ├── 01/
│       │   │   └── 15/
│       │   │       ├── {uuid}-thumb.jpg
│       │   │       └── {uuid}-thumb.png
```

**Storage Design Principles:**

- Files stored outside web root for security
- Hierarchical directory structure (year/month/day) to prevent filesystem limitations
- UUID-based filenames to prevent conflicts and directory traversal
- Separate thumbnail storage for images
- Configurable storage backends (local filesystem, cloud storage)

## Components and Interfaces

### Database Schema

```sql
CREATE TABLE task_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    original_filename VARCHAR(255) NOT NULL,
    stored_filename VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    uploaded_by UUID REFERENCES users(id),
    has_thumbnail BOOLEAN DEFAULT FALSE,
    thumbnail_path TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_task_attachments_task_id ON task_attachments(task_id);
CREATE INDEX idx_task_attachments_uploaded_by ON task_attachments(uploaded_by);
```

### Go Models

```go
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

type UploadRequest struct {
    Files []multipart.FileHeader `form:"files"`
}

type AttachmentResponse struct {
    Attachments []TaskAttachment `json:"attachments"`
    TotalSize   int64            `json:"total_size"`
    Count       int              `json:"count"`
}
```

### API Endpoints

```
POST   /api/tasks/{taskId}/attachments     - Upload files
GET    /api/tasks/{taskId}/attachments     - List attachments
GET    /api/attachments/{id}/download      - Download file
GET    /api/attachments/{id}/thumbnail     - Get thumbnail
DELETE /api/attachments/{id}               - Delete attachment
```

### Frontend Components

**AttachmentManager.svelte** - Main component integrated into TaskDetails

- File upload with drag-and-drop
- Attachment list with thumbnails
- Progress indicators
- Error handling

**AttachmentUploader.svelte** - Upload interface component

- Drag-and-drop zone
- File selection
- Upload progress
- Validation feedback

**AttachmentList.svelte** - Display component

- File list with metadata
- Thumbnail previews
- Download/delete actions
- Responsive design

## Data Models

### File Validation Rules

```go
type FileValidationConfig struct {
    MaxFileSize      int64    // 10MB default
    MaxTotalSize     int64    // 100MB per task default
    AllowedTypes     []string // MIME types
    AllowedExtensions []string // File extensions
    BlockedExtensions []string // Dangerous extensions
}

var DefaultConfig = FileValidationConfig{
    MaxFileSize:  10 * 1024 * 1024, // 10MB
    MaxTotalSize: 100 * 1024 * 1024, // 100MB
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
}
```

### Security Model

**Access Control:**

- Users can only upload/download attachments for tasks they have access to
- Task access follows existing project ownership rules
- Admin users have full access to all attachments

**File Security:**

- Files stored outside web root
- Secure filename generation using UUIDs
- MIME type validation
- File content scanning (basic)
- Security headers on downloads

## Correctness Properties

_A property is a characteristic or behavior that should hold true across all valid executions of a system-essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees._

### Property Reflection

After analyzing the acceptance criteria, several properties can be consolidated:

- Properties 1.1 and 1.5 both relate to file type validation and can be combined
- Properties 2.1 and 2.5 both relate to attachment display information and can be combined
- Properties 4.1 and 4.2 both relate to file cleanup and can be combined
- Properties 6.1, 6.3, and 6.4 all relate to API behavior and can be combined

### Core Properties

**Property 1: File Type Validation**
_For any_ uploaded file, the system should accept it if and only if its extension and MIME type are in the allowed lists and not in the blocked lists
**Validates: Requirements 1.1, 1.5**

**Property 2: File Size Validation**
_For any_ uploaded file or set of files, the system should reject uploads that exceed individual file size limits or total task attachment limits
**Validates: Requirements 1.2, 4.4**

**Property 3: Upload Processing Consistency**
_For any_ batch of uploaded files, each file should be processed individually with success/failure reported per file, and valid files should be stored with corresponding database records
**Validates: Requirements 1.3, 1.4**

**Property 4: Attachment Display Completeness**
_For any_ task with attachments, the display should include all attachment metadata (filename, size, upload date, uploader name) and thumbnails for image files
**Validates: Requirements 2.1, 2.2, 2.5**

**Property 5: Download Access Control**
_For any_ download request, the system should verify the user has access to the parent task before serving the file with appropriate security headers
**Validates: Requirements 2.3, 3.1, 3.5**

**Property 6: Permission-Based UI**
_For any_ user viewing attachments, delete options should be available if and only if the user has appropriate permissions for the parent task
**Validates: Requirements 2.4**

**Property 7: Secure File Storage**
_For any_ uploaded file, the system should store it with a UUID-based filename in a structured directory hierarchy outside the web root
**Validates: Requirements 3.2, 4.3, 4.5**

**Property 8: File Cleanup Consistency**
_For any_ deleted task or attachment, all associated physical files (including thumbnails) should be immediately removed from storage
**Validates: Requirements 4.1, 4.2**

**Property 9: UI Reactivity**
_For any_ attachment operation (upload, delete), the UI should update without page refresh and provide appropriate progress indicators and feedback
**Validates: Requirements 5.2, 5.3, 5.4, 5.5**

**Property 10: API Consistency**
_For any_ API operation, the system should return appropriate HTTP status codes, descriptive error messages, and maintain referential integrity between tasks and attachments
**Validates: Requirements 6.1, 6.2, 6.3, 6.4, 6.5**

## Error Handling

### File Upload Errors

- **Invalid file type**: Return 400 with specific MIME type/extension error
- **File too large**: Return 413 with size limit information
- **Total size exceeded**: Return 413 with total limit information
- **Storage failure**: Return 500 with generic error message
- **Malicious content**: Return 400 with security violation message

### Download Errors

- **File not found**: Return 404
- **Access denied**: Return 403 with permission error
- **Corrupted file**: Return 500 with file integrity error

### Database Errors

- **Task not found**: Return 404
- **Attachment not found**: Return 404
- **Foreign key violation**: Return 400 with relationship error

### Frontend Error Handling

- Network errors: Retry mechanism with exponential backoff
- Upload failures: Per-file error display with retry options
- Validation errors: Real-time feedback with clear messaging
- Progress tracking: Cancellation support and error recovery

## Testing Strategy

### Dual Testing Approach

The testing strategy employs both unit tests and property-based tests as complementary approaches:

**Unit Tests:**

- Specific examples demonstrating correct behavior
- Edge cases and error conditions
- Integration points between components
- API endpoint functionality with known inputs

**Property-Based Tests:**

- Universal properties across all inputs using randomized data
- Comprehensive input coverage through generation
- Minimum 100 iterations per property test
- Each test tagged with: **Feature: task-attachments, Property {number}: {property_text}**

### Property-Based Testing Configuration

Using **Testify** with **gopter** for Go property-based testing:

- Generate random file data, names, sizes, and types
- Test file validation across all possible inputs
- Verify storage consistency with random file sets
- Test access control with random user/task combinations
- Validate UI behavior with random attachment states

### Test Categories

**File Handling Tests:**

- Upload validation with various file types and sizes
- Storage path generation and file organization
- Thumbnail generation for image files
- File cleanup on deletion

**Security Tests:**

- Access control verification
- Malicious filename handling
- MIME type validation
- Directory traversal prevention

**API Integration Tests:**

- Multipart upload handling
- Error response formatting
- Authentication and authorization
- Database transaction integrity

**Frontend Tests:**

- Component rendering with various attachment states
- Drag-and-drop functionality
- Progress indicator behavior
- Error message display
