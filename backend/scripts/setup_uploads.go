package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// setupUploads creates the uploads directory structure with proper permissions
func main() {
	baseUploadPath := os.Getenv("UPLOAD_PATH")
	if baseUploadPath == "" {
		baseUploadPath = "./uploads"
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(baseUploadPath)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	// Create directory structure
	directories := []string{
		filepath.Join(absPath, "attachments"),
		filepath.Join(absPath, "thumbnails"),
	}

	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", dir, err)
		}
		fmt.Printf("Created directory: %s\n", dir)
	}

	// Create .gitkeep files to ensure directories are tracked in git
	gitkeepFiles := []string{
		filepath.Join(absPath, "attachments", ".gitkeep"),
		filepath.Join(absPath, "thumbnails", ".gitkeep"),
	}

	for _, file := range gitkeepFiles {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			f, err := os.Create(file)
			if err != nil {
				log.Printf("Warning: Failed to create .gitkeep file %s: %v", file, err)
				continue
			}
			f.Close()
			fmt.Printf("Created .gitkeep file: %s\n", file)
		}
	}

	// Set proper permissions (readable/writable by owner, readable by group)
	for _, dir := range directories {
		if err := os.Chmod(dir, 0755); err != nil {
			log.Printf("Warning: Failed to set permissions for %s: %v", dir, err)
		}
	}

	fmt.Printf("Upload directory structure initialized successfully at: %s\n", absPath)
	fmt.Println("Directory permissions set to 0755 (rwxr-xr-x)")
	fmt.Println("Ensure the application has write access to this directory")
}