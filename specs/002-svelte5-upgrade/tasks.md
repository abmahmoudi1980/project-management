# Tasks: Upgrade Frontend to Svelte 5

**Input**: Design documents from `/specs/002-svelte5-upgrade/`
**Prerequisites**: plan.md, spec.md, research.md, quickstart.md

**Tests**: No automated tests in this project - manual testing workflow will be used.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Web app**: `frontend/src/`, `frontend/src/components/`, `frontend/src/stores/`

---

## Phase 1: Setup (Dependency Upgrades)

**Purpose**: Upgrade npm packages and verify basic build compatibility

- [ ] T001 Backup current package-lock.json to frontend/package-lock.json.backup
- [ ] T002 Update svelte dependency to ^5.0.0 in frontend/package.json
- [ ] T003 Update @sveltejs/vite-plugin-svelte dependency to ^5.0.0 in frontend/package.json
- [ ] T004 Run npm install in frontend/ to update dependencies
- [ ] T005 Verify build succeeds with npm run build in frontend/
- [ ] T006 Verify dev server starts with npm run dev in frontend/

---

## Phase 2: Foundational (N/A for this feature)

**Note**: This feature has no foundational blocking tasks. User Story 1 can begin immediately after Phase 1.

---

## Phase 3: User Story 1 - Core Application Functionality Preserved (Priority: P1) 🎯 MVP

**Goal**: Migrate all 8 Svelte components to Svelte 5 runes syntax while maintaining 100% functionality

**Independent Test**: Run complete application workflow (create project → add tasks → log time) and verify identical behavior to Svelte 4 version. Check browser console for zero errors.

### Component Migration - Core App

- [ ] T007 [US1] Migrate App.svelte to Svelte 5 runes in frontend/src/App.svelte (convert selectedProject to $state(null), keep onMount)
- [ ] T008 [US1] Test App.svelte: verify project list displays, selection works, task list appears on selection

### Component Migration - Modal

- [ ] T009 [P] [US1] Migrate Modal.svelte to Svelte 5 runes in frontend/src/components/Modal.svelte (convert props to $props(), any internal state to $state)
- [ ] T010 [US1] Test Modal.svelte: verify modal opens/closes, backdrop click works, ESC key works

### Component Migration - Forms

- [ ] T011 [P] [US1] Migrate ProjectForm.svelte to Svelte 5 runes in frontend/src/components/ProjectForm.svelte (convert props to $props(), form state to $state)
- [ ] T012 [US1] Test ProjectForm: verify form opens, fields accept input, submission creates project, form closes

- [ ] T013 [P] [US1] Migrate TaskForm.svelte to Svelte 5 runes in frontend/src/components/TaskForm.svelte (convert props to $props(), form state to $state)
- [ ] T014 [US1] Test TaskForm: verify form opens, all fields work including date picker, submission creates task

- [ ] T015 [P] [US1] Migrate TimeLogForm.svelte to Svelte 5 runes in frontend/src/components/TimeLogForm.svelte (convert props to $props(), form state to $state)
- [ ] T016 [US1] Test TimeLogForm: verify hours field accepts numbers, description works, submission creates time log

### Component Migration - Lists

- [ ] T017 [P] [US1] Migrate ProjectList.svelte to Svelte 5 runes in frontend/src/components/ProjectList.svelte (convert props to $props(), local state to $state)
- [ ] T018 [US1] Test ProjectList: verify list displays, click selection works, events dispatch correctly

- [ ] T019 [P] [US1] Migrate TaskList.svelte to Svelte 5 runes in frontend/src/components/TaskList.svelte (convert props to $props(), local state to $state)
- [ ] T020 [US1] Test TaskList: verify tasks display, task details show, time logs display, edit/delete work

### Component Migration - Date Picker

- [ ] T021 [US1] Migrate JalaliDatePicker.svelte to Svelte 5 runes in frontend/src/components/JalaliDatePicker.svelte (convert props to $props(), calendar state to $state, derived values to $derived)
- [ ] T022 [US1] Test JalaliDatePicker thoroughly: calendar opens, dates selectable, updates form, Jalali conversion works, calendar closes

### Store Verification

- [ ] T023 [P] [US1] Verify projectStore.js remains compatible in frontend/src/stores/projectStore.js (no changes needed - verify $projects syntax works in components)
- [ ] T024 [P] [US1] Verify taskStore.js remains compatible in frontend/src/stores/taskStore.js (no changes needed - verify $tasks syntax works in components)
- [ ] T025 [P] [US1] Verify timeLogStore.js remains compatible in frontend/src/stores/timeLogStore.js (no changes needed - verify $timeLogs syntax works in components)

