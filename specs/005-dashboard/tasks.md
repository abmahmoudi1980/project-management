# Tasks: Project Manager Dashboard

**Input**: Design documents from `/specs/005-dashboard/`  
**Branch**: `005-dashboard`  
**Prerequisites**: plan.md ✅, spec.md ✅, research.md ✅, data-model.md ✅, contracts/ ✅

**Note**: Tests are NOT included (not explicitly requested in feature specification, project uses manual testing). Tasks follow strict checklist format with file paths.

**Organization**: Dashboard is a single, unified feature with no independent user stories to partition. Tasks are organized by implementation layers (database → backend → frontend) following the established architecture pattern.

---

## Format: `- [ ] [ID] [P?] Description with file path`

- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions
- All tasks sequentially numbered T001-T045
- No [Story] label (single feature, not multiple stories)

---

## Phase 1: Setup & Project Configuration

**Purpose**: Prepare project structure and dependencies for dashboard implementation

- [ ] T001 Create feature branch `005-dashboard` (already done: verify with `git branch --show-current`)
- [ ] T002 [P] Create spec directory structure in `/specs/005-dashboard/` (already done: verify with `ls -la specs/005-dashboard/`)
- [ ] T003 Review and understand existing backend structure in `/backend/` (handlers, services, repositories, models)
- [ ] T004 Review and understand existing frontend structure in `/frontend/src/` (components, stores, lib)
- [ ] T005 Ensure PostgreSQL is running and accessible (`psql -U your_user -d your_database`)
- [ ] T006 Verify Go version >= 1.21 with `go version`
- [ ] T007 Verify Node.js version >= 18 with `node --version`
- [ ] T008 Review AGENTS.md for project conventions and patterns

---

## Phase 2: Foundational - Database Setup (BLOCKING)

**Purpose**: Create database tables and indexes that all backend components depend on

**⚠️ CRITICAL**: No backend work can begin until database migration is complete

- [ ] T009 Create migration file `backend/migration/005_add_dashboard_meetings.sql` with SQL from data-model.md
  - Include: meetings table, meeting_attendees junction table, 8 indexes
  - Reference: data-model.md Database Migration section
- [ ] T010 Run migration with `cd backend && go run ./cmd/migrate`
- [ ] T011 Verify migration success:
  - Connect to database: `psql -U user -d database`
  - Check tables: `\dt meetings*` and `\dt meeting_attendees*`
  - Check indexes: `\di idx_meetings*` and `\di idx_meeting_attendees*`
  - Check indexes on existing tables: `\di idx_projects_status*`, `\di idx_tasks_assigned*`
- [ ] T012 Document any migration issues or deviations in implementation notes

**Checkpoint**: Database schema ready - all tables and indexes created successfully

---

## Phase 3: Backend - Models & Data Structures

**Purpose**: Define Go structs representing dashboard data

- [ ] T013 Create `backend/models/meeting.go` with Meeting and MeetingAttendee structs
  - Fields: id (UUID), title, description, meeting_date, duration_minutes, project_id, created_by, timestamps
  - Follow existing struct patterns (PascalCase fields, json tags, pgx types)
  - Reference: data-model.md Entity Definitions section
- [ ] T014 Create `backend/models/dashboard.go` with Dashboard DTOs:
  - DashboardResponse, DashboardStatistics, StatValue, ProjectCard, TaskSummary
  - Follow existing model patterns and conventions
  - Reference: data-model.md Dashboard Aggregate Data section
- [ ] T015 Verify models compile: `cd backend && go build ./models`
- [ ] T016 Review models for consistency with existing patterns

---

## Phase 4: Backend - Repositories (Data Access)

**Purpose**: Create data access layer for dashboard and meeting queries

