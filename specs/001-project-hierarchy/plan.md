# Implementation Plan: Project Hierarchy (Parent / Child Projects)

**Branch**: `001-project-hierarchy` | **Date**: 2026-06-07 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-project-hierarchy/spec.md`

## Summary

Add parent/child relationships to the existing Project entity so that:

- A new project can be created with an optional parent (existing project) or no parent (top-level).
- An existing project can be re-parented to a different parent (or to no parent).
- The project list view renders a hierarchical tree (top-level roots with nested children, expand/collapse).
- Each project's detail page lists its direct child projects.
- The system prevents circular references, prevents deletion of projects that have children, and treats orphan rows as top-level.

The change is localized to:
- `backend/migration/` — one new SQL file adding `parent_id UUID NULL` to `projects` with a self-FK and an index.
- `backend/models/project.go` — add `ParentID *uuid.UUID`, add it to create/update request DTOs, add a new `ProjectNode` tree DTO and a `ProjectChildrenList` DTO.
- `backend/repositories/project_repository.go` — add `ListTree`, `ListChildren(parentID)`, cycle-check helpers, set/clear `parent_id` on create/update, and a guarded delete.
- `backend/services/project_service.go` — add validation (parent exists, not a descendant of self), build the tree from the flat list, decide delete blocking.
- `backend/handlers/project_handler.go` — add `GET /api/projects/tree`, `GET /api/projects/:id/children`, accept `parent_id` on create/update.
- `backend/routes/routes.go` — register the two new GETs.
- `frontend/src/lib/api.js` — new `getProjectTree`, `getProjectChildren`, pass `parent_id` on create/update.
- `frontend/src/components/` — replace the flat project list with a tree component; add a "Child Projects" panel to the project detail page; add a parent-picker to the create/edit form.

The plan keeps the existing layered architecture (handlers → services → repositories), pgx-only data access, Svelte 5 runes on the frontend, and Tailwind styling.

## Technical Context

**Language/Version**: Go 1.25+ (backend, per `backend/AGENTS.md`); JavaScript with Svelte 5 runes (frontend, per root `AGENTS.md`)
**Primary Dependencies**: Fiber v2 (HTTP), pgx v5 (PostgreSQL driver), `github.com/google/uuid`; Svelte 5, Tailwind CSS, `jalali-moment` (frontend). No new external dependencies are introduced by this feature.
**Storage**: PostgreSQL — one schema change: add nullable `parent_id UUID` to the existing `projects` table with a self-referential foreign key `ON DELETE SET NULL` and a B-tree index on `parent_id`.
**Testing**: Manual end-to-end testing only (no test framework currently exists in the repo). Verification will be performed by exercising the flows described in `quickstart.md`.
**Target Platform**: Web application — Linux/Windows server for the API, modern evergreen browsers for the SPA.
**Project Type**: Web application (backend REST API + Svelte SPA). Confirmed by user input and existing `backend/` + `frontend/` directories.
**Performance Goals**: Tree render of ≥50 projects / ≥3 nesting levels in <2s (SC-002); per-node expand/collapse <200ms (SC-006).
**Constraints**: No breaking changes to existing project endpoints' JSON shape — `parent_id` is added as an optional field. Existing projects with `parent_id = NULL` are top-level and must continue to work.
**Scale/Scope**: Single workspace, tens-to-hundreds of projects. No multi-tenancy, no cross-workspace tree. Tree is fully client-renderable (no need for server-side pagination in the foreseeable scope).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Pre-research check** (above). The current `.specify/memory/constitution.md` is the stock template — no project-specific principles or gates have been ratified yet. With no ratified gates, the Constitution Check is vacuously satisfied. The following conventional gates are still respected by this plan because they are encoded in the existing `AGENTS.md` files and the existing code:

| # | Gate (from AGENTS.md) | Status | Evidence in this plan |
|---|----------------------|--------|-----------------------|
| 1 | Backend layering: handlers → services → repositories, never skip | PASS | New endpoint logic lives in `services/project_service.go` and `repositories/project_repository.go`; the handler only wires HTTP. |
| 2 | Backend uses pgx directly, not `database/sql` | PASS | All new repository code uses pgx v5 (`pgxpool`, `Query`/`QueryRow`/`Exec`). |
| 3 | Backend uses UUID primary keys, snake_case columns, PascalCase structs, JSON tags on every field | PASS | `parent_id` is `uuid NULL`; struct field is `ParentID *uuid.UUID` with `json:"parent_id"`. |
| 4 | No DB connection exposed outside `config/`; no hardcoded values | PASS | New repository methods take `pgxpool.Pool` via constructor, same as existing repos. |
| 5 | Never store plaintext secrets / never log secrets | N/A | This feature introduces no secrets. |
| 6 | Svelte 5: `$props()`, `$state`, `$derived`, `$effect` — no `export let`, no TypeScript | PASS | New components use Svelte 5 runes; project remains in JavaScript. |
| 7 | Tailwind for styling, `jalali-moment` for Persian dates | PASS | Tree indentation, expand/collapse chevrons, and child list use Tailwind utilities; any dates reuse the existing `jalali-moment` helper. |
| 8 | Migrations are forward-only SQL files under `backend/migration/` with a numeric prefix | PASS | New file is `009_add_project_parent_id.sql`. |
| 9 | Default admin credentials must be changed after first login | N/A | This feature does not touch the admin user. |
| 10 | Dashboard / new endpoints follow the single-endpoint-aggregation pattern only when serving a dashboard view | PASS | The new `GET /api/projects/tree` is a single endpoint, but it is **not** the dashboard — it is a list view. The dashboard already has its own endpoint and is not affected. |

**Result**: No violations. No complexity-tracking entries required.

### Post-design re-check (after Phase 1)

Re-evaluated against the artefacts produced (`research.md`, `data-model.md`, `contracts/*.md`, `quickstart.md`):

| # | Gate | Post-design status | Evidence |
|---|------|-------------------|----------|
| 1 | Layering: handlers → services → repositories | PASS | All new endpoint logic lives in services and repositories; the handler only wires HTTP. Documented in `data-model.md` §3–4. |
| 2 | pgx only, no `database/sql` | PASS | New repository methods all use pgx v5 (`pgxpool.Pool`); no new SQL driver introduced. |
| 3 | UUID primary keys, snake_case columns, PascalCase structs, JSON tags | PASS | `parent_id` is `uuid NULL`; Go field is `ParentID *uuid.UUID` with `json:"parent_id"`. |
| 4 | No DB connection outside `config/` | PASS | New repository methods follow the existing constructor pattern that takes a `pgxpool.Pool`. |
| 5 | No plaintext secrets / no logged secrets | N/A | This feature introduces no secrets. |
| 6 | Svelte 5 runes (`$props`, `$state`, `$derived`, `$effect`), no `export let`, no TypeScript | PASS | New components (`ProjectTree.svelte`, `ProjectChildrenList.svelte`) are described in `data-model.md` §5 using runes only. |
| 7 | Tailwind for styling, `jalali-moment` for Persian dates | PASS | Tree indentation, chevron, and child-list panels use Tailwind utilities; any dates reuse the existing helper. |
| 8 | Migrations are forward-only SQL files with a numeric prefix | PASS | New file is `backend/migration/009_add_project_parent_id.sql`, idempotent, additive (nullable column + FK + index). |
| 9 | Single-endpoint-aggregation pattern is reserved for dashboard views | PASS | The new `GET /api/projects/tree` and `GET /api/projects/:id/children` are list/detail endpoints, not the dashboard. The dashboard endpoint and its 30-second refresh pattern are untouched. |
| 10 | Default admin credentials must be changed after first login | N/A | This feature does not touch the admin user. |

**Post-design result**: All gates still pass. No new violations introduced by the design. No complexity-tracking entries required. The plan is ready for `/speckit.tasks`.

## Project Structure

### Documentation (this feature)

```text
specs/001-project-hierarchy/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   ├── api-projects.md
│   └── api-projects-tree.md
├── checklists/
│   └── requirements.md  # Already produced by /speckit.specify
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

The project is a web application (frontend + backend), so Option 2 applies. No new top-level directories are introduced; the change is additive within the existing layout.

```text
backend/
├── handlers/
│   └── project_handler.go          # + GET /api/projects/tree, + GET /api/projects/:id/children, accept parent_id on create/update
├── services/
│   └── project_service.go          # + cycle validation, + tree build, + delete blocking
├── repositories/
│   └── project_repository.go       # + ListTree, ListChildren, SetParent, HasChildren, CountDescendantsOf
├── models/
│   └── project.go                  # + ParentID, + ProjectNode (tree DTO), + ProjectChildrenList
├── migration/
│   └── 009_add_project_parent_id.sql  # NEW: parent_id column, FK, index
├── routes/
│   └── routes.go                   # + register the two new GET routes
└── ...

frontend/
└── src/
    ├── components/
    │   ├── ProjectTree.svelte      # NEW: tree view with expand/collapse
    │   ├── ProjectChildrenList.svelte  # NEW: child projects panel
    │   ├── ProjectForm.svelte      # MODIFIED: add parent picker
    │   └── ...
    ├── lib/
    │   └── api.js                  # + getProjectTree, getProjectChildren, pass parent_id on create/update
    └── ...
```

**Structure Decision**: Web application (Option 2). No new modules or top-level directories. All new code lives inside the existing `backend/` and `frontend/src/` trees.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations. This section is intentionally empty.
