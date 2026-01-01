# Research: Dashboard Implementation

**Feature**: 005-dashboard  
**Date**: 2026-01-01  
**Purpose**: Resolve technical unknowns and identify best practices for dashboard implementation

---

## Research Questions

### 1. Dashboard Data Aggregation Strategy

**Question**: How should we efficiently aggregate statistics (active projects count, pending tasks count, etc.) without impacting database performance?

**Research Findings**:
- **Current Database**: PostgreSQL with existing tables (projects, tasks, users, sessions)
- **Query Patterns**: Dashboard requires COUNT queries across multiple tables with filters
- **Options Evaluated**:
  1. **Real-time aggregation**: Execute COUNT queries on each dashboard load
  2. **Cached aggregation**: Store pre-computed counts in Redis/memory cache
  3. **Materialized views**: PostgreSQL materialized views refreshed periodically
  4. **Database triggers**: Update counters on insert/update/delete

**Decision**: Real-time aggregation with selective indexing

**Rationale**:
- System is small-to-medium scale (based on existing codebase context)
- PostgreSQL COUNT with proper indexes is fast for expected data volume
- No additional infrastructure needed (no Redis)
- Counts remain accurate in real-time
- Can optimize later if performance issues arise

**Implementation Notes**:
- Add indexes on filtered columns: `projects.status`, `tasks.status`, `tasks.assigned_user_id`, `users.active`
- Single API endpoint `/api/dashboard/stats` returns all statistics
- Backend aggregates data before sending to frontend

**Alternatives Considered**:
- **Cached aggregation**: Rejected due to added complexity of cache invalidation and no current caching infrastructure
- **Materialized views**: Rejected as overkill for current scale, adds refresh complexity

---

### 2. Real-time Dashboard Updates

**Question**: How should the dashboard refresh when data changes (30-second auto-refresh vs WebSocket vs Server-Sent Events)?

**Research Findings**:
- **Current Architecture**: REST API with Fiber v2 backend
- **Options Evaluated**:
  1. **Polling (30-second interval)**: Frontend makes periodic GET requests
  2. **WebSocket**: Bidirectional real-time connection
  3. **Server-Sent Events (SSE)**: Unidirectional server push
  4. **Long polling**: Keep connection open until data changes

**Decision**: Simple polling with 30-second interval

**Rationale**:
- Matches existing architecture (REST API pattern)
- No WebSocket infrastructure exists in current codebase
- Dashboard doesn't require sub-second updates
- 30-second delay acceptable for project management use case
- Minimal backend changes required
- Lower server resource usage for expected user count

**Implementation Notes**:
- Use Svelte `setInterval` in `$effect` hook
- Cancel interval on component unmount
- Show subtle loading indicator during refresh
- Debounce rapid user-triggered refreshes

**Alternatives Considered**:
- **WebSocket**: Rejected as over-engineered for current needs, requires new backend infrastructure
- **SSE**: Rejected for similar reasons, not supported by existing backend framework patterns

---

### 3. Progress Calculation for Projects

**Question**: How should project completion percentage be calculated?

**Research Findings**:
- **Existing Schema**: Tasks have `status` field (likely enum with values like: todo, in_progress, done, blocked)
- **Options Evaluated**:
  1. **Task count ratio**: (completed tasks / total tasks) × 100
  2. **Weighted by estimated hours**: Sum(completed_hours) / Sum(total_hours) × 100
  3. **Manual override**: Project manager sets percentage manually
  4. **Hybrid**: Calculate automatically, allow manual override

**Decision**: Task count ratio (completed tasks / total tasks)

**Rationale**:
- Simple and transparent calculation
- No additional time tracking fields needed
- Matches common project management practice
- Easy for users to understand
- Consistent with existing schema (tasks have status field)

**Implementation Notes**:
- Count tasks WHERE status = 'completed' or equivalent
- Store calculation in API response, not in database
- Handle division by zero (projects with no tasks show 0%)
- Consider adding `manual_progress_override` field to projects table for future enhancement

