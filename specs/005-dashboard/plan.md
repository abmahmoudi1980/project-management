# Implementation Plan: Dashboard

**Branch**: `005-dashboard` | **Date**: 2026-01-01 | **Spec**: [spec.md](spec.md)  
**Input**: Feature specification from `/specs/005-dashboard/spec.md`

---

## Summary

Implement a project manager dashboard that displays real-time statistics (active projects, pending tasks, team members, upcoming deadlines), recent project cards with progress tracking, user's priority task list, and next upcoming team meeting. The dashboard serves as the primary landing page after login, providing comprehensive project oversight in a single view. Technical approach uses existing Go/Fiber backend with new dashboard aggregation endpoint, meetings entity with PostgreSQL storage, Svelte 5 frontend with 30-second polling for updates, and reusable components following established patterns.

---

## Technical Context

**Language/Version**: Go 1.21+ (backend), JavaScript ES6+ (frontend)  
**Primary Dependencies**: Fiber v2 (HTTP), pgx v5 (PostgreSQL driver), Svelte 5 (UI), Tailwind CSS (styling), jalali-moment (Persian dates), Lucide (icons)  
**Storage**: PostgreSQL 14+ with new `meetings` and `meeting_attendees` tables, indexes on existing tables for dashboard queries  
**Testing**: Manual testing (no automated tests currently), curl/Postman for API testing  
**Target Platform**: Web application (desktop and mobile responsive), modern browsers (Chrome, Firefox, Safari, Edge)  
**Project Type**: Web application (existing architecture: backend/ and frontend/ directories)  
**Performance Goals**: Dashboard loads in <2 seconds, supports 100 concurrent users, auto-refresh in <500ms without UI disruption  
**Constraints**: Must integrate with existing auth (JWT cookies), maintain layered architecture (handlers→services→repositories), use existing database schema patterns (UUIDs, snake_case), no new infrastructure (no Redis, WebSocket, or caching)  
**Scale/Scope**: Single organization scale (expected 10-50 users, 100-500 projects, 1000-5000 tasks), 4 stat cards, up to 4 project cards, up to 5 task items, 1 meeting card, single `/api/dashboard` endpoint

---

## Constitution Check

*No constitution file found in `.specify/memory/constitution.md`. Skipping constitution validation. This is acceptable for projects without explicit architectural constitutions.*

**Status**: ✅ PASSED (no constitution to validate against)

**Assumptions**:
- Following existing codebase patterns as documented in `AGENTS.md`
- Layered architecture (handlers→services→repositories) is implicit standard
- No violations detected based on existing patterns

---

## Project Structure

### Documentation (this feature)

```text
specs/005-dashboard/
├── spec.md              # Feature specification (COMPLETE)
├── plan.md              # This file (COMPLETE)
├── research.md          # Phase 0 output (COMPLETE)
├── data-model.md        # Phase 1 output (COMPLETE)
├── quickstart.md        # Phase 1 output (COMPLETE)
├── contracts/           # Phase 1 output (COMPLETE)
│   └── api-contract.md  # REST API specifications
├── checklists/          # Quality validation
│   └── requirements.md  # Spec quality checklist
└── tasks.md             # Phase 2 output (PENDING - use /speckit.tasks command)
```

### Source Code (repository root)

