---
description: "Task list for Enhance Entities with Redmine Fields feature"
---

# Tasks: Enhance Entities with Redmine Fields

**Input**: Design documents from `/specs/001-enhance-entities-with-redmine-fields/`
**Prerequisites**: plan.md (required), spec.md (required for user stories)

**Tests**: Tests are NOT included in this task list (no test framework configured per plan.md)

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

This is a web application with:

- Backend: `backend/` (Go + Fiber + PostgreSQL)
- Frontend: `frontend/` (Svelte + Vite + Tailwind CSS)

---

## Phase 1: Setup (Database Migration)

**Purpose**: Add new columns to existing tables without data loss

- [ ] T001 Create migration directory in backend/migration/
- [ ] T002 Create migration script backend/migration/001_enhance_entities.sql with all ALTER TABLE statements
- [ ] T003 Add identifier, homepage, is_public columns to projects table in backend/migration/001_enhance_entities.sql
- [ ] T004 Add description, assignee_id, author_id, category, start_date, due_date, estimated_hours, done_ratio columns to tasks table in backend/migration/001_enhance_entities.sql
- [ ] T005 Add user_id, activity_type, project_id columns to time_logs table in backend/migration/001_enhance_entities.sql
- [ ] T006 Add indexes (idx_projects_identifier, idx_tasks_assignee_id, idx_tasks_due_date, idx_time_logs_user_id, idx_time_logs_activity_type) in backend/migration/001_enhance_entities.sql
- [ ] T007 Add constraints (unique identifier, done_ratio check 0-100) in backend/migration/001_enhance_entities.sql
- [ ] T008 Run migration script and verify schema changes with \d commands

**Checkpoint**: Database schema ready - backend model updates can now begin

---

## Phase 2: User Story 1 - Enhanced Project Management (Priority: P1) 🎯 MVP

**Goal**: Enable projects to have unique identifiers, homepage URLs, and public visibility settings for better organization and sharing

**Independent Test**: Create a project with identifier "proj-001", homepage "https://example.com", and is_public=true. Verify it's stored correctly. Try creating another project with same identifier and verify validation error.

### Backend Models for User Story 1

- [ ] T009 [P] [US1] Update Project struct in backend/models/project.go to add Identifier (string), Homepage (\*string), IsPublic (bool) fields with JSON tags
- [ ] T010 [P] [US1] Update CreateProjectRequest struct in backend/models/project.go to include identifier, homepage, is_public fields
- [ ] T011 [P] [US1] Update UpdateProjectRequest struct in backend/models/project.go to include identifier, homepage, is_public fields

### Backend Repositories for User Story 1

- [ ] T012 [US1] Update Create function in backend/repositories/project_repository.go INSERT statement to include identifier, homepage, is_public columns
- [ ] T013 [US1] Update GetAll function in backend/repositories/project_repository.go SELECT statement to include identifier, homepage, is_public columns
- [ ] T014 [US1] Update GetByID function in backend/repositories/project_repository.go SELECT statement to include identifier, homepage, is_public columns
- [ ] T015 [US1] Update Update function in backend/repositories/project_repository.go to handle identifier, homepage, is_public columns

### Backend Services for User Story 1

- [ ] T016 [US1] Add ValidateProjectIdentifier function in backend/services/project*service.go to check uniqueness and regex pattern ^[a-zA-Z0-9*-]+$
- [ ] T017 [US1] Add ValidateHomepageURL function in backend/services/project_service.go to validate URL format if provided
- [ ] T018 [US1] Update CreateProject function in backend/services/project_service.go to call validation for identifier and homepage
- [ ] T019 [US1] Update UpdateProject function in backend/services/project_service.go to call validation for identifier and homepage

### Backend Handlers for User Story 1

- [ ] T020 [US1] Update CreateProject handler in backend/handlers/project_handler.go to parse identifier, homepage, is_public from request body
- [ ] T021 [US1] Update UpdateProject handler in backend/handlers/project_handler.go to parse identifier, homepage, is_public from request body
- [ ] T022 [US1] Add error handling in backend/handlers/project_handler.go for validation errors (duplicate identifier, invalid URL)

### Frontend for User Story 1