- [ ] T017 Create `backend/repositories/dashboard_repository.go`:
  - NewDashboardRepository(db) constructor
  - GetStatistics(ctx, userID, userRole) → DashboardStatistics
  - GetRecentProjects(ctx, userID, userRole, limit) → []ProjectCard
  - GetUserTasks(ctx, userID, limit) → []TaskSummary
  - Reference: plan.md Phase 4 Backend Services section for query logic
- [ ] T018 Create `backend/repositories/meeting_repository.go`:
  - NewMeetingRepository(db) constructor
  - CreateMeeting(ctx, meeting) error
  - GetNextMeetingForUser(ctx, userID) → *MeetingWithAttendees
  - GetMeetingByID(ctx, meetingID) → *MeetingWithAttendees
  - AddAttendees(ctx, meetingID, userIDs) error
  - ListMeetings(ctx, userID, from, to, limit, offset) → ([]Meeting, error)
  - Reference: data-model.md Entity Definitions for repository methods
- [ ] T019 [P] Add indexes on existing tables in repositories if not in migration:
  - Verify: idx_projects_status_updated, idx_tasks_assigned_status, idx_tasks_due_date, idx_users_active
  - Add to dashboard queries WHERE/ORDER BY clauses if missing
- [ ] T020 Test repository methods compile and build: `cd backend && go build ./repositories`
- [ ] T021 Review repository queries for SQL injection prevention (use parameterized queries)

---

## Phase 5: Backend - Services (Business Logic)

**Purpose**: Implement business logic and validation for dashboard and meetings

- [ ] T022 Create `backend/services/dashboard_service.go`:
  - NewDashboardService(dashboardRepo, meetingRepo) constructor
  - GetDashboardData(ctx, userID, userRole) → *DashboardResponse
  - Calls 4 repository methods and aggregates into DashboardResponse
  - Handles nil/empty responses gracefully
  - Reference: plan.md Data Flow Dashboard Load Sequence
- [ ] T023 Create `backend/services/meeting_service.go`:
  - NewMeetingService(meetingRepo, userRepo) constructor
  - GetNextMeetingForUser(ctx, userID) → *MeetingWithAttendees
  - CreateMeeting(ctx, userID, input) → *MeetingWithAttendees
  - CreateMeetingInput struct with validation rules
  - validateMeetingInput(input) function
  - Verify attendees exist before creating
  - Add meeting creator to attendees automatically
  - Reference: data-model.md Validation Rules section
- [ ] T024 Implement validation functions:
  - Meeting title: 1-200 characters
  - Description: 0-5000 characters
  - Duration: 1-1440 minutes
  - Meeting date: must be in future
  - At least 1 attendee required
- [ ] T025 Implement error handling:
  - Return descriptive error messages
  - Wrap errors with context when needed
  - Handle pgx.ErrNoRows gracefully
- [ ] T026 Test service methods compile: `cd backend && go build ./services`
- [ ] T027 Review services for business logic correctness

---

## Phase 6: Backend - Handlers (HTTP Layer)

**Purpose**: Create HTTP endpoints for dashboard and meetings

- [ ] T028 Create `backend/handlers/dashboard_handler.go`:
  - NewDashboardHandler(dashboardService) constructor
  - GetDashboard(c *fiber.Ctx) error handler
  - Extract user from JWT (c.Locals("user"))
  - Call service.GetDashboardData(ctx, user.ID, user.Role)
  - Return JSON response with proper error handling
  - HTTP status: 200 OK or 500 Internal Server Error
  - Reference: api-contract.md GET /api/dashboard section
- [ ] T029 Create `backend/handlers/meeting_handler.go`:
  - NewMeetingHandler(meetingService) constructor
  - GetNextMeeting(c *fiber.Ctx) - GET /api/meetings/next
  - CreateMeeting(c *fiber.Ctx) - POST /api/meetings
  - ListMeetings(c *fiber.Ctx) - GET /api/meetings
  - GetMeeting(c *fiber.Ctx) - GET /api/meetings/:id
  - Parse JSON body, validate input, handle errors
  - Reference: api-contract.md Meeting Endpoints sections