```text
backend/
├── cmd/
│   └── migrate/         # Database migration runner
├── config/              # Database & JWT configuration
├── handlers/            # HTTP request handlers
│   ├── auth_handler.go
│   ├── comment_handler.go
│   ├── dashboard_handler.go     # NEW - Dashboard endpoint
│   ├── meeting_handler.go       # NEW - Meeting CRUD
│   ├── project_handler.go
│   ├── task_handler.go
│   ├── timelog_handler.go
│   └── user_handler.go
├── middleware/          # Auth & other middleware
├── migration/           # SQL migration files
│   ├── 001_enhance_entities.sql
│   ├── 002_add_user_authentication.sql
│   ├── 003_add_comments.sql
│   └── 005_add_dashboard_meetings.sql  # NEW - Meetings tables
├── models/              # Domain entities (structs)
│   ├── comment.go
│   ├── dashboard.go     # NEW - Dashboard DTOs
│   ├── meeting.go       # NEW - Meeting models
│   ├── password_reset_token.go
│   ├── project.go
│   ├── session.go
│   ├── task.go
│   ├── timelog.go
│   └── user.go
├── repositories/        # Data access layer
│   ├── comment_repository.go
│   ├── dashboard_repository.go  # NEW - Dashboard queries
│   ├── meeting_repository.go    # NEW - Meeting data access
│   ├── password_reset_repository.go
│   ├── project_repository.go
│   ├── session_repository.go
│   ├── task_repository.go
│   ├── timelog_repository.go
│   └── user_repository.go
├── routes/              # Route registration
│   └── routes.go        # MODIFY - Add dashboard & meeting routes
├── services/            # Business logic layer
│   ├── auth_service.go
│   ├── comment_service.go
│   ├── dashboard_service.go     # NEW - Dashboard aggregation
│   ├── email_service.go
│   ├── meeting_service.go       # NEW - Meeting logic & validation
│   ├── project_service.go
│   ├── task_service.go
│   ├── timelog_service.go
│   └── user_service.go
└── main.go              # MODIFY - Wire new dependencies

frontend/
└── src/
    ├── components/      # Svelte 5 components
    │   ├── AdvancedTaskSearch.svelte
    │   ├── Avatar.svelte              # NEW - User avatar with initials
    │   ├── CommentList.svelte
    │   ├── Dashboard.svelte           # NEW - Main dashboard page
    │   ├── MeetingCard.svelte         # NEW - Upcoming meeting widget
    │   ├── ProjectCard.svelte         # NEW - Project summary card
    │   ├── StatCard.svelte            # NEW - Statistics card widget
    │   ├── TaskListItem.svelte        # NEW - Task list item with checkbox
    │   └── ... (existing components)
    ├── lib/             # Utilities & API clients
    │   ├── api/
    │   │   ├── dashboard.js   # NEW - Dashboard API calls
    │   │   └── meetings.js    # NEW - Meetings API calls
    │   └── utils.js           # MODIFY - Add Jalali date formatter
    └── stores/          # State management
        └── (existing stores)
```

**Structure Decision**: Using existing web application structure with `backend/` and `frontend/` directories. Dashboard is a new feature within the existing architecture, not a separate application. Follows established layered backend pattern (handlers→services→repositories) and component-based frontend pattern.

---

## Complexity Tracking

> **Constitution validation not applicable (no constitution file found)**

No complexity violations to justify. Implementation follows existing patterns documented in codebase.

---

## Implementation Phases

### Phase 0: Research & Discovery ✅ COMPLETE

**Output**: [research.md](research.md)

**Key Decisions Made**:
1. **Data aggregation**: Real-time COUNT queries with indexes (no caching/materialized views)
2. **Real-time updates**: 30-second polling (no WebSocket/SSE)
3. **Progress calculation**: Task count ratio (completed/total × 100)
4. **Task priority**: 4 levels (Critical=4, High=3, Medium=2, Low=1)
5. **Role-based filtering**: Backend service layer enforces visibility
6. **Meetings entity**: New table required (not existing feature)
7. **Jalali dates**: Client-side formatting with jalali-moment
8. **Avatars**: Generated from initials (no uploads for MVP)
9. **Change indicators**: 7-day lookback comparison
10. **Notifications**: Out of scope (bell icon present but non-functional)

**Technology Stack Confirmed**:
- Backend: Go + Fiber v2 + pgx + PostgreSQL
- Frontend: Svelte 5 + Tailwind CSS + Lucide icons + jalali-moment
- Auth: JWT in httpOnly cookies (existing)
- No new infrastructure dependencies

---

### Phase 1: Data Design & Contracts ✅ COMPLETE

**Outputs**:
- [data-model.md](data-model.md) - Entity definitions, relationships, validation rules
- [contracts/api-contract.md](contracts/api-contract.md) - REST API specifications
- [quickstart.md](quickstart.md) - Implementation guide
- [AGENTS.md](../../AGENTS.md) - Updated context file