- [ ] T023 [P] [US1] Update projectStore in frontend/src/stores/projectStore.js to include identifier, homepage, isPublic in state
- [ ] T024 [P] [US1] Update createProject and updateProject functions in frontend/src/lib/api.js to send identifier, homepage, is_public fields
- [ ] T025 [US1] Add identifier input field (text) to ProjectForm.svelte with validation for alphanumeric, underscore, hyphen only
- [ ] T026 [US1] Add homepage input field (URL type) to ProjectForm.svelte with optional validation
- [ ] T027 [US1] Add is_public checkbox to ProjectForm.svelte
- [ ] T028 [US1] Update ProjectList.svelte to display identifier, homepage link (if present), and public badge
- [ ] T029 [US1] Add client-side validation feedback in ProjectForm.svelte for invalid identifiers and URLs

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. Projects can be created with identifier, homepage, and public visibility.

---

## Phase 3: User Story 2 - Enhanced Task Management (Priority: P1) 🎯 MVP

**Goal**: Track detailed task information including assignee, category, dates, estimates, and progress for better task management

**Independent Test**: Create a task with description, category "Backend", start_date, due_date, estimated_hours=8.5, and done_ratio=50. Verify all fields are stored. Try setting done_ratio=150 and verify validation error. Try setting due_date before start_date and verify validation error.

### Backend Models for User Story 2

- [ ] T030 [P] [US2] Update Task struct in backend/models/task.go to add Description (string), AssigneeID (*uuid.UUID), AuthorID (*uuid.UUID), Category (*string), StartDate (*time.Time), DueDate (*time.Time), EstimatedHours (*float64), DoneRatio (int) fields with JSON tags
- [ ] T031 [P] [US2] Update CreateTaskRequest struct in backend/models/task.go to include all new task fields
- [ ] T032 [P] [US2] Update UpdateTaskRequest struct in backend/models/task.go to include all new task fields

### Backend Repositories for User Story 2

- [ ] T033 [US2] Update Create function in backend/repositories/task_repository.go INSERT statement to include description, assignee_id, author_id, category, start_date, due_date, estimated_hours, done_ratio columns
- [ ] T034 [US2] Update GetAll function in backend/repositories/task_repository.go SELECT statement to include all new task columns
- [ ] T035 [US2] Update GetByID function in backend/repositories/task_repository.go SELECT statement to include all new task columns
- [ ] T036 [US2] Update GetByProjectID function in backend/repositories/task_repository.go SELECT statement to include all new task columns
- [ ] T037 [US2] Update Update function in backend/repositories/task_repository.go to handle all new task columns with NULL handling for optional fields

### Backend Services for User Story 2

- [ ] T038 [US2] Add ValidateTaskDates function in backend/services/task_service.go to ensure due_date >= start_date when both provided
- [ ] T039 [US2] Add ValidateDoneRatio function in backend/services/task_service.go to ensure 0 <= done_ratio <= 100
- [ ] T040 [US2] Add ValidateEstimatedHours function in backend/services/task_service.go to ensure estimated_hours >= 0 if provided
- [ ] T041 [US2] Update CreateTask function in backend/services/task_service.go to call all validation functions
- [ ] T042 [US2] Update UpdateTask function in backend/services/task_service.go to call all validation functions

### Backend Handlers for User Story 2

- [ ] T043 [US2] Update CreateTask handler in backend/handlers/task_handler.go to parse description, assignee_id, author_id, category, start_date, due_date, estimated_hours, done_ratio from request body
- [ ] T044 [US2] Update UpdateTask handler in backend/handlers/task_handler.go to parse all new task fields from request body
- [ ] T045 [US2] Add error handling in backend/handlers/task_handler.go for validation errors (invalid dates, invalid done_ratio, invalid estimated_hours)

### Frontend for User Story 2

- [ ] T046 [P] [US2] Update taskStore in frontend/src/stores/taskStore.js to include description, assigneeId, authorId, category, startDate, dueDate, estimatedHours, doneRatio in state
- [ ] T047 [P] [US2] Update createTask and updateTask functions in frontend/src/lib/api.js to send all new task fields
- [ ] T048 [US2] Add description textarea to TaskForm.svelte with proper styling
- [ ] T049 [US2] Add category input field (text or select dropdown) to TaskForm.svelte
- [ ] T050 [US2] Add start_date and due_date date picker inputs to TaskForm.svelte with date range validation
- [ ] T051 [US2] Add estimated_hours number input to TaskForm.svelte with min=0 validation
- [ ] T052 [US2] Add done_ratio slider or number input (0-100) to TaskForm.svelte with range validation
- [ ] T053 [US2] Update TaskList.svelte to display category, dates, estimated hours, and progress bar based on done_ratio
- [ ] T054 [US2] Add client-side validation feedback in TaskForm.svelte for invalid date ranges and numeric values