**Alternatives Considered**:
- **Weighted by hours**: Rejected because time tracking may not be consistently used
- **Manual override**: Rejected for MVP to keep implementation simple, can add later

---

### 4. Task Priority and Sorting

**Question**: What priority levels should be supported and how should tasks be sorted on the dashboard?

**Research Findings**:
- **Existing Schema**: Tasks likely have priority field (check migration files)
- **Industry Standards**: Most PM tools use 4-5 priority levels
- **Options Evaluated**:
  1. **4 levels**: Critical, High, Medium, Low
  2. **3 levels**: High, Medium, Low
  3. **5 levels**: Critical, High, Medium, Low, Trivial
  4. **Numeric (1-10)**: Flexible but less user-friendly

**Decision**: 4 priority levels (Critical, High, Medium, Low)

**Rationale**:
- Matches reference design (dashboard.html) which shows these exact levels
- Critical level separates urgent issues from normal high-priority work
- Widely recognized and easy to understand
- Maps to numeric values for sorting (Critical=4, High=3, Medium=2, Low=1)

**Sorting Algorithm**:
1. **Primary sort**: Priority (Critical > High > Medium > Low)
2. **Secondary sort**: Due date (earliest first)
3. **Tertiary sort**: Created date (oldest first)

**Implementation Notes**:
- Ensure tasks table has priority ENUM or INTEGER field
- Dashboard query: `ORDER BY priority DESC, due_date ASC, created_at ASC LIMIT 5`
- Frontend displays priority as colored badges matching reference design

**Alternatives Considered**:
- **3 levels**: Rejected as lacking distinction for truly urgent issues
- **Numeric priority**: Rejected as less user-friendly, harder to understand at a glance

---

### 5. Role-Based Data Filtering

**Question**: How should dashboard data be filtered based on user roles (Admin vs Project Manager vs Team Member)?

**Research Findings**:
- **Existing Auth**: Feature 003 implemented JWT-based authentication with user roles
- **Repository Pattern**: Backend uses repositories for data access with filtering
- **Options Evaluated**:
  1. **Backend filtering**: API enforces role-based access at service/repository level
  2. **Frontend filtering**: Fetch all data, filter in UI
  3. **Hybrid**: Backend filters critical data, frontend for display preferences

**Decision**: Backend filtering at service layer

**Rationale**:
- Security: Never send unauthorized data to client
- Performance: Reduce payload size
- Consistency: Matches existing architecture (services layer handles business logic)
- Scalability: Works as data grows

**Filtering Rules**:
- **Admin**: See all projects and all team members' tasks
- **Project Manager**: See projects they manage + assigned projects, team tasks for their projects
- **Team Member**: See only projects they're assigned to, only their own tasks

**Implementation Notes**:
- Add `user_id` and `user_role` parameters to service methods
- Use existing `user_repository` to fetch user's project assignments
- Filter in repository layer: `WHERE project_id IN (user_project_ids)`
- "Your Tasks" section always filters by `assigned_user_id = current_user_id`

**Alternatives Considered**:
- **Frontend filtering**: Rejected due to security concerns and performance
- **No filtering (show everything)**: Rejected as violates principle of least privilege

---

### 6. Meeting Data Integration

**Question**: Does the system already have meetings/calendar functionality, or is this a new entity?

**Research Findings**:
- **Schema Review**: No existing `meetings` or `calendar_events` table in migration files
- **Existing Features**: Features 001-004 cover projects, tasks, users, search - no meetings
- **Reference Design**: dashboard.html shows "Team Meeting" card with title, description, time, attendees

**Decision**: Add meetings as new entity

**Rationale**:
- Dashboard spec includes team meeting card as a feature requirement
- Natural fit for project management system
- Relatively simple entity (fewer fields than projects/tasks)

**Meeting Entity Design** (for data-model.md):
```
Meeting:
  - id: UUID (primary key)
  - title: VARCHAR(200)
  - description: TEXT
  - meeting_date: TIMESTAMP
  - duration_minutes: INTEGER
  - project_id: UUID (foreign key, optional)
  - created_by: UUID (foreign key to users)
  - created_at: TIMESTAMP
  - updated_at: TIMESTAMP

MeetingAttendee (junction table):
  - meeting_id: UUID (foreign key)
  - user_id: UUID (foreign key)
  - PRIMARY KEY (meeting_id, user_id)
```

