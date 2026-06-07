# Data Model: Project Hierarchy

**Feature**: 001-project-hierarchy
**Date**: 2026-06-07
**Source spec**: [spec.md](./spec.md)
**Source research**: [research.md](./research.md)

This document describes the persistent and in-memory data structures that implement the project hierarchy feature. It is the source of truth for the migration, the Go models, and the API contracts.

---

## 1. Persistent model

### 1.1 Table: `projects` (modified)

The existing `projects` table gains one new column, one new constraint, and one new index. No existing column is renamed, dropped, or has its type changed. All existing rows are valid under the new schema (they all get `parent_id = NULL`).

| Column      | Type           | Nullability | Default | Notes                                                                                 |
| ----------- | -------------- | ----------- | ------- | ------------------------------------------------------------------------------------- |
| `parent_id` | `UUID`         | NULL        | —       | Self-FK to `projects(id) ON DELETE SET NULL`. NULL means top-level (root) project.   |

The full set of new DDL is in the migration file `backend/migration/009_add_project_parent_id.sql` and is reproduced here for reference:

```sql
ALTER TABLE projects
    ADD COLUMN IF NOT EXISTS parent_id UUID;

ALTER TABLE projects
    ADD CONSTRAINT fk_projects_parent
    FOREIGN KEY (parent_id) REFERENCES projects(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_projects_parent_id ON projects(parent_id);

COMMENT ON COLUMN projects.parent_id IS
    'Optional parent project. NULL means top-level. ON DELETE SET NULL keeps orphans accessible.';
```

### 1.2 Invariants enforced by the database

- A project's `parent_id`, if set, refers to a row that exists in `projects(id)`.
- If the referenced parent row is deleted, every child's `parent_id` is set to `NULL` automatically (FR-008 / R1 / R6).
- `parent_id = id` is **not** prevented at the DB level. The application layer is the source of truth for cycle prevention (R2) because the check is cheaper and produces a better error message than a deferred constraint trigger.

### 1.3 Invariants enforced by the application

- **No cycles.** Setting `parent_id = X` is rejected if X is the project itself or any of its descendants (R2). Implemented in `services/project_service.go` with a recursive CTE before the UPDATE.
- **Parent must exist** at the moment of insert/update. A `parent_id` that does not resolve to an existing row is rejected with a 400.
- **No deletion with children.** `DELETE FROM projects WHERE id = $1` is rejected by the service layer if any row has `parent_id = $1` (FR-007 / R7). The error message names the count of direct children.

---

## 2. Go domain model

All changes are in `backend/models/project.go`. Existing fields and JSON tags are preserved verbatim; new fields are additive and optional (use `omitempty` so existing clients that do not know about `parent_id` are unaffected).

### 2.1 `Project` (modified)

```go
type Project struct {
    ID          uuid.UUID  `json:"id"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      string     `json:"status"`
    Identifier  string     `json:"identifier"`
    Homepage    *string    `json:"homepage,omitempty"`
    IsPublic    bool       `json:"is_public"`
    UserID      *uuid.UUID `json:"user_id,omitempty"`
    CreatedBy   *uuid.UUID `json:"created_by,omitempty"`
    ParentID    *uuid.UUID `json:"parent_id,omitempty"` // NEW
    StartDate   *time.Time `json:"start_date,omitempty"`
    DueDate     *time.Time `json:"due_date,omitempty"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`
}
```

`ParentID` is a `*uuid.UUID` (nullable) — `nil` means top-level project.

### 2.2 `CreateProjectRequest` (modified)

```go
type CreateProjectRequest struct {
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      string     `json:"status"`
    Identifier  string     `json:"identifier"`
    Homepage    *string    `json:"homepage,omitempty"`
    IsPublic    bool       `json:"is_public"`
    ParentID    *uuid.UUID `json:"parent_id,omitempty"` // NEW
    StartDate   *time.Time `json:"start_date,omitempty"`
    DueDate     *time.Time `json:"due_date,omitempty"`
}
```

`ParentID == nil` → top-level. Validation order: validate title/identifier/URL → check that `ParentID` (if set) refers to an existing project → check that the request does not introduce a cycle (only relevant on update; on create the project has no descendants yet, so cycle is impossible).

### 2.3 `UpdateProjectRequest` (modified)

```go
type UpdateProjectRequest struct {
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      string     `json:"status"`
    Identifier  string     `json:"identifier"`
    Homepage    *string    `json:"homepage,omitempty"`
    IsPublic    bool       `json:"is_public"`
    ParentID    *uuid.UUID `json:"parent_id,omitempty"` // NEW — may be nil to clear parent
    StartDate   *time.Time `json:"start_date,omitempty"`
    DueDate     *time.Time `json:"due_date,omitempty"`
}
```

A non-nil `ParentID` sets the new parent. A nil `ParentID` (or omitted from the JSON) clears the parent and makes the project a top-level project.

### 2.4 `ProjectChildrenList` (new)

Response shape for `GET /api/projects/:id/children`.

```go
type ProjectChildrenList struct {
    ProjectID uuid.UUID `json:"project_id"`
    Count     int       `json:"count"`
    Children  []Project `json:"children"`
}
```

`Children` is sorted by `title ASC` for a stable UI. The list contains **direct** children only — grandchildren are not inlined (FR-005, US3 acceptance scenario 3).

### 2.5 `ProjectTree` (new)

Response shape for `GET /api/projects/tree`. The server returns a flat list of all projects in the user's accessible scope; the client builds the tree by grouping on `parent_id`. The shape intentionally mirrors `Project` so the existing project row UI can be reused.

```go
type ProjectTree struct {
    Nodes  []ProjectTreeNode `json:"nodes"`
    Total  int               `json:"total"`
}

