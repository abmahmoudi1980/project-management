# Tasks: Task Search Feature

**Input**: Design documents from `/specs/004-task-search/`  
**Prerequisites**: ✅ plan.md, ✅ spec.md, ✅ research.md, ✅ data-model.md, ✅ contracts/API.md  
**Branch**: `004-task-search`  
**Date**: 2026-01-01

## Format: `[ID] [P?] [Story?] Description with file path`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: User story (US1, US2, US3)
- Exact file paths included in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Frontend environment and base component setup

- [ ] T001 Verify Svelte 5 development environment (`npm run dev` works in `frontend/`)
- [ ] T002 Review existing JalaliDatePicker.svelte component structure in `frontend/src/components/JalaliDatePicker.svelte`
- [ ] T003 [P] Understand TaskList.svelte component structure in `frontend/src/components/TaskList.svelte`
- [ ] T004 [P] Review taskStore.js state management in `frontend/src/stores/taskStore.js`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core filtering infrastructure that user stories depend on

**⚠️ CRITICAL**: User story implementation cannot begin until this phase is complete

- [ ] T005 Create utility functions for Jalali date conversion in `frontend/src/lib/dateUtils.js`:
  - `jalaliStringToDate(jalaliStr)`: Convert "1403/10/15" → JavaScript Date (Gregorian)
  - `dateToJalaliString(date)`: Convert Date → "1403/10/15" format
  - Use existing `jalali-moment` library
  - Add proper error handling for invalid inputs

- [ ] T006 Create filter utility functions in `frontend/src/lib/filterUtils.js`:
  - `evaluateTextFilter(task, searchTerm)`: Check if title or description contains search term (case-insensitive)
  - `evaluateDateFilter(task, dateField, fromDate, toDate)`: Check if task date within range
  - `evaluateAllFilters(task, filters)`: Apply all filters with AND logic
  - Include JSDoc comments for clarity

- [ ] T007 Create TaskSearch.svelte component skeleton in `frontend/src/components/TaskSearch.svelte`:
  - Define filter state using Svelte 5 runes: `text`, `start_date_from`, `start_date_to`, `due_date_from`, `due_date_to`
  - Create `clearFilters()` function
  - Export filters via component props or event dispatcher
  - Do NOT include styling yet (deferred to Phase 3)

**Checkpoint**: Filtering infrastructure ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Filter Tasks by Text Search (Priority: P1) 🎯 MVP

**Goal**: Enable users to search for tasks by entering keywords that match task title or description

**Independent Test**: Enter "API" in search field → verify only tasks with "API" in title/description shown; change to "xyz" → verify "No tasks found" message

### Implementation for User Story 1

- [ ] T008 [P] [US1] Complete TaskSearch component text input field in `frontend/src/components/TaskSearch.svelte`:
  - Add text input HTML element with placeholder "جستجو در وظایف..."
  - Bind input value to `filters.text`
  - Use Tailwind CSS classes for styling (deferred to T015)

- [ ] T009 [P] [US1] Update TaskList.svelte to add text filter state in `frontend/src/components/TaskList.svelte`:
  - Add `let filters = $state({ text: '', ... })` reactive variable
  - Import `evaluateTextFilter` from `frontend/src/lib/filterUtils.js`

- [ ] T010 [US1] Create derived `filteredTasks` in TaskList.svelte using `$derived.by()`:
  - Filter `$tasks.tasks` array using `evaluateAllFilters(task, filters)`
  - Real-time updates as filters change (leverages Svelte 5 reactivity)

- [ ] T011 [US1] Update task iteration loop in `frontend/src/components/TaskList.svelte`:
  - Change `{#each $tasks.tasks || [] as task}` to `{#each filteredTasks as task}`
  - Maintains all existing task rendering logic

- [ ] T012 [US1] Add "No tasks found" message in TaskList.svelte:
  - Display when `filteredTasks.length === 0 && $tasks.tasks.length > 0`
  - Show message: "هیچ وظیفه‌ای یافت نشد"
  - Position above task list

- [ ] T013 [US1] Update result count display in TaskList.svelte:
  - Change from `{$tasks.total}` to show filtered results: `{filteredTasks.length} / {$tasks.total}`
  - Display in task list header

- [ ] T014 [US1] Integrate TaskSearch component in TaskList.svelte:
  - Import TaskSearch component: `import TaskSearch from './TaskSearch.svelte'`
  - Bind filters: `<TaskSearch bind:filters />`
  - Position above current task list (before task iteration loop)

- [ ] T015 [P] [US1] Style TaskSearch text input field in `frontend/src/components/TaskSearch.svelte`:
  - Use Tailwind CSS classes matching project style (see TaskList.svelte for reference)
  - Mobile responsive: full width on small screens
  - RTL support for Persian text
  - Focus states and hover effects

