# Data Model: Dashboard

**Feature**: 005-dashboard  
**Date**: 2026-01-01  
**Status**: Design Phase

---

## Overview

The dashboard feature introduces one new entity (**Meeting**) and leverages three existing entities (Project, Task, User). This document defines the data structures, relationships, and validation rules.

---

## Entity Definitions

### New Entity: Meeting

Represents a scheduled team meeting with attendees.

**Table**: `meetings`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | UUID | PRIMARY KEY, NOT NULL | Unique identifier |
| title | VARCHAR(200) | NOT NULL | Meeting title/subject |
| description | TEXT | NULL | Detailed meeting description |
| meeting_date | TIMESTAMP WITH TIME ZONE | NOT NULL | Scheduled date and time (UTC) |
| duration_minutes | INTEGER | NOT NULL, DEFAULT 60 | Meeting duration in minutes |
| project_id | UUID | NULL, FOREIGN KEY → projects(id) | Optional associated project |
| created_by | UUID | NOT NULL, FOREIGN KEY → users(id) | User who created the meeting |
| created_at | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Record creation timestamp |
| updated_at | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | Last update timestamp |

**Indexes**:
```sql
CREATE INDEX idx_meetings_date ON meetings(meeting_date);
CREATE INDEX idx_meetings_project ON meetings(project_id) WHERE project_id IS NOT NULL;
CREATE INDEX idx_meetings_created_by ON meetings(created_by);
```

**Validation Rules**:
- `title`: 1-200 characters, required
- `description`: 0-5000 characters, optional
- `meeting_date`: Must be in the future when creating
- `duration_minutes`: Must be > 0 and <= 1440 (24 hours)
- `project_id`: Must exist in projects table if provided

---

### New Entity: MeetingAttendee

Junction table for many-to-many relationship between meetings and users.

**Table**: `meeting_attendees`

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| meeting_id | UUID | NOT NULL, FOREIGN KEY → meetings(id) ON DELETE CASCADE | Meeting identifier |
| user_id | UUID | NOT NULL, FOREIGN KEY → users(id) ON DELETE CASCADE | Attendee user identifier |
| response_status | VARCHAR(20) | NOT NULL, DEFAULT 'pending' | Attendance status |
| added_at | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT NOW() | When attendee was added |

**Primary Key**: `(meeting_id, user_id)`

**Indexes**:
```sql
CREATE INDEX idx_meeting_attendees_user ON meeting_attendees(user_id);
CREATE INDEX idx_meeting_attendees_meeting ON meeting_attendees(meeting_id);
```

**Validation Rules**:
- `response_status`: ENUM('pending', 'accepted', 'declined', 'maybe')
- Composite primary key prevents duplicate attendees per meeting

---

### Existing Entity: Project (Reference)

Used for dashboard project cards.

**Table**: `projects` (already exists)

**Fields Used by Dashboard**:
- `id`: UUID - Project identifier
- `name`: VARCHAR - Project name (displayed in cards)
- `description`: TEXT - Project description
- `status`: VARCHAR/ENUM - Current status (Planning, In Progress, On Track, Review, Completed)
- `start_date`: TIMESTAMP - Project start date
- `due_date`: TIMESTAMP - Project deadline (shown on dashboard)
- `updated_at`: TIMESTAMP - Last update time (for "recent" sorting)
- `created_at`: TIMESTAMP - Creation time

**Dashboard-Specific Queries**:
```sql
-- Recent projects (up to 4, ordered by last update)
SELECT id, name, status, due_date, updated_at
FROM projects
WHERE status IN ('Planning', 'In Progress', 'On Track', 'Review')
  AND id IN (user_visible_project_ids)
ORDER BY updated_at DESC
LIMIT 4;
```

**Note**: Client/organization name from spec will use project `description` or a future `client_name` field (not in current schema, defaultto project name for now).

---

### Existing Entity: Task (Reference)

Used for dashboard task list.

**Table**: `tasks` (already exists)

**Fields Used by Dashboard**:
- `id`: UUID - Task identifier
- `title`: VARCHAR - Task name
- `project_id`: UUID - Parent project
- `assigned_user_id`: UUID - User assigned to task
- `status`: VARCHAR/ENUM - Current status (todo, in_progress, done, blocked, etc.)
- `priority`: INTEGER or ENUM - Task priority (1=Low, 2=Medium, 3=High, 4=Critical)
- `due_date`: TIMESTAMP - Task deadline
- `created_at`: TIMESTAMP - Creation time (for sorting)