- [ ] T030 Implement error response formatting:
  - Consistent error format: {error, message, field}
  - Proper HTTP status codes (400, 401, 403, 404, 500)
  - Reference: api-contract.md Error Handling section
- [ ] T031 Test handlers compile: `cd backend && go build ./handlers`
- [ ] T032 Review handlers for HTTP concern separation

---

## Phase 7: Backend - Routes & Wiring

**Purpose**: Register new endpoints and wire dependencies

- [ ] T033 Update `backend/routes/routes.go`:
  - Create dashboardHandler and meetingHandler instances
  - Register GET /api/dashboard with auth middleware
  - Register GET /api/meetings/next with auth middleware
  - Register POST /api/meetings with auth middleware
  - Register GET /api/meetings with auth middleware
  - Register GET /api/meetings/:id with auth middleware
  - Reference: plan.md Phase 7 Backend Handlers & Routes section
- [ ] T034 Update `backend/main.go`:
  - Initialize DashboardRepository with db pool
  - Initialize MeetingRepository with db pool
  - Initialize DashboardService with repositories
  - Initialize MeetingService with repositories
  - Pass to SetupRoutes() function
  - Verify all dependencies injected correctly
- [ ] T035 Build backend: `cd backend && go build`
- [ ] T036 Verify backend starts: `cd backend && go run main.go` (CTRL+C to stop)
  - Check for startup errors in console
  - Verify "Server started on :3000" message

**Checkpoint**: Backend API endpoints ready and running

---

## Phase 8: Backend - API Testing

**Purpose**: Verify API contracts work correctly

- [ ] T037 Test dashboard endpoint with curl:
  ```bash
  curl -X GET http://localhost:3000/api/dashboard \
    -H "Cookie: access_token=YOUR_JWT_TOKEN" | jq .
  ```
  - Verify response has statistics, recent_projects, user_tasks, next_meeting
  - Verify all fields match api-contract.md structure
- [ ] T038 Test meeting creation:
  ```bash
  curl -X POST http://localhost:3000/api/meetings \
    -H "Cookie: access_token=YOUR_JWT_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"title":"Test","meeting_date":"2026-01-15T10:00:00Z","duration_minutes":60,"attendee_ids":["UUID"]}'
  ```
  - Verify 201 Created response
  - Verify meeting_id returned
- [ ] T039 Test get next meeting:
  ```bash
  curl -X GET http://localhost:3000/api/meetings/next \
    -H "Cookie: access_token=YOUR_JWT_TOKEN"
  ```
  - Verify 200 OK if meeting exists
  - Verify 204 No Content if no meetings
- [ ] T040 Test role-based filtering:
  - Create JWT for different roles (admin, project_manager, team_member)
  - Verify dashboard returns different data per role
  - Document any filtering issues found

---

## Phase 9: Frontend - API Clients

**Purpose**: Create JavaScript functions to call backend endpoints

- [ ] T041 Create `frontend/src/lib/api/dashboard.js`:
  - getDashboardData() → Promise<DashboardResponse>
  - Fetch GET /api/dashboard with credentials: 'include'
  - Return parsed JSON
  - Reference: api-contract.md GET /api/dashboard section
- [ ] T042 Create `frontend/src/lib/api/meetings.js`:
  - getNextMeeting() → Promise<Meeting|null>
  - createMeeting(data) → Promise<Meeting>
  - Fetch with proper headers and body
  - Handle 204 No Content response
  - Reference: api-contract.md Meeting Endpoints sections
- [ ] T043 Update `frontend/src/lib/utils.js`:
  - Add formatJalaliDate(date, format) function
  - Use jalali-moment library
  - Formats: 'full' (jYYYY/jMM/jDD), 'short' (jMM/jDD), 'time' (HH:mm)
  - Handle special cases: "امروز" (today), "فردا" (tomorrow)
  - Reference: quickstart.md Utility Functions section