**Checkpoint**: User Story 1 complete - text search fully functional and independently testable

---

## Phase 4: User Story 2 - Filter Tasks by Date Range (Priority: P1)

**Goal**: Enable users to filter tasks by start and due date ranges using Jalali date picker

**Independent Test**: Select start date range 1403/01/01 to 1403/10/31 → verify only tasks with start_date in that range shown; clear dates → verify all tasks shown again

### Implementation for User Story 2

- [ ] T016 [P] [US2] Create date input fields in TaskSearch.svelte in `frontend/src/components/TaskSearch.svelte`:
  - Add HTML labels: "تاریخ شروع (از)" (Start Date From), "تاریخ شروع (تا)" (Start Date To)
  - Add HTML labels: "مهلت (از)" (Due Date From), "مهلت (تا)" (Due Date To)
  - Create Svelte state for displaying dates: `start_date_from_display`, `start_date_to_display`, etc.
  - Do NOT style yet (deferred to T018)

- [ ] T017 [P] [US2] Integrate JalaliDatePicker in TaskSearch.svelte:
  - Import JalaliDatePicker: `import JalaliDatePicker from './JalaliDatePicker.svelte'`
  - Add 4 instances (start_from, start_to, due_from, due_to)
  - Bind display values: `bind:value={start_date_from_display}`
  - Handle change events to convert Jalali → Gregorian using `jalaliStringToDate()`

- [ ] T018 [P] [US2] Implement date change handlers in TaskSearch.svelte:
  - `handleStartDateFromChange(jalaliStr)`: Parse and store as `filters.start_date_from`
  - `handleStartDateToChange(jalaliStr)`: Parse and store as `filters.start_date_to`
  - `handleDueDateFromChange(jalaliStr)`: Parse and store as `filters.due_date_from`
  - `handleDueDateToChange(jalaliStr)`: Parse and store as `filters.due_date_to`
  - Handle invalid dates (clear field, prevent filtering)

- [ ] T019 [US2] Verify date filtering logic in filterUtils.js `evaluateDateFilter()`:
  - Test filtering with various date ranges
  - Handle null dates (no filter)
  - Handle tasks without dates (excluded from range filtering)

- [ ] T020 [P] [US2] Style date input section in TaskSearch.svelte:
  - Use Tailwind CSS for consistent layout
  - Mobile: Stack 4 date fields vertically
  - Tablet: 2x2 grid
  - Desktop: 4 fields in row
  - RTL support for Persian labels

- [ ] T021 [US2] Test date range filtering:
  - Manual test: Apply start date range → verify filtering works
  - Manual test: Apply due date range → verify filtering works
  - Manual test: Clear date fields → verify all tasks shown

**Checkpoint**: User Story 2 complete - date range filtering fully functional and independently testable

---

## Phase 5: User Story 3 - Combined Search Filters (Priority: P2)

**Goal**: Enable users to combine text search and date filters for precise task discovery

**Independent Test**: Enter text "API" + select date ranges → verify only tasks matching ALL criteria shown; change one filter → verify results update immediately

### Implementation for User Story 3

- [ ] T022 [US3] Verify combined filter logic in filterUtils.js:
  - `evaluateAllFilters(task, filters)` applies text AND date filters
  - Confirm AND logic (all must match): textMatch && startMatch && dueMatch
  - Verify with acceptance scenario: text="API", start_date_from=date, due_date_to=date → only matching tasks shown

- [ ] T023 [US3] Test combined filter scenarios in TaskList.svelte:
  - Manual test: Text + start date range active → verify filtering
  - Manual test: Text + due date range active → verify filtering
  - Manual test: Text + both date ranges active → verify filtering
  - Manual test: Modify one filter while others active → verify results update immediately (real-time)

- [ ] T024 [US3] Verify "No tasks found" with combined filters:
  - Apply filters that match no tasks
  - Confirm "No tasks found" message displays
  - Test clearing individual filters one-by-one → results should update

- [ ] T025 [P] [US3] Create filter summary display in TaskSearch.svelte:
  - Show active filter count if > 0
  - Example: "فیلترهای فعال: 3" (Active filters: 3)
  - Help users understand what's filtered

**Checkpoint**: User Story 3 complete - combined filtering fully functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Refinement and feature completion across all user stories

- [ ] T026 Implement comprehensive "Clear Filters" button in TaskSearch.svelte:
  - Add button that resets all filter state to initial values
  - Styling consistent with project (see TaskList.svelte for reference)
  - Only show button if `hasActiveFilters()` returns true
  - Position at end of filter controls