**Dashboard-Specific Queries**:
```sql
-- User's top 5 priority tasks
SELECT t.id, t.title, t.priority, t.due_date, p.name as project_name
FROM tasks t
JOIN projects p ON t.project_id = p.id
WHERE t.assigned_user_id = $1
  AND t.status != 'done'
ORDER BY t.priority DESC, t.due_date ASC, t.created_at ASC
LIMIT 5;
```

**Priority Mapping**:
- `4` or `'critical'` → "Critical" (red badge)
- `3` or `'high'` → "High" (orange badge)
- `2` or `'medium'` → "Medium" (blue badge)
- `1` or `'low'` → "Low" (gray badge)

---

### Existing Entity: User (Reference)

Used for avatars, assignments, and authentication.

**Table**: `users` (already exists)

**Fields Used by Dashboard**:
- `id`: UUID - User identifier
- `email`: VARCHAR - Email address
- `full_name`: VARCHAR - User's full name (for avatar initials)
- `role`: VARCHAR/ENUM - User role (admin, project_manager, team_member)
- `active`: BOOLEAN - Whether user is active

**Avatar Generation Logic** (in frontend):
```javascript
// Extract initials from full_name: "John Doe" → "JD"
function getInitials(fullName) {
  const parts = fullName.trim().split(' ');
  if (parts.length === 1) return parts[0].substring(0, 2).toUpperCase();
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
}

// Generate consistent color from user ID
function getAvatarColor(userId) {
  const colors = ['indigo', 'purple', 'blue', 'pink', 'orange', 'green', 'red', 'teal'];
  const hash = userId.split('').reduce((acc, char) => acc + char.charCodeAt(0), 0);
  return colors[hash % colors.length];
}
```

---

## Entity Relationships

```
┌─────────────┐
│   User      │
└──────┬──────┘
       │
       │ created_by (meetings)
       │ assigned_user_id (tasks)
       │ user_id (meeting_attendees)
       │
       ├────────────────────────────┐
       │                            │
       ▼                            ▼
┌─────────────┐            ┌──────────────────┐
│  Meeting    │────────────│ MeetingAttendee  │
└──────┬──────┘  1:N       └──────────────────┘
       │                            │ N:1
       │ project_id (optional)      │
       │                            │
       ▼                            ▼
┌─────────────┐                   User
│  Project    │
└──────┬──────┘
       │ 1:N
       │
       ▼
┌─────────────┐
│    Task     │
└──────┬──────┘
       │ N:1
       │
       ▼
      User (assigned_user_id)
```

**Relationship Descriptions**:
1. **User → Meeting**: One user creates many meetings (1:N)
2. **Meeting → MeetingAttendee**: One meeting has many attendees (1:N)
3. **User → MeetingAttendee**: One user attends many meetings (1:N)
4. **Project → Meeting**: One project may have many meetings (1:N, optional)
5. **Project → Task**: One project has many tasks (1:N, existing)
6. **User → Task**: One user is assigned many tasks (1:N, existing)

---

## Dashboard Aggregate Data

The dashboard displays several computed values that don't have dedicated tables.

### Statistics Payload

**Endpoint**: `GET /api/dashboard/stats`

**Response Structure**:
```json
{
  "active_projects": {
    "current": 12,
    "previous": 10,
    "change": 2
  },
  "pending_tasks": {
    "current": 48,
    "previous": 53,
    "change": -5
  },
  "team_members": {
    "current": 24,
    "previous": 20,
    "change": 4
  },
  "upcoming_deadlines": {
    "current": 7,
    "previous": 5,
    "change": 2
  }
}
```

**Calculation Logic**:

```sql
-- Active projects (current)
SELECT COUNT(*) FROM projects
WHERE status IN ('In Progress', 'On Track')
  AND id IN (user_visible_project_ids);

-- Active projects (7 days ago)
SELECT COUNT(*) FROM projects
WHERE status IN ('In Progress', 'On Track')
  AND id IN (user_visible_project_ids)
  AND (updated_at <= NOW() - INTERVAL '7 days'
       OR created_at <= NOW() - INTERVAL '7 days');

-- Pending tasks (current)
SELECT COUNT(*) FROM tasks
WHERE status != 'done'
  AND project_id IN (user_visible_project_ids);

-- Team members (current - for admins only)
SELECT COUNT(*) FROM users WHERE active = true;

-- Upcoming deadlines (next 7 days)
SELECT COUNT(*) FROM (
  SELECT id FROM tasks WHERE due_date BETWEEN NOW() AND NOW() + INTERVAL '7 days'
  UNION ALL
  SELECT id FROM projects WHERE due_date BETWEEN NOW() AND NOW() + INTERVAL '7 days'
) AS deadlines;
```