type ProjectTreeNode struct {
    Project
    // Children is left empty in the wire response — the client folds the flat
    // list into a tree by `parent_id` and ignores this field. It is present
    // in the type so a future server-side nested response can be added
    // without a breaking change.
    Children []ProjectTreeNode `json:"children,omitempty"`
}
```

The flat list is sorted by `parent_id NULLS FIRST, title ASC` so root projects come first and each parent's children are contiguous, which makes the client's grouping step a single linear pass.

---

## 3. Repository contract

New methods on `repositories/project_repository.go`. All take a `context.Context` and a `pgxpool.Pool` (already used by the rest of the repository), per the project's pgx-only convention.

| Method | Purpose | Returns |
| ------ | ------- | ------- |
| `ListTree(ctx) ([]Project, error)` | Fetch every project the caller is allowed to see, sorted for tree grouping. | Flat list of `Project` with `parent_id` populated. |
| `ListChildren(ctx, parentID uuid.UUID) ([]Project, error)` | Fetch direct children of one project, sorted by title. | `[]Project` (may be empty). |
| `HasChildren(ctx, parentID uuid.UUID) (bool, int, error)` | Used by the service to decide whether to block a delete. | `(hasChildren, count, err)`. |
| `IsDescendantOf(ctx, ancestorID, candidateID uuid.UUID) (bool, error)` | Returns `true` if `candidateID` is a descendant of `ancestorID`, including being the same row. Used by the cycle check. | `(bool, error)`. Implemented with `WITH RECURSIVE descendants AS (...)`. |

Modified methods:

| Method | Change |
| ------ | ------ |
| `Create(ctx, req CreateProjectRequest, userID *uuid.UUID) (*Project, error)` | Insert `parent_id` from `req.ParentID` if set. |
| `Update(ctx, id uuid.UUID, req UpdateProjectRequest) (*Project, error)` | Update `parent_id` from `req.ParentID`. `nil` clears the parent. |
| `Delete(ctx, id uuid.UUID) error` | No change at the repository layer. Blocking is done in the service layer using `HasChildren` before calling this. |

---

## 4. Service-layer behaviour

New validation rules live in `services/project_service.go` and are applied in this order:

1. **On `CreateProject`**: if `req.ParentID` is set, call `repo.GetByID(req.ParentID)`; return `ErrValidation` with message "parent project not found" if it does not exist. No cycle check is needed on create (the new project has no descendants).
2. **On `UpdateProject`**: if `req.ParentID` is set, call `repo.IsDescendantOf(targetID, req.ParentID)`. If `true`, return `ErrValidation` with message "cannot set a project as its own descendant". If `req.ParentID` is `nil`, no check is required (clearing is always safe).
3. **On `DeleteProject`**: call `repo.HasChildren(id)`. If `hasChildren`, return `ErrValidation` with message "cannot delete project with N sub-projects; re-parent or delete them first". The count is included so the user knows what to do.

All errors surface as HTTP 400 with `{"error": "<message>"}`, matching the existing handler convention (`project_handler.go:53,89`).

---

## 5. Frontend model (Svelte 5)

### 5.1 `ProjectTree` component state

```js
// frontend/src/components/ProjectTree.svelte
let { projects = [] } = $props(); // flat list of projects with parent_id
let expanded = $state(new Set());   // Set<string> of project IDs that are open

let tree = $derived(buildTree(projects, expanded));

function buildTree(projects, expanded) {
  // 1. Group projects by parent_id (parent_id === null → roots).
  // 2. For each project, decide whether to render its children based on `expanded`.
  // 3. Return a nested structure: [{ project, children: [...] }, ...]
}
```

The first level is expanded by default on mount (`$effect(() => { expanded = new Set(projects.filter(p => p.parent_id === null).map(p => p.id)) })`).

### 5.2 `ProjectChildrenList` component

Receives `projectId` as a prop, fetches `GET /api/projects/:id/children` once on mount, and renders a Tailwind list. Empty state renders the message "No sub-projects yet".

### 5.3 Parent picker

Added to the existing project form. When editing, the picker excludes the current project and its descendants (so a cycle is impossible from the UI) and provides a "no parent (top-level)" clear button. On create, the picker excludes nothing because the new project has no descendants yet.

---

## 6. Summary of schema and code surface

| Layer | Surface | Type of change |
| ----- | ------- | -------------- |
| Database | `projects.parent_id` column, FK, index | New (migration 009) |
| Go model | `Project.ParentID`, `CreateProjectRequest.ParentID`, `UpdateProjectRequest.ParentID` | Modified (additive) |
| Go model | `ProjectChildrenList`, `ProjectTree`/`ProjectTreeNode` | New |
| Repository | `ListTree`, `ListChildren`, `HasChildren`, `IsDescendantOf` | New |
| Repository | `Create`, `Update` | Modified (write `parent_id`) |
| Service | `CreateProject`, `UpdateProject`, `DeleteProject` | Modified (validation rules) |
| Handler | `GetAllProjects`, `CreateProject`, `UpdateProject` | Modified (pass `parent_id` through) |
| Handler | `GetProjectTree`, `GetProjectChildren` | New |
| Routes | `GET /api/projects/tree`, `GET /api/projects/:id/children` | New |
| Frontend | `api.js` — `getProjectTree`, `getProjectChildren`, `parent_id` on create/update | Modified (additive) |
| Frontend | `ProjectTree.svelte`, `ProjectChildrenList.svelte` | New |
| Frontend | `ProjectForm.svelte` (or equivalent) | Modified (add parent picker) |
| Frontend | project list page | Modified (replace flat list with `<ProjectTree>`) |
| Frontend | project detail page | Modified (mount `<ProjectChildrenList>`) |

No new external dependency is introduced at any layer.
