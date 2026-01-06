# Implementation Plan: Task Attachments

## Overview

This implementation plan breaks down the task attachments feature into discrete coding steps that build incrementally. Each task focuses on specific components while maintaining integration with the existing Go backend architecture and Svelte 5 frontend.

## Tasks

- [x] 1. Set up database schema and core models

  - Create migration file for task_attachments table with proper indexes
  - Define TaskAttachment model in Go with JSON tags
  - Add attachment-related request/response structs
  - _Requirements: 6.5_

- [ ]\* 1.1 Write property test for database schema

  - **Property 10: API Consistency - Database integrity portion**
  - **Validates: Requirements 6.5**

- [x] 2. Implement file storage infrastructure

  - [x] 2.1 Create file storage service with directory hierarchy

    - Implement secure file path generation using UUID and date structure
    - Create directory creation and file storage functions
    - Add file cleanup utilities for deletion
    - _Requirements: 3.2, 4.3, 4.5_

  - [ ]\* 2.2 Write property test for secure file storage

    - **Property 7: Secure File Storage**
    - **Validates: Requirements 3.2, 4.3, 4.5**

  - [x] 2.3 Implement file validation service

    - Create MIME type and extension validation
    - Add file size validation (individual and total)
    - Implement basic malicious content detection
    - _Requirements: 1.1, 1.2, 4.4_

  - [ ]\* 2.4 Write property tests for file validation
    - **Property 1: File Type Validation**
    - **Property 2: File Size Validation**
    - **Validates: Requirements 1.1, 1.2, 1.5, 4.4**

- [-] 3. Create attachment repository layer

  - [x] 3.1 Implement attachment repository with CRUD operations

    - Create database operations for attachment metadata
    - Add methods for listing attachments by task
    - Implement cascade deletion support
    - _Requirements: 4.1, 6.5_

  - [ ]\* 3.2 Write property test for file cleanup
    - **Property 8: File Cleanup Consistency**
    - **Validates: Requirements 4.1, 4.2**

- [x] 4. Implement attachment service layer

  - [x] 4.1 Create attachment service with business logic

    - Integrate file validation with storage operations
    - Implement batch upload processing with individual error handling
    - Add thumbnail generation for image files
    - _Requirements: 1.3, 1.4, 2.2_

  - [ ]\* 4.2 Write property test for upload processing

    - **Property 3: Upload Processing Consistency**
    - **Validates: Requirements 1.3, 1.4**

  - [x] 4.3 Add access control to attachment service

    - Implement task permission verification
    - Add user-based attachment access checks
    - _Requirements: 3.1, 2.4_

  - [ ]\* 4.4 Write property test for access control
    - **Property 5: Download Access Control**
    - **Property 6: Permission-Based UI**
    - **Validates: Requirements 2.3, 2.4, 3.1, 3.5**

- [x] 5. Create attachment handler layer

  - [x] 5.1 Implement upload handler with multipart support

    - Create POST endpoint for file uploads
    - Add proper error handling and response formatting
    - Implement progress tracking support
    - _Requirements: 6.1, 6.4_

  - [x] 5.2 Implement download and list handlers

    - Create secure download endpoint with proper headers
    - Add attachment listing endpoint
    - Implement thumbnail serving endpoint
    - _Requirements: 6.2, 6.3, 3.5_

  - [x] 5.3 Implement delete handler

    - Create DELETE endpoint for attachments
    - Ensure proper cleanup of files and database records
    - _Requirements: 4.2_

  - [ ]\* 5.4 Write property test for API consistency
    - **Property 10: API Consistency**
    - **Validates: Requirements 6.1, 6.2, 6.3, 6.4, 6.5**

- [x] 6. Integrate attachment routes

  - [x] 6.1 Add attachment routes to main router
    - Register all attachment endpoints
    - Apply authentication middleware
    - Update main.go to initialize attachment components
    - _Requirements: 6.1, 6.2, 6.3_

- [x] 7. Checkpoint - Backend API Testing

  - Ensure all backend tests pass, verify API endpoints work correctly, ask the user if questions arise.

- [ ] 8. Create frontend attachment components

  - [ ] 8.1 Create AttachmentUploader component

    - Implement drag-and-drop file upload interface
    - Add file validation feedback and progress indicators
    - Create upload cancellation functionality
    - _Requirements: 5.2, 5.3, 5.5_

  - [ ]\* 8.2 Write property test for UI reactivity (upload)

    - **Property 9: UI Reactivity - Upload portion**
    - **Validates: Requirements 5.2, 5.3, 5.5**

  - [ ] 8.3 Create AttachmentList component

    - Display attachments with metadata (filename, size, date, uploader)
    - Implement thumbnail previews for images
    - Add download and delete action buttons
    - _Requirements: 2.1, 2.2, 2.5_

  - [ ]\* 8.4 Write property test for attachment display

    - **Property 4: Attachment Display Completeness**
    - **Validates: Requirements 2.1, 2.2, 2.5**

  - [ ] 8.5 Create AttachmentManager component
    - Integrate uploader and list components
    - Handle state management and API communication
    - Implement real-time UI updates after operations
    - _Requirements: 5.4_

- [ ] 9. Integrate attachments into TaskDetails

  - [ ] 9.1 Add attachment section to TaskDetails component

    - Integrate AttachmentManager into existing task view
    - Update task details layout to accommodate attachments
    - Ensure responsive design consistency
    - _Requirements: 5.1_

  - [ ] 9.2 Update API client for attachment operations
    - Add attachment endpoints to api.js
    - Implement multipart upload support
    - Add error handling for attachment operations
    - _Requirements: 6.1, 6.2, 6.3_

- [ ]\* 9.3 Write property test for complete UI reactivity

  - **Property 9: UI Reactivity**
  - **Validates: Requirements 5.2, 5.3, 5.4, 5.5**

- [ ] 10. Add configuration and environment setup

  - [ ] 10.1 Add attachment configuration

    - Create configuration for file size limits and allowed types
    - Add environment variables for storage paths
    - Update Docker configuration if needed
    - _Requirements: 1.1, 1.2, 4.4_

  - [ ] 10.2 Create uploads directory structure
    - Set up secure file storage directories
    - Ensure proper permissions and security
    - _Requirements: 3.4_

- [ ] 11. Final integration and testing

  - [ ] 11.1 End-to-end integration testing

    - Test complete upload-to-download workflow
    - Verify all error scenarios work correctly
    - Test with various file types and sizes
    - _Requirements: All_

  - [ ]\* 11.2 Write integration tests for complete workflow
    - Test full attachment lifecycle
    - Verify error handling across all components
    - _Requirements: All_

- [ ] 12. Final checkpoint - Complete system verification
  - Ensure all tests pass, verify complete attachment functionality works end-to-end, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Property tests validate universal correctness properties using randomized inputs
- Unit tests validate specific examples and edge cases
- Backend tasks (1-7) can be completed independently of frontend tasks (8-11)
- Configuration tasks (10) should be completed before final testing (11-12)
