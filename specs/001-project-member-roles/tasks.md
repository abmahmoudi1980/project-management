# Tasks: Project Member Role Assignment

**Input**: Design documents from `/specs/001-project-member-roles/`  
**Branch**: `001-project-member-roles`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/openapi.yaml, quickstart.md

**Tests**: Automated test tasks are not included because the specification does not explicitly request TDD or new automated test coverage.

**Organization**: Tasks are grouped by user story so each story can be implemented and validated independently.

## Phase 1: Setup (Project Initialization)

**Purpose**: Prepare implementation scaffolding and file structure for this feature.

- [x] T001 Create migration file scaffold in `backend/migration/008_add_project_members.sql`
- [x] T002 [P] Create model file scaffold in `backend/models/project_member.go`
- [x] T003 [P] Create repository file scaffold in `backend/repositories/project_member_repository.go`
- [x] T004 [P] Create service file scaffold in `backend/services/project_member_service.go`
- [x] T005 [P] Create handler file scaffold in `backend/handlers/project_member_handler.go`
- [x] T006 [P] Create frontend API client file scaffold in `frontend/src/lib/api/projectMembers.js`
- [x] T007 [P] Create frontend member panel file scaffold in `frontend/src/components/ProjectMembersPanel.svelte`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Deliver core schema and shared backend wiring required by all user stories.

**⚠️ CRITICAL**: No user story work can start until this phase is complete.

- [ ] T008 Implement `project_roles` and `project_members` tables with constraints and indexes in `backend/migration/008_add_project_members.sql`
- [ ] T009 Seed predefined active roles (`owner`, `manager`, `contributor`, `viewer`) in `backend/migration/008_add_project_members.sql`
- [ ] T010 Run and verify migration using `backend/cmd/migrate/main.go`
- [ ] T011 Define shared DTO/entity structs (role/member request-response types) in `backend/models/project_member.go`
- [ ] T012 Add shared constructor wiring for project member repository/service/handler in `backend/main.go`

**Checkpoint**: Database schema and shared wiring are ready for story implementation.

---

## Phase 3: User Story 1 - Add Member With Role (Priority: P1) 🎯 MVP

**Goal**: Admin can add a project member with a predefined role and view project members.

**Independent Test**: Call add-member API for a project with valid `user_id` and `project_role_id`, then list members and verify assigned role and duplicate prevention.

### Implementation for User Story 1

- [ ] T013 [US1] Implement repository method to insert project membership in `backend/repositories/project_member_repository.go`
- [ ] T014 [US1] Implement repository method to list project members with role join in `backend/repositories/project_member_repository.go`
- [ ] T015 [US1] Implement service method for add-member workflow (project exists, duplicate check, persistence) in `backend/services/project_member_service.go`
- [ ] T016 [US1] Implement service method to fetch project member list in `backend/services/project_member_service.go`
- [ ] T017 [US1] Implement handler for `POST /api/projects/:projectId/members` in `backend/handlers/project_member_handler.go`
- [ ] T018 [US1] Implement handler for `GET /api/projects/:projectId/members` in `backend/handlers/project_member_handler.go`
- [ ] T019 [US1] Register member list/add routes in `backend/routes/routes.go`
- [ ] T020 [US1] Add frontend API functions `addProjectMember` and `getProjectMembers` in `frontend/src/lib/api/projectMembers.js`
- [ ] T021 [US1] Implement basic members panel (list + add submit flow) in `frontend/src/components/ProjectMembersPanel.svelte`
- [ ] T022 [US1] Mount member panel in project view in `frontend/src/App.svelte`
- [ ] T023 [US1] Verify MVP scenario and duplicate-member rejection via quickstart flow in `specs/001-project-member-roles/quickstart.md`

**Checkpoint**: US1 is independently functional and demonstrable as MVP.

---

## Phase 4: User Story 2 - Select From System Users (Priority: P2)

**Goal**: Admin selects members from eligible active system users for a specific project.

**Independent Test**: Open add-member flow and confirm eligible users are listed from system users, excluding existing project members and inactive users.

### Implementation for User Story 2