**Checkpoint**: At this point, User Story 2 should be fully functional and testable independently. Tasks can be created with full metadata including dates, estimates, and progress tracking.

---

## Phase 4: User Story 3 - Enhanced Time Logging (Priority: P2)

**Goal**: Log time with activity types and user attribution for more accurate and categorizable time tracking

**Independent Test**: Create a time log on a task with user_id and activity_type="Development". Create another with activity_type="Testing". Filter time logs by activity type and verify correct results.

### Backend Models for User Story 3

- [ ] T055 [P] [US3] Update TimeLog struct in backend/models/timelog.go to add UserID (*uuid.UUID), ActivityType (*string) fields with JSON tags
- [ ] T056 [P] [US3] Update CreateTimeLogRequest struct in backend/models/timelog.go to include user_id, activity_type fields
- [ ] T057 [P] [US3] Update UpdateTimeLogRequest struct in backend/models/timelog.go to include user_id, activity_type fields

### Backend Repositories for User Story 3

- [ ] T058 [US3] Update Create function in backend/repositories/timelog_repository.go INSERT statement to include user_id, activity_type columns
- [ ] T059 [US3] Update GetAll function in backend/repositories/timelog_repository.go SELECT statement to include user_id, activity_type columns
- [ ] T060 [US3] Update GetByID function in backend/repositories/timelog_repository.go SELECT statement to include user_id, activity_type columns
- [ ] T061 [US3] Update GetByTaskID function in backend/repositories/timelog_repository.go SELECT statement to include user_id, activity_type columns
- [ ] T062 [US3] Update Update function in backend/repositories/timelog_repository.go to handle user_id, activity_type columns with NULL handling

### Backend Services for User Story 3

- [ ] T063 [US3] Add ValidateActivityType function in backend/services/timelog_service.go to validate activity_type against allowed values (optional: Development, Testing, Design, Documentation, etc.)
- [ ] T064 [US3] Update CreateTimeLog function in backend/services/timelog_service.go to handle user_id and activity_type
- [ ] T065 [US3] Update UpdateTimeLog function in backend/services/timelog_service.go to handle user_id and activity_type

### Backend Handlers for User Story 3

- [ ] T066 [US3] Update CreateTimeLog handler in backend/handlers/timelog_handler.go to parse user_id, activity_type from request body
- [ ] T067 [US3] Update UpdateTimeLog handler in backend/handlers/timelog_handler.go to parse user_id, activity_type from request body

### Frontend for User Story 3

- [ ] T068 [P] [US3] Update timeLogStore in frontend/src/stores/timeLogStore.js to include userId, activityType in state
- [ ] T069 [P] [US3] Update createTimeLog and updateTimeLog functions in frontend/src/lib/api.js to send user_id, activity_type fields
- [ ] T070 [US3] Add activity_type select dropdown to TimeLogForm.svelte with predefined options (Development, Testing, Design, Documentation, etc.)
- [ ] T071 [US3] Add user_id selector to TimeLogForm.svelte (or auto-populate from current user if authentication exists)
- [ ] T072 [US3] Update time log display components to show activity_type badge and user information

**Checkpoint**: At this point, User Story 3 should be fully functional and testable independently. Time logs can be created with activity types and user attribution.

---

## Phase 5: User Story 4 - Project-Level Time Tracking (Priority: P3)

**Goal**: Log time at the project level (not just task level) for tracking overhead and administrative time

**Independent Test**: Create a time log with project_id but no task_id. Verify it's stored correctly. View project time reports and verify both project-level and task-level time are aggregated.

### Backend Models for User Story 4

- [ ] T073 [P] [US4] Update TimeLog struct in backend/models/timelog.go to add ProjectID (\*uuid.UUID) field with JSON tag
- [ ] T074 [P] [US4] Update CreateTimeLogRequest struct in backend/models/timelog.go to include project_id field
- [ ] T075 [P] [US4] Update UpdateTimeLogRequest struct in backend/models/timelog.go to include project_id field

### Backend Repositories for User Story 4

- [ ] T076 [US4] Update Create function in backend/repositories/timelog_repository.go INSERT statement to include project_id column
- [ ] T077 [US4] Update GetAll function in backend/repositories/timelog_repository.go SELECT statement to include project_id column
- [ ] T078 [US4] Update GetByID function in backend/repositories/timelog_repository.go SELECT statement to include project_id column
- [ ] T079 [US4] Add GetByProjectID function in backend/repositories/timelog_repository.go to retrieve project-level time logs
- [ ] T080 [US4] Update Update function in backend/repositories/timelog_repository.go to handle project_id column with NULL handling

