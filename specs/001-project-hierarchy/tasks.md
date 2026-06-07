---
description: "Task list for Project Hierarchy (001-project-hierarchy)"
---

# Tasks: Project Hierarchy (Parent / Child Projects)

**Input**: Design documents from `/specs/001-project-hierarchy/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/, quickstart.md

**Tests**: The spec explicitly states the project uses manual end-to-end testing only — no test framework currently exists. **No test tasks are generated.** Validation follows the steps in `quickstart.md`.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story. US1 and US2 share priority P1 and form the MVP together; US3 and US4 build on top of them.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: Which user story this task belongs to (e.g. US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

This is a **Web app**: backend lives in `backend/`, frontend lives in `frontend/src/`. All paths below use the existing layered structure (handlers → services → repositories) on the backend and Svelte 5 runes on the frontend.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm a clean baseline before making changes.

- [x] T001 Run `go build` from `backend/` and confirm the project compiles cleanly before any changes

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Schema change, model field, and write paths that MUST be complete before ANY user story can be implemented. Every user story reads `parent_id` on the project and (where relevant) writes it.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [x] T002 Create migration `backend/migration/009_add_project_parent_id.sql` with `ALTER TABLE projects ADD COLUMN parent_id UUID`, the `fk_projects_parent` self-FK (`ON DELETE SET NULL`), the `idx_projects_parent_id` index, and a `COMMENT ON COLUMN` (idempotent: `IF NOT EXISTS` everywhere, per the pattern in `005_add_dashboard_meetings.sql`)
- [x] T003 Add `ParentID *uuid.UUID` field with `json:"parent_id,omitempty"` to `Project`, `CreateProjectRequest`, and `UpdateProjectRequest` in `backend/models/project.go`
- [x] T004 Update `ProjectRepository.GetAll` and `GetByID` in `backend/repositories/project_repository.go` to `SELECT parent_id` and `Scan` it into the new `ParentID` field
- [x] T005 Update `ProjectRepository.Create` in `backend/repositories/project_repository.go` to include `parent_id` in the `INSERT … VALUES` and the `RETURNING` clause
- [x] T006 Update `ProjectRepository.Update` in `backend/repositories/project_repository.go` to include `parent_id` in the `UPDATE … SET` and the `RETURNING` clause

**Checkpoint**: Foundation ready. The `projects` table has the new column, the Go model carries `parent_id` through every read and write. User story implementation can now begin.

---

## Phase 3: User Story 1 - Create a Project with a Parent (Priority: P1) 🎯 MVP

**Goal**: A user can create a new project and either pick an existing project as its parent or leave the parent blank to create a top-level project. The system rejects self-parent at the database / handler boundary.

**Independent Test**: Create a new project with a parent selected; verify it appears nested under that parent in the project list. Create another with no parent; verify it appears as a top-level project. Attempting to create a project whose parent equals its own ID is impossible on create (the ID is server-generated), so the cycle check is exercised in US4.

### Implementation for User Story 1

- [x] T007 [US1] In `ProjectService.CreateProject` (`backend/services/project_service.go`), after the existing identifier/URL validation and before `repo.Create`, if `req.ParentID != nil` call `repo.GetByID(ctx, *req.ParentID)`; return `ErrValidation` with message `"parent project not found"` if the result is `nil`
- [x] T008 [US1] Update `createProject` in `frontend/src/lib/api.js` to forward `parent_id` from the caller's payload (no other change to the function signature)
- [x] T009 [US1] Add a parent-picker field to `ProjectForm.svelte` (`frontend/src/components/ProjectForm.svelte`): type-ahead / autocomplete over the existing `getProjects()` list, with a "No parent (top-level)" clear button. In create mode, the picker lists all projects with no exclusion. Store the chosen value in the form's reactive state and include it in the submitted payload as `parent_id` (or omit when the clear button is used)

**Checkpoint**: US1 is fully functional. A user can create a top-level or sub-project from the form. The "sub-project appears nested under the parent" part is only visible end-to-end after US2 is wired in, but the data is correct from this checkpoint onwards.

---

## Phase 4: User Story 2 - Browse Projects in a Tree Structure (Priority: P1) 🎯 MVP

**Goal**: The project list page renders projects as a hierarchical tree (top-level roots with nested children, expand/collapse controls) and every project name is a link to its detail page.

**Independent Test**: With a workspace containing a multi-level hierarchy, the list page renders the tree with correct indentation, expand/collapse works, and the deepest descendants are reachable without leaving the page. All four acceptance scenarios in `spec.md` US2 pass.

### Implementation for User Story 2

- [ ] T010 [P] [US2] Add `ProjectTree` (`{Nodes []ProjectTreeNode; Total int}`) and `ProjectTreeNode` (embeds `Project`; adds `Children []ProjectTreeNode` with `json:"children,omitempty"`) DTOs to `backend/models/project.go`
- [ ] T011 [P] [US2] Add `ProjectRepository.ListTree(ctx) ([]Project, error)` in `backend/repositories/project_repository.go` returning every project the caller can see, ordered `ORDER BY parent_id NULLS FIRST, title ASC` (uses `idx_projects_parent_id`)
- [ ] T012 [US2] Add `ProjectService.GetProjectTree(ctx) (*ProjectTree, error)` in `backend/services/project_service.go` that calls `repo.ListTree` and wraps the result in `ProjectTree{Nodes, Total: len(Nodes)}` — sorting is already done in the repository
- [ ] T013 [US2] Add `ProjectHandler.GetProjectTree` in `backend/handlers/project_handler.go`: extract the user from the auth context, call `service.GetProjectsByUser` (or a new service method that reuses it for the role-based filter), and return the projects wrapped in a flat `ProjectTree` payload. JSON: `{"nodes": [...], "total": N}`
- [ ] T014 [US2] Register `GET /api/projects/tree` in `backend/routes/routes.go`, inside the existing `projects := api.Group("/projects", middleware.RequireAuth)` group. **Order matters**: register `/tree` BEFORE `/:id` so Fiber does not capture `"tree"` as an `:id`
- [ ] T015 [US2] Add `getProjectTree` to `frontend/src/lib/api.js` returning the parsed `{nodes, total}` object
- [ ] T016 [US2] Create `frontend/src/components/ProjectTree.svelte`: accepts `let { projects = [] } = $props()` (flat list with `parent_id`); holds `let expanded = $state(new Set())`; uses `$derived` to group by `parent_id` into a tree; renders recursively with a chevron button (`aria-expanded`) per parent; first level is expanded by default via `$effect`
- [ ] T017 [US2] Replace the body of `frontend/src/components/ProjectList.svelte` so it calls `getProjectTree()` on mount and renders `<ProjectTree {projects} />` instead of the flat list. Keep the empty state and the existing "New Project" CTA

**Checkpoint**: US1 and US2 are both functional — a user can create sub-projects and see them nested in the tree. This is the **MVP**.

---

## Phase 5: User Story 3 - View Child Projects on a Project's Detail Page (Priority: P2)

**Goal**: On a project's detail page, a "Child Projects" section lists the project's direct children (with empty state when none). Each child name links to that child's detail page.

**Independent Test**: Open a project that has children — the section lists them with working links. Open a project with no children — the section shows the empty-state message. Open a non-existent project — the page renders a "project not found" state.

### Implementation for User Story 3

- [ ] T018 [P] [US3] Add `ProjectChildrenList` (`{ProjectID uuid.UUID; Count int; Children []Project}`) DTO to `backend/models/project.go`
- [ ] T019 [P] [US3] Add `ProjectRepository.ListChildren(ctx, parentID uuid.UUID) ([]Project, error)` in `backend/repositories/project_repository.go` running `SELECT … FROM projects WHERE parent_id = $1 ORDER BY title ASC` (uses `idx_projects_parent_id`)
- [ ] T020 [US3] Add `ProjectService.GetProjectChildren(ctx, id uuid.UUID) (*ProjectChildrenList, error)` in `backend/services/project_service.go`: call `repo.GetByID` first, return `ErrNotFound` (`models.ErrNotFound`) if the parent does not exist so the handler can return 404; otherwise call `repo.ListChildren` and wrap in `ProjectChildrenList`
- [ ] T021 [US3] Add `ProjectHandler.GetProjectChildren` in `backend/handlers/project_handler.go`: parse `:id` as UUID, call the service, return 400 on bad UUID, 404 on `ErrNotFound`, 500 on other errors, 200 with the DTO on success
- [ ] T022 [US3] Register `GET /api/projects/:id/children` in `backend/routes/routes.go` inside the existing projects group, after the existing `/:id` route (Fiber routes are matched in order, and `/:id/children` is unambiguous against `/:id`)
- [ ] T023 [P] [US3] Add `getProjectChildren(id)` to `frontend/src/lib/api.js` returning the parsed `ProjectChildrenList`
- [ ] T024 [P] [US3] Create `frontend/src/components/ProjectChildrenList.svelte`: accepts `let { projectId } = $props()`; on mount fetches `getProjectChildren(projectId)`; renders the populated list with each name as a link to that child's detail page; renders the empty-state message ("No sub-projects yet") when `count === 0`; renders the "project not found" state when the fetch returns 404
- [ ] T025 [US3] Mount `<ProjectChildrenList {projectId} />` in `frontend/src/App.svelte` inside the `{:else if selectedProject}` branch, between the project header (title + description) and the existing `<TaskList project={selectedProject} />` element. Pass `selectedProject.id` as the `projectId` prop

**Checkpoint**: US1, US2, US3 are all functional independently. The user can drill from the tree → parent detail → child detail.

---

## Phase 6: User Story 4 - Re-parent an Existing Project (Priority: P3)

**Goal**: A user can edit a project and change its parent (to a different project, or to none for top-level). The system rejects self/descendant parents with a clear error. The system also rejects deletion of a project that has children.

**Independent Test**: Edit an existing project, change its parent, save — verify the project now appears under the new parent in the tree and no longer under the old parent. Clear the parent — verify the project becomes top-level. Attempt to set a project as its own descendant (via direct API) — verify the 400 response. Attempt to delete a project with children — verify the 400 response.

### Implementation for User Story 4

- [ ] T026 [US4] Add `ProjectRepository.IsDescendantOf(ctx, ancestorID, candidateID uuid.UUID) (bool, error)` in `backend/repositories/project_repository.go` using a recursive CTE: `WITH RECURSIVE descendants AS (SELECT id, parent_id FROM projects WHERE id = $1 UNION ALL SELECT p.id, p.parent_id FROM projects p JOIN descendants d ON p.parent_id = d.id) SELECT EXISTS (SELECT 1 FROM descendants WHERE id = $2)`. Returns `true` if `candidateID == ancestorID` or any transitive descendant
- [ ] T027 [US4] Add `ProjectRepository.HasChildren(ctx, parentID uuid.UUID) (bool, int, error)` in `backend/repositories/project_repository.go` running `SELECT COUNT(*) FROM projects WHERE parent_id = $1` and returning `(count > 0, count, nil)`
- [ ] T028 [US4] In `ProjectService.UpdateProject` (`backend/services/project_service.go`), when `req.ParentID != nil`, call `repo.IsDescendantOf(ctx, id, *req.ParentID)`; if `true`, return `ErrValidation` with message `"cannot set a project as its own descendant"`. If `req.ParentID == nil`, no check is needed (clearing is always safe)
- [ ] T029 [US4] In `ProjectService.DeleteProject` (`backend/services/project_service.go`), call `repo.HasChildren(ctx, id)` first; if `hasChildren`, return `ErrValidation` with message `"cannot delete project with N sub-projects; re-parent or delete them first"` (substitute the actual `count`). Otherwise call `repo.Delete` as today
- [ ] T030 [US4] Update `ProjectHandler.DeleteProject` in `backend/handlers/project_handler.go` to return `400` with the service's error message when `DeleteProject` returns `ErrValidation` (the current code returns `500` on any error — replace with a `errors.Is(err, models.ErrValidation)` check). Also return `404` (not `500`) when the service reports the project does not exist
- [ ] T031 [US4] Update `updateProject` in `frontend/src/lib/api.js` to forward `parent_id` from the caller's payload
- [ ] T032 [US4] Extend the parent picker added in T009 (in `frontend/src/components/ProjectForm.svelte`) to support edit mode: when editing an existing project, fetch the current project tree, exclude the project itself and all of its descendants from the picker options, and submit the chosen value (or `null` to clear) as `parent_id` in the update payload

**Checkpoint**: All four user stories are functional independently. US4 is the only one that depends on US1's form (it shares the parent picker) — they cannot run truly in parallel in the implementation, but US4 is testable end-to-end on its own once T032 is done.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories.

- [ ] T033 [P] Run the end-to-end walkthrough in `specs/001-project-hierarchy/quickstart.md` §5 (US1, US2, US3, US4) and §6 (edge cases) and confirm every step matches the expected outcome
- [ ] T034 [P] Verify the success criteria in `specs/001-project-hierarchy/quickstart.md` §7 (SC-001 through SC-007) and tick them off in the quickstart
- [ ] T035 [P] Update `AGENTS.md` to record the new Svelte components (`ProjectTree.svelte`, `ProjectChildrenList.svelte`) and bump the component count line in the OVERVIEW section
- [ ] T036 [P] Add a short entry to the `## Recent Changes` section of `AGENTS.md` (repository root) describing the project-hierarchy feature

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately.
- **Foundational (Phase 2)**: Depends on Setup completion — **BLOCKS all user stories**. T002 is the only one that requires a database round-trip; T003–T006 are code-only and can be done in any order, but the recommended execution order is the one listed above.
- **User Stories (Phase 3+)**: All depend on the Foundational phase completion.
  - US1 (P1) → US2 (P1) → US3 (P2) → US4 (P3) is the recommended sequential order. US1 and US2 form the MVP and are both P1.
  - US4 technically depends on US1's parent picker in the form (T009 / T032 share `ProjectForm.svelte`), so US4 cannot fully start until US1's form work is in place.
  - US2 and US3 do not depend on each other and can be developed in parallel.