**Key Artifacts**:

**New Tables**:
1. `meetings` - Meeting information (id, title, description, meeting_date, duration_minutes, project_id, created_by, timestamps)
2. `meeting_attendees` - Junction table (meeting_id, user_id, response_status, added_at)

**New Indexes**:
- `idx_meetings_date` - For upcoming meeting queries
- `idx_meeting_attendees_user` - For user meeting lookups
- `idx_projects_status_updated` - For recent projects query
- `idx_tasks_assigned_status` - For user task queries
- `idx_tasks_priority_due` - For priority sorting

**API Endpoints**:
1. `GET /api/dashboard` - All dashboard data (single request)
2. `PATCH /api/tasks/:id` - Complete task (existing, used by dashboard)
3. `GET /api/meetings/next` - Next upcoming meeting
4. `POST /api/meetings` - Create meeting
5. `GET /api/meetings` - List all meetings
6. `GET /api/meetings/:id` - Meeting details

**Data Structures**:
- `DashboardResponse` - Aggregated dashboard payload
- `DashboardStatistics` - 4 stat values with changes
- `ProjectCard` - Project summary with progress & team
- `TaskSummary` - Task item with priority & project
- `MeetingWithAttendees` - Meeting with attendee list

---

### Phase 2: Implementation Planning (Next Step)

**Pending**: Run `/speckit.tasks` command to generate `tasks.md`

**Tasks File Will Include**:
- [ ] Database migration (create tables, indexes)
- [ ] Backend models (Meeting, MeetingAttendee, Dashboard DTOs)
- [ ] Backend repositories (DashboardRepository, MeetingRepository)
- [ ] Backend services (DashboardService, MeetingService)
- [ ] Backend handlers (DashboardHandler, MeetingHandler)
- [ ] Backend routes (register new endpoints)
- [ ] Backend dependency wiring (main.go)
- [ ] Frontend API clients (dashboard.js, meetings.js)
- [ ] Frontend utilities (Jalali date formatter)
- [ ] Frontend components (Avatar, StatCard, ProjectCard, TaskListItem, MeetingCard)
- [ ] Frontend dashboard page (Dashboard.svelte)
- [ ] Frontend routing (add /dashboard route)
- [ ] Manual testing (API endpoints, UI interactions, auto-refresh)
- [ ] Documentation (update README, inline code comments)

**Estimated Implementation Time**: 12-16 hours (see quickstart.md for breakdown)

---

## Architecture Decisions

### Backend Architecture

**Layered Pattern** (following existing convention):
```
HTTP Request
    ↓
Handler (validation, HTTP concerns)
    ↓
Service (business logic, coordination)
    ↓
Repository (data access, SQL queries)
    ↓
PostgreSQL Database
```

**Key Principles**:
1. Handlers don't call repositories directly
2. Services don't handle HTTP concerns
3. Repositories don't contain business logic
4. Each layer has clear responsibility

**Dashboard-Specific**:
- `DashboardHandler`: Extracts user from JWT, calls service, returns JSON
- `DashboardService`: Coordinates multiple repository calls, aggregates data
- `DashboardRepository`: Executes SQL queries for statistics and project cards
- `MeetingRepository`: Handles meeting data access separately

### Frontend Architecture

**Component Hierarchy**:
```
Dashboard.svelte (page)
    ├── StatCard.svelte (x4 - statistics)
    ├── ProjectCard.svelte (x4 - recent projects)
    ├── TaskListItem.svelte (x5 - user tasks)
    └── MeetingCard.svelte (x1 - next meeting)
        └── Avatar.svelte (user avatars)
```

**State Management**:
- Dashboard data: Local component state with `$state()`
- Auto-refresh: `$effect()` hook with 30-second interval
- Task completion: Optimistic UI update, remove after 2 seconds
- No global store needed (dashboard is self-contained)

**API Communication**:
- Single `getDashboardData()` call on mount and every 30 seconds
- Separate `completeTask(id)` for task checkboxes
- Fetch with `credentials: 'include'` for JWT cookies

### Database Design