---

### Project Progress Calculation

**Formula**: `progress = (completed_tasks / total_tasks) * 100`

**Implementation**:
```sql
SELECT
  p.id,
  p.name,
  COUNT(t.id) as total_tasks,
  COUNT(CASE WHEN t.status = 'done' THEN 1 END) as completed_tasks,
  CASE
    WHEN COUNT(t.id) = 0 THEN 0
    ELSE ROUND((COUNT(CASE WHEN t.status = 'done' THEN 1 END)::NUMERIC / COUNT(t.id)) * 100)
  END as progress_percentage
FROM projects p
LEFT JOIN tasks t ON t.project_id = p.id
WHERE p.id = $1
GROUP BY p.id, p.name;
```

**Returned in Project Card Payload**:
```json
{
  "id": "uuid",
  "name": "Website Redesign",
  "client": "Acme Corp",
  "status": "In Progress",
  "progress": 75,
  "due_date": "2023-10-24T00:00:00Z",
  "team_members": [
    {"id": "uuid1", "full_name": "User One"},
    {"id": "uuid2", "full_name": "User Two"},
    {"id": "uuid3", "full_name": "User Three"}
  ]
}
```

---

### Team Members for Projects

Projects need to show team member avatars. This requires a project-user relationship.

**Options**:
1. **Infer from task assignments**: Users assigned to any task in project are team members
2. **Explicit project_members table**: Junction table for project assignments
3. **Project.owner_id only**: Show only project owner

**Decision**: Infer from task assignments (Option 1)

**Query**:
```sql
-- Get unique users assigned to tasks in a project
SELECT DISTINCT u.id, u.full_name
FROM users u
JOIN tasks t ON t.assigned_user_id = u.id
WHERE t.project_id = $1
LIMIT 3;
```

**Rationale**: No additional table needed, automatically reflects active participation.

---

## Database Migration

**Migration File**: `backend/migration/005_add_dashboard_meetings.sql`

```sql
-- Create meetings table
CREATE TABLE IF NOT EXISTS meetings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  title VARCHAR(200) NOT NULL,
  description TEXT,
  meeting_date TIMESTAMP WITH TIME ZONE NOT NULL,
  duration_minutes INTEGER NOT NULL DEFAULT 60 CHECK (duration_minutes > 0 AND duration_minutes <= 1440),
  project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
  created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create meeting_attendees junction table
CREATE TABLE IF NOT EXISTS meeting_attendees (
  meeting_id UUID NOT NULL REFERENCES meetings(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  response_status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (response_status IN ('pending', 'accepted', 'declined', 'maybe')),
  added_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  PRIMARY KEY (meeting_id, user_id)
);

-- Indexes for meetings
CREATE INDEX idx_meetings_date ON meetings(meeting_date);
CREATE INDEX idx_meetings_project ON meetings(project_id) WHERE project_id IS NOT NULL;
CREATE INDEX idx_meetings_created_by ON meetings(created_by);

-- Indexes for meeting_attendees
CREATE INDEX idx_meeting_attendees_user ON meeting_attendees(user_id);
CREATE INDEX idx_meeting_attendees_meeting ON meeting_attendees(meeting_id);

-- Indexes for dashboard queries on existing tables
CREATE INDEX IF NOT EXISTS idx_projects_status_updated ON projects(status, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_tasks_assigned_status ON tasks(assigned_user_id, status);
CREATE INDEX IF NOT EXISTS idx_tasks_due_date ON tasks(due_date) WHERE due_date IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_tasks_priority_due ON tasks(priority DESC, due_date ASC);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(active) WHERE active = true;

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_meetings_updated_at
  BEFORE UPDATE ON meetings
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();
```

---

## Validation Rules Summary

### Meeting Validation
- **Title**: Required, 1-200 characters
- **Description**: Optional, max 5000 characters
- **Meeting Date**: Required, must be valid timestamp
- **Duration**: Required, 1-1440 minutes
- **Project ID**: Optional, must exist if provided
- **Created By**: Required, must be valid user ID

