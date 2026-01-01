# Research Findings: Task Search Feature

**Phase**: 0 (Research & Discovery)  
**Date**: 2026-01-01  
**Status**: Complete  
**Link**: [plan.md](plan.md)

## Overview

This document contains research findings addressing unknowns identified in the implementation plan. All findings are based on direct code examination and project conventions.

---

## Research Topic 1: JalaliDatePicker Component Structure

**Unknown**: How is the existing JalaliDatePicker component structured?

**Finding**: 
- **Location**: `frontend/src/components/JalaliDatePicker.svelte`
- **Architecture**: Standalone component using Svelte 5 runes
- **Props**: 
  - `value` (bindable string): Date in format YYYY/MM/DD or YYYY-MM-DD
  - `placeholder` (string): Display text, defaults to "مثال: 1403/10/10" (Jalali example)
  - `error` (boolean): For error state display
- **Key Features**:
  - Converts between Gregorian (YYYY-MM-DD) and Jalali (YYYY/MM/DD) formats
  - Uses `jalali-moment` library for date conversion
  - Provides interactive calendar UI with month/year navigation
  - Auto-formats input with slashes (e.g., "14031010" → "1403/10/10")
  - Validates date input before conversion
- **State Management**:
  - Uses Svelte 5 reactive variables: `$state()`, `$derived()`, `$effect()`
  - Dispatches custom events via `createEventDispatcher()`
- **Integration**: Can be directly reused in TaskSearch component with minimal adaptation

**Decision**: ✅ **REUSE existing JalaliDatePicker**
- Component is well-designed and handles all date format conversions needed
- Supports both Gregorian and Jalali input formats
- Already integrated with `jalali-moment` library used project-wide

---

## Research Topic 2: Task Date Fields Format

**Unknown**: What is the current date format stored in task.start_date and task.due_date?

**Finding**:
- **Backend Storage** (Go models in `backend/models/task.go`):
  - `StartDate: *time.Time` (Go time.Time type)
  - `DueDate: *time.Time` (Go time.Time type)
  - JSON serialization: ISO 8601 format (YYYY-MM-DDThh:mm:ssZ)
  - Both fields are optional (pointers, can be nil)

- **Frontend Handling** (in `TaskList.svelte`):
  - Used with helper function: `formatJalaliDate(task.start_date)`
  - Display example: "شروع: 1403/10/15" (Jalali format)
  - Compared directly for due date validation: `new Date(task.due_date) < new Date()`

- **API Contract**:
  - Tasks received from API contain ISO 8601 dates: `"2025-01-15T00:00:00Z"`
  - JalaliDatePicker expects YYYY-MM-DD or YYYY/MM/DD format
  - Conversion handled by JalaliDatePicker component using `jalali-moment`

**Decision**: ✅ **Use ISO 8601 dates for backend filtering (if implemented)**
- Frontend receives dates as ISO strings and passes to JalaliDatePicker
- JalaliDatePicker converts to/from Jalali format
- Client-side filtering can parse both formats using `jalali-moment`

---

## Research Topic 3: TaskList Filtering Capability

**Unknown**: Does TaskList already have filtering capability that can be extended?

**Finding**:
- **Current Architecture** (in `TaskList.svelte`):
  - Uses `$tasks` store from `taskStore.js` 
  - Directly iterates over `$tasks.tasks` array
  - No built-in filtering logic currently exists
  - Pagination implemented via intersection observer for infinite scroll
  - Data is fetched from API: `api.tasks.getByProject(projectId, page, pageSize)`

- **Current Flow**:
  1. TaskList requests tasks from store
  2. Store calls API with `projectId` and pagination params
  3. TaskList renders returned tasks directly (no client-side filtering)
  4. No query parameters for search/filter sent to API

- **TaskStore State** (in `frontend/src/stores/taskStore.js`):
  - Maintains `tasks[]`, `currentPage`, `pageSize`, `total`, `hasMore`
  - No search/filter state currently exists

**Decision**: ✅ **Implement client-side filtering in TaskList (MVP)**
- Add TaskSearch component above task list
- Add filter state to TaskList component using Svelte 5 runes
- Filter `$tasks.tasks` in memory before rendering
- No API changes required for MVP
- Optional Phase 2: Backend filtering via query parameters for performance at scale

---

## Research Topic 4: Task Data Flow Architecture

**Unknown**: How are tasks currently passed to TaskList component?

**Finding**:
- **Data Source Chain**:
  - `App.svelte` → displays ProjectList/Project selection
  - Selected project passed to `TaskList.svelte` as prop
  - `TaskList.svelte` calls `tasks.load(project.id)` on mount
  - Store fetches from API: `GET /api/tasks?project_id=<id>&page=<page>&page_size=<size>`
  - Tasks stored in `$tasks.tasks` (reactive store)
  - UI re-renders on store update

- **Component Props** (TaskList.svelte):
  - Receives: `{ project }` via `let { project } = $props()`
  - Data: Actually comes from reactive store subscription

- **Reactive Pattern**:
  - Uses Svelte store subscriptions with `$` prefix: `$tasks.tasks`, `$tasks.total`
  - Store updates trigger automatic component re-renders

**Decision**: ✅ **Add filtering state to TaskList component**
- Create reactive variables for search filters using Svelte 5 runes
- Use `$derived()` to compute filtered task list
- Pass filtered list to rendering loop instead of `$tasks.tasks`
- Maintain original store data (no mutation)

