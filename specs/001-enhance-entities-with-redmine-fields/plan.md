# Implementation Plan: Enhance Entities with Redmine Fields

**Branch**: `001-enhance-entities-with-redmine-fields` | **Date**: 2025-12-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-enhance-entities-with-redmine-fields/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

This feature enhances the existing Project, Task, and TimeLog entities by adding fields inspired by Redmine's data model. The goal is to make entities more complete for project management without adding Redmine's full complexity (workflows, custom fields, plugins, etc.).

**Primary requirements:**
- Add `identifier`, `homepage`, and `is_public` to projects
- Add `description`, `assignee_id`, `author_id`, `category`, `start_date`, `due_date`, `estimated_hours`, and `done_ratio` to tasks
- Add `user_id`, `activity_type`, and optional `project_id` to time logs

**Technical approach:** Incremental database migration using PostgreSQL `ALTER TABLE` statements, followed by updates to Go models, repositories, services, handlers, and frontend components.

## Technical Context

**Language/Version**: Go 1.21+ (backend), Node.js 18+ with Svelte 4 (frontend)
**Primary Dependencies**: Fiber (Go web framework), pgx (PostgreSQL driver), Svelte, Vite, Tailwind CSS
**Storage**: PostgreSQL 12+ with UUID primary keys
**Testing**: No test framework currently configured (NEEDS CLARIFICATION - add testing in future)
**Target Platform**: Linux server (backend), Web browser (frontend)
**Project Type**: web application (backend + frontend)
**Performance Goals**: <200ms p95 for API responses, maintain current performance levels
**Constraints**: Must maintain backward compatibility - existing data must remain intact after migration
**Scale/Scope**: Small-to-medium scale (< 1000 concurrent users), < 10K projects/tasks

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

[No constitution file exists - skipping constitution check. To add this, create `.specify/memory/constitution.md`]

## Project Structure

### Documentation (this feature)

```text
specs/001-enhance-entities-with-redmine-fields/
├── plan.md              # This file
├── spec.md              # Feature specification
└── migration.sql         # Database migration script (to be created)
```

### Source Code (repository root)

```text
backend/
├── models/
│   ├── project.go        # UPDATE: Add new fields
│   ├── task.go          # UPDATE: Add new fields
│   └── timelog.go      # UPDATE: Add new fields
├── repositories/
│   ├── project_repository.go    # UPDATE: Handle new fields
│   ├── task_repository.go      # UPDATE: Handle new fields
│   └── timelog_repository.go  # UPDATE: Handle new fields
├── services/
│   ├── project_service.go       # UPDATE: Validation for new fields
│   ├── task_service.go         # UPDATE: Validation for new fields
│   └── timelog_service.go     # UPDATE: Handle project_id option
├── handlers/
│   ├── project_handler.go       # UPDATE: Process new fields
│   ├── task_handler.go         # UPDATE: Process new fields
│   └── timelog_handler.go     # UPDATE: Handle project_id option
├── routes/
│   └── routes.go              # NO CHANGE - same endpoints
├── main.go                   # NO CHANGE
└── migration/
    └── 001_enhance_entities.sql   # NEW: Database migration

frontend/
├── src/
│   ├── components/
│   │   ├── ProjectForm.svelte       # UPDATE: Add new fields
│   │   ├── ProjectList.svelte       # UPDATE: Display new fields
│   │   ├── TaskForm.svelte         # UPDATE: Add new fields
│   │   └── TaskList.svelte         # UPDATE: Display new fields
│   ├── stores/
│   │   ├── projectStore.js          # UPDATE: Handle new fields
│   │   ├── taskStore.js            # UPDATE: Handle new fields
│   │   └── timeLogStore.js        # UPDATE: Handle new fields
│   └── lib/
│       └── api.js                 # UPDATE: Pass new fields
├── package.json
└── vite.config.js
```

**Structure Decision**: This is a web application with a Go backend and Svelte frontend. The feature modifies existing entities across the full stack: database schema, backend models/services/handlers, and frontend components/stores.

## Phase-by-Phase Implementation

### Phase 0: Database Migration

