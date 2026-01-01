# Feature Specification: Task Search

**Feature Branch**: `004-task-search`  
**Created**: 2026-01-01  
**Status**: Draft  
**Input**: User description: "Add a search panel on top of tasks to get a string and dates to search title, description and start_date and due_date, user should be able to enter Jalali date on dates input with the Jalali date picker that exists on the project"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Filter Tasks by Text Search (Priority: P1)

A project manager wants to quickly find specific tasks by searching for keywords in task titles and descriptions without needing to manually scroll through a potentially long task list.

**Why this priority**: This is the core search functionality that addresses the primary use case of finding tasks. It's fundamental to the feature and provides immediate value with minimal dependencies.

**Independent Test**: Can be fully tested by entering search terms in a text field and verifying that only tasks containing those keywords in title or description are displayed.

**Acceptance Scenarios**:

1. **Given** a task list with multiple tasks, **When** the user enters "API" in the text search field, **Then** only tasks with "API" in the title or description are displayed
2. **Given** a task list with multiple tasks, **When** the user enters a partial keyword like "auth", **Then** tasks containing "authentication", "authorize", etc. are displayed (case-insensitive)
3. **Given** a user searching for a keyword, **When** no tasks match the search term, **Then** a message indicating "No tasks found" is displayed
4. **Given** a user with an active search, **When** they clear the search field, **Then** all tasks are displayed again

---

### User Story 2 - Filter Tasks by Date Range (Priority: P1)

A project manager wants to filter tasks by their start date and due date to find tasks in a specific time period, enabling better project planning and deadline management.

**Why this priority**: Date filtering is essential for time-based project planning and is equally important as text search for task discovery.

**Independent Test**: Can be fully tested by selecting start and due date ranges using the Jalali date picker and verifying that only tasks within those date ranges are displayed.

**Acceptance Scenarios**:

1. **Given** a task list with tasks having different start dates, **When** the user selects a start date range, **Then** only tasks with start_date within that range are displayed
2. **Given** a task list with tasks having different due dates, **When** the user selects a due date range, **Then** only tasks with due_date within that range are displayed
3. **Given** a user applying both start and due date filters, **When** both filters are active, **Then** only tasks matching both date conditions are displayed
4. **Given** a user with a date filter active, **When** they clear the date fields, **Then** all tasks are displayed again (date filter removed)

---

### User Story 3 - Combined Search Filters (Priority: P2)

A project manager wants to combine text and date filters to narrow down their search results to a very specific set of tasks.

**Why this priority**: This provides advanced filtering capability for power users. It builds on the foundation of text and date filtering and adds flexibility for complex queries.

**Independent Test**: Can be fully tested by applying both text search and date filters simultaneously and verifying that results match all filter criteria.

**Acceptance Scenarios**:

1. **Given** a task list, **When** the user enters a search term AND selects both start and due date ranges, **Then** only tasks matching all three filter criteria are displayed
2. **Given** filters applied, **When** the user modifies one filter, **Then** results update immediately to reflect the change
3. **Given** combined filters with no matching results, **When** filters are active, **Then** a message is displayed and user can modify filters or clear them

---

### Edge Cases

- What happens when a task has no description? (Text search should still work with title only)
- How does the system handle invalid Jalali date inputs? (Should show an error message and prevent filtering)
- What if a task's start date is after its due date? (System should still display the task if it falls within selected ranges)
- How are overdue tasks handled in date filtering? (Should be included based on their actual dates)
- What happens when the task list is empty? (Search panel should still be available and functional)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST display a search panel positioned at the top of the tasks list with input fields for text search and date filters
- **FR-002**: System MUST provide a text input field that accepts a search string and filters tasks by matching against task title and description (case-insensitive)
- **FR-003**: System MUST provide two date input fields (Start Date and Due Date) that use the existing Jalali date picker component for date selection
- **FR-004**: System MUST filter tasks by start_date when a Start Date range is selected
- **FR-005**: System MUST filter tasks by due_date when a Due Date range is selected
- **FR-006**: System MUST support combining text search with date filters - all filters must be applied simultaneously (AND logic)
- **FR-007**: System MUST update the task list in real-time as the user modifies any search filter
- **FR-008**: System MUST display all tasks when search filters are cleared or empty
- **FR-009**: System MUST accept Jalali date format in the date input fields (matching the existing Jalali date picker format used in the project)
- **FR-010**: System MUST display a message when no tasks match the applied filters
- **FR-011**: System MUST allow users to clear individual filters independently

### Key Entities

- **Task**: Contains title, description, start_date, and due_date properties that are used for filtering
- **Search Filter State**: Maintains the current values of text search, start date, and due date filters applied to the task list

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can locate a specific task by text search within 2 seconds of entering a search term
- **SC-002**: Task list updates to show filtered results within 500ms of any filter change
- **SC-003**: 90% of users successfully complete their first task search without additional help
- **SC-004**: Date picker Jalali calendar displays correctly and users can select dates without errors
- **SC-005**: Combined filters (text + dates) work correctly 100% of the time when multiple filters are applied
- **SC-006**: Search feature reduces task discovery time by at least 50% compared to manual scrolling

## Assumptions

- The existing Jalali date picker component (referenced as `JalaliDatePicker.svelte`) can be reused or adapted for the search panel
- Date filtering uses exact date matching (tasks within the selected range are shown)
- Text search is case-insensitive and uses substring matching
- The TaskList component structure allows insertion of a search panel above the current task list display
- Tasks without a description or without start/due dates should still be searchable by available fields