### Backend Services for User Story 4

- [ ] T081 [US4] Update ValidateTimeLog function in backend/services/timelog_service.go to allow either task_id OR project_id (at least one required)
- [ ] T082 [US4] Add GetProjectTimeLogs function in backend/services/timelog_service.go to aggregate project and task-level time logs
- [ ] T083 [US4] Update CreateTimeLog function in backend/services/timelog_service.go to handle optional project_id

### Backend Handlers for User Story 4

- [ ] T084 [US4] Update CreateTimeLog handler in backend/handlers/timelog_handler.go to parse project_id from request body
- [ ] T085 [US4] Update UpdateTimeLog handler in backend/handlers/timelog_handler.go to parse project_id from request body
- [ ] T086 [US4] Add GetProjectTimeLogs handler in backend/handlers/timelog_handler.go for project time aggregation endpoint

### Backend Routes for User Story 4

- [ ] T087 [US4] Add GET /projects/:id/timelogs route in backend/routes/routes.go to retrieve project-level time logs

### Frontend for User Story 4

- [ ] T088 [P] [US4] Update timeLogStore in frontend/src/stores/timeLogStore.js to include projectId in state
- [ ] T089 [P] [US4] Update createTimeLog and updateTimeLog functions in frontend/src/lib/api.js to send project_id field
- [ ] T090 [US4] Update TimeLogForm.svelte to allow selecting project OR task (make task optional when project is selected)
- [ ] T091 [US4] Add project time report view to display aggregated project and task-level time logs
- [ ] T092 [US4] Add validation in TimeLogForm.svelte to ensure at least one of project_id or task_id is provided

**Checkpoint**: At this point, User Story 4 should be fully functional and testable independently. Time can be logged at project level and aggregated reports work correctly.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T093 [P] Add data migration script to auto-generate identifiers for existing projects in backend/migration/001_enhance_entities.sql
- [ ] T094 [P] Update README.md with new field descriptions and usage examples
- [ ] T095 Code cleanup and refactoring across backend and frontend files
- [ ] T096 Add comments and documentation to new validation functions
- [ ] T097 [P] Update API documentation (if exists) with new field specifications
- [ ] T098 Performance testing with large datasets to ensure indexes are effective
- [ ] T099 Manual testing of all user stories end-to-end
- [ ] T100 Run full application smoke test following existing workflows

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately. Database migration MUST complete before any backend work.
- **User Stories (Phase 2-5)**: All depend on Phase 1 (database migration) completion
  - User stories can proceed in parallel (if staffed) once database is ready
  - Or sequentially in priority order (US1 → US2 → US3 → US4)
  - US1 and US2 are both P1 priority and recommended for MVP
  - US3 and US4 add additional value but not blocking
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Phase 1 database migration - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Phase 1 database migration - No dependencies on other stories (independent from US1)
- **User Story 3 (P2)**: Can start after Phase 1 database migration - No dependencies on other stories (works with existing tasks)
- **User Story 4 (P3)**: Can start after Phase 1 database migration - No dependencies on other stories (extends time logging)

**KEY INSIGHT**: All user stories are independent after database migration. Each enhances a different entity and can be implemented/tested separately.

### Within Each User Story

Standard flow for each story:

1. Models first (can be parallel if multiple entities)
2. Repositories next (depend on models)
3. Services next (depend on repositories)
4. Handlers next (depend on services)
5. Frontend stores and API (can be parallel)
6. Frontend components last (depend on stores/API)

### Parallel Opportunities

**Phase 1 (Setup)**:

- T003, T004, T005 can be written in parallel (different SQL sections)
- T006, T007 can be done in parallel (different SQL sections)

**After Phase 1 completes, ALL user stories can start in parallel**:

- Team Member A: User Story 1 (Projects)
- Team Member B: User Story 2 (Tasks)
- Team Member C: User Story 3 (Time Logs - activity types)
- Team Member D: User Story 4 (Time Logs - project level)

**Within User Story 1**:

- T009, T010, T011 (models) can run in parallel
- T023, T024 (frontend stores/API) can run in parallel

**Within User Story 2**:

- T030, T031, T032 (models) can run in parallel
- T046, T047 (frontend stores/API) can run in parallel