- [ ] T044 Test API clients load without errors:
  - Import in browser console to verify syntax
  - Check for missing dependencies (jalali-moment)

---

## Phase 10: Frontend - Reusable Components

**Purpose**: Create Svelte 5 components with runes

- [ ] T045 [P] Create `frontend/src/components/Avatar.svelte`:
  - Props: user object, size (sm|md|lg)
  - Extract initials from user.full_name
  - Generate color from user.id hash
  - Render: div with initials, background color, rounded shape
  - Reference: quickstart.md Step 7 Avatar component
- [ ] T046 [P] Create `frontend/src/components/StatCard.svelte`:
  - Props: title, value, change, icon, iconColor
  - Display: icon, value, title, change indicator
  - Change color: green (+), red (-), gray (0)
  - Reference: quickstart.md Step 7 StatCard component
- [ ] T047 [P] Create `frontend/src/components/ProjectCard.svelte`:
  - Props: project object, onclick callback
  - Display: status badge, name, client, progress bar, team avatars, deadline
  - Status colors: planning (gray), in progress (blue), on track (green), review (purple)
  - Team avatars: max 3 visible, "+N" if more
  - Reference: quickstart.md Step 7 ProjectCard component
- [ ] T048 [P] Create `frontend/src/components/TaskListItem.svelte`:
  - Props: task object, onComplete callback
  - Display: checkbox, title, project name, priority badge, due date
  - Priority colors: critical (red), high (orange), medium (blue), low (gray)
  - Checkbox triggers onComplete(taskId)
  - Completed state: strikethrough + fade out
  - Reference: quickstart.md Step 7 TaskListItem component
- [ ] T049 [P] Create `frontend/src/components/MeetingCard.svelte`:
  - Props: meeting object (nullable)
  - Display: title, description, time, attendee avatars (max 3)
  - Gradient background: indigo-600 to purple-700
  - Hidden if meeting is null
  - Reference: quickstart.md Step 7 MeetingCard component
- [ ] T050 Test all components render without errors:
  - Import each in a test Svelte file
  - Check browser console for errors
  - Verify props are defined with $props()

---

## Phase 11: Frontend - Dashboard Page

**Purpose**: Create main dashboard page component

- [ ] T051 Create `frontend/src/components/Dashboard.svelte`:
  - State: dashboardData, loading, error
  - $effect hook: calls loadDashboard() on mount and every 30 seconds
  - loadDashboard() function: calls getDashboardData(), handles errors
  - handleTaskComplete(taskId) function: calls completeTask(), removes task after 2 seconds
  - navigateToProject(projectId) function: navigates to project details
  - Render: statistics grid, projects grid, tasks list, meeting card
  - Reference: quickstart.md Step 8 Frontend Dashboard Page section
- [ ] T052 Add error handling in Dashboard.svelte:
  - Show error banner if loading fails
  - Don't block UI for partial failures
  - Allow auto-refresh to retry every 30 seconds
- [ ] T053 Implement auto-refresh:
  - $effect creates setInterval for 30 seconds
  - Cleanup: return function that clears interval
  - No loading spinner (subtle update)
  - Reference: quickstart.md Data Flow section
- [ ] T054 Test Dashboard.svelte renders:
  - No console errors
  - All child components render
  - Auto-refresh interval working (check developer tools)

---

## Phase 12: Frontend - Routing & Integration

**Purpose**: Add dashboard to app routing and layout

- [ ] T055 Update `frontend/src/App.svelte`:
  - Add route for /dashboard pointing to Dashboard.svelte component
  - Make /dashboard the default route after login
  - Verify routing works with browser navigation
- [ ] T056 Update `frontend/src/components/Navigation.svelte` (if exists):
  - Add "Dashboard" link in sidebar/menu
  - Highlight as active when on /dashboard
  - Point to /dashboard route
