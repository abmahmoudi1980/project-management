# API Contracts: Task Search Feature

**Phase**: 1 (Design)  
**Date**: 2026-01-01  
**Status**: MVP (No API Changes)  
**Link**: [data-model.md](../data-model.md)

## Overview

The MVP implementation of task search is **frontend-only** and requires **NO new API endpoints** or contract changes. All task data is fetched via the existing task list API, and filtering is performed client-side.

This document outlines:
1. **MVP** (Current): Existing API usage (no changes)
2. **Phase 2** (Optional): Future backend filtering endpoint

---

## MVP: Existing API Usage

### Existing Endpoint: Get Tasks by Project

**Endpoint**: `GET /api/tasks?project_id=<id>&page=<page>&page_size=<size>`

**Request**:
```http
GET /api/tasks?project_id=550e8400-e29b-41d4-a716-446655440000&page=1&page_size=10 HTTP/1.1
Host: localhost:3000
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Response** (200 OK):
```json
{
  "tasks": [
    {
      "id": "task-uuid-1",
      "project_id": "project-uuid",
      "title": "Implement User Authentication",
      "description": "Create login endpoint with JWT",
      "start_date": "2024-03-15T00:00:00Z",
      "due_date": "2024-04-15T00:00:00Z",
      "priority": "High",
      "category": "Backend",
      "completed": false,
      "assignee_id": "user-uuid",
      "author_id": "user-uuid",
      "estimated_hours": 8.5,
      "done_ratio": 25,
      "created_at": "2024-02-20T10:30:00Z",
      "updated_at": "2024-03-01T14:22:00Z"
    },
    {
      "id": "task-uuid-2",
      "project_id": "project-uuid",
      "title": "Write API Documentation",
      "description": "Document all endpoints with examples",
      "start_date": "2024-03-20T00:00:00Z",
      "due_date": "2024-04-01T00:00:00Z",
      "priority": "Medium",
      "category": "Documentation",
      "completed": false,
      "assignee_id": null,
      "author_id": "user-uuid",
      "estimated_hours": null,
      "done_ratio": 0,
      "created_at": "2024-02-25T09:15:00Z",
      "updated_at": "2024-03-10T16:45:00Z"
    }
  ],
  "total": 42,
  "page": 1,
  "page_size": 10,
  "has_more": true
}
```

**Usage in MVP**:
- Frontend calls this endpoint to fetch all tasks for a project
- Response array stored in `$tasks.tasks` reactive store
- Frontend applies search filters in-memory to `$tasks.tasks`
- Filtered results displayed without additional API calls

**Contract Status**: ✅ **UNCHANGED** (no modifications needed)

---

## Phase 2 (Optional): Backend Search Endpoint

### Proposed Endpoint: Search Tasks

**Purpose**: Server-side filtering for improved performance at scale

**Status**: Optional enhancement (not implemented in MVP)

**Endpoint**: `GET /api/tasks/search?project_id=<id>&search=<text>&start_date_from=<date>&start_date_to=<date>&due_date_from=<date>&due_date_to=<date>&page=<page>&page_size=<size>`

**Request**:
```http
GET /api/tasks/search?project_id=550e8400-e29b-41d4-a716-446655440000&search=API&start_date_from=2024-03-01&start_date_to=2024-03-31&page=1&page_size=10 HTTP/1.1
Host: localhost:3000
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Query Parameters**:

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `project_id` | UUID | Yes | Project to search tasks in |
| `search` | string | No | Search term (matches title or description) |
| `start_date_from` | ISO 8601 | No | Minimum start date (inclusive) |
| `start_date_to` | ISO 8601 | No | Maximum start date (inclusive) |
| `due_date_from` | ISO 8601 | No | Minimum due date (inclusive) |
| `due_date_to` | ISO 8601 | No | Maximum due date (inclusive) |
| `page` | integer | No | Page number (default: 1) |
| `page_size` | integer | No | Results per page (default: 10, max: 100) |

**Response** (200 OK):
```json
{
  "tasks": [
    {
      "id": "task-uuid-1",
      "project_id": "project-uuid",
      "title": "Implement User Authentication API",
      "description": "Create login endpoint with JWT",
      "start_date": "2024-03-15T00:00:00Z",
      "due_date": "2024-04-15T00:00:00Z",
      "priority": "High",
      "category": "Backend",
      "completed": false,
      "assignee_id": "user-uuid",
      "author_id": "user-uuid",
      "estimated_hours": 8.5,
      "done_ratio": 25,
      "created_at": "2024-02-20T10:30:00Z",
      "updated_at": "2024-03-01T14:22:00Z"
    }
  ],
  "total": 12,
  "page": 1,
  "page_size": 10,
  "has_more": true
}
```

