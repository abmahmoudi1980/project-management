# Data Model: Task Search Feature

**Phase**: 1 (Design)  
**Date**: 2026-01-01  
**Status**: Design Complete  
**Link**: [research.md](research.md)

## Overview

This document defines the entities, state structures, and data transformations for the task search feature. No new database tables or backend entities are required for the MVP.

---

## Entities

### Task (Existing)

**Purpose**: Represents a project task - used for filtering in search results

**Fields**:
```
{
  id: UUID                    // Unique identifier
  project_id: UUID            // Parent project
  title: string               // Task name (searchable)
  description: string         // Task details (searchable)
  start_date: DateTime|null   // Task start (filterable)
  due_date: DateTime|null     // Task deadline (filterable)
  priority: string            // High | Medium | Low
  completed: boolean          // Completion status
  category: string|null       // Task category (optional)
  assignee_id: UUID|null      // Assigned user
  done_ratio: int             // 0-100 completion percentage
  estimated_hours: float|null // Time estimate
  created_at: DateTime        // Creation timestamp
  updated_at: DateTime        // Last modification timestamp
}
```

**Validation Rules**:
- `title`: Non-empty string (1-255 characters)
- `start_date`: Valid ISO 8601 date, optional
- `due_date`: Valid ISO 8601 date, optional
- Date consistency: `start_date` can be any valid date (may be > due_date)

**Notes**:
- Task entity is unchanged; existing API continues to serve full Task objects
- All fields available for display; only title, description, start_date, due_date used for filtering

---

### Search Filter State (New)

**Purpose**: Maintains user-input search filters at component level

**Location**: `TaskList.svelte` or `TaskSearch.svelte` component state

**Fields**:
```
{
  text: string                // Keyword search string (empty = no filter)
  start_date_from: Date|null  // Start date range begin (Gregorian Date object)
  start_date_to: Date|null    // Start date range end (Gregorian Date object)
  due_date_from: Date|null    // Due date range begin (Gregorian Date object)
  due_date_to: Date|null      // Due date range end (Gregorian Date object)
}
```

**Initial State**:
```
{
  text: "",
  start_date_from: null,
  start_date_to: null,
  due_date_from: null,
  due_date_to: null
}
```

**Constraints**:
- Text field: Max 255 characters
- Date fields: Must be valid ISO 8601 dates (from JalaliDatePicker)
- No requirement that `from_date <= to_date` (system handles invalid ranges gracefully)

**State Management** (Svelte 5):
```javascript
let filters = $state({
  text: "",
  start_date_from: null,
  start_date_to: null,
  due_date_from: null,
  due_date_to: null
});
```

---

## Data Transformations

### Text Search Filter

**Input**: 
- Filter state: `filters.text` (user-entered string)
- Task: `{ title: string, description: string }`

**Process**:
```
searchTerm = filters.text.toLowerCase().trim()
IF searchTerm is empty:
  ✓ Include task (no text filter active)
ELSE:
  IF task.title.toLowerCase().includes(searchTerm) 
     OR task.description.toLowerCase().includes(searchTerm):
    ✓ Include task
  ELSE:
    ✗ Exclude task
```

**Output**: Boolean (include/exclude)

**Examples**:
- Search: "API", Task Title: "Build API authentication" → ✓ Match
- Search: "api", Task Desc: "Create API endpoints" → ✓ Match (case-insensitive)
- Search: "test", Task: "Task Details", Desc: "Testing unit" → ✓ Match (substring)
- Search: "xyz", Task: "Build X-ray API", Desc: "Other" → ✗ No Match

---

### Date Range Filter (Start Date)

**Input**:
- Filter state: `filters.start_date_from` and `filters.start_date_to` (Date objects)
- Task: `{ start_date: ISO 8601 string | null }`