- [ ] T027 [P] Implement individual filter clear functionality:
  - Allow clearing text filter (click on text input, select all, delete)
  - Allow clearing start date range (clear one or both fields)
  - Allow clearing due date range (clear one or both fields)
  - Each clear updates results immediately

- [ ] T028 Verify accessibility and keyboard navigation:
  - All inputs have labels (for screen readers)
  - Tab order is logical (left-to-right, top-to-bottom)
  - Focus states visible
  - RTL layout respected

- [ ] T029 Test mobile responsiveness:
  - Test on mobile viewport (< 768px)
  - Search panel should stack vertically
  - Input fields should be full width
  - JalaliDatePicker should be usable on mobile

- [ ] T030 [P] Add error handling for edge cases:
  - Empty task list: Search panel still appears and functional
  - Task without description: Text search still works (searches title only)
  - Task without dates: Date filters exclude task (expected behavior)
  - Invalid Jalali date input: JalaliDatePicker handles (no changes needed)

- [ ] T031 Update TaskList.svelte imports and structure:
  - Remove test/console.log statements
  - Clean up any unused variables
  - Ensure consistent code formatting

- [ ] T032 Test complete user journeys (manual):
  - Journey 1: Search for "API" → find matching task → verify result
  - Journey 2: Filter by date range → find tasks in period → verify result
  - Journey 3: Combine text + date filters → find specific tasks → verify result
  - Journey 4: Clear filters → see all tasks again → verify result
  - Journey 5: Apply filters → no results → clear → results return → verify

- [ ] T033 [P] Verify performance meets success criteria:
  - **SC-001**: Search returns results within 2 seconds ✅ (client-side, < 100ms typical)
  - **SC-002**: UI updates within 500ms of filter change ✅ (Svelte 5 reactivity, instant)
  - **SC-004**: Jalali date picker works correctly ✅ (reused component)
  - **SC-005**: Combined filters work 100% ✅ (AND logic validated)

- [ ] T034 Validate against quickstart.md implementation guide:
  - Ensure all 5 implementation steps completed
  - Verify component structure matches guide
  - Confirm filtering functions match design
  - Check styling matches Tailwind patterns in project

- [ ] T035 [P] Document any deviations from specification:
  - If changes made, update comments in code
  - Note any UX improvements made during implementation
  - Record any technical decisions not in spec

- [ ] T036 Final verification against acceptance scenarios:
  - US1 Scenario 1: Enter "API" → only API tasks shown ✅
  - US1 Scenario 2: Enter "auth" → case-insensitive match ✅
  - US1 Scenario 3: No matches → "No tasks found" ✅
  - US1 Scenario 4: Clear search → all tasks shown ✅
  - US2 Scenario 1: Select start date range → filtered ✅
  - US2 Scenario 2: Select due date range → filtered ✅
  - US2 Scenario 3: Both date ranges → filtered ✅
  - US2 Scenario 4: Clear dates → all shown ✅
  - US3 Scenario 1: Text + dates → all match ✅
  - US3 Scenario 2: Modify filter → immediate update ✅
  - US3 Scenario 3: No matches → message shown ✅

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies → Start immediately
- **Foundational (Phase 2)**: Depends on Setup → BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase
  - **US1 (Phase 3)**: Can start after Foundational
  - **US2 (Phase 4)**: Can start after Foundational (independent from US1)
  - **US3 (Phase 5)**: Depends on US1 + US2 (combines both features)
- **Polish (Phase 6)**: Depends on all user stories complete

### Parallel Opportunities

**Phase 1 (Setup)** - All tasks can run in sequence (informational):
- Understanding existing code is fastest in linear order

**Phase 2 (Foundational)** - Tasks T005-T007 can be parallelized:
- T005 (date utils) independent from T006 (filter utils) and T007 (component)
- Suggested approach: Parallel on different developers
  - Dev A: T005
  - Dev B: T006
  - Dev C: T007

**Phase 3 (US1)** - Tasks can be parallelized:
- T008 (text input) and T009 (filter state) parallel → T010 (derived) → T011 (iteration) → T012 (message) → T013 (count) → T014 (integration) → T015 (styling)
- Suggested teams:
  - Dev A: T008, T009, T015 (text input & styling)
  - Dev B: T010-T014 (filtering logic & integration)

**Phase 4 (US2)** - Can run in parallel with Phase 3 (different files):
- T016-T018 can be parallel with Phase 3 tasks
- Suggested: Dev C on Phase 4 while Dev A+B complete Phase 3

**Phase 5 (US3)** - Must wait for Phase 3 + 4:
- Tasks are validation/testing of combined logic
- Can start once Phase 4 T018 completes

