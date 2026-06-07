# Phase 0 — Research: Project Hierarchy

**Feature**: 001-project-hierarchy
**Date**: 2026-06-07

This document records the technical decisions made for the project hierarchy feature, the rationale, and the alternatives that were considered. All "NEEDS CLARIFICATION" items from the Technical Context have been resolved by reference to the existing codebase (root `AGENTS.md`, `backend/AGENTS.md`, the project source) and to common patterns used elsewhere in the system.

---

## R1. Storage layout for the parent reference

**Decision**: Add a single nullable column `parent_id UUID` to the existing `projects` table, with a self-referential foreign key to `projects(id) ON DELETE SET NULL` and a B-tree index on `parent_id`.

**Rationale**:
- The hierarchy is a forest (a project has at most one parent; many top-level projects are allowed). The classic SQL representation of a tree in a single table is a self-referential FK on a nullable `parent_id` — this is what the spec assumes and what every consumer of the data (tree view, detail page, dashboard) will read.
- `ON DELETE SET NULL` is the correct behaviour for FR-008 (orphan handling): if a parent row is removed by any out-of-band path, every child row is automatically promoted to a top-level row, so no data becomes inaccessible.
- A B-tree index on `parent_id` is required because both new endpoints (`ListChildren` and the tree builder's "all projects grouped by parent" query) are filtered by `parent_id` and must stay fast as the project count grows. SC-002 (50 projects / 3 levels in <2s) is easy to meet with this index, even without it the query plan is `Index Scan using idx_projects_parent_id`.

**Alternatives considered**:
- *Adjacency list with a separate `project_relations` table* (m:n parents). Rejected because the spec explicitly assumes a single parent (see Assumptions in `spec.md`). Adds complexity for no benefit.
- *Materialized path (`/A/B/C/`)* or *nested sets* (lft/rgt). Rejected because the operations required by the spec are all `O(1)` reads of "children of X" and `O(1)` writes of "set parent of X", which the adjacency list does in a single UPDATE. Materialized paths and nested sets are appropriate for deep, read-mostly trees that need fast subtree queries — neither is required here.
- *Closure table* (`project_ancestors`). Rejected because it adds a second table to keep in sync on every write, and the only "is X a descendant of Y?" question (used for cycle detection) is solvable with a single recursive CTE without persisting ancestors.

## R2. Cycle prevention strategy

**Decision**: In the service layer, before persisting a `parent_id` change, run a recursive CTE that returns `TRUE` if and only if the proposed parent is a descendant of the project being moved. If true, reject the request with HTTP 400 and a clear message.

**Rationale**:
- Cycles cannot be created by `ON DELETE SET NULL` (the only DB-level cascade), so the only code path that can introduce a cycle is `UpdateProject` setting a new `parent_id`. Validating at the service layer keeps the repository thin and gives us a clean `ErrValidation` to map to 400 in the handler.
- A single recursive CTE (`WITH RECURSIVE descendants AS (...) SELECT 1 FROM descendants WHERE id = $1`) is portable PostgreSQL, runs in milliseconds even for very deep trees, and requires no extra table.
- The check must run on the **proposed** state, not the current state — i.e., we check "is the proposed parent a descendant of the project in the database right now?" That is exactly what the CTE does.

**Alternatives considered**:
- *App-level walk of `parent_id` chain in Go* (load project, load its parent, load its parent, …). Rejected: requires N round-trips and is unbounded by depth. Recursive CTE is one round-trip.
- *DB trigger that raises an exception on cycle*. Rejected: hides validation behind an opaque PG error code and forces the handler to translate it to a 400 with a useful message. Service-layer check is more debuggable and keeps the error vocabulary consistent with the rest of the codebase (every other validation error in `project_service.go` already returns a 400 with a string message).

## R3. Tree-shape representation for the client

**Decision**: The server returns a single flat list of all projects (the same `Project` shape already used by `GET /api/projects`), and the client builds the tree in the browser by `Map<parentID, Project[]>` grouping. A new convenience endpoint `GET /api/projects/tree` returns the same flat list under the `nodes` key plus an empty `children` slot on each item so the client can decide whether to fold it in or not.

**Rationale**:
- The dashboard and the existing `GetProjectsByUser` already operate on flat lists. Reusing the same shape avoids introducing a second, parallel "tree DTO" that has to be kept in sync with `Project`.
- The expected scale (≤ a few hundred projects, single workspace) makes client-side tree building free.
- Returning the flat list, rather than a pre-nested tree, lets the client reuse the existing project card / row UI for each node in the tree without a special-cased renderer.

**Alternatives considered**:
- *Server returns a pre-nested tree of `ProjectNode { …, children: ProjectNode[] }`*. Considered because it would let the client render with a single recursive component. Rejected because it duplicates the `Project` shape, makes the response harder to consume for non-tree views, and offers no measurable benefit at the expected scale. Kept available as a future option if the project count ever grows past ~1000.
- *Server returns `{roots: ProjectNode[]}` only (drops top-level projects that have a parent from the response)*. Rejected because it would force the client to make a second call to fetch a child's subtree when expanding a node.

## R4. Tree rendering and expand/collapse in Svelte 5

**Decision**: A new `ProjectTree.svelte` component holds a `Set<uuid>` of expanded node IDs in `$state`, recursively renders each node, and uses a single `<button aria-expanded={isOpen}>` for the chevron. Initial state expands only the first level by default.

**Rationale**:
- Svelte 5 runes (`$state`, `$derived`) are the project's mandated pattern (root `AGENTS.md`).
- A `Set` is the right shape for expand state because membership check is O(1) and the size is bounded by the number of projects in the tree.
- Defaulting to "first level expanded, deeper levels collapsed" matches the common project-management UX and keeps the initial render fast even with a large tree (SC-002).
- `aria-expanded` on the chevron button is included for accessibility (existing project page already uses semantic HTML).

**Alternatives considered**:
- *Persist expand state in `localStorage`*. Rejected as out of scope for this feature; the spec does not require it. Can be added later as a small follow-up.
- *Render all nodes flat with CSS indent based on depth (no recursion)*. Rejected because expand/collapse of nested nodes is much harder to express without recursive components.

## R5. Parent picker in the create/edit form

**Decision**: Reuse the same autocomplete-style text input pattern used elsewhere in the project form (the existing form already has a similar picker for project members and identifier suggestions). The picker shows project title + identifier, supports a "no parent (top-level)" clear button, and excludes the current project and its descendants when editing (so the user cannot accidentally pick a child as the new parent).

**Rationale**:
- Consistent UX with the rest of the form.
- Excluding self-and-descendants in the picker is the client-side companion of the server-side cycle check (R2) — the user is guided away from an invalid choice, and the server still rejects it as a defence in depth.
- Empty selection is explicitly supported with a clear button; that maps directly to FR-001 and FR-002 ("leave the parent unset").

**Alternatives considered**:
- *Native `<select>` element with all projects listed*. Rejected because it scales poorly past a few dozen projects and does not support type-ahead search.
- *Separate "Move to..." modal opened from the project detail page*. Rejected because the spec keeps re-parenting inside the normal edit flow (US4 acceptance scenarios are written against the edit form, not a modal).

## R6. Index, migration safety, and rollback

**Decision**: Migration is `009_add_project_parent_id.sql` and does the following in one transaction:
1. `ALTER TABLE projects ADD COLUMN IF NOT EXISTS parent_id UUID;`
2. `ALTER TABLE projects ADD CONSTRAINT fk_projects_parent FOREIGN KEY (parent_id) REFERENCES projects(id) ON DELETE SET NULL;`
3. `CREATE INDEX IF NOT EXISTS idx_projects_parent_id ON projects(parent_id);`
4. `COMMENT ON COLUMN projects.parent_id IS 'Optional parent project. NULL means top-level project. ON DELETE SET NULL keeps orphaned rows accessible.';`

**Rationale**:
- `IF NOT EXISTS` on the column and index makes the migration idempotent, which matches the pattern used in earlier migrations (`005_add_dashboard_meetings.sql` uses the same `IF NOT EXISTS` clauses).
- `ON DELETE SET NULL` is required for FR-008.
- The index is required for the new query patterns (R1).
- The migration is forward-only and adds a nullable column with no `NOT NULL` and no default, so it is safe to run on a populated `projects` table — every existing row gets `parent_id = NULL` and therefore becomes a top-level project, which preserves all current behaviour.

**Alternatives considered**:
- *Backfill existing rows from a "project group" or "module" column*. Rejected because no such column exists in the current schema; the existing `projects` table is flat.
- *Splitting the migration into two files (column, then FK)*. Rejected: the FK is what gives us orphan-handling (R1), and running them in one transaction is simpler and equally safe.

## R7. Delete blocking for projects with children

**Decision**: In the service layer, before calling `DELETE`, run `SELECT 1 FROM projects WHERE parent_id = $1 LIMIT 1`. If a row is found, return `models.ErrValidation` with a message naming the count of direct children and instructing the user to delete or re-parent them first.

**Rationale**:
- This is the simplest expression of FR-007 and it runs in O(1) thanks to `idx_projects_parent_id`.
- It runs in the service layer (not the repository) so the message can be enriched ("3 sub-projects must be removed or re-parented first") rather than a generic "constraint violation".
- Combined with R6's `ON DELETE SET NULL`, this means a *direct* delete of a project with children is rejected with a clear 400, but a *cascading* delete (e.g., from a future admin tool that drops the FK) would not silently destroy data — children would simply become top-level.

**Alternatives considered**:
- *Let the database do it with a trigger or a `CHECK` constraint*. Rejected because it produces opaque error codes and bypasses the existing `ErrValidation` pattern used throughout the service layer.
- *Cascade delete (`ON DELETE CASCADE`) of children*. Explicitly rejected by the spec (FR-007) — the user must consciously re-parent or remove children.

## R8. Authorization

**Decision**: No new authorization rules. The existing `RequireAuth` middleware applies to all new endpoints (`GET /api/projects/tree`, `GET /api/projects/:id/children`, `PUT /api/projects/:id`) and uses the same `GetProjectsByUser` access pattern. Re-parenting inherits the same authorization as updating any other project field.

**Rationale**:
- The spec's Assumptions section explicitly says "Any user who has permission to create or edit a project is also allowed to set or change its parent to any other project; finer-grained 'who can be a parent of whom' rules are out of scope for this feature."
- Adding role-based "who can be a parent of whom" would be speculative and would require touching the role/permission model, which is a separate concern.

**Alternatives considered**:
- *Restrict re-parenting to the project owner or an admin*. Rejected: out of scope per the spec; can be added later as a separate feature.

---

## Summary of resolved unknowns

All Technical Context items have been resolved. There are no remaining `NEEDS CLARIFICATION` markers. The plan in `plan.md` can proceed to Phase 1 (data model, contracts, quickstart) without further research.