**Query Optimization Strategy**:
1. **Statistics**: 4 separate COUNT queries (fast with indexes)
2. **Recent Projects**: Single query with JOIN for tasks (calculate progress)
3. **User Tasks**: Single query with JOIN for project names
4. **Next Meeting**: Single query with JOIN for attendees
5. **Change Indicators**: Additional WHERE clauses for 7-day-ago counts

**Index Strategy**:
- Status + updated_at for recent projects
- assigned_user_id + status for user tasks
- priority + due_date for task sorting
- meeting_date for upcoming meetings
- user_id for meeting attendees

**No Triggers/Materialized Views**: Keeping queries simple and direct for now. Can optimize later if performance issues arise.

### Security Architecture

**Authentication**:
- Existing JWT middleware applies to all dashboard endpoints
- User identity extracted from JWT token in cookie
- No additional auth logic needed

**Authorization**:
- Service layer filters data by user role
- Admins: See all projects, all tasks
- Project Managers: See managed projects, team tasks
- Team Members: See assigned projects, own tasks only

**Data Filtering**:
```go
// Pseudo-code for authorization in service
func (s *DashboardService) GetDashboardData(userID, userRole) {
    visibleProjectIDs := getVisibleProjects(userID, userRole)
    stats := calculateStats(visibleProjectIDs)
    projects := getRecentProjects(visibleProjectIDs, limit=4)
    tasks := getUserTasks(userID, limit=5) // Always own tasks
    meeting := getNextMeeting(userID) // Attendee or creator
    return DashboardResponse{...}
}
```

### Performance Architecture

**Backend**:
- Connection pooling: pgxpool (default 4 connections, configurable)
- Query optimization: Strategic indexes, LIMIT clauses
- Single endpoint: Reduce HTTP round trips
- No N+1 queries: Use JOINs for related data

**Frontend**:
- Polling interval: 30 seconds (not too aggressive)
- No spinner during refresh: Subtle update, no flicker
- Optimistic updates: Task completion instant, sync on next refresh
- Lazy icon initialization: Lucide icons created after mount

**Database**:
- Indexes cover WHERE, ORDER BY, and JOIN columns
- Small result sets (4 projects, 5 tasks, 1 meeting)
- No full table scans with proper indexes
- UTC timestamps (no timezone conversion in DB)

---

## Data Flow

### Dashboard Load Sequence

1. **User navigates to `/dashboard`**
   - Frontend: Dashboard.svelte mounts
   - `$effect()` triggers `loadDashboard()`

2. **Frontend fetches data**
   - `getDashboardData()` → `GET /api/dashboard`
   - JWT cookie sent automatically

3. **Backend processes request**
   - Middleware: Validates JWT, extracts user
   - Handler: Calls `DashboardService.GetDashboardData(userID, role)`
   - Service: Calls 4 repository methods in parallel/sequence
   - Repositories: Execute SQL queries with user filtering
   - Database: Returns results

4. **Backend aggregates response**
   - Service: Combines statistics, projects, tasks, meeting
   - Handler: Returns JSON payload
   - HTTP: Sends response to frontend

5. **Frontend renders**
   - Parse JSON response
   - Update `dashboardData` state
   - Svelte re-renders components
   - Lucide creates icons

6. **Auto-refresh cycle**
   - Wait 30 seconds
   - Repeat steps 2-5
   - No loading spinner (subtle update)

### Task Completion Flow

1. **User clicks task checkbox**
   - TaskListItem: `handleCheck()` called
   - Optimistic: Set `completed = true` (instant UI feedback)

2. **Frontend updates backend**
   - `completeTask(taskId)` → `PATCH /api/tasks/:id`
   - Body: `{ status: 'done' }`

3. **Backend processes**
   - Handler: Validates task ownership
   - Service: Updates task status
   - Repository: `UPDATE tasks SET status = 'done' WHERE id = $1`

4. **Frontend post-update**
   - Wait 2 seconds (show strikethrough)
   - Remove task from list
   - Next auto-refresh syncs statistics

### Meeting Display Flow