**Phase 6 (Polish)** - Tasks mostly parallel:
- All marked [P] can run in parallel
- Non-parallel tasks are sequential cleanup at end

### Estimated Timeline (Single Developer)

| Phase | Tasks | Est. Time | Notes |
|-------|-------|-----------|-------|
| Phase 1 (Setup) | T001-T004 | 30 min | Review existing code |
| Phase 2 (Foundational) | T005-T007 | 60 min | Core utilities |
| Phase 3 (US1) | T008-T015 | 90 min | Text search complete |
| Phase 4 (US2) | T016-T021 | 90 min | Date filtering complete |
| Phase 5 (US3) | T022-T025 | 30 min | Validation only |
| Phase 6 (Polish) | T026-T036 | 60 min | Refinement & testing |
| **Total** | 36 tasks | **360 min (6 hours)** | Development + testing |

### Estimated Timeline (3 Developers - Parallel)

| Phase | Tasks | Timeline | Assignment |
|-------|-------|----------|-----------|
| Phase 1 | T001-T004 | 30 min | Dev A (lead) |
| Phase 2 | T005-T007 | 60 min | Dev A, B, C in parallel |
| Phase 3 + 4 | T008-T021 | 90 min | Dev A (US1), Dev B (US2) in parallel |
| Phase 5 | T022-T025 | 30 min | Dev A or B |
| Phase 6 | T026-T036 | 60 min | All 3 (parallel) |
| **Total** | 36 tasks | **180 min (3 hours)** | Wall clock time |

---

## Task Grouping by File

### frontend/src/components/TaskSearch.svelte (New)
- T007 (skeleton), T008 (text input), T016 (date fields), T017 (JalaliDatePicker), T018 (handlers), T020 (styling), T026 (clear button), T025 (summary)

### frontend/src/components/TaskList.svelte (Modified)
- T009 (filter state), T010 (derived), T011 (iteration), T012 (message), T013 (count), T014 (integration), T031 (cleanup)

### frontend/src/lib/dateUtils.js (New)
- T005 (utility functions)

### frontend/src/lib/filterUtils.js (New)
- T006 (filter functions)

### Testing & Verification (All)
- T019 (test filtering), T021 (test dates), T022 (verify combined), T023 (test scenarios), T024 (no results), T027 (clear), T028 (accessibility), T029 (mobile), T030 (edge cases), T032 (journeys), T033 (performance), T034 (validation), T035 (deviations), T036 (acceptance)

---

## Acceptance Criteria Mapping

Each user story task list satisfies its acceptance scenarios:

**User Story 1 (Text Search)**:
- [US1 Scenario 1] ← T011 (iteration), T010 (filtering), T008 (input field)
- [US1 Scenario 2] ← T006 (case-insensitive matching in filterUtils)
- [US1 Scenario 3] ← T012 (No tasks found message)
- [US1 Scenario 4] ← T013 (count update), T026 (clear button)

**User Story 2 (Date Range)**:
- [US2 Scenario 1] ← T016 (date fields), T019 (filtering logic), T010 (derived)
- [US2 Scenario 2] ← T016 (due date field), T019 (filtering logic)
- [US2 Scenario 3] ← T022 (combined logic), T006 (AND logic)
- [US2 Scenario 4] ← T027 (individual clear), T021 (test clear)

**User Story 3 (Combined)**:
- [US3 Scenario 1] ← T022 (verify combined), T023 (test scenarios)
- [US3 Scenario 2] ← T010 (real-time update via $derived)
- [US3 Scenario 3] ← T024 (no results message), T025 (filter summary)

---

## Success Criteria Verification

| Criteria | Verified By | Task(s) |
|----------|-------------|---------|
| SC-001: 2 sec search response | T033 | Client-side filtering instant |
| SC-002: 500ms UI update | T033 | Svelte 5 $derived() reactive |
| SC-003: 90% user success | T032 | Manual journey testing |
| SC-004: Jalali picker works | T017 | Reused component validation |
| SC-005: Combined 100% correct | T022-T023 | AND logic verification |
| SC-006: 50% time improvement | T032 | Journey testing (subjective) |

---

## No Backend Changes Required

✅ **Confirmed**: All tasks are frontend-only
- No backend handlers created
- No database changes
- No API endpoint modifications
- Existing `/api/tasks?project_id=<id>` endpoint unchanged
- Zero impact to other features

---

## Ready for Implementation

All 36 tasks are specific, actionable, and ready for development. Each task includes:
- ✅ Clear description
- ✅ Exact file paths
- ✅ Dependencies identified
- ✅ Parallelization markers [P] where applicable
- ✅ User story mapping [US1, US2, US3]
- ✅ Acceptance criteria mapping
- ✅ Estimated effort (180 min parallel, 360 min serial)