- **Polish (Phase 7)**: Depends on all desired user stories being complete.

### User Story Dependencies

- **US1 (P1)**: Can start after Foundational (Phase 2). No dependencies on other stories.
- **US2 (P1)**: Can start after Foundational (Phase 2). No dependencies on other stories. **Shares the MVP with US1.**
- **US3 (P2)**: Can start after Foundational (Phase 2). No dependencies on other stories.
- **US4 (P3)**: Can start after Foundational (Phase 2). Shares the parent-picker component with US1 (same file: `ProjectForm.svelte`), so concurrent edits to that file are not safe.

### Within Each User Story

- No tests (manual verification only — see `quickstart.md`).
- Models before services.
- Services before handlers.
- Handlers before routes.
- Backend before frontend.

### Parallel Opportunities

- **Phase 2**: T003 (model) and T004 (repo Get* reads) touch different files but T004 depends on the new `ParentID` field in T003 — they cannot be parallel. T005 (Create) and T006 (Update) touch the same file and depend on T004; they can be done in any order but not in parallel with each other.
- **Phase 4 (US2)**: T010 (model) and T011 (repo) are in different files with no dependency between them — true parallel work. T015 (api.js) and T016 (ProjectTree.svelte) are different files with no dependency — true parallel work. The middle steps (T012–T014) are sequential because each layer depends on the previous.
- **Phase 5 (US3)**: T018 (model) and T019 (repo) are parallel. T023 (api.js) and T024 (ProjectChildrenList.svelte) are parallel.
- **Phase 7 (Polish)**: All four tasks are independent — full parallelism.