1. **Backend fetches next meeting**
   - Service: `MeetingRepository.GetNextMeetingForUser(userID)`
   - SQL: `WHERE meeting_date > NOW() AND (creator OR attendee) ORDER BY meeting_date ASC LIMIT 1`

2. **Two scenarios**:
   - **Meeting exists**: Include in dashboard response
   - **No meeting**: `next_meeting: null`

3. **Frontend rendering**:
   - `{#if meeting}` - Show MeetingCard
   - `{:else}` - Hide card completely
   - Display: Title, description, time, up to 3 avatars

---

## API Design

### Dashboard Endpoint Design

**Why Single Endpoint?**
- Reduce HTTP round trips (4 separate requests → 1)
- Atomic data snapshot (consistent view)
- Simpler frontend logic (one API call)
- Better caching potential (future enhancement)

**Payload Size**: Typically 10-50KB
- Statistics: ~200 bytes
- 4 Projects: ~2-4KB (depends on team size)
- 5 Tasks: ~1-2KB
- 1 Meeting: ~500 bytes
- JSON overhead: ~500 bytes

**Endpoint**: `GET /api/dashboard`

**Response Structure**:
```json
{
  "statistics": { /* 4 stat objects */ },
  "recent_projects": [ /* up to 4 project objects */ ],
  "user_tasks": [ /* up to 5 task objects */ ],
  "next_meeting": { /* meeting object or null */ }
}
```

### Meeting Endpoints Design

**Separate CRUD**: Meetings have full lifecycle management
- Create: `POST /api/meetings`
- Read one: `GET /api/meetings/:id`
- Read list: `GET /api/meetings`
- Read next: `GET /api/meetings/next` (dashboard-specific)
- Update: `PATCH /api/meetings/:id` (future)
- Delete: `DELETE /api/meetings/:id` (future)

**Dashboard Uses**: Only `GET /api/meetings/next`
- Optimized query for dashboard use case
- Returns 204 if no meetings (distinguishes from error)
- Includes full attendee list (not just 3 visible)

---

## Error Handling

### Backend Error Strategy

**Validation Errors** (400 Bad Request):
```go
if len(input.Title) == 0 {
    return c.Status(400).JSON(fiber.Map{
        "error": "Bad Request",
        "message": "Title is required",
        "field": "title",
    })
}
```

**Authorization Errors** (403 Forbidden):
```go
if task.AssignedUserID != user.ID {
    return c.Status(403).JSON(fiber.Map{
        "error": "Forbidden",
        "message": "You can only update tasks assigned to you",
    })
}
```

**Not Found Errors** (404):
```go
if errors.Is(err, pgx.ErrNoRows) {
    return c.Status(404).JSON(fiber.Map{
        "error": "Not Found",
        "message": "Task not found",
    })
}
```

**Database Errors** (500 Internal Server Error):
```go
if err != nil {
    log.Error("Database error", err)
    return c.Status(500).JSON(fiber.Map{
        "error": "Internal Server Error",
        "message": "Failed to fetch dashboard data",
    })
}
```

### Frontend Error Strategy

**Network Errors**:
```javascript
try {
    dashboardData = await getDashboardData();
} catch (err) {
    error = err.message;
    // Show error banner, don't break UI
}
```

**Graceful Degradation**:
- If statistics fail: Show "—" instead of numbers
- If projects fail: Show empty state message
- If tasks fail: Show empty task list
- If meeting fails: Hide meeting card
- Never block entire dashboard for partial failure

**Retry Strategy**:
- Auto-refresh will retry every 30 seconds
- No manual retry button needed
- Failed updates show alert, don't break UI

---

## Testing Strategy

### Manual Testing Focus

**No automated tests currently** (per AGENTS.md), focus on manual validation.

**Critical Paths to Test**:
1. Dashboard load (fresh login)
2. Statistics accuracy (verify counts match data)
3. Project progress calculation (matches task completion)
4. Task priority sorting (correct order)
5. Task completion (checkbox → API → UI update)
6. Auto-refresh (30-second interval, no flicker)
7. Role-based filtering (different users see different data)
8. Empty states (no projects, no tasks, no meetings)
9. Jalali date formatting (today, tomorrow, dates)
10. Responsive layout (mobile, tablet, desktop)