### Integration Testing

- [ ] T026 [US1] Run full user workflow test: Create project → verify appears in list
- [ ] T027 [US1] Run full user workflow test: Select project → add task → verify task appears
- [ ] T028 [US1] Run full user workflow test: Select task → log time → verify time log appears
- [ ] T029 [US1] Run full user workflow test: Edit project/task/time log → verify updates persist
- [ ] T030 [US1] Run full user workflow test: Delete time log → task → project → verify deletions work
- [ ] T031 [US1] Verify browser console shows zero errors and zero warnings during all workflows

**Checkpoint**: User Story 1 complete - all functionality works identically to Svelte 4 version

---

## Phase 4: User Story 2 - Developer Experience Improvements (Priority: P2)

**Goal**: Verify all components use Svelte 5 runes properly and build process is clean

**Independent Test**: Code review confirms all components use runes syntax, build completes with zero warnings

### Code Quality Verification

- [ ] T032 [P] [US2] Code review App.svelte: verify uses $state for reactive variables, $derived for computed values
- [ ] T033 [P] [US2] Code review Modal.svelte: verify uses $props() for component props
- [ ] T034 [P] [US2] Code review ProjectForm.svelte: verify uses $state for form state, $props() for props
- [ ] T035 [P] [US2] Code review TaskForm.svelte: verify uses $state for form state, $props() for props
- [ ] T036 [P] [US2] Code review TimeLogForm.svelte: verify uses $state for form state, $props() for props
- [ ] T037 [P] [US2] Code review ProjectList.svelte: verify uses $props() and $state appropriately
- [ ] T038 [P] [US2] Code review TaskList.svelte: verify uses $props() and $state appropriately
- [ ] T039 [P] [US2] Code review JalaliDatePicker.svelte: verify uses $state, $derived, and $props appropriately

### Build Process Verification

- [ ] T040 [US2] Run npm run build in frontend/ and verify zero errors
- [ ] T041 [US2] Run npm run build in frontend/ and verify zero warnings
- [ ] T042 [US2] Verify build output has no deprecated syntax warnings
- [ ] T043 [US2] Run npm run dev in frontend/ and verify dev server starts without warnings

**Checkpoint**: Developer experience improved - all components follow Svelte 5 best practices

---

## Phase 5: User Story 3 - Build Performance Optimization (Priority: P3)

**Goal**: Measure and verify build performance improvements from Svelte 5

**Independent Test**: Compare build times and bundle sizes before/after upgrade

### Performance Measurement

- [ ] T044 [P] [US3] Measure production build time with time npm run build in frontend/
- [ ] T045 [P] [US3] Measure bundle size with ls -lh frontend/dist/ after build
- [ ] T046 [P] [US3] Document bundle size in specs/002-svelte5-upgrade/performance-metrics.md
- [ ] T047 [US3] Test HMR speed: make change in component, measure time to reflect in browser
- [ ] T048 [US3] Compare metrics to baseline (Svelte 4 metrics if available)

### Build Configuration Optimization

- [ ] T049 [US3] Review vite.config.js in frontend/vite.config.js for any Svelte 5-specific optimizations
- [ ] T050 [US3] Verify Svelte plugin options are optimal for production in frontend/vite.config.js
- [ ] T051 [US3] Run npm run preview in frontend/ to test production build locally

**Checkpoint**: Performance metrics documented, build configuration optimized

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final cleanup and documentation

- [ ] T052 [P] Update AGENTS.md with final Svelte 5 migration status
- [ ] T053 [P] Update specs/002-svelte5-upgrade/quickstart.md with actual migration experience notes
- [ ] T054 Run final verification: npm run build && npm run dev in frontend/
- [ ] T055 Git commit all changes with descriptive message
- [ ] T056 Create pull request or merge to main branch per team workflow

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: N/A - skipped for this feature
- **User Story 1 (Phase 3)**: Depends on Setup (Phase 1) completion
- **User Story 2 (Phase 4)**: Depends on User Story 1 completion (needs migrated code to review)
- **User Story 3 (Phase 5)**: Depends on User Story 1 completion (needs working build to measure)
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Phase 1 - No dependencies on other stories
- **User Story 2 (P2)**: Depends on User Story 1 (needs code to review)
- **User Story 3 (P3)**: Depends on User Story 1 (needs working build to measure)