---

## Parallel Examples

### Phase 4 (US2) — two parallel waves

Wave 1 (after Foundational):
```bash
# In parallel (different files, no dependency):
Task T010: "Add ProjectTree and ProjectTreeNode DTOs to backend/models/project.go"
Task T011: "Add ProjectRepository.ListTree in backend/repositories/project_repository.go"
```

Wave 2 (after T010, T011, T012):
```bash
# In parallel (different files, no dependency):
Task T015: "Add getProjectTree to frontend/src/lib/api.js"
Task T016: "Create ProjectTree.svelte in frontend/src/components/ProjectTree.svelte"
```

### Phase 5 (US3) — two parallel waves

Wave 1:
```bash
# In parallel:
Task T018: "Add ProjectChildrenList DTO to backend/models/project.go"
Task T019: "Add ProjectRepository.ListChildren in backend/repositories/project_repository.go"
```

Wave 2 (after T020–T022):
```bash
# In parallel:
Task T023: "Add getProjectChildren to frontend/src/lib/api.js"
Task T024: "Create ProjectChildrenList.svelte in frontend/src/components/ProjectChildrenList.svelte"
```

### Phase 7 (Polish) — all in parallel

```bash
Task T033: "Run quickstart.md end-to-end walkthrough"
Task T034: "Verify success criteria SC-001..SC-007"
Task T035: "Update AGENTS.md component count"
Task T036: "Add Recent Changes entry to AGENTS.md"
```

