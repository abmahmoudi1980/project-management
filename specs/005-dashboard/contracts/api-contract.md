# API Contract: Dashboard

**Feature**: 005-dashboard  
**Base URL**: `/api`  
**Authentication**: Required (JWT token in httpOnly cookie)  
**Date**: 2026-01-01

---

## Endpoints

### 1. Get Dashboard Data

Retrieves all dashboard information in a single request.

**Endpoint**: `GET /api/dashboard`

**Authentication**: Required

**Query Parameters**: None

**Request Headers**:
```
Cookie: access_token=<jwt_token>
```

**Success Response** (200 OK):
```json
{
  "statistics": {
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
  },
  "recent_projects": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Website Redesign",
      "client": "Acme Corp",
      "status": "In Progress",
      "progress": 75,
      "due_date": "2023-10-24T00:00:00Z",
      "updated_at": "2023-10-13T15:30:00Z",
      "team_members": [
        {
          "id": "660e8400-e29b-41d4-a716-446655440001",
          "full_name": "Alice Johnson"
        },
        {
          "id": "660e8400-e29b-41d4-a716-446655440002",
          "full_name": "Bob Smith"
        },
        {
          "id": "660e8400-e29b-41d4-a716-446655440003",
          "full_name": "Carol White"
        }
      ],
      "total_members": 5
    }
  ],
  "user_tasks": [
    {
      "id": "770e8400-e29b-41d4-a716-446655440010",
      "title": "Design Homepage Mockups",
      "project_name": "Website Redesign",
      "project_id": "550e8400-e29b-41d4-a716-446655440000",
      "priority": 3,
      "priority_label": "High",
      "due_date": "2023-10-14T00:00:00Z",
      "status": "in_progress",
      "created_at": "2023-10-10T09:00:00Z"
    }
  ],
  "next_meeting": {
    "id": "880e8400-e29b-41d4-a716-446655440020",
    "title": "Team Meeting",
    "description": "Weekly sync with the design team.",
    "meeting_date": "2023-10-13T10:00:00Z",
    "duration_minutes": 60,
    "attendees": [
      {
        "id": "660e8400-e29b-41d4-a716-446655440001",
        "full_name": "Alice Johnson"
      },
      {
        "id": "660e8400-e29b-41d4-a716-446655440002",
        "full_name": "Bob Smith"
      },
      {
        "id": "660e8400-e29b-41d4-a716-446655440003",
        "full_name": "Carol White"
      }
    ],
    "total_attendees": 5
  }
}
```

**Field Descriptions**:

**statistics**:
- `current` (integer): Current count
- `previous` (integer): Count from 7 days ago
- `change` (integer): Difference (current - previous), can be negative

**recent_projects** (array, max 4 items):
- `id` (string, UUID): Project unique identifier
- `name` (string): Project name
- `client` (string): Client or organization name (defaults to project name if not available)
- `status` (string): One of: "Planning", "In Progress", "On Track", "Review", "Completed"
- `progress` (integer, 0-100): Completion percentage based on task ratio
- `due_date` (string, ISO 8601): Project deadline in UTC
- `updated_at` (string, ISO 8601): Last update timestamp in UTC
- `team_members` (array, max 3 items): Team member objects with id and full_name
- `total_members` (integer): Total count of all team members (including those not shown)

**user_tasks** (array, max 5 items):
- `id` (string, UUID): Task unique identifier
- `title` (string): Task name
- `project_name` (string): Name of parent project
- `project_id` (string, UUID): Parent project ID
- `priority` (integer, 1-4): Priority level (1=Low, 2=Medium, 3=High, 4=Critical)
- `priority_label` (string): Human-readable priority ("Low", "Medium", "High", "Critical")
- `due_date` (string, ISO 8601, nullable): Task deadline in UTC (null if not set)
- `status` (string): Current task status
- `created_at` (string, ISO 8601): Task creation timestamp in UTC