### Meeting Attendee Validation
- **Meeting ID**: Required, must exist
- **User ID**: Required, must exist
- **Response Status**: Must be one of: pending, accepted, declined, maybe
- **Uniqueness**: One attendee record per user per meeting

### Business Rules
1. **Meeting visibility**: Users can only see meetings they're invited to OR created
2. **Dashboard meeting card**: Show only NEXT meeting (meeting_date > NOW, earliest first)
3. **Meeting in past**: Can still be displayed if meeting_date is within last 24 hours
4. **No attendees**: Meeting valid even with 0 attendees (draft state)
5. **Creator is attendee**: System automatically adds creator to attendees list

---

## API Payload Examples

### GET /api/dashboard

Returns all dashboard data in a single response to minimize requests.

**Response**:
```json
{
  "statistics": {
    "active_projects": { "current": 12, "previous": 10, "change": 2 },
    "pending_tasks": { "current": 48, "previous": 53, "change": -5 },
    "team_members": { "current": 24, "previous": 20, "change": 4 },
    "upcoming_deadlines": { "current": 7, "previous": 5, "change": 2 }
  },
  "recent_projects": [
    {
      "id": "uuid-1",
      "name": "Website Redesign",
      "client": "Acme Corp",
      "status": "In Progress",
      "progress": 75,
      "due_date": "2023-10-24T00:00:00Z",
      "team_members": [
        { "id": "user-1", "full_name": "User One" },
        { "id": "user-2", "full_name": "User Two" },
        { "id": "user-3", "full_name": "User Three" }
      ],
      "total_members": 5
    }
  ],
  "user_tasks": [
    {
      "id": "task-1",
      "title": "Design Homepage Mockups",
      "project_name": "Website Redesign",
      "project_id": "uuid-1",
      "priority": 3,
      "priority_label": "High",
      "due_date": "2023-10-14T00:00:00Z",
      "status": "in_progress"
    }
  ],
  "next_meeting": {
    "id": "meeting-1",
    "title": "Team Meeting",
    "description": "Weekly sync with the design team.",
    "meeting_date": "2023-10-13T10:00:00Z",
    "duration_minutes": 60,
    "attendees": [
      { "id": "user-1", "full_name": "User One" },
      { "id": "user-2", "full_name": "User Two" },
      { "id": "user-3", "full_name": "User Three" }
    ],
    "total_attendees": 5
  }
}
```

**Status Codes**:
- `200 OK`: Success
- `401 Unauthorized`: Not authenticated
- `500 Internal Server Error`: Server error

---

## Data Access Patterns

### Read Operations (Primary Use Case)

**Dashboard Load**:
1. GET /api/dashboard → Fetch all dashboard data
2. Frontend renders components from single payload
3. Auto-refresh every 30 seconds (same GET request)

**Query Optimization**:
- Use JOINs to minimize round trips
- Add WHERE clauses for user's visible projects/tasks
- Leverage indexes for sorting (ORDER BY updated_at, priority, due_date)
- LIMIT results (4 projects, 5 tasks, 1 meeting)

### Write Operations (Task Completion)

**Mark Task Complete**:
1. PATCH /api/tasks/:id → Update task status to 'done'
2. Frontend optimistically updates UI
3. Next auto-refresh syncs with server state

**Access Control**:
- Users can only complete tasks assigned to them
- Verify `assigned_user_id = current_user_id` in backend

---

## Future Considerations

**Not in MVP, but data model supports**:
1. **Recurring meetings**: Add `recurrence_rule` field to meetings table
2. **Meeting notes**: Add `notes` TEXT field or separate meeting_notes table
3. **Project client field**: Add `client_name` VARCHAR to projects table
4. **Manual progress override**: Add `manual_progress` INTEGER to projects table
5. **Task time tracking**: Add `estimated_hours`, `actual_hours` to tasks
6. **Notification preferences**: Add meeting notification settings
7. **Calendar integration**: Add external_calendar_id for syncing

---

## Notes

- All timestamps stored as UTC in PostgreSQL (`TIMESTAMP WITH TIME ZONE`)
- Frontend converts to Jalali calendar for display using `jalali-moment`
- UUIDs generated by database (`gen_random_uuid()`)
- Soft deletes not used; meetings cascade delete attendees
- Priority values: 1=Low, 2=Medium, 3=High, 4=Critical (consistent with existing task schema assumption)
- Status values: Use existing enum values from projects/tasks tables
