# API Contract: Project CRUD Endpoints (modified)

**Feature**: 001-project-hierarchy
**Scope**: Existing `POST /api/projects`, `PUT /api/projects/:id`, `DELETE /api/projects/:id`, plus `GET /api/projects` and `GET /api/projects/:id`. All existing endpoints gain an optional `parent_id` field on request/response. The response shape is additive — clients that do not know about `parent_id` are unaffected.

**Authentication**: All endpoints require a valid session (Bearer token or httpOnly cookie), enforced by the existing `RequireAuth` middleware. No new auth rules are added (see `research.md` R8).

**Content type**: `application/json` for all request and response bodies.

**Error format**: `{"error": "<human-readable message>"}` with an appropriate HTTP status, matching the existing handler convention (`backend/handlers/project_handler.go:53,89,103`).

---

## 1. `GET /api/projects` (modified)

Returns the same flat list of projects as today, but every project in the response now includes a `parent_id` field (omitted when the project is top-level).

**Response 200** — array of `Project` (see §5):

```json
[
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
  }
]
```

> Note: `parent_id` is rendered as `null` for top-level projects in this document for clarity; in the actual JSON it is omitted (`omitempty` on `*uuid.UUID`).

**Errors**:
- `401 Unauthorized` — missing/invalid session.

---

## 2. `GET /api/projects/:id` (modified)

Same as before, but the response `Project` now includes `parent_id` (omitted when top-level).

**Response 200** — `Project` (see §5).

**Errors**:
- `400 Bad Request` — `id` is not a valid UUID.
- `401 Unauthorized`.
- `404 Not Found` — no project with that id.

---

## 3. `POST /api/projects` (modified)

Creates a new project.

**Request body** — `CreateProjectRequest` (see §5):

```json
{
  "title": "Mobile App",
  "description": "iOS and Android client",
  "status": "active",
  "identifier": "MA",
  "is_public": false,
  "parent_id": "8d2…-…-…-…-………",   // optional, omit for top-level
  "start_date": "2026-03-01T00:00:00Z",
  "due_date": "2026-09-30T00:00:00Z"
}
```

**Response 201** — the created `Project` (see §5), with `parent_id` populated or omitted.

**Errors**:
- `400 Bad Request` — invalid body, missing required fields, or:
  - `"parent project not found"` — `parent_id` is set but does not resolve.
- `401 Unauthorized`.

**Validation rules** (existing + new):
- `title` is required and non-empty.
- `identifier` is required, unique, and follows the existing identifier pattern.
- `homepage`, if set, must be a syntactically valid URL (existing rule).
- `parent_id`, if set, must refer to an existing project (NEW).
- On create, no cycle check is required because the new project has no descendants yet (R2).

---

## 4. `PUT /api/projects/:id` (modified)

Updates an existing project. The body shape is `UpdateProjectRequest` (see §5). All fields are sent on every update; the server replaces the stored values for the fields it understands and ignores unknown fields.

**Request body** — example re-parenting a project:

```json
{
  "title": "Mobile App",
  "description": "iOS and Android client",
  "status": "active",
  "identifier": "MA",
  "is_public": false,
  "parent_id": "9ab…-…-…-…-………",   // NEW parent
  "start_date": "2026-03-01T00:00:00Z",
  "due_date": "2026-09-30T00:00:00Z"
}
```

To **clear the parent** (make the project top-level), send `"parent_id": null` in the body.

**Response 200** — the updated `Project` (see §5).

**Errors**:
- `400 Bad Request` — invalid body, missing fields, or:
  - `"parent project not found"` — `parent_id` is set but does not resolve.
  - `"cannot set a project as its own descendant"` — `parent_id` is set to a project that is the target or a descendant of the target (cycle prevention, FR-006, R2).
- `401 Unauthorized`.
- `404 Not Found` — `:id` does not exist.

**Validation rules** (existing + new):
- All create-time rules still apply, except that a cycle check is now performed.
- The cycle check is `IsDescendantOf(targetID, req.ParentID)` — i.e., the proposed new parent must not be a descendant of the project being updated.

---

## 5. `DELETE /api/projects/:id` (modified)

Deletes a project. New behaviour: a project that has one or more child projects cannot be deleted (FR-007, R7). The user must re-parent or delete the children first.

**Response 204** — no body, on success.

**Errors**:
- `400 Bad Request`:
  - `"cannot delete project with N sub-projects; re-parent or delete them first"` — the project has direct children.
- `401 Unauthorized`.
- `404 Not Found` — `:id` does not exist (existing behaviour, surfaced as a 500 in the current code; see the note below).

> **Note on existing behaviour**: The current `DeleteProject` handler returns `500` on any error (`project_handler.go:103`). As part of this feature, the "not found" path is normalised to `404` to match the rest of the handler, since the service layer now distinguishes "not found" from "blocked by children". This is a small, scope-aligned fix; the contract here reflects the new behaviour.

---

## 6. Shared types

### 6.1 `Project`

```json
{
  "id": "uuid",
  "title": "string",
  "description": "string",
  "status": "string",
  "identifier": "string",
  "homepage": "string?",
  "is_public": "bool",
  "user_id": "uuid?",
  "created_by": "uuid?",
  "parent_id": "uuid? (NEW, omitted for top-level)",
  "start_date": "datetime?",
  "due_date": "datetime?",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### 6.2 `CreateProjectRequest` / `UpdateProjectRequest`

Same as `Project` minus `id`, `user_id`, `created_by`, `created_at`, `updated_at`. On update, sending `"parent_id": null` clears the parent; omitting the field leaves the existing parent unchanged.

---

## 7. Backward compatibility

- `parent_id` is `omitempty` on every response, so clients that do not read it see no change.
- `parent_id` is `omitempty` on every request, so clients that do not send it see no change.
- The only new error messages on existing endpoints are the two new 400s on POST/PUT and the new 400 on DELETE; clients that previously saw a 500 on delete will now see a 400 with a clearer message, which is a strict improvement.