**next_meeting** (object, nullable):
- `id` (string, UUID): Meeting unique identifier
- `title` (string): Meeting title
- `description` (string, nullable): Meeting description
- `meeting_date` (string, ISO 8601): Scheduled meeting time in UTC
- `duration_minutes` (integer): Meeting duration
- `attendees` (array, max 3 items): Attendee objects with id and full_name
- `total_attendees` (integer): Total count of all attendees (including those not shown)
- **null** if no meetings scheduled within next 7 days

**Error Responses**:

401 Unauthorized:
```json
{
  "error": "Unauthorized",
  "message": "Authentication required"
}
```

500 Internal Server Error:
```json
{
  "error": "Internal Server Error",
  "message": "Failed to fetch dashboard data"
}
```

**Notes**:
- Response filtered by user's role and project assignments
- Admins see all projects; team members see only assigned projects
- Tasks always filtered by logged-in user's assignments
- Meetings filtered to show only those where user is attendee or creator
- All timestamps in UTC; client converts to Jalali calendar
- Response size typically 10-50KB depending on data volume

---

### 2. Complete Task (Existing Endpoint)

Marks a task as complete. Used by dashboard task list checkboxes.

**Endpoint**: `PATCH /api/tasks/:id`

**Authentication**: Required

**URL Parameters**:
- `id` (UUID, required): Task identifier

**Request Headers**:
```
Cookie: access_token=<jwt_token>
Content-Type: application/json
```

**Request Body**:
```json
{
  "status": "done"
}
```

**Success Response** (200 OK):
```json
{
  "id": "770e8400-e29b-41d4-a716-446655440010",
  "title": "Design Homepage Mockups",
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "assigned_user_id": "660e8400-e29b-41d4-a716-446655440001",
  "status": "done",
  "priority": 3,
  "due_date": "2023-10-14T00:00:00Z",
  "created_at": "2023-10-10T09:00:00Z",
  "updated_at": "2023-10-13T16:45:00Z"
}
```

**Error Responses**:

401 Unauthorized:
```json
{
  "error": "Unauthorized",
  "message": "Authentication required"
}
```

403 Forbidden:
```json
{
  "error": "Forbidden",
  "message": "You can only update tasks assigned to you"
}
```

404 Not Found:
```json
{
  "error": "Not Found",
  "message": "Task not found"
}
```

400 Bad Request:
```json
{
  "error": "Bad Request",
  "message": "Invalid status value"
}
```

**Notes**:
- Users can only update tasks assigned to them (`assigned_user_id` check)
- Status must be a valid enum value from tasks table
- Dashboard should optimistically update UI before receiving response

---

## Meeting Endpoints (New)

### 3. Get User's Next Meeting

Retrieves the next upcoming meeting for the authenticated user.

**Endpoint**: `GET /api/meetings/next`

**Authentication**: Required

**Query Parameters**: None

**Request Headers**:
```
Cookie: access_token=<jwt_token>
```

**Success Response** (200 OK):
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440020",
  "title": "Team Meeting",
  "description": "Weekly sync with the design team.",
  "meeting_date": "2023-10-13T10:00:00Z",
  "duration_minutes": 60,
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "created_by": "660e8400-e29b-41d4-a716-446655440001",
  "attendees": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "full_name": "Alice Johnson",
      "response_status": "accepted"
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440002",
      "full_name": "Bob Smith",
      "response_status": "pending"
    }
  ]
}
```

**Success Response - No Meeting** (204 No Content):
```
(Empty response body)
```

**Error Responses**:

401 Unauthorized:
```json
{
  "error": "Unauthorized",
  "message": "Authentication required"
}
```

**Notes**:
- Returns only the NEXT meeting where user is attendee or creator
- Meeting must be scheduled in the future (meeting_date > NOW)
- Returns 204 if no meetings found within next 7 days
- Attendees array includes full list (not limited to 3)

---

### 4. Create Meeting

Creates a new meeting.

**Endpoint**: `POST /api/meetings`

**Authentication**: Required

**Request Headers**:
```
Cookie: access_token=<jwt_token>
Content-Type: application/json
```

**Request Body**:
```json
{
  "title": "Team Meeting",
  "description": "Weekly sync with the design team.",
  "meeting_date": "2023-10-13T10:00:00Z",
  "duration_minutes": 60,
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "attendee_ids": [
    "660e8400-e29b-41d4-a716-446655440001",
    "660e8400-e29b-41d4-a716-446655440002"
  ]
}
```

**Field Validation**:
- `title` (string, required): 1-200 characters
- `description` (string, optional): 0-5000 characters
- `meeting_date` (string, ISO 8601, required): Must be future timestamp
- `duration_minutes` (integer, required): 1-1440 minutes
- `project_id` (string, UUID, optional): Must exist in projects table
- `attendee_ids` (array of UUIDs, required): At least 1 user, all must exist

**Success Response** (201 Created):
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440020",
  "title": "Team Meeting",
  "description": "Weekly sync with the design team.",
  "meeting_date": "2023-10-13T10:00:00Z",
  "duration_minutes": 60,
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "created_by": "660e8400-e29b-41d4-a716-446655440099",
  "created_at": "2023-10-10T14:30:00Z",
  "updated_at": "2023-10-10T14:30:00Z",
  "attendees": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "full_name": "Alice Johnson",
      "response_status": "pending"
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440002",
      "full_name": "Bob Smith",
      "response_status": "pending"
    }
  ]
}
```

