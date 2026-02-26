# Research: Project Member Role Assignment

**Feature**: 001-project-member-roles  
**Date**: 2026-02-26  
**Purpose**: Resolve technical decisions for adding project members from system users with predefined project roles.

## Decision 1: Membership persistence model

**Decision**: Introduce a dedicated junction table `project_members` (`project_id`, `user_id`, `project_role_id`, `joined_at`) with a unique constraint on (`project_id`, `user_id`).

**Rationale**:
- Current schema has no project membership relation.
- A dedicated table cleanly supports many-to-many user↔project assignment.
- Unique constraint enforces “no duplicate membership” requirement at DB level.
- Fits existing PostgreSQL + repository SQL pattern.

**Alternatives considered**:
- Add `members` JSON/array field on `projects` (rejected: poor queryability and integrity).
- Reuse task assignments as implicit project membership (rejected: semantically incorrect and incomplete).

## Decision 2: Predefined project roles representation

**Decision**: Add `project_roles` catalog table and reference it from `project_members.project_role_id`.

**Rationale**:
- Requirement calls for predefined roles selectable by admin.
- Table-backed catalog supports active/inactive roles and future expansion without schema rewrites.
- Avoids hardcoding role strings in multiple backend/frontend places.

**Alternatives considered**:
- `CHECK` constraint enum-like text field on membership (rejected: inflexible for role lifecycle changes).
- Reuse `users.role` (`admin/user`) for project role (rejected: system auth role != project membership role).

## Decision 3: Eligible user selection behavior

**Decision**: Eligible users for a project are active users not already present in `project_members` for that project.

**Rationale**:
- Aligns with requirement to select users from system accounts while preventing duplicates.
- Uses existing `users.is_active` field and join/not-exists query patterns.
- Keeps validation enforceable both in API and database constraints.

**Alternatives considered**:
- Include inactive users (rejected: conflicts with “eligible user” semantics).
- Allow duplicate inserts then dedupe in UI (rejected: poor integrity and race-prone).

## Decision 4: API design and authorization

**Decision**: Add project-scoped REST endpoints under `/api/projects/{projectId}/members` and enforce admin-only mutating operations via existing middleware.

**Rationale**:
- Resource nesting naturally models project membership ownership.
- Existing auth middleware already supports `RequireRole("admin")`.
- Keeps route consistency with current project/task endpoint style.

**Alternatives considered**:
- Global `/api/project-members` endpoints (rejected: less discoverable and weaker resource context).
- Frontend-only authorization checks (rejected: security risk).

## Decision 5: Layer integration pattern

**Decision**: Add dedicated handler/service/repository files for project members and roles; do not bypass layers.

**Rationale**:
- Matches mandatory architecture in project conventions.
- Simplifies testability and future changes.
- Keeps project feature boundaries clean.

**Alternatives considered**:
- Add SQL directly in handlers (rejected: violates architecture rules).
- Merge logic into existing `project_repository` only (rejected: bloats existing component and weakens separation).

## Decision 6: Migration strategy

**Decision**: Add a new migration file `008_add_project_members.sql` that creates `project_roles`, seeds default roles, creates `project_members`, and indexes.

**Rationale**:
- Existing repository uses incremental SQL migrations.
- Seeder in migration ensures predefined roles are available immediately.
- Non-destructive migration aligns with current approach.

**Alternatives considered**:
- Runtime auto-create roles on API startup (rejected: hidden side effects).
- Manual SQL setup outside migrations (rejected: non-repeatable environments).