**Within User Story 3**:

- T055, T056, T057 (models) can run in parallel
- T068, T069 (frontend stores/API) can run in parallel

**Within User Story 4**:

- T073, T074, T075 (models) can run in parallel
- T088, T089 (frontend stores/API) can run in parallel

**Phase 6 (Polish)**:

- T093, T094, T097 can run in parallel (different files)

---

## Parallel Example: User Story 2 (Tasks)

```bash
# Launch all models for User Story 2 together:
Task: "Update Task struct in backend/models/task.go to add all new fields"
Task: "Update CreateTaskRequest struct in backend/models/task.go"
Task: "Update UpdateTaskRequest struct in backend/models/task.go"

# Launch frontend stores and API together:
Task: "Update taskStore in frontend/src/stores/taskStore.js"
Task: "Update createTask and updateTask functions in frontend/src/lib/api.js"
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2)

Recommended MVP includes both P1 priority stories for maximum value:

1. Complete Phase 1: Database migration (CRITICAL - blocks all work)
2. Complete Phase 2: User Story 1 (Enhanced Projects)
3. **VALIDATE US1**: Test project creation with identifier, homepage, public flag
4. Complete Phase 3: User Story 2 (Enhanced Tasks)
5. **VALIDATE US2**: Test task creation with dates, estimates, progress
6. **STOP and VALIDATE MVP**: Test both stories independently and together
7. Deploy/demo if ready

**MVP Scope**: 54 tasks (T001-T054)
**MVP Value**: Complete enhanced project and task management with metadata tracking

### Incremental Delivery

1. Phase 1: Database migration → Schema ready
2. Phase 2: User Story 1 → Test independently → Deploy/Demo (Projects enhanced!)
3. Phase 3: User Story 2 → Test independently → Deploy/Demo (Tasks enhanced!)
4. Phase 4: User Story 3 → Test independently → Deploy/Demo (Time tracking enhanced!)
5. Phase 5: User Story 4 → Test independently → Deploy/Demo (Project-level time tracking!)
6. Phase 6: Polish → Final release

Each story adds value without breaking previous stories.

### Parallel Team Strategy

With multiple developers:

1. Complete Phase 1 (Database) together
2. Once Phase 1 is done:
   - Developer A: User Story 1 (Projects) - T009 through T029
   - Developer B: User Story 2 (Tasks) - T030 through T054
   - Developer C: User Story 3 (Time Logs - activity) - T055 through T072
   - Developer D: User Story 4 (Time Logs - project) - T073 through T092
3. Stories complete and integrate independently
4. Team completes Phase 6 (Polish) together

---

## Summary

**Total Tasks**: 100 tasks

- Phase 1 (Setup): 8 tasks
- Phase 2 (User Story 1 - Projects): 21 tasks
- Phase 3 (User Story 2 - Tasks): 25 tasks
- Phase 4 (User Story 3 - Time Logs Activity): 18 tasks
- Phase 5 (User Story 4 - Time Logs Project): 20 tasks
- Phase 6 (Polish): 8 tasks

**Task Distribution by Story**:

- User Story 1 (Enhanced Projects): 21 tasks
- User Story 2 (Enhanced Tasks): 25 tasks
- User Story 3 (Enhanced Time Logging): 18 tasks
- User Story 4 (Project-Level Time Tracking): 20 tasks
- Infrastructure (Setup + Polish): 16 tasks

**Parallel Opportunities**: 27 tasks marked [P] can run in parallel within their phase

**MVP Recommendation**: Complete Phase 1 + User Stories 1 & 2 (54 tasks total)

**Estimated Timeline**:

- Phase 1: 0.5 hours
- User Story 1: 2 hours
- User Story 2: 2.5 hours
- User Story 3: 1.5 hours
- User Story 4: 1.5 hours
- Polish: 1 hour
- **Total: ~9 hours** (sequential) or **~4 hours** (with 4-person parallel team)

---

## Notes

- [P] tasks = different files, no dependencies - safe to parallelize
- [Story] label maps task to specific user story for traceability and independent delivery
- Each user story is independently completable and testable - no cross-story blocking dependencies
- Database migration (Phase 1) is the only true blocker - must complete before any backend work
- All validation follows spec.md requirements (FR-001 through FR-017)
- No tests included per plan.md (no test framework configured)
- Commit after each task or logical group of related tasks
- Stop at any checkpoint to validate story independently before proceeding