**Error Responses**:

401 Unauthorized:
```json
{
  "error": "Unauthorized",
  "message": "Authentication required"
}
```

400 Bad Request:
```json
{
  "error": "Bad Request",
  "message": "Title is required and must be 1-200 characters",
  "field": "title"
}
```

404 Not Found:
```json
{
  "error": "Not Found",
  "message": "Project not found",
  "field": "project_id"
}
```

**Notes**:
- Meeting creator automatically added to attendees with status "accepted"
- All other attendees start with status "pending"
- Meeting date must be in future when creating
- Project ID optional; meeting can exist without project association

---

### 5. List All Meetings

Retrieves all meetings for the authenticated user (creator or attendee).

**Endpoint**: `GET /api/meetings`

**Authentication**: Required

**Query Parameters**:
- `from` (ISO 8601 date, optional): Start date filter (default: NOW)
- `to` (ISO 8601 date, optional): End date filter (default: +30 days)
- `limit` (integer, optional): Max results (default: 50, max: 100)
- `offset` (integer, optional): Pagination offset (default: 0)

**Request Headers**:
```
Cookie: access_token=<jwt_token>
```

**Success Response** (200 OK):
```json
{
  "meetings": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440020",
      "title": "Team Meeting",
      "description": "Weekly sync with the design team.",
      "meeting_date": "2023-10-13T10:00:00Z",
      "duration_minutes": 60,
      "project_id": "550e8400-e29b-41d4-a716-446655440000",
      "created_by": "660e8400-e29b-41d4-a716-446655440001",
      "attendee_count": 5
    }
  ],
  "total": 12,
  "limit": 50,
  "offset": 0
}
```

**Error Responses**:

401 Unauthorized:
```json
{
  "error": "Unauthorized",
  "message": "Authentication required"
}
```

400 Bad Request:
```json
{
  "error": "Bad Request",
  "message": "Invalid date format"
}
```

**Notes**:
- Returns only meetings where user is creator or attendee
- Ordered by meeting_date ASC (chronological)
- Attendee details not included in list view (use GET /api/meetings/:id for details)

---

### 6. Get Meeting Details

Retrieves full details for a specific meeting.

**Endpoint**: `GET /api/meetings/:id`

**Authentication**: Required

**URL Parameters**:
- `id` (UUID, required): Meeting identifier

**Request Headers**:
```
Cookie: access_token=<jwt_token>
```

**Success Response** (200 OK):
```json
{
  "id": "880e8400-e29b-41d4-a716-446655440020",
  "title": "Team Meeting",
  "description": "Weekly sync with the design team.",
  "meeting_date": "2023-10-13T10:00:00Z",
  "duration_minutes": 60,
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "project_name": "Website Redesign",
  "created_by": "660e8400-e29b-41d4-a716-446655440001",
  "created_at": "2023-10-10T14:30:00Z",
  "updated_at": "2023-10-10T14:30:00Z",
  "attendees": [
    {
      "id": "660e8400-e29b-41d4-a716-446655440001",
      "full_name": "Alice Johnson",
      "email": "alice@example.com",
      "response_status": "accepted",
      "added_at": "2023-10-10T14:30:00Z"
    },
    {
      "id": "660e8400-e29b-41d4-a716-446655440002",
      "full_name": "Bob Smith",
      "email": "bob@example.com",
      "response_status": "pending",
      "added_at": "2023-10-10T14:30:00Z"
    }
  ]
}
```

