# Quickstart: Project Hierarchy

**Feature**: 001-project-hierarchy
**Audience**: Developer running the feature end-to-end on a local workstation.

This quickstart walks through the steps a developer follows to bring the project-hierarchy feature up locally, exercise the four user stories from `spec.md`, and verify the success criteria.

---

## 1. Prerequisites

- PostgreSQL reachable on the URL expected by the backend (see `backend/AGENTS.md` for env-var loading via `godotenv`).
- Go 1.25+ on the PATH.
- Node 20+ on the PATH (for the Svelte 5 frontend).
- A working copy of this repository on branch `001-project-hierarchy`.

## 2. Bring up the database

```bash
# From repo root
cd backend
go run ./cmd/migrate        # applies all migrations, including the new 009
```

The migration `009_add_project_parent_id.sql` is idempotent. Running it on a database that already has the column is a no-op. On a fresh database it:

- Adds the `parent_id` column to `projects`.
- Adds the self-referential FK with `ON DELETE SET NULL`.
- Creates `idx_projects_parent_id`.
- Records a `COMMENT ON COLUMN` for documentation.

Verify with:

```bash
psql "$DATABASE_URL" -c "\d projects" | grep parent_id
```

You should see `parent_id | uuid |` in the column list and a foreign-key constraint named `fk_projects_parent` at the bottom of the table description.

## 3. Bring up the backend

```bash
# Same shell, still in backend/
go run main.go              # listens on :3000
```

Smoke-test the new endpoint:

```bash
# Log in first to get a session cookie (existing flow, e.g. POST /api/auth/login).
curl -b cookies.txt http://localhost:3000/api/projects/tree
# Expect: {"nodes":[...],"total":N}
```

If `nodes` is empty, the workspace has no projects yet — that is fine, the tree endpoint still returns 200 with `total: 0`.

## 4. Bring up the frontend

```bash
# New shell
cd frontend
npm install
npm run dev                  # listens on :5173
```

Open `http://localhost:5173` and log in.

## 5. Walk the user stories

### US1 — Create a project with a parent (P1)

1. Navigate to **Projects** → **New Project**.
2. Fill in title, identifier, description.
3. In the new **Parent Project** field, start typing the title of an existing project and pick one from the autocomplete.
4. Submit.

**Expect**: the new project is created. The success page (or the project list) shows it nested under the chosen parent.

### US2 — Browse the project list as a tree (P1)

1. Navigate to **Projects**.
2. The page renders a tree, not a flat list. Top-level projects are at the root; child projects are visually indented.
3. Click the chevron on a parent → its children become visible. Click again → they collapse.
4. Click any project name → you land on its detail page.

**Expect**: all four acceptance scenarios in `spec.md` US2 pass.

### US3 — See child projects on a project detail page (P2)

1. Open a project that has at least one child.
2. Scroll to the **Child Projects** section. The section lists the direct children, each name is a link.
3. Open a project that has no children.
4. The **Child Projects** section is rendered with the empty-state message "No sub-projects yet".

**Expect**: both populated and empty paths render correctly.

### US4 — Re-parent an existing project (P3)

1. Open any project's **Edit** form.
2. Change the **Parent Project** field to a different project (or click "No parent (top-level)").
3. Submit.
4. The project now appears under the new parent in the tree, and the old parent's **Child Projects** section no longer lists it.

**Cycle attempt**: while editing a project, try to set its parent to one of its own children. The picker excludes descendants, so the UI prevents the choice; if you bypass the UI and `PUT` directly, the server returns `400` with `{"error":"cannot set a project as its own descendant"}`.

## 6. Verify the edge cases

| Edge case | How to trigger | Expected behaviour |
| --- | --- | --- |
| Delete a project with children | `DELETE /api/projects/:id` on a parent that has children | `400` with `{"error":"cannot delete project with N sub-projects; re-parent or delete them first"}`. The handler returns 400, not 500. |
| Orphan handling (out-of-band) | Manually `DELETE FROM projects WHERE id = '<parent-id>'` in psql | All children get `parent_id = NULL` automatically (the FK uses `ON DELETE SET NULL`). They reappear as top-level projects in the tree on next page load. |
| Project with no children | Open such a project | **Child Projects** section is rendered with the empty-state message, not hidden. |
| Deeply nested tree | Create a chain A → B → C → D | All four appear in the tree; expanding every level reaches D in two clicks per level (SC-007). |
| No projects at all | Fresh database | Tree page renders an empty state with a call-to-action to create the first project. |
| Detail page of a deleted project | Open a project, then delete it from another tab | Detail page renders a "project not found" state. |

## 7. Verify the success criteria

| SC | How to verify | Pass criterion |
| -- | ------------- | -------------- |
| SC-001 | Time creating a sub-project with a parent selected and all required fields filled. | ≤ 60 seconds wall-clock. |
| SC-002 | Seed ≥ 50 projects in a ≥ 3-level hierarchy. Reload the list page. | First paint of the tree in < 2 seconds. |
| SC-003 | Spot-check 10 parent-child pairs. | Every child appears under its parent in the tree, and every parent that has children lists them on the detail page. |
| SC-004 | Try to set a project's parent to one of its descendants (UI and raw API). | 100% rejected with `400` and the `cannot set a project as its own descendant` message. |
| SC-005 | Try to `DELETE` a project with children. | 100% rejected with `400` and the `cannot delete project with N sub-projects` message. |
| SC-006 | Click expand on a parent with many children, then collapse. | Per-node response < 200 ms (subjectively instant; client-side). |
| SC-007 | Navigate from the tree to a deep descendant. | ≤ 2 clicks per level. |

## 8. Rollback

If the feature needs to be rolled back:

1. Revert the code changes.
2. The migration is **additive** (a nullable column with a self-FK and an index). It can be left in place without breaking the application, because the application code is the only thing that writes to `parent_id`. To remove the column:

   ```sql
   ALTER TABLE projects DROP CONSTRAINT IF EXISTS fk_projects_parent;
   DROP INDEX IF EXISTS idx_projects_parent_id;
   ALTER TABLE projects DROP COLUMN IF EXISTS parent_id;
   ```

   This is destructive — every project's parent assignment is lost. Do it only if rolling back the feature entirely.