**Implementation Notes**:
- Dashboard shows only NEXT upcoming meeting (meeting_date > NOW, ORDER BY meeting_date ASC LIMIT 1)
- Filter by attendees: WHERE user_id IN (SELECT user_id FROM meeting_attendees WHERE meeting_id = meetings.id)
- Meeting card hidden if no meetings in next 7 days

**Alternatives Considered**:
- **Skip meetings feature**: Rejected because spec explicitly requires it
- **External calendar integration**: Rejected as too complex for MVP, can add later

---

### 7. Persian/Jalali Date Formatting

**Question**: How should dates be displayed in Jalali calendar format?

**Research Findings**:
- **Existing Dependency**: `jalali-moment` already in frontend package.json (from AGENTS.md)
- **Options Evaluated**:
  1. **Client-side formatting**: Convert dates in Svelte components
  2. **Server-side formatting**: Backend returns Jalali dates
  3. **Hybrid**: Server sends UTC, client converts to Jalali

**Decision**: Client-side formatting with jalali-moment

**Rationale**:
- Leverage existing dependency
- Keep backend timezone-agnostic (UTC timestamps)
- Supports internationalization (easy to add Gregorian option later)
- Frontend controls display format
- Matches existing pattern in codebase

**Implementation Notes**:
- Create utility function `formatJalaliDate(date, format)` in `frontend/src/lib/utils.js`
- Use format: `jYYYY/jMM/jDD` for dates, `jMM jDD` for dashboard compact dates
- Display "امروز" (today), "فردا" (tomorrow) for relative dates
- Store all dates as UTC timestamps in database

**Alternatives Considered**:
- **Server-side formatting**: Rejected to keep API timezone-independent
- **New library**: Rejected, jalali-moment is sufficient and already installed

---

### 8. Avatar Display Strategy

**Question**: How should user avatars be displayed (generated vs uploaded)?

**Research Findings**:
- **User Table**: Likely has name/email fields, may not have avatar URL
- **Options Evaluated**:
  1. **Generated avatars**: Use initials on colored backgrounds
  2. **Uploaded avatars**: Users upload profile pictures
  3. **Gravatar**: Use email hash to fetch from Gravatar
  4. **Hybrid**: Support uploads, fallback to generated

**Decision**: Generated avatars with initials (MVP)

**Rationale**:
- No file upload infrastructure needed
- Consistent with reference design (shows initials: JD, U1, U2)
- Fast implementation
- No storage costs
- Can add upload capability later

**Implementation Notes**:
- Extract initials from user's full name (first letter of first and last name)
- Generate consistent background color from user ID hash
- Use 8 predefined colors (indigo, purple, blue, pink, orange, green, red, teal)
- Create reusable Svelte component: `<Avatar user={user} size="sm|md|lg" />`

**Alternatives Considered**:
- **Uploaded avatars**: Rejected for MVP due to file storage complexity
- **Gravatar**: Rejected as may not have images for all users, requires external service

---

### 9. Change Indicators Calculation

**Question**: How should "+2" / "-5" change indicators be calculated (7 days ago comparison)?

**Research Findings**:
- **Spec Requirement**: "Change indicators shall calculate difference from same period 7 days ago"
- **Options Evaluated**:
  1. **Real-time comparison**: Query data for 7 days ago on each request
  2. **Daily snapshots**: Store daily counts, compare to 7-day-old snapshot
  3. **No indicators**: Show only current counts
  4. **Rolling average**: Show trend over time

**Decision**: Real-time comparison with 7-day lookback

**Rationale**:
- Simple implementation for MVP
- No additional tables needed
- Accurate to the moment
- Acceptable query performance with proper indexes

