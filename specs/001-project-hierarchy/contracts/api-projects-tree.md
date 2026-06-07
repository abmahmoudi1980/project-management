# API Contract: New Endpoints (project tree and child list)

**Feature**: 001-project-hierarchy
**Scope**: Two new GET endpoints that support the tree view and the "Child Projects" panel on the project detail page.

**Authentication**: Both endpoints require a valid session, enforced by the existing `RequireAuth` middleware. They use the same role-based access pattern as `GET /api/projects` (the existing `GetProjectsByUser` logic decides which projects the caller can see).

**Content type**: `application/json`.

**Error format**: `{"error": "<human-readable message>"}` with an appropriate HTTP status, matching the existing handler convention.

---

## 1. `GET /api/projects/tree`

Returns every project the caller is allowed to see, as a flat list sorted for tree grouping on the client (`parent_id NULLS FIRST, title ASC`). The client is responsible for folding the flat list into a tree by `parent_id`. See `research.md` R3 for why the server does not return a pre-nested structure.

**Response 200** — `ProjectTree`:

```json
{
  "nodes": [
    {
      "id": "8d2…",
      "title": "Product Development",
      "description": "All product work",
      "status": "active",
      "identifier": "PD",
      "is_public": false,
      "parent_id": null,
      "created_at": "2026-01-15T10:00:00Z",
      "updated_at": "2026-01-20T10:00:00Z"
    },
    {
      "id": "1f4…",
      "title": "Mobile App",
      "description": "iOS and Android",
      "status": "active",
      "identifier": "MA",
      "is_public": false,
      "parent_id": "8d2…",
      "created_at": "2026-02-01T10:00:00Z",
      "updated_at": "2026-02-01T10:00:00Z"
    },
    {
      "id": "c11…",
      "title": "iOS Client",
      "status": "active",
      "identifier": "IOS",
      "is_public": false,
      "parent_id": "1f4…",
      "created_at": "2026-02-10T10:00:00Z",
      "updated_at": "2026-02-10T10:00:00Z"
    }
  ],
  "total": 3
}
```

- The `nodes` list contains **every** project in the caller's scope, including children of children.
- `parent_id` is omitted for top-level projects.
- `total == nodes.length`.
- Sort order: `parent_id NULLS FIRST, title ASC` — root projects first, siblings grouped together.

**Errors**:
- `401 Unauthorized` — missing/invalid session.

**Performance**: Single query with an `ORDER BY parent_id NULLS FIRST, title ASC`. Uses the existing `idx_projects_parent_id` index (R6). Targets <2s for 50 projects / 3 levels (SC-002).

---

## 2. `GET /api/projects/:id/children`

Returns the **direct** children of a single project, sorted by `title ASC`. Used by the "Child Projects" section on the project detail page (FR-005, US3). Grandchildren are not included; the user navigates to a child to see its own children.

**Response 200** — `ProjectChildrenList`:

```json
{
  "project_id": "8d2…",
  "count": 2,
  "children": [
    {
      "id": "1f4…",
      "title": "Mobile App",
      "description": "iOS and Android",
      "status": "active",
      "identifier": "MA",
      "is_public": false,
      "parent_id": "8d2…",
      "created_at": "2026-02-01T10:00:00Z",
      "updated_at": "2026-02-01T10:00:00Z"
    },
    {
      "id": "5b7…",
      "title": "Web Platform",
      "status": "active",
      "identifier": "WEB",
      "is_public": false,
      "parent_id": "8d2…",
      "created_at": "2026-02-05T10:00:00Z",
      "updated_at": "2026-02-05T10:00:00Z"
    }
  ]
}
```

- `count == children.length`.
- `children` may be empty: in that case the section on the detail page renders an empty-state message ("No sub-projects yet"). The endpoint itself still returns 200 with `count: 0` and `children: []` — the 404 path is reserved for "the parent project itself does not exist".

**Errors**:
- `400 Bad Request` — `:id` is not a valid UUID.
- `401 Unauthorized`.
- `404 Not Found` — no project with that id (so the caller can render a "project not found" state on the detail page, per the Edge Cases in `spec.md`).

**Performance**: Single `SELECT … WHERE parent_id = $1 ORDER BY title ASC` against `idx_projects_parent_id`. No recursion needed.

---

## 3. Notes for the frontend

- The list page should call `getProjectTree()` once on mount and then use the local expand/collapse state. No need to refetch the tree to expand a node (R3).
- The detail page should call `getProjectChildren(projectId)` on mount. No need to subscribe to changes from other users for this iteration — a manual reload is sufficient. (Realtime updates are out of scope for this feature; the existing 30-second dashboard refresh pattern is the precedent for that kind of behaviour but is not applied here.)
- A project that exists but has no children returns 200 with `count: 0` — the frontend must render an empty state in that case.
- A project that does not exist returns 404 — the frontend must render a "project not found" state.