### Within User Story 1 (Sequential & Parallel Opportunities)

**Sequential Order**:
1. Core App migration (T007-T008) - MUST be first (other components depend on it)
2. Modal migration (T009-T010) - Should be early (used by forms)
3. Forms & Lists can proceed in parallel after Modal
4. Date Picker (T021-T022) - Can be parallel but test thoroughly
5. Store verification (T023-T025) - Can be parallel once components migrated
6. Integration testing (T026-T031) - MUST be after all migrations complete

**Parallel Opportunities within US1**:
- After T008 completes:
  - T009 (Modal) can run
  - T011 (ProjectForm) can run
  - T013 (TaskForm) can run
  - T015 (TimeLogForm) can run
  - T017 (ProjectList) can run
  - T019 (TaskList) can run
- T023, T024, T025 (Store verifications) can all run in parallel
- T032-T039 (Code reviews in US2) can all run in parallel

---

## Parallel Example: User Story 1 Component Migrations

```bash
# After App.svelte is migrated and tested (T007-T008), launch these in parallel:

# Developer 1 or Agent 1:
Task: "Migrate Modal.svelte to Svelte 5 runes in frontend/src/components/Modal.svelte"

# Developer 2 or Agent 2:
Task: "Migrate ProjectForm.svelte to Svelte 5 runes in frontend/src/components/ProjectForm.svelte"
Task: "Migrate TaskForm.svelte to Svelte 5 runes in frontend/src/components/TaskForm.svelte"

# Developer 3 or Agent 3:
Task: "Migrate ProjectList.svelte to Svelte 5 runes in frontend/src/components/ProjectList.svelte"
Task: "Migrate TaskList.svelte to Svelte 5 runes in frontend/src/components/TaskList.svelte"

# After forms are migrated, launch date picker:
Task: "Migrate JalaliDatePicker.svelte to Svelte 5 runes in frontend/src/components/JalaliDatePicker.svelte"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (dependency upgrades) - ~10 minutes
2. Complete Phase 3: User Story 1 (all component migrations) - ~2-3 hours
3. **STOP and VALIDATE**: Test all workflows, verify zero regressions
4. Deploy/demo if ready

**This gives you**: Fully functional Svelte 5 application with all features working

### Incremental Delivery

1. Complete Setup (Phase 1) → Dependencies upgraded
2. Complete User Story 1 (Phase 3) → Full functionality in Svelte 5 ✓ (MVP!)
3. Complete User Story 2 (Phase 4) → Code quality verified ✓
4. Complete User Story 3 (Phase 5) → Performance measured ✓
5. Polish (Phase 6) → Documentation complete ✓

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup (Phase 1) together - 10 minutes
2. Developers split User Story 1 component migrations:
   - Dev A: App + Modal + ProjectForm (T007-T012)
   - Dev B: TaskForm + TimeLogForm (T013-T016)
   - Dev C: ProjectList + TaskList (T017-T020)
   - Dev D: JalaliDatePicker + Store verification (T021-T025)
3. Team reconvenes for integration testing (T026-T031)
4. User Stories 2 and 3 can proceed quickly (mostly verification)

**Timeline with parallel work**: 1-2 hours instead of 3-4 hours sequential

---

## Task Summary

**Total Tasks**: 56 tasks

### Tasks per Phase:
- **Phase 1 (Setup)**: 6 tasks
- **Phase 2 (Foundational)**: 0 tasks (N/A)
- **Phase 3 (User Story 1)**: 25 tasks
- **Phase 4 (User Story 2)**: 12 tasks
- **Phase 5 (User Story 3)**: 8 tasks
- **Phase 6 (Polish)**: 5 tasks

### Tasks per User Story:
- **User Story 1**: 25 tasks (component migrations + testing)
- **User Story 2**: 12 tasks (code quality verification)
- **User Story 3**: 8 tasks (performance measurement)

### Parallelizable Tasks: 24 tasks marked [P]

### MVP Scope (User Story 1 only): 31 tasks (Phase 1 + Phase 3)

---

## Notes

- [P] tasks = different files, no dependencies
- [US1], [US2], [US3] labels map tasks to specific user stories
- Each user story is independently completable and testable
- No automated tests - manual testing workflow used
- Stop after each component migration to test individually
- Commit frequently (after each component or logical group)
- All components must work identically to Svelte 4 version
- Stores require zero changes (backward compatible)