**Implementation Notes**:
- API calculates both current count and 7-day-ago count in single transaction
- Add timestamp filters: `created_at <= NOW() - INTERVAL '7 days'` for historical count
- Return JSON: `{ "current": 12, "previous": 10, "change": +2 }`
- Frontend displays change with color coding (green for positive, red for negative)

**Alternatives Considered**:
- **Daily snapshots**: Rejected as adds complexity, requires scheduled job
- **No indicators**: Rejected as spec explicitly requires them
- **Rolling average**: Rejected as too complex for current need

---

### 10. Notification Badge Count

**Question**: What should the notification bell badge show?

**Research Findings**:
- **Reference Design**: Shows red dot badge on notification bell
- **Existing Schema**: No notifications table currently exists
- **Options Evaluated**:
  1. **Build notifications system**: Create notifications table and logic
  2. **Show simple count**: Count of tasks due today/overdue
  3. **Mock for now**: Show placeholder, implement notifications later
  4. **Remove badge**: Don't show until notifications exist

**Decision**: Remove notification badge (out of scope for dashboard)

**Rationale**:
- Notifications is a separate feature requiring significant work
- Dashboard spec doesn't detail notification requirements
- Better to not show than to show misleading information
- Can add in future feature (006-notifications)

**Implementation Notes**:
- Keep bell icon in header for UI consistency
- Remove red badge dot
- Make icon clickable but show "Coming soon" message
- Document in "Out of Scope" that notification functionality is separate feature

**Alternatives Considered**:
- **Quick notification system**: Rejected as would be incomplete and misleading
- **Show task count**: Rejected as notifications ≠ tasks, would confuse users

---

## Technology Stack Summary

Based on existing codebase (from AGENTS.md) and research above:

**Backend**:
- Language: Go
- Framework: Fiber v2
- Database: PostgreSQL with pgx driver
- Authentication: JWT (access 15min, refresh 7d) in httpOnly cookies
- Password Hashing: bcrypt (cost 10)
- Architecture: Handlers → Services → Repositories (strict layering)

**Frontend**:
- Framework: Svelte 5 with runes
- Styling: Tailwind CSS
- Icons: Lucide icons library
- Date Formatting: jalali-moment
- Language: JavaScript (not TypeScript)
- State Management: Svelte stores

**Database**:
- Primary Keys: UUID
- Naming: snake_case for columns, PascalCase for Go structs
- Migrations: SQL files in `backend/migration/`

---

## Implementation Approach

Based on research findings, the dashboard implementation will:

1. **Create new API endpoint**: `/api/dashboard` returning all dashboard data in single response
2. **Add meetings entity**: New table + handler + service + repository following existing patterns
3. **Leverage existing entities**: Projects, tasks, users tables already exist
4. **Create dashboard component**: New Svelte 5 component with runes
5. **Use existing patterns**: Follow established architecture in features 001-004
6. **No new infrastructure**: Use existing database, no caching/WebSocket/Redis

**Performance Strategy**:
- Add database indexes for dashboard queries
- Single API call loads all dashboard data
- 30-second polling for updates
- Lazy load avatars if needed

**Security Strategy**:
- Reuse existing JWT middleware
- Filter data by user role at service layer
- Never send unauthorized data to client

---

## Open Questions / Decisions Deferred

1. **Number format preference**: Persian-Indic digits (۱۲۳) vs Western Arabic (123) - awaiting user input, defaulting to Western for consistency with codebase
2. **Project creation modal**: Spec mentions "New Project" button but modal is out of scope - will navigate to existing project creation page or show coming soon
3. **Task creation modal**: Similar to above, out of scope for dashboard feature
4. **Search functionality**: Header has search bar, but search implementation is in feature 004 (task search) - dashboard will link to existing search

---

## Next Steps

With research complete, proceed to:
1. **Phase 1**: Create data-model.md with entity definitions
2. **Phase 1**: Create API contracts in /contracts/
3. **Phase 1**: Create quickstart.md with development guide
4. **Phase 1**: Update AGENTS.md context file
5. **Phase 2**: Generate tasks.md with implementation checklist (separate command: `/speckit.tasks`)