- [ ] T024 [US2] Implement repository query for eligible users by project in `backend/repositories/project_member_repository.go`
- [ ] T025 [US2] Implement service method for eligible-user search/filter in `backend/services/project_member_service.go`
- [ ] T026 [US2] Implement handler for `GET /api/projects/:projectId/members/eligible-users` in `backend/handlers/project_member_handler.go`
- [ ] T027 [US2] Register eligible-users route in `backend/routes/routes.go`
- [ ] T028 [US2] Add `getEligibleProjectUsers` API function in `frontend/src/lib/api/projectMembers.js`
- [ ] T029 [US2] Replace raw user-id input with eligible-user selector in `frontend/src/components/ProjectMembersPanel.svelte`
- [ ] T030 [US2] Validate eligible-user filtering behavior (active only, non-members only) in `specs/001-project-member-roles/quickstart.md`

**Checkpoint**: US2 is independently functional and testable on top of MVP.

---

## Phase 5: User Story 3 - Use Predefined Roles Consistently (Priority: P3)

**Goal**: Only predefined active project roles are assignable and consistently displayed.

**Independent Test**: Role dropdown shows predefined roles only, member add with invalid role is rejected, and saved members display role labels.

### Implementation for User Story 3

- [ ] T031 [US3] Implement repository query for active predefined roles in `backend/repositories/project_member_repository.go`
- [ ] T032 [US3] Enforce active predefined role validation in add-member service flow in `backend/services/project_member_service.go`
- [ ] T033 [US3] Implement handler for `GET /api/project-roles` in `backend/handlers/project_member_handler.go`
- [ ] T034 [US3] Register project-roles route in `backend/routes/routes.go`
- [ ] T035 [US3] Add `getProjectRoles` API function in `frontend/src/lib/api/projectMembers.js`
- [ ] T036 [US3] Replace raw role-id input with predefined-role selector in `frontend/src/components/ProjectMembersPanel.svelte`
- [ ] T037 [US3] Render role display names in project member list UI in `frontend/src/components/ProjectMembersPanel.svelte`
- [ ] T038 [US3] Validate invalid-role rejection and predefined-role-only selection in `specs/001-project-member-roles/quickstart.md`

**Checkpoint**: US3 is independently functional and all story requirements are satisfied.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final hardening and consistency checks across all stories.

- [ ] T039 [P] Align API contract examples with implemented payloads in `specs/001-project-member-roles/contracts/openapi.yaml`
- [ ] T040 [P] Align data model notes with final schema/query decisions in `specs/001-project-member-roles/data-model.md`
- [ ] T041 Run backend build verification in `backend/main.go`
- [ ] T042 Run frontend build verification in `frontend/package.json`
- [ ] T043 Perform end-to-end quickstart validation and record final notes in `specs/001-project-member-roles/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1 (Setup)**: no dependencies.
- **Phase 2 (Foundational)**: depends on Phase 1; blocks all story work.
- **Phase 3 (US1)**: depends on Phase 2 completion.
- **Phase 4 (US2)**: depends on Phase 2 completion and integrates with US1 member panel/API.
- **Phase 5 (US3)**: depends on Phase 2 completion and integrates with US1 add flow.
- **Phase 6 (Polish)**: depends on all targeted stories being complete.

### User Story Dependencies

- **US1 (P1)**: first deliverable (MVP), no dependency on other stories.
- **US2 (P2)**: functionally independent but builds on US1 UI/API touchpoints.
- **US3 (P3)**: functionally independent but builds on US1 add-member workflow.

### Recommended Execution Order

1. Complete Setup + Foundational.
2. Deliver US1 as MVP and validate independently.
3. Deliver US2 (eligible selection) and validate independently.
4. Deliver US3 (role consistency) and validate independently.
5. Run polish tasks.

---

## Parallel Execution Examples

### User Story 1

- Run T013 and T014 in parallel (same repository file but separate methods after initial scaffold).
- Run T020 and T021 in parallel (API client and component implementation).

### User Story 2

- Run T024 and T028 in parallel (backend eligible-user query + frontend API wrapper).

### User Story 3

- Run T031 and T035 in parallel (backend roles query + frontend roles API wrapper).

---

## Implementation Strategy

### MVP First

- Deliver Phases 1–3 only (through T023) to provide immediate value: admin can add members with roles and see project membership.

### Incremental Delivery

- Add US2 next to improve usability with system-user selection.
- Add US3 last to harden role governance and consistency.

### Team Parallelization

- After foundational completion, backend and frontend tasks within each story can be split across developers.
- Use `[P]` tasks first to maximize throughput with minimal file conflicts.