**API Testing with curl**:
```bash
# Get JWT token first
TOKEN=$(curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"Admin123!"}' \
  -c cookies.txt | jq -r .token)

# Test dashboard endpoint
curl -X GET http://localhost:3000/api/dashboard \
  -b cookies.txt | jq .

# Test task completion
curl -X PATCH http://localhost:3000/api/tasks/TASK_UUID \
  -b cookies.txt \
  -H "Content-Type: application/json" \
  -d '{"status":"done"}' | jq .

# Test meeting creation
curl -X POST http://localhost:3000/api/meetings \
  -b cookies.txt \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Meeting",
    "meeting_date": "2026-01-15T10:00:00Z",
    "duration_minutes": 60,
    "attendee_ids": ["USER_UUID"]
  }' | jq .
```

### Browser Testing Checklist

**Chrome DevTools**:
- [ ] Network tab: Verify /api/dashboard called every 30 seconds
- [ ] Network tab: Verify JWT cookie sent with requests
- [ ] Console: No JavaScript errors
- [ ] Performance tab: Dashboard loads in <2 seconds
- [ ] Responsive mode: Test 320px, 768px, 1024px, 1920px widths

**User Flows**:
- [ ] Admin sees all projects
- [ ] Team member sees only assigned projects
- [ ] Task completion updates counter
- [ ] Completed task fades out after 2 seconds
- [ ] Clicking project navigates to detail page
- [ ] Auto-refresh doesn't interrupt user interaction
- [ ] Empty states display correctly
- [ ] Long project/task names truncate properly

---

## Deployment Considerations

### Environment Variables

**Backend** (`backend/.env`):
```
DATABASE_URL=postgres://user:pass@host:5432/dbname
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRY=15m
REFRESH_EXPIRY=168h
PORT=3000
ENV=production
```

**Frontend** (`frontend/.env.production`):
```
VITE_API_BASE_URL=https://api.yourdomain.com
```

### Database Migration

**Production Deployment**:
1. Backup database before migration
2. Run migration: `go run ./cmd/migrate`
3. Verify tables created: `\dt meetings*`
4. Verify indexes created: `\di idx_meetings*`
5. Test with sample data

**Rollback Plan**:
```sql
-- If needed, drop new tables
DROP TABLE IF EXISTS meeting_attendees CASCADE;
DROP TABLE IF EXISTS meetings CASCADE;

-- Drop indexes on existing tables
DROP INDEX IF EXISTS idx_projects_status_updated;
DROP INDEX IF EXISTS idx_tasks_assigned_status;
-- ... (other new indexes)
```

### Performance Monitoring

**Metrics to Track**:
- Dashboard load time (target: <2 seconds)
- API response time (target: <500ms)
- Database query duration
- Concurrent user count
- Error rate

**Logging**:
```go
// Add to dashboard handler
log.Info("Dashboard request",
    "user_id", user.ID,
    "role", user.Role,
    "duration", time.Since(start),
)
```

### Scaling Considerations

**Current Scale** (MVP):
- 10-50 users
- 100-500 projects
- 1000-5000 tasks
- 100-500 meetings

**Future Scale** (if needed):
- Add Redis caching for statistics (5-minute TTL)
- Use materialized views for expensive aggregations
- Implement WebSocket for real-time updates
- Add CDN for static assets
- Database read replicas for high read load

---

## Future Enhancements

**Not in MVP** (documented in spec.md "Out of Scope"):

1. **Dashboard Customization**
   - User-configurable widgets
   - Drag-and-drop layout
   - Widget visibility toggles

2. **Advanced Visualization**
   - Charts and graphs (burn-down, velocity)
   - Timeline view
   - Gantt charts

3. **Filtering & Searching**
   - Filter projects by status
   - Search within dashboard
   - Date range selectors

4. **Notifications**
   - In-app notification center
   - Email notifications
   - Desktop push notifications

5. **Export & Sharing**
   - PDF export
   - Share dashboard link
   - Email reports

6. **Mobile App**
   - Native iOS/Android apps
   - Offline support
   - Push notifications