**Error Responses**:

401 Unauthorized:
```json
{
  "error": "Unauthorized",
  "message": "Authentication required"
}
```

403 Forbidden:
```json
{
  "error": "Forbidden",
  "message": "You don't have permission to view this meeting"
}
```

404 Not Found:
```json
{
  "error": "Not Found",
  "message": "Meeting not found"
}
```

**Notes**:
- User must be creator or attendee to view meeting details
- Includes full attendee list with email addresses
- Project name included if meeting associated with project

---

## Data Types

### Common Types

**UUID**: String in format `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`

**ISO 8601 Timestamp**: String in format `YYYY-MM-DDTHH:MM:SSZ` (UTC timezone)

**Priority Values**:
- `1` = "Low"
- `2` = "Medium"
- `3` = "High"
- `4` = "Critical"

**Project Status Values**:
- "Planning"
- "In Progress"
- "On Track"
- "Review"
- "Completed"

**Task Status Values**:
- "todo"
- "in_progress"
- "done"
- "blocked"

**Meeting Response Status Values**:
- "pending"
- "accepted"
- "declined"
- "maybe"

---

## Rate Limiting

**Dashboard Endpoint**:
- Rate: 120 requests per minute per user
- Recommended polling: Every 30 seconds
- Burst limit: 10 requests in 5 seconds

**Meeting Endpoints**:
- Rate: 60 requests per minute per user
- No burst limits

**Rate Limit Headers** (included in all responses):
```
X-RateLimit-Limit: 120
X-RateLimit-Remaining: 115
X-RateLimit-Reset: 1697204400
```

---

## Authentication

All endpoints require authentication via JWT token in httpOnly cookie.

**Cookie Name**: `access_token`

**Token Expiration**: 15 minutes (access token), 7 days (refresh token)

**Refresh Flow**: Use existing `/api/auth/refresh` endpoint

**Logout**: Use existing `/api/auth/logout` endpoint

---

## CORS

**Allowed Origins**: Same origin only (no CORS needed for initial implementation)

**Future**: If frontend and backend on different domains, add CORS headers

---

## Versioning

**Current Version**: v1 (implicit, no version in URL)

**Breaking Changes**: Will introduce `/api/v2/` prefix if needed

**Deprecation**: 6-month notice for deprecated endpoints

---

## Error Handling

All error responses follow consistent format:

```json
{
  "error": "ErrorType",
  "message": "Human-readable description",
  "field": "optional_field_name",
  "code": "optional_error_code"
}
```

**HTTP Status Codes**:
- `200 OK`: Success with body
- `201 Created`: Resource created
- `204 No Content`: Success with no body
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

---

## Testing

**Contract Tests Required**:
1. GET /api/dashboard returns all fields in correct format
2. GET /api/dashboard filters data by user role
3. PATCH /api/tasks/:id completes task successfully
4. PATCH /api/tasks/:id fails for tasks not assigned to user
5. POST /api/meetings creates meeting with attendees
6. GET /api/meetings/next returns nearest future meeting
7. GET /api/meetings/next returns 204 if no meetings

**Integration Tests Required**:
1. Dashboard statistics calculate correctly
2. Project progress calculation matches task completion ratio
3. Task priority sorting works correctly
4. Meeting visibility respects attendee relationships
5. Rate limiting enforces request limits

---

## Notes

- All endpoints use existing JWT authentication middleware
- Dashboard endpoint combines multiple queries for efficiency
- Meeting endpoints follow RESTful conventions
- Timestamps always in UTC; client handles timezone conversion
- Pagination uses offset/limit pattern (not cursor-based)
- No WebSocket or real-time push notifications in MVP