- [ ] T057 Update frontend package.json if needed:
  - Verify jalali-moment is installed: `npm ls jalali-moment`
  - Verify lucide-svelte (or lucide) is installed for icons
  - Run `npm install` if any dependencies missing
- [ ] T058 Test frontend navigation:
  - Start dev server: `cd frontend && npm run dev`
  - Navigate to /dashboard
  - Dashboard loads without errors
  - No console errors in browser

---

## Phase 13: Frontend - Styling & Responsive Design

**Purpose**: Ensure dashboard looks correct and is responsive

- [ ] T059 Verify Tailwind CSS classes used in components:
  - Check color classes: text-indigo-600, bg-blue-100, etc. match Tailwind docs
  - Grid layouts: grid-cols-1, grid-cols-4, gap-6
  - Responsive: md:grid-cols-2, lg:grid-cols-4
  - Reference: dashboard.html reference design for layout
- [ ] T060 Test responsive design:
  - Desktop (1920px): All 4 stat cards in row, 2x2 projects grid
  - Tablet (768px): 2 stat cards per row, stacked projects
  - Mobile (320px): Stacked everything, single column
  - Use browser devtools responsive mode
- [ ] T061 Verify Lucide icons render:
  - Icons appear with correct colors
  - Icons scale properly with size
  - Icons missing → fallback symbols work
- [ ] T062 Test dark text on light backgrounds (accessibility):
  - Color contrast ratios pass WCAG AA
  - White icons on colored backgrounds visible
  - All text readable

---

## Phase 14: End-to-End Testing

**Purpose**: Verify complete feature works from backend to frontend

- [ ] T063 Start backend server:
  ```bash
  cd backend
  go run main.go
  ```
  - Verify "Server started on :3000" message
- [ ] T064 Start frontend dev server:
  ```bash
  cd frontend
  npm run dev
  ```
  - Verify "Local: http://localhost:5173/" message
- [ ] T065 Test complete user flows:
  - **Flow 1**: Login → Navigate to /dashboard → Dashboard loads with data
  - **Flow 2**: Dashboard displays → Auto-refresh every 30 seconds → New data appears
  - **Flow 3**: Click task checkbox → Task marks complete → Task list updates
  - **Flow 4**: Click project card → Navigate to project details
  - **Flow 5**: Task completion updates "Pending Tasks" stat → Verify counter decreases
- [ ] T066 Verify dashboard data accuracy:
  - Active Projects count matches projects with status In Progress or On Track
  - Pending Tasks count matches incomplete tasks
  - Team Members count matches active users
  - Upcoming Deadlines count matches tasks/projects due in next 7 days
  - Project progress = completed_tasks / total_tasks × 100
  - Task sorting: priority (Critical > High > Medium > Low) → due date
- [ ] T067 Verify role-based filtering:
  - **Admin**: See all projects, all teams' tasks (no filtering)
  - **Project Manager**: See managed projects, team members' tasks for those projects
  - **Team Member**: See assigned projects, only own tasks
  - Create test accounts for each role, verify filtering
- [ ] T068 Test edge cases:
  - No projects exist: Show empty state
  - No tasks assigned: Show "You're all caught up!" message
  - No meetings scheduled: Hide meeting card (not show empty card)
  - Long project names: Truncate with ellipsis (text-truncate)
  - Overdue tasks: Highlight in red with "Overdue" badge
  - 50+ tasks: Show only top 5 with "View All" link
- [ ] T069 Browser compatibility testing:
  - Chrome: Works correctly
  - Firefox: Works correctly
  - Safari: Works correctly
  - Edge: Works correctly
  - Mobile Safari (if available): Works correctly
- [ ] T070 Performance validation:
  - Dashboard loads in < 2 seconds (browser Network tab)
  - API response time < 500ms (curl with timing)
  - Auto-refresh doesn't cause flickering or lag
  - No memory leaks (developer tools Performance/Memory)

