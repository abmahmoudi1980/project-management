# Implementation Plan: Task Search

**Branch**: `004-task-search` | **Date**: 2026-01-01 | **Spec**: [spec.md](spec.md)  
**Input**: Feature specification from `/specs/004-task-search/spec.md`

## Summary

This feature adds a search panel above the task list to enable users to filter tasks by text (title/description) and date range (start_date/due_date) using Jalali date picker. The implementation is frontend-focused (Svelte 5 component) with optional backend enhancement for server-side filtering. The feature integrates with the existing Jalali date picker component and Svelte 5 reactive state management.

## Technical Context

**Language/Version**: JavaScript (ES6+) / Svelte 5  
**Primary Dependencies**: Svelte 5 (runes), `jalali-moment` (existing for Persian date handling), Tailwind CSS  
**Storage**: N/A (client-side filtering in MVP; backend search optional for P2)  
**Testing**: Manual testing (no test framework currently in use per AGENTS.md)  
**Target Platform**: Web browser (desktop and mobile-responsive)  
**Project Type**: Web application (frontend + backend available but MVP is frontend-only)  
**Performance Goals**: 500ms max filter response time (SC-002), real-time UI updates (FR-007)  
**Constraints**: Must reuse existing JalaliDatePicker component, Svelte 5 patterns, maintain accessibility  
**Scale/Scope**: Single task list view enhancement, 11 functional requirements, 3 user stories

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Details |
|-----------|--------|---------|
| **Scope Clarity** | ✅ PASS | Feature scope is narrowly defined: search/filter only, no CRUD operations |
| **Technology Stack Alignment** | ✅ PASS | Uses existing Svelte 5 + Jalali date picker; no new major dependencies required |
| **Architectural Fit** | ✅ PASS | Frontend component enhancement; layered backend available if needed; no DB schema changes required |
| **Testability** | ✅ PASS | User scenarios are independently testable; each filter type can be verified separately |
| **Reuse/DRY** | ✅ PASS | Leverages existing JalaliDatePicker component; follows established Svelte 5 patterns |

**No Constitution Violations**: Feature aligns with project architecture and principles.

## Project Structure

### Documentation (this feature)

```text
specs/004-task-search/
├── spec.md              # Feature specification (complete)
├── plan.md              # This file (Phase 1+ planning)
├── research.md          # Phase 0 output (research & discovery)
├── data-model.md        # Phase 1 output (entity definitions)
├── quickstart.md        # Phase 1 output (developer quickstart)
├── contracts/           # Phase 1 output (API contracts if needed)
└── checklists/
    └── requirements.md  # Quality validation checklist
```

### Source Code (repository root) - Frontend Focus

```text
frontend/src/
├── components/
│   ├── TaskSearch.svelte        # NEW: Search panel component (PR-001)
│   ├── TaskList.svelte          # MODIFY: Integrate TaskSearch above list
│   ├── JalaliDatePicker.svelte  # EXISTING: Reuse for date selection
│   └── [13 other existing components]
├── stores/
│   └── taskStore.js             # MODIFY: Add search state management if needed
└── lib/
    └── api.js                   # EXISTING: API client
```

**Backend** (Optional for Phase 2):
```text
backend/handlers/
└── task_handler.go              # MODIFY: Add search/filter query parameters

backend/services/
└── task_service.go              # MODIFY: Add filtering logic (server-side)

backend/repositories/
└── task_repository.go           # MODIFY: Add filtered query method
```

**Structure Decision**: 
- **Phase 1 (MVP)**: Frontend-only client-side filtering with new `TaskSearch.svelte` component
- **Phase 2 (Optional)**: Backend enhancement for server-side filtering (if performance requires)
- **Integration point**: TaskSearch state → TaskList filtering → Filtered UI rendering

## Implementation Approach

### Phase 1: Frontend Search Component (MVP)

1. **Create TaskSearch.svelte component**:
   - Text input field for keyword search
   - Two Jalali date picker fields (start date, due date)
   - Clear/Reset button
   - Real-time filter state management using Svelte 5 runes
   - Responsive design with Tailwind CSS

2. **Integrate with TaskList.svelte**:
   - Place TaskSearch above current TaskList
   - Pass filter state to TaskList
   - Implement client-side filtering logic
   - Display "No tasks found" message when appropriate

3. **Use Svelte 5 Patterns**:
   - State: `let filters = $state({ text: '', startDate: null, dueDate: null })`
   - Effects: `$effect()` for filter updates
   - Props: `let { tasks } = $props()` for data passing

### Phase 2: Optional Backend Enhancement

1. **Add API endpoint** (if needed for performance):
   - `/api/tasks/search?text=...&startDate=...&dueDate=...`
   - Handler → Service → Repository layering

2. **Database filtering**:
   - Implement SQL WHERE clauses for title/description/dates
   - Return filtered results from database

## Complexity Tracking

No Constitutional violations requiring justification. Feature is straightforward frontend enhancement.

## Known Unknowns (Research Phase)

- [ ] How is the existing JalaliDatePicker component structured? (Review JalaliDatePicker.svelte)
- [ ] What is the current date format stored in task.start_date and task.due_date? (Check task model)
- [ ] Does TaskList already have filtering capability that can be extended? (Review TaskList.svelte)
- [ ] How are tasks currently passed to TaskList component? (Check App.svelte and stores)
- [ ] Should date ranges be inclusive or exclusive? (Clarify with acceptance scenarios)

**Next Steps**: 
1. Run Phase 0 research to resolve unknowns
2. Generate research.md with findings
3. Proceed to Phase 1 design (data-model.md, contracts/)
4. Update agent context with new component