**Process**:
```
IF both filters are null:
  ✓ Include task (no start_date filter active)
ELSE:
  IF task.start_date is null:
    ✗ Exclude task (cannot match range without date)
  ELSE:
    taskDate = parseDate(task.start_date) → JavaScript Date
    IF filters.start_date_from !== null:
      IF taskDate < filters.start_date_from:
        ✗ Exclude task
    IF filters.start_date_to !== null:
      IF taskDate > filters.start_date_to:
        ✗ Exclude task
    ✓ Include task (within range)
```

**Output**: Boolean (include/exclude)

**Examples**:
- Filter: From 1403/01/01 to 1403/10/31
- Task start_date: "2024-03-20" → ✓ Within range
- Task start_date: "2024-12-31" → ✗ Outside range
- Task start_date: null → ✗ No date to match
- Filter: Only From 1403/01/01 (no To)
- Task start_date: "2024-06-15" → ✓ On or after From
- Task start_date: "2023-12-31" → ✗ Before From

---

### Date Range Filter (Due Date)

**Input**:
- Filter state: `filters.due_date_from` and `filters.due_date_to` (Date objects)
- Task: `{ due_date: ISO 8601 string | null }`

**Process**: (Same as Start Date filter above)

**Output**: Boolean (include/exclude)

---

### Combined Filters (AND Logic)

**Input**:
- All filter state fields
- Task with all properties

**Process**:
```
textMatch = evaluateTextFilter(task)
startDateMatch = evaluateStartDateFilter(task)
dueDateMatch = evaluateDueDateFilter(task)

IF textMatch AND startDateMatch AND dueDateMatch:
  ✓ Include task in results
ELSE:
  ✗ Exclude task from results
```

**Output**: Boolean (include/exclude)

**Example**:
- Filter: text="API", start_date_from="2024-03-01", due_date_to="2024-06-30"
- Task 1: title="API Auth", start="2024-03-15", due="2024-04-30"
  - textMatch=✓, startDateMatch=✓, dueDateMatch=✓ → **✓ INCLUDE**
- Task 2: title="Database", start="2024-03-15", due="2024-04-30"
  - textMatch=✗ → **✗ EXCLUDE** (fails text filter)
- Task 3: title="API Docs", start="2024-02-15", due="2024-04-30"
  - textMatch=✓, startDateMatch=✗ → **✗ EXCLUDE** (fails start date filter)

---

## Derived State

### Filtered Task List

**Definition** (Svelte 5 `$derived`):
```javascript
let filteredTasks = $derived.by(() => {
  return ($tasks.tasks || []).filter(task => {
    const textMatch = evaluateTextFilter(filters, task);
    const startMatch = evaluateStartDateFilter(filters, task);
    const dueMatch = evaluateDueDateFilter(filters, task);
    return textMatch && startMatch && dueMatch;
  });
});
```

**Purpose**: Real-time computed list of tasks matching all active filters

**Usage**: Render loop `{#each filteredTasks as task}`

**Reactivity**: Automatically updates when:
- `filters` state changes
- `$tasks.tasks` store updates
- Any task property changes

---

### Result Count

**Definition** (Svelte 5 `$derived`):
```javascript
let resultCount = $derived(filteredTasks.length);
```

**Purpose**: Display number of matching tasks

**Usage**: 
- "3 وظیفه" (3 tasks) - update from `$tasks.total` to `resultCount`
- Show "No tasks found" when `resultCount === 0`

---

### Has Active Filters

**Definition**:
```javascript
let hasActiveFilters = $derived.by(() => {
  return filters.text.trim() !== "" ||
         filters.start_date_from !== null ||
         filters.start_date_to !== null ||
         filters.due_date_from !== null ||
         filters.due_date_to !== null;
});
```

**Purpose**: Determine if any filters are active

**Usage**: 
- Show "Clear Filters" button only if `hasActiveFilters === true`
- Adjust messaging when no results shown

---

## State Transitions

### Clear All Filters