**Goal**: Add new columns to existing tables without data loss.

**Steps**:
1. Create `backend/migration/001_enhance_entities.sql`
2. Add columns to `projects` table:
   - `identifier VARCHAR(50) UNIQUE`
   - `homepage VARCHAR(255)`
   - `is_public BOOLEAN DEFAULT FALSE`
3. Add columns to `tasks` table:
   - `description TEXT`
   - `assignee_id UUID` (foreign key to users - optional for now)
   - `author_id UUID` (foreign key to users - optional for now)
   - `category VARCHAR(100)`
   - `start_date DATE`
   - `due_date DATE`
   - `estimated_hours DECIMAL(10,2)`
   - `done_ratio INTEGER DEFAULT 0 CHECK (done_ratio >= 0 AND done_ratio <= 100)`
4. Add columns to `time_logs` table:
   - `user_id UUID` (foreign key to users - optional for now)
   - `activity_type VARCHAR(50)`
   - `project_id UUID` (foreign key to projects - optional)
5. Add indexes for performance:
   - `idx_projects_identifier` on `projects(identifier)`
   - `idx_tasks_assignee_id` on `tasks(assignee_id)`
   - `idx_tasks_due_date` on `tasks(due_date)`
   - `idx_time_logs_user_id` on `time_logs(user_id)`
   - `idx_time_logs_activity_type` on `time_logs(activity_type)`

**Validation**:
- Run migration: `psql -U postgres -d project_management -f backend/migration/001_enhance_entities.sql`
- Verify all columns added: `\d projects`, `\d tasks`, `\d time_logs`

### Phase 1: Backend Models

**Goal**: Update Go struct definitions to include new fields.

**Files**:
- `backend/models/project.go`
- `backend/models/task.go`
- `backend/models/timelog.go`

**Changes**:
- Add new fields with JSON tags
- Update `Create*Request` and `Update*Request` structs
- Ensure nullable fields use pointer types (`*string`, `*int`, etc.)

