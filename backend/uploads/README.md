# Uploads Directory

This directory contains user-uploaded files for the task attachments feature.

## Directory Structure

```
uploads/
├── attachments/          # Original uploaded files
│   └── YYYY/MM/DD/      # Organized by upload date
│       └── {uuid}-{filename}
└── thumbnails/          # Generated thumbnails for images
    └── YYYY/MM/DD/      # Organized by upload date
        └── {uuid}-thumb.{ext}
```

## Security Considerations

1. **Location**: This directory should be outside the web root to prevent direct access
2. **Permissions**: Directory permissions should be 0755 (rwxr-xr-x)
3. **File Access**: Files are served through the API with proper authentication
4. **File Names**: All files use UUID-based names to prevent conflicts and directory traversal
5. **Validation**: All uploads are validated for type, size, and content before storage

## Configuration

The upload path can be configured via environment variables:

- `UPLOAD_PATH`: Base directory for uploads (default: ./uploads)
- `MAX_FILE_SIZE`: Maximum individual file size in bytes (default: 10MB)
- `MAX_TOTAL_SIZE`: Maximum total attachments per task in bytes (default: 100MB)
- `THUMBNAIL_SIZE`: Thumbnail size in pixels (default: 200px)

## Maintenance

- Old files are automatically cleaned up when tasks or attachments are deleted
- The directory structure is created automatically by the application
- Use `go run scripts/setup_uploads.go` to manually initialize the directory structure

## Docker Deployment

In Docker environments:

- The uploads directory is mounted as a volume for persistence
- Path is typically set to `/app/uploads` inside the container
- Ensure proper volume permissions in docker-compose.yml
