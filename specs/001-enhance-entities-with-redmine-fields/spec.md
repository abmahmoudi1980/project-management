# Feature Specification: Enhance Entities with Redmine Fields

**Feature Branch**: `001-enhance-entities-with-redmine-fields`
**Created**: 2025-12-30
**Status**: Draft
**Input**: User description: "I want to add more features to this project. Take a look at Redmine. I don't want to add all the Redmine features to this project, but I want to make the current entities more complete and add the necessary fields to them."

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Enhanced Project Management (Priority: P1)

As a project manager, I want to create projects with additional metadata (identifier, homepage URL, public visibility) so that I can better organize and share project information.

**Why this priority**: This is foundational - all other enhancements build on having more complete project metadata. These fields are essential for project identification and sharing.

**Independent Test**: Can be fully tested by creating projects with the new fields and verifying they are stored and retrieved correctly, without needing tasks or time logs.

**Acceptance Scenarios**:

1. **Given** a new project creation form, **When** user enters a unique identifier and optional homepage URL, **Then** the project is created with these fields saved to the database
2. **Given** an existing project, **When** user updates the is_public flag, **Then** the project's visibility status changes accordingly
3. **Given** a project with an identifier, **When** user tries to create another project with the same identifier, **Then** the system returns a validation error

---

### User Story 2 - Enhanced Task Management (Priority: P1)

As a developer, I want to track detailed task information including assignee, category, dates, estimates, and progress so that I can manage work more effectively.

**Why this priority**: Task management is core to the application. These fields provide essential project management capabilities without adding complexity like workflows or custom fields.

**Independent Test**: Can be fully tested by creating and updating tasks with the new fields and verifying all data is persisted and accessible.

**Acceptance Scenarios**:

1. **Given** a task in a project, **When** user sets an assignee, category, and start/due dates, **Then** all fields are saved and displayed correctly
2. **Given** a task with estimated hours, **When** user updates the done ratio (0-100%), **Then** the progress is calculated and displayed
3. **Given** a task, **When** user adds a description, **Then** the description is saved and rendered properly
4. **Given** a task without a due date, **When** user filters tasks by date range, **Then** the task is excluded from results

---

### User Story 3 - Enhanced Time Logging (Priority: P2)

As a team member, I want to log time with activity types and user attribution so that time tracking is more accurate and categorizable.

**Why this priority**: Time tracking is valuable but less critical than project and task management. This enhancement adds useful categorization without requiring workflow changes.

**Independent Test**: Can be fully tested by creating time logs with different activity types and verifying they are stored and filterable.

**Acceptance Scenarios**:

1. **Given** a task, **When** user logs time with an activity type (e.g., Development, Testing), **Then** the time entry is saved with the activity type
2. **Given** multiple time entries, **When** user filters by activity type, **Then** only entries matching that type are shown
3. **Given** a time entry, **When** user updates the user_id, **Then** the entry shows the new user

---

### User Story 4 - Project-Level Time Tracking (Priority: P3)

As a project manager, I want to log time at the project level (not just task level) so that I can track overhead and administrative time.

**Why this priority**: Nice-to-have enhancement that provides more flexibility. Doesn't block core functionality.

**Independent Test**: Can be fully tested by creating time logs with project_id instead of task_id.

**Acceptance Scenarios**:

1. **Given** a project, **When** user logs time without specifying a task, **Then** the time entry is associated with the project only
2. **Given** project-level time entries, **When** viewing project time reports, **Then** both project and task time are aggregated

---

### Edge Cases

- What happens when a task's done_ratio is set to 100% but the task is not marked as completed? → Both fields can exist independently (completion vs progress)
- What happens when a project identifier contains invalid characters? → Validation should only allow alphanumeric, underscores, and hyphens
- What happens when a time entry has both project_id and task_id? → Task ID should take precedence for association
- What happens when estimated_hours is negative? → Validation should reject negative values

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow projects to have a unique identifier (VARCHAR 50)
- **FR-002**: System MUST allow projects to have an optional homepage URL (VARCHAR 255)
- **FR-003**: System MUST allow projects to have an is_public flag (BOOLEAN, default: false)
- **FR-004**: System MUST validate that project identifiers are unique across all projects
- **FR-005**: System MUST allow tasks to have a description (TEXT, optional)
- **FR-006**: System MUST allow tasks to have an assignee_id (UUID, references users table)
- **FR-007**: System MUST allow tasks to have an author_id (UUID, references users table, default: creator)
- **FR-008**: System MUST allow tasks to have a category (VARCHAR 100, optional)
- **FR-009**: System MUST allow tasks to have start_date (DATE, optional)
- **FR-010**: System MUST allow tasks to have due_date (DATE, optional)
- **FR-011**: System MUST allow tasks to have estimated_hours (DECIMAL, optional, must be >= 0)
- **FR-012**: System MUST allow tasks to have done_ratio (INTEGER, 0-100, default: 0)
- **FR-013**: System MUST validate that due_date >= start_date when both are provided
- **FR-014**: System MUST validate that done_ratio is between 0 and 100
- **FR-015**: System MUST allow time logs to have a user_id (UUID, references users table)
- **FR-016**: System MUST allow time logs to have an activity_type (VARCHAR 50, optional, e.g., "Development", "Testing", "Design")
- **FR-017**: System MUST allow time logs to be associated with a project (project_id, optional) instead of just a task

### Key Entities *(include if feature involves data)*

- **Project**: Enhanced with identifier (unique), homepage URL, and is_public flag for better project organization and sharing
- **Task**: Enhanced with description, assignee, author, category, start/due dates, estimated hours, and done_ratio for comprehensive task management
- **TimeLog**: Enhanced with user_id, activity_type, and optional project_id for more detailed time tracking

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All new fields are correctly added to the database schema with appropriate constraints
- **SC-002**: Backend API endpoints support all new fields in create/update operations
- **SC-003**: Frontend forms display and validate all new fields appropriately
- **SC-004**: Existing functionality (projects, tasks, time logs) continues to work without breaking changes
- **SC-005**: Database migration script cleanly adds new columns without data loss
- **SC-006**: All validation rules are enforced (unique identifiers, date ranges, numeric ranges)

### Technical Outcomes

- **TC-001**: Database migration runs successfully with no errors
- **TC-002**: All existing data remains intact after migration
- **TC-003**: New Go models include all additional fields with proper JSON tags
- **TC-004**: Repository functions properly handle new fields in CRUD operations
- **TC-005**: Service layer validates new fields according to requirements
- **TC-006**: Frontend stores and components display new fields correctly