---

## Implementation Strategy

### MVP First (US1 + US2 only)

1. Complete Phase 1: Setup (T001).
2. Complete Phase 2: Foundational (T002–T006).
3. Complete Phase 3: US1 (T007–T009).
4. Complete Phase 4: US2 (T010–T017).
5. **STOP and VALIDATE**: Walk US1 and US2 in `quickstart.md` §5. A user can create sub-projects and see them nested in a tree. The "appears under the parent" verification of US1's acceptance scenario is satisfied by US2's tree view.
6. Deploy / demo the MVP.

### Incremental Delivery

1. Setup + Foundational → foundation ready.
2. Add US1 → independent verification (data only) → ready.
3. Add US2 → tree view live → **MVP**, deploy / demo.
4. Add US3 → child list on detail page → deploy / demo.
5. Add US4 → re-parent + delete blocking → deploy / demo.
6. Polish → AGENTS.md updated, success criteria ticked off.

### Parallel Team Strategy

With multiple developers:

1. Together: Setup + Foundational.
2. After Foundational:
   - Developer A: US2 (tree view).
   - Developer B: US1 (create with parent) — must finish T009 before Developer C starts US4 (shared file).
   - Developer C: US3 (child list) — fully independent of A and B.
3. After US1's T009 lands: Developer C (or a new developer) starts US4.
4. After all four stories: Developer D runs the Polish phase.

---

## Notes

- All file paths are relative to the repository root.
- No new external dependency is introduced by any task — all work uses the libraries already in `backend/go.mod` and `frontend/package.json`.
- US4 is the only user story that touches the DELETE endpoint's behaviour. The current handler returns `500` on any error; T030 normalises this to `400` (validation) / `404` (not found). This is a small, scope-aligned fix and is required for FR-007 to surface a useful error.
- T014 is order-sensitive: `GET /api/projects/tree` must be registered before `GET /api/projects/:id` in `routes.go`, otherwise Fiber will capture the literal string `"tree"` as the `:id` parameter and the tree endpoint becomes unreachable. The task description calls this out explicitly.
- T022 is unambiguous against T014's registration because `:id/children` has a fixed suffix that does not collide with any other route.
- T032 modifies the same file (and the same picker component) added by T009. Concurrent work on T009 and T032 is **not** safe; they must be serialised. The dependency is captured in the "User Story Dependencies" section above.
- Stop at any checkpoint to validate the story independently against `spec.md` and `quickstart.md`.