---

## Phase 15: Documentation & Finalization

**Purpose**: Document implementation and prepare for deployment

- [ ] T071 Update `backend/AGENTS.md`:
  - Add dashboard feature details
  - Document new handlers, services, repositories
  - Update component count in frontend section
  - Reference: AGENTS.md template sections
- [ ] T072 Update root `README.md`:
  - Add dashboard to feature list
  - Document how to access dashboard (/dashboard route)
  - Add quickstart section for dashboard feature
- [ ] T073 Add code comments:
  - Backend service: Explain data aggregation logic
  - Frontend Dashboard component: Explain auto-refresh pattern
  - API client functions: Explain error handling
- [ ] T074 Document any deviations:
  - If Persian number formatting not implemented (Western used instead)
  - If notification badge removed (out of scope)
  - If any features skipped or altered from spec
- [ ] T075 Run final verification:
  - `cd backend && go build` - compiles without errors
  - `cd frontend && npm run build` - builds without errors
  - No compiler warnings or linting errors
  - All tasks in this checklist are complete

---

## Phase 16: Git Commit & Review

**Purpose**: Prepare feature for merge

- [ ] T076 [P] Stage all backend changes:
  ```bash
  cd backend
  git add models/meeting.go models/dashboard.go \
           repositories/dashboard_repository.go repositories/meeting_repository.go \
           services/dashboard_service.go services/meeting_service.go \
           handlers/dashboard_handler.go handlers/meeting_handler.go \
           routes/routes.go main.go migration/005_add_dashboard_meetings.sql
  git commit -m "feat(dashboard): add backend models, repositories, services, handlers"
  ```
- [ ] T077 [P] Stage all frontend changes:
  ```bash
  cd frontend
  git add src/components/Avatar.svelte src/components/Dashboard.svelte \
           src/components/MeetingCard.svelte src/components/ProjectCard.svelte \
           src/components/StatCard.svelte src/components/TaskListItem.svelte \
           src/lib/api/dashboard.js src/lib/api/meetings.js \
           src/lib/utils.js src/App.svelte
  git commit -m "feat(dashboard): add frontend components and API clients"
  ```
- [ ] T078 [P] Stage documentation changes:
  ```bash
  git add AGENTS.md README.md specs/005-dashboard/
  git commit -m "docs(dashboard): update project documentation"
  ```
- [ ] T079 Verify git history:
  ```bash
  git log --oneline -5
  ```
  - All feature commits present
  - Messages are clear and descriptive