7. **Integrations**
   - Calendar integration (Google, Outlook)
   - Slack notifications
   - Jira import/export

**Priority for Next Iterations**:
1. Automated tests (unit + integration)
2. Caching layer (Redis)
3. WebSocket real-time updates
4. Dashboard customization
5. Data visualization

---

## Risk Mitigation

### Technical Risks

**Risk**: Database performance degrades with scale  
**Mitigation**: Indexes cover all dashboard queries, monitor query performance, add caching if needed

**Risk**: Auto-refresh causes excessive server load  
**Mitigation**: 30-second interval is reasonable, rate limiting can be added, backend caching reduces DB load

**Risk**: Frontend state management becomes complex  
**Mitigation**: Keep dashboard self-contained, no global state, simple polling pattern

**Risk**: JWT token expiration during dashboard use  
**Mitigation**: Existing refresh token mechanism handles this transparently

### Implementation Risks

**Risk**: Backend layering violation  
**Mitigation**: Code review checklist, follow quickstart.md strictly, reference existing handlers

**Risk**: Svelte 5 runes syntax errors  
**Mitigation**: Use $props(), $state(), $derived() consistently, reference existing components

**Risk**: SQL injection vulnerabilities  
**Mitigation**: pgx parameterized queries prevent injection, never string concatenation for SQL

**Risk**: Authorization bypass  
**Mitigation**: Service layer always filters by user role, never trust frontend, test with different roles

### Timeline Risks

**Risk**: Implementation takes longer than 12-16 hours  
**Mitigation**: Follow quickstart.md step-by-step, test incrementally, ask for help if blocked

**Risk**: Database migration fails in production  
**Mitigation**: Test migration on staging first, have rollback plan ready, backup before migration

---

## Success Criteria

**Technical**:
- [ ] All API endpoints return correct data structure
- [ ] Dashboard loads in <2 seconds
- [ ] No console errors in browser
- [ ] Auto-refresh works without flicker
- [ ] Role-based filtering enforced
- [ ] Database migration successful

**Functional** (from spec.md):
- [ ] 4 statistics display with change indicators
- [ ] Up to 4 recent projects with progress bars
- [ ] Up to 5 user tasks sorted by priority
- [ ] Next meeting card appears if meeting exists
- [ ] Task checkbox marks task complete
- [ ] Clicking project card navigates to details
- [ ] Jalali dates display correctly
- [ ] Empty states handle gracefully

**User Experience**:
- [ ] Dashboard provides value within 30 seconds of viewing
- [ ] Users can identify priority tasks immediately
- [ ] Project health visible at a glance
- [ ] No confusion about what actions are available
- [ ] Responsive on mobile and desktop

**Code Quality**:
- [ ] Backend follows layered architecture
- [ ] Frontend uses Svelte 5 runes correctly
- [ ] No SQL injection vulnerabilities
- [ ] Errors handled gracefully
- [ ] Code follows existing patterns

---

## Next Steps

1. **Review this plan** - Ensure all stakeholders understand approach
2. **Run `/speckit.tasks` command** - Generate detailed task checklist in `tasks.md`
3. **Begin implementation** - Follow [quickstart.md](quickstart.md) step-by-step
4. **Test incrementally** - Validate each layer (DB → backend → frontend)
5. **Deploy to staging** - Test with real data and multiple users
6. **Collect feedback** - Iterate based on user needs
7. **Deploy to production** - Following deployment checklist above

---

## Questions & Clarifications

**Resolved**:
- ✅ Number format preference: Default to Western numerals (123) for consistency
- ✅ Meetings entity: Confirmed as new feature (not existing)
- ✅ Notification badge: Removed from MVP (coming soon message)
- ✅ Project creation modal: Link to existing page (not in dashboard scope)

**Open** (non-blocking):
- User preference for Persian vs Western numerals (can be added later as setting)
- Chart/graph requirements for future iterations
- Mobile app timeline and requirements

---

**Implementation Ready**: ✅ YES

This plan provides complete technical guidance for implementing the dashboard feature. Proceed to `/speckit.tasks` command to generate detailed implementation checklist.