**Response Codes**:

| Code | Description |
|------|-------------|
| 200 | Success - results returned (may be empty array) |
| 400 | Invalid query parameters (e.g., malformed dates) |
| 401 | Unauthorized - JWT missing or invalid |
| 403 | Forbidden - user lacks project access |
| 500 | Server error |

**Error Response Example** (400 Bad Request):
```json
{
  "error": "invalid date format for start_date_from - expected ISO 8601",
  "details": {
    "parameter": "start_date_from",
    "received": "01/03/1403"
  }
}
```

**Implementation Notes** (if Phase 2 proceeds):
1. **Backend Handler** (`backend/handlers/task_handler.go`):
   - Add new handler function `SearchTasks(c *fiber.Ctx)`
   - Parse query parameters (with validation)
   - Call service layer with filter criteria

2. **Service Layer** (`backend/services/task_service.go`):
   - Add filtering logic using parameters
   - Validate date ranges
   - Implement pagination

3. **Repository Layer** (`backend/repositories/task_repository.go`):
   - Implement `SearchTasks(projectID, search, startFrom, startTo, dueFrom, dueTo, page, pageSize)`
   - Build dynamic SQL WHERE clauses:
     ```sql
     WHERE project_id = $1
       AND (title ILIKE $2 OR description ILIKE $2)
       AND (start_date >= $3 AND start_date <= $4)
       AND (due_date >= $5 AND due_date <= $6)
     ORDER BY due_date, priority
     LIMIT $7 OFFSET $8
     ```

4. **Route Registration** (`backend/routes/routes.go`):
   - Add route: `router.Get("/api/tasks/search", authMiddleware, taskSearchHandler)`

---

## MVP vs Phase 2 Comparison

| Aspect | MVP (Frontend) | Phase 2 (Backend) |
|--------|----------------|-------------------|
| **Location** | Browser memory | Database server |
| **Performance** | O(n) on client | O(1) database index |
| **Network** | Single call for all data | Single call with filters |
| **Scalability** | ~100-1000 tasks OK | Unlimited tasks |
| **Implementation** | Svelte component state | Go handler + service + repo |
| **Latency** | Instant for small lists | Network + query time |
| **Browser Memory** | Higher (all tasks loaded) | Lower (only results) |
| **Effort** | 2-3 hours | 4-6 hours |

**When to Upgrade to Phase 2**:
- Typical project has 5000+ tasks
- User complaints about UI lag with large task lists
- Database query optimization becomes necessary
- Team requests server-side filtering for other clients (mobile app, etc.)

---

## Frontend Integration (MVP)

### Current Usage

**File**: `frontend/src/stores/taskStore.js`

```javascript
const load = async (projectId, reset = false) => {
  // Makes request to existing endpoint
  const response = await api.tasks.getByProject(projectId, 1, pageSize);
  // Response contains all fields needed for frontend filtering
};
```

**No changes required** - store continues using existing endpoint.

---

## No Breaking Changes

✅ MVP implementation requires:
- No API modifications
- No contract changes
- No endpoint additions
- Backward compatible with existing clients

---

## Future-Ready Design

Contract is designed for potential Phase 2 upgrades:
- Query parameter names match field names in Task model
- Date formats use ISO 8601 (consistent with existing API)
- Response structure identical to existing endpoint
- Pagination pattern same as current API

Frontend code can upgrade to Phase 2 without changes:
```javascript
// MVP: Uses client-side filtering
const response = await api.tasks.getByProject(projectId, 1, pageSize);

// Phase 2: Sends filters to backend (same response structure)
const response = await api.tasks.search(projectId, {
  search: "API",
  start_date_from: "2024-03-01",
  // ...
});
```

---

## Summary

**MVP**: Reuses existing `/api/tasks?project_id=<id>` endpoint
**Phase 2**: Adds `/api/tasks/search` endpoint with query parameters
**Breaking Changes**: None
**Migration Path**: Gradual upgrade from client-side to server-side filtering