**From**: Any filter state with active filters
**To**: 
```
{
  text: "",
  start_date_from: null,
  start_date_to: null,
  due_date_from: null,
  due_date_to: null
}
```
**Trigger**: User clicks "Clear Filters" button
**Result**: `filteredTasks` updates to show all tasks from `$tasks.tasks`

### Clear Individual Filter

**Examples**:
- Clear text: `filters.text = ""`
- Clear start date range: `filters.start_date_from = null; filters.start_date_to = null`
- Clear due date: `filters.due_date_from = null`

**Trigger**: User clears input field or date picker
**Result**: `filteredTasks` updates, applying remaining active filters

### Update Text Filter

**From**: `{ text: "old search" }`
**To**: `{ text: "new search" }`
**Trigger**: User types in search input
**Debounce**: Optional - consider 300ms debounce for performance

### Update Date Filter

**From**: `{ start_date_from: null }`
**To**: `{ start_date_from: Date(2024-03-01) }`
**Trigger**: User selects date in JalaliDatePicker
**Result**: `filteredTasks` immediately updates with new date constraint

---

## No Backend Changes (MVP)

**Important**: The MVP implementation requires NO changes to:
- Database schema
- Backend Task model
- API endpoints
- Data serialization

All filtering happens in the frontend using data already available from the existing `/api/tasks?project_id=<id>&page=<page>` endpoint.

**Future Enhancement** (Phase 2):
If scaling requires server-side filtering, add optional query parameters:
```
GET /api/tasks?project_id=<id>&search=<text>&start_date_from=<date>&start_date_to=<date>&due_date_from=<date>&due_date_to=<date>&page=<page>&page_size=<size>
```

---

## Validation Rules

### Text Filter
- Max length: 255 characters
- Allow: Any UTF-8 characters (Persian, English, symbols)
- Validation: None required (substring matching handles all inputs)

### Date Filters
- Format: ISO 8601 (from JalaliDatePicker)
- Validation: Handled by JalaliDatePicker component (invalid dates rejected)
- No cross-field validation (from > to is allowed, results in no matches)

### Combined Filters
- No mutual exclusivity rules
- All filters use AND logic (not OR)
- Clearing one filter doesn't clear others

---

## Error Handling

### Invalid Date Input
- **Source**: User types invalid Jalali date in date picker
- **Handler**: JalaliDatePicker component validation (existing)
- **Result**: Date field cleared or reverted to previous valid value
- **UI**: Show error message on date picker

### Empty Task List
- **Cause**: No tasks in project OR all tasks filtered out
- **Handler**: Check `filteredTasks.length === 0`
- **Result**: Display "No tasks found" message instead of task list

### Missing Task Fields
- **Title missing**: Never occurs (required field)
- **Description missing**: Treat as empty string (search still works)
- **Dates missing**: Task excluded from date filter results (correct behavior)

---

## Example Complete Filter Evaluation

```
Input Task:
{
  id: "task-123",
  title: "Implement User Authentication API",
  description: "Create endpoints for login, logout, and token refresh",
  start_date: "2024-03-15T00:00:00Z",
  due_date: "2024-04-15T00:00:00Z",
  category: "Backend",
  priority: "High"
}

Active Filters:
{
  text: "API",
  start_date_from: 2024-03-01,
  start_date_to: 2024-03-31,
  due_date_from: null,
  due_date_to: null
}

Evaluation:
1. textMatch = "API".toLowerCase().includes("api") → true
            OR "Create endpoints...".toLowerCase().includes("api") → false
            → Result: TRUE

2. startMatch = 2024-03-15 >= 2024-03-01 AND 2024-03-15 <= 2024-03-31
            → Result: TRUE

3. dueMatch = (null = no filter)
            → Result: TRUE

4. Combined = TRUE AND TRUE AND TRUE
            → **INCLUDE TASK IN RESULTS**
```

---

## Ready for Component Design

Data model is complete. Proceed to:
- Create `TaskSearch.svelte` component with filter state
- Create `TaskList.svelte` modifications for filtered rendering
- Create API client updates (if needed)