- [ ] T080 Create feature summary:
  - Total commits: 3 (backend, frontend, docs)
  - Files changed: ~20 backend, ~8 frontend, docs
  - Lines of code: ~800 backend, ~600 frontend
  - Database tables: 2 new (meetings, meeting_attendees)
  - API endpoints: 5 new (/api/dashboard, /api/meetings/*)
  - Components: 6 new (Avatar, Dashboard, MeetingCard, ProjectCard, StatCard, TaskListItem)

---

## Dependencies & Execution Order

### Phase Dependencies

1. **Phase 1 Setup**: No dependencies ✅
2. **Phase 2 Database**: Depends on Phase 1 ✅
3. **Phase 3 Models**: Depends on Phase 2 ✅
4. **Phase 4 Repositories**: Depends on Phase 3 ✅
5. **Phase 5 Services**: Depends on Phase 4 ✅
6. **Phase 6 Handlers**: Depends on Phase 5 ✅
7. **Phase 7 Routes**: Depends on Phase 6 ✅
8. **Phase 8 Backend Testing**: Depends on Phase 7 ✅
9. **Phase 9 Frontend API Clients**: Depends on Phase 8 ✅
10. **Phase 10 Frontend Components**: Can start after Phase 9 ✅
11. **Phase 11 Dashboard Page**: Depends on Phase 10 ✅
12. **Phase 12 Routing**: Depends on Phase 11 ✅
13. **Phase 13 Styling**: Depends on Phase 12 ✅
14. **Phase 14 E2E Testing**: Depends on Phase 13 ✅
15. **Phase 15 Documentation**: Depends on Phase 14 ✅
16. **Phase 16 Git Commit**: Depends on Phase 15 ✅

### Parallel Opportunities

**Phase 1 Setup**: All tasks independent
- T002, T003, T004, T005, T006, T007, T008 → Can run in parallel

**Phase 3 Models**: Create two files in parallel
- T013 (meeting.go) and T014 (dashboard.go) → Can run in parallel

**Phase 4 Repositories**: Create two files in parallel
- T017 (dashboard_repository.go) and T018 (meeting_repository.go) → Can run in parallel
- T019 (indexes) → Can check in parallel with above

**Phase 10 Components**: Create six components in parallel
- T045 (Avatar), T046 (StatCard), T047 (ProjectCard), T048 (TaskListItem), T049 (MeetingCard) → Can run in parallel
- Each component is independent

**Phase 16 Git Commit**: Three commits in parallel (if using multiple terminals)
- T076 (backend), T077 (frontend), T078 (docs) → Can stage in parallel

### Recommended Sequential Order

For single developer:

1. T001-T008: Setup phase (30 min)
2. T009-T012: Database setup (45 min) **BLOCKING**
3. T013-T040: Backend complete (6-7 hours)
   - Can pause after T036 for backend testing
4. T041-T044: Frontend API clients (45 min)
5. T045-T062: Frontend components & styling (4-5 hours)
6. T063-T070: Complete E2E testing (2 hours)
7. T071-T080: Documentation & cleanup (1 hour)

**Total estimated time**: 12-16 hours

---

## Parallel Example: Multiple Developers

If team can work in parallel after Phase 2:

```
Developer A (Backend):        Developer B (Frontend):
T013 meeting.go              T041 dashboard.js API client
T014 dashboard.go            T043 utils.js (Jalali dates)
T017 dashboard_repository    T045 Avatar.svelte
T018 meeting_repository      T046 StatCard.svelte
T022 dashboard_service       T047 ProjectCard.svelte
T023 meeting_service         T048 TaskListItem.svelte
T028 dashboard_handler       T049 MeetingCard.svelte
T029 meeting_handler         T051 Dashboard.svelte
T033 routes.go               T055 App.svelte routing
T034 main.go wiring
T037-T040 API testing

After Phase 2 Database complete (2 hours):
- Developer A: 5 hours on backend
- Developer B: 4 hours on frontend (can start after T041)
- Both: 2 hours on testing & documentation

Merge: All done in ~7 hours total (vs 16 sequential)
```

---

## Implementation Strategy

### MVP First (Recommended)

1. **Phase 1**: Setup (30 min)
2. **Phase 2**: Database (45 min)
3. **Phase 3-8**: Complete backend (7 hours)
   - STOP and test with curl → Verify API works
4. **Phase 9-14**: Complete frontend (5 hours)
   - Test with browser → Verify dashboard works
5. **Demo**: Show working dashboard to stakeholders
6. **Phase 15-16**: Polish & cleanup (2 hours)

**Checkpoint at end of Phase 8**: Fully functional API, testable with curl  
**Checkpoint at end of Phase 14**: Complete feature, ready for production

### Incremental Delivery

This feature is atomic (single feature, not multiple stories), but could break at checkpoints:

- **After Phase 8**: Stable backend API, no frontend yet
- **After Phase 11**: Dashboard page working with mock data (if no backend)
- **After Phase 14**: Feature ready for staging/production testing

---

## Success Criteria (Verification Checklist)

### Database Verification
- [ ] `\dt meetings` shows table created
- [ ] `\dt meeting_attendees` shows junction table created
- [ ] `\di idx_meetings*` shows all meeting indexes
- [ ] No migration errors in console

### Backend Verification
- [ ] `go build` compiles without errors
- [ ] Server starts: `go run main.go`
- [ ] GET /api/dashboard returns 200 with correct structure
- [ ] POST /api/meetings creates meeting with 201
- [ ] GET /api/meetings/next returns meeting or 204
- [ ] Role-based filtering works (different roles see different data)

### Frontend Verification
- [ ] Dashboard page loads at http://localhost:5173/dashboard
- [ ] All components render without errors
- [ ] Statistics show with correct numbers
- [ ] Projects display with progress bars
- [ ] Tasks show with checkboxes
- [ ] Meeting card displays (if meeting exists) or hidden (if none)
- [ ] Task checkbox marks task complete
- [ ] Auto-refresh works every 30 seconds

### Performance Verification
- [ ] Dashboard loads in < 2 seconds
- [ ] API response time < 500ms
- [ ] No console errors or warnings
- [ ] No memory leaks detected

### User Flow Verification
- [ ] Login → Dashboard → See statistics
- [ ] Click project → Navigate to project details
- [ ] Check task → Task marked complete, counter updates
- [ ] Wait 30s → Dashboard auto-refreshes with new data
- [ ] Different role → Different data displayed

---

## Notes

- Each task has a specific file path for implementation
- [P] tasks can be parallelized (different files, no blocking dependencies)
- No [Story] label because this is a single, unified feature (no independent user stories)
- Tests are NOT included (project uses manual testing)
- Follow existing code patterns and conventions from AGENTS.md
- Commit after logical groups of tasks
- Stop at any checkpoint to validate feature independently
- Reference quickstart.md for detailed implementation code examples
- Reference data-model.md for schema details
- Reference api-contract.md for API specifications

---

## Quick Reference: File Changes Summary

### New Files (Backend)
- `backend/migration/005_add_dashboard_meetings.sql` - Database migration
- `backend/models/meeting.go` - Meeting and MeetingAttendee structs
- `backend/models/dashboard.go` - Dashboard DTOs
- `backend/repositories/dashboard_repository.go` - Dashboard queries
- `backend/repositories/meeting_repository.go` - Meeting data access
- `backend/services/dashboard_service.go` - Dashboard aggregation
- `backend/services/meeting_service.go` - Meeting business logic
- `backend/handlers/dashboard_handler.go` - Dashboard HTTP endpoint
- `backend/handlers/meeting_handler.go` - Meeting HTTP endpoints

### Modified Files (Backend)
- `backend/routes/routes.go` - Add 5 new routes
- `backend/main.go` - Wire new dependencies

### New Files (Frontend)
- `frontend/src/components/Avatar.svelte` - User avatar component
- `frontend/src/components/Dashboard.svelte` - Main dashboard page
- `frontend/src/components/MeetingCard.svelte` - Meeting widget
- `frontend/src/components/ProjectCard.svelte` - Project card
- `frontend/src/components/StatCard.svelte` - Statistics card
- `frontend/src/components/TaskListItem.svelte` - Task item
- `frontend/src/lib/api/dashboard.js` - Dashboard API client
- `frontend/src/lib/api/meetings.js` - Meetings API client

### Modified Files (Frontend)
- `frontend/src/lib/utils.js` - Add Jalali date formatter
- `frontend/src/App.svelte` - Add /dashboard route
- `frontend/src/components/Navigation.svelte` - Add dashboard link (if exists)

### Documentation Files
- `AGENTS.md` - Update with dashboard feature details
- `README.md` - Add dashboard to feature list
- `specs/005-dashboard/tasks.md` - This file

---

**Status**: ✅ Tasks ready for implementation

Start with Phase 1 (Setup), then Phase 2 (Database), then proceed sequentially through backend phases, then frontend phases. Use parallel opportunities to speed up development if team capacity allows.