**Example (project.go)**:
```go
type Project struct {
    ID          uuid.UUID  `json:"id"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      string     `json:"status"`
    Identifier  string     `json:"identifier"`    // NEW
    Homepage    *string    `json:"homepage,omitempty"`    // NEW
    IsPublic    bool       `json:"is_public"`         // NEW
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}
```

### Phase 2: Backend Repositories

**Goal**: Update SQL queries to handle new fields in CRUD operations.

**Files**:
- `backend/repositories/project_repository.go`
- `backend/repositories/task_repository.go`
- `backend/repositories/timelog_repository.go`

**Changes**:
- Update `INSERT` statements to include new fields
- Update `SELECT` statements to include new fields
- Update `UPDATE` statements to include new fields
- Handle NULL values for optional fields

### Phase 3: Backend Services

**Goal**: Add validation logic for new fields.

**Files**:
- `backend/services/project_service.go`
- `backend/services/task_service.go`

**Validation rules**:
- **Projects**:
  - `identifier` must be unique
  - `identifier` must match regex: `^[a-zA-Z0-9_-]+$`
  - `homepage` must be valid URL format if provided
- **Tasks**:
  - `due_date` >= `start_date` (if both provided)
  - `estimated_hours` >= 0
  - `done_ratio` between 0 and 100

### Phase 4: Backend Handlers

**Goal**: Process new fields from HTTP requests.

**Files**:
- `backend/handlers/project_handler.go`
- `backend/handlers/task_handler.go`
- `backend/handlers/timelog_handler.go`

**Changes**:
- Update `c.BodyParser()` calls to capture new fields
- Pass new fields to service layer
- Return validation errors for invalid data

### Phase 5: Frontend API

**Goal**: Update API client to send new fields.

**File**: `frontend/src/lib/api.js`

**Changes**:
- Update `createProject`, `updateProject` functions to send `identifier`, `homepage`, `is_public`
- Update `createTask`, `updateTask` functions to send all new task fields
- Update `createTimeLog`, `updateTimeLog` functions to send `user_id`, `activity_type`, `project_id`

### Phase 6: Frontend Stores

**Goal**: Update state management to handle new fields.

**Files**:
- `frontend/src/stores/projectStore.js`
- `frontend/src/stores/taskStore.js`
- `frontend/src/stores/timeLogStore.js`

**Changes**:
- Include new fields in store state
- Ensure new fields are passed through to API

### Phase 7: Frontend Components

**Goal**: Update UI to display and capture new fields.

**Files**:
- `frontend/src/components/ProjectForm.svelte`
- `frontend/src/components/ProjectList.svelte`
- `frontend/src/components/TaskForm.svelte`
- `frontend/src/components/TaskList.svelte`
- `frontend/src/components/TimeLogForm.svelte`

**Changes**:
- Add form fields for new inputs
- Display new fields in list views
- Add validation feedback for invalid inputs
- Use appropriate input types (date pickers, select dropdowns, checkboxes)

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A | N/A | N/A |

## Migration Strategy

### Data Migration

**Current data**: Existing projects, tasks, and time logs have no values for new fields.

**Default values**:
- `projects.identifier`: Auto-generate from title (lowercase, replace spaces with hyphens)
- `projects.homepage`: NULL
- `projects.is_public`: FALSE
- `tasks.description`: NULL
- `tasks.assignee_id`: NULL
- `tasks.author_id`: NULL
- `tasks.category`: NULL
- `tasks.start_date`: NULL
- `tasks.due_date`: NULL
- `tasks.estimated_hours`: NULL
- `tasks.done_ratio`: 0
- `time_logs.user_id`: NULL
- `time_logs.activity_type`: "General"
- `time_logs.project_id`: NULL (keep existing task_id association)

**Backward compatibility**: All new fields are optional (nullable) or have default values. Existing data remains valid and functional.

## Rollback Plan

If issues arise after deployment:

1. **Database**: Revert by dropping new columns:
   ```sql
   ALTER TABLE projects DROP COLUMN identifier, DROP COLUMN homepage, DROP COLUMN is_public;
   ALTER TABLE tasks DROP COLUMN description, DROP COLUMN assignee_id, DROP COLUMN author_id,
                DROP COLUMN category, DROP COLUMN start_date, DROP COLUMN due_date,
                DROP COLUMN estimated_hours, DROP COLUMN done_ratio;
   ALTER TABLE time_logs DROP COLUMN user_id, DROP COLUMN activity_type, DROP COLUMN project_id;
   ```
2. **Code**: Revert Git commit to pre-feature state
3. **Frontend**: Revert to previous build

## Testing Checklist

- [ ] Database migration runs without errors
- [ ] All existing data is accessible after migration
- [ ] Create project with identifier, homepage, is_public → success
- [ ] Create two projects with same identifier → validation error
- [ ] Create task with all new fields → success
- [ ] Create task with due_date < start_date → validation error
- [ ] Create task with done_ratio > 100 → validation error
- [ ] Update task progress from 0% to 50% → success
- [ ] Create time log with activity_type → success
- [ ] Create time log with project_id (no task) → success
- [ ] Frontend forms display all new fields
- [ ] Frontend validation shows appropriate error messages
- [ ] List views display new fields correctly

## Dependencies & Blocking

**No external dependencies** - this feature only uses existing Go and Node.js packages.

**Blocking requirements**:
- None - this is a standalone enhancement

## Timeline Estimates

- Phase 0 (Database): 0.5 hours
- Phase 1 (Models): 0.5 hours
- Phase 2 (Repositories): 1 hour
- Phase 3 (Services): 1 hour
- Phase 4 (Handlers): 1 hour
- Phase 5 (Frontend API): 0.5 hours
- Phase 6 (Frontend Stores): 0.5 hours
- Phase 7 (Frontend Components): 2 hours

**Total estimated time**: 7 hours

## Next Steps

After this plan is reviewed and approved:
1. Create database migration script (`backend/migration/001_enhance_entities.sql`)
2. Update backend models, repositories, services, handlers
3. Update frontend API, stores, and components
4. Test all changes manually (or add automated tests if time permits)
5. Create PR and merge to main branch