---

## Research Topic 5: Date Range Filtering Behavior

**Unknown**: Should date ranges be inclusive or exclusive?

**Finding**:
- **User Stories** (from spec.md):
  - "tasks with start_date within that range"
  - "tasks with due_date within that range"
  - Implies inclusive matching on both boundaries

- **Acceptance Scenario**:
  - User selects "start date from 1403/01/01 to 1403/31/10"
  - System shows "tasks with start_date within that range"
  - No explicit exclusion mentioned

- **Date Comparison Examples** (from existing code in TaskList.svelte):
  - Overdue detection: `new Date(task.due_date) < new Date()` (strict less-than)
  - Display: Shows actual dates as-is

**Decision**: ✅ **Use INCLUSIVE date range matching**
- Implementation: `startDate <= task.start_date <= endDate`
- Applies to both start_date and due_date ranges
- Aligns with user expectations ("within that range" = inclusive)
- Easier to understand and predict for users
- Example: Selecting "1403/01/01 to 1403/10/31" shows all tasks on those days

---

## Research Topic 6: Task Fields for Text Search

**Unknown**: Which task fields should text search cover?

**Finding**:
- **Specification Requirements** (from spec.md FR-002):
  - "filters tasks by matching against task title and description"
  - Case-insensitive matching required
  - Substring matching (partial keywords)

- **Available Task Fields**:
  - `title` (string): Task name - PRIMARY target
  - `description` (string): Task details - PRIMARY target
  - `category` (optional string): Task category
  - `priority` (string): High/Medium/Low
  - NOT specified for search: priority, category, assignee, etc.

- **Current Frontend Display** (TaskList.svelte):
  - Title shown prominently
  - Description shown with `line-clamp-2` (visible to users)
  - Other metadata (category, dates, assignee) shown separately

**Decision**: ✅ **Search only title and description (per spec)**
- Implementation: Check if title OR description contains search string
- Case-insensitive: Convert both to lowercase before comparing
- Substring matching: Use `string.includes(searchTerm.toLowerCase())`
- Future: Can extend to other fields if needed

---

## Research Topic 7: Empty State Handling

**Unknown**: How should the UI handle empty results?

**Finding**:
- **Specification Requirements** (from spec.md FR-010):
  - "System MUST display a message when no tasks match the applied filters"
  
- **Acceptance Scenarios**:
  - Story 1: "a message indicating 'No tasks found' is displayed"
  - Story 2: "a message is displayed and user can modify filters or clear them"

- **Current UI Pattern** (in TaskList.svelte):
  - Shows task count at top: `{$tasks.total} وظیفه` (tasks in Persian)
  - Task list uses `{#each $tasks.tasks || [] as task}` pattern
  - No current empty state message

- **Suggested Message** (in Persian):
  - Text: "هیچ وظیفه‌ای یافت نشد" (No tasks found)
  - Placement: Where task list would be rendered
  - Action: Show "Clear Filters" button

**Decision**: ✅ **Display "No tasks found" message**
- Show when filtered task count === 0
- Display below search panel
- Include "Clear Filters" button to reset search state
- Maintain original message style from project (Persian RTL layout)

---

## Implementation Notes

### Client-Side Filtering Algorithm (MVP)

```
1. Get search text from input field (if any)
2. Get start_date and due_date from date pickers (if any)
3. For each task in $tasks.tasks:
   - Check text filter: title.includes(text) OR description.includes(text)
   - Check start_date filter: task.start_date >= startDate AND task.start_date <= endDate
   - Check due_date filter: task.due_date >= dueDate AND task.due_date <= endDate
   - Include task if ALL active filters match (AND logic)
4. Return filtered array to render loop
5. Display "No tasks found" if result array is empty
```

### Performance Considerations

- **MVP Target**: < 500ms filter response (SC-002)
- **Client-side filtering**: O(n) complexity with 10-100 tasks = negligible
- **Scaling concern**: If task list grows to 1000+ items, implement backend filtering
- **Phase 2 option**: Add `?search=&start_date=&due_date=` query parameters to API

### Jalali Date Conversion Flow

```
User Input (Jalali picker) → "1403/10/15"
                          ↓
        JalaliDatePicker converts to ISO → "2024-12-31"
                          ↓
        Store in filter state → "2024-12-31"
                          ↓
        JavaScript Date comparison → new Date("2024-12-31")
                          ↓
        Compare with task ISO dates → task.start_date < userDate < task.due_date
```

---

## Decisions Summary

| Unknown | Decision | Rationale |
|---------|----------|-----------|
| JalaliDatePicker reuse | ✅ Reuse as-is | Already handles all required conversions |
| Date format handling | ✅ ISO 8601 backend, Jalali UI | Matches existing patterns |
| Filtering location | ✅ Client-side MVP | Simpler, faster for typical task counts |
| Date range behavior | ✅ Inclusive on both ends | User-friendly, intuitive |
| Text search fields | ✅ Title + Description only | Per specification |
| Empty state | ✅ "No tasks found" message | Required by spec, improves UX |

---

## Ready for Phase 1 Design

All research unknowns have been resolved. Proceed to:
1. ✅ Technical context fully understood
2. ✅ Component architecture identified
3. ✅ Integration points mapped
4. ✅ Data flow validated
5. → Next: Create data-model.md, contracts/, and quickstart.md
