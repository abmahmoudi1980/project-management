# Agent Context File

**Last Updated**: 2025-12-30  
**Branch**: 002-svelte5-upgrade  
**Purpose**: Provide AI coding agents with project-specific context

## Technology Stack

### Backend
- **Language**: Go 1.21+
- **Framework**: Fiber (Go web framework)
- **Database**: PostgreSQL 12+ (using pgx driver)
- **Architecture**: Layered (handlers → services → repositories → models)

### Frontend  
- **Language**: JavaScript (ES6+)
- **Framework**: Svelte 5.0.0 (upgraded from Svelte 4.2.0) ✅
- **Build Tool**: Vite 6.0.0 (upgraded from 5.0.0) ✅
- **Vite Plugin**: @sveltejs/vite-plugin-svelte 5.0.0 (upgraded from 3.0.0) ✅
- **Styling**: Tailwind CSS 3.4.0
- **State Management**: Svelte stores (writable/readable)
- **Date Library**: jalali-moment 3.3.11 (Persian calendar)

### Development
- **Node.js**: 18+
- **npm**: Package management
- **Git**: Version control

## Project Structure

```
project-management/
├── backend/
│   ├── models/       # Data structures
│   ├── repositories/ # Database operations
│   ├── services/     # Business logic
│   ├── handlers/     # HTTP handlers
│   ├── routes/       # Route definitions
│   ├── config/       # Configuration
│   └── migration/    # Database migrations
├── frontend/
│   ├── src/
│   │   ├── components/  # Svelte components (8 components)
│   │   ├── stores/      # State management (3 stores)
│   │   └── lib/         # Utilities (API client)
│   └── package.json
└── specs/
    └── [###-feature-name]/
        ├── spec.md       # Feature specification
        ├── plan.md       # Implementation plan
        ├── research.md   # Technical research
        ├── quickstart.md # Quick reference
        └── tasks.md      # Task breakdown
```

## Current Feature Context

### Feature 002: Svelte 5 Upgrade (COMPLETED) ✅

**Status**: Implementation Complete  
**Branch**: `002-svelte5-upgrade`

**Objective**: Upgrade frontend from Svelte 4.2.0 to Svelte 5 with full runes migration while maintaining 100% backward compatibility.

**Key Changes**:
- Svelte: 4.2.0 → 5.0.0 ✅
- @sveltejs/vite-plugin-svelte: 3.0.0 → 5.0.0 ✅
- Vite: 5.0.0 → 6.0.0 ✅
- All components migrated to runes syntax ✅
- Props: `export let` → `$props()` ✅
- State: `let x = 0` → `let x = $state(0)` ✅
- Derived: `$: computed = x * 2` → `let computed = $derived(x * 2)` ✅
- Effects: `$: { sideEffect() }` → `$effect(() => { sideEffect() })` ✅
- Stores remain compatible (no changes needed) ✅

**Components Migrated** (8 total):
1. ✅ App.svelte - Main application component
2. ✅ Modal.svelte - Reusable modal dialog
3. ✅ ProjectForm.svelte - Project creation/editing
4. ✅ ProjectList.svelte - Project listing
5. ✅ TaskForm.svelte - Task creation/editing
6. ✅ TaskList.svelte - Task listing
7. ✅ TimeLogForm.svelte - Time log entry
8. ✅ JalaliDatePicker.svelte - Persian calendar date picker

**Stores** (3 total - no migration needed):
1. ✅ projectStore.js - Project state management
2. ✅ taskStore.js - Task state management
3. ✅ timeLogStore.js - Time log state management

**Additional Changes**:
- Fixed nested button HTML validation issue in ProjectList.svelte
- Vite upgraded to 6.0.0 to satisfy peer dependencies

### Feature 001: Enhance Entities with Redmine Fields (COMPLETED)

**Status**: Implementation in progress  
**Branch**: `001-enhance-entities-with-redmine-fields`

**Objective**: Add Redmine-inspired fields to Project, Task, and TimeLog entities.

**Database Changes**:
- Projects: identifier, homepage, is_public
- Tasks: description, assignee_id, author_id, category, start_date, due_date, estimated_hours, done_ratio
- Time Logs: user_id, activity_type, project_id (optional)

## Svelte 5 Key Patterns

### Component Structure
```svelte
<script>
  // Props (replaces export let)
  let { title = 'Default', count, onClick } = $props();
  
  // Local state (reactive)
  let localValue = $state(0);
  
  // Computed values (replaces $:)
  let doubled = $derived(localValue * 2);
  
  // Side effects (replaces $: blocks)
  $effect(() => {
    console.log('localValue changed:', localValue);
  });
  
  // Lifecycle (unchanged)
  import { onMount } from 'svelte';
  onMount(() => {
    // Initialization
  });
  
  // Events (unchanged)
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
</script>

<!-- Template (unchanged) -->
<div>Content</div>
```

### Store Usage (Unchanged)
```svelte
<script>
  import { myStore } from './stores/myStore';
</script>

{#each $myStore as item}
  <div>{item.name}</div>
{/each}
```

## Development Commands

### Frontend
```bash
cd frontend
npm install          # Install dependencies
npm run dev          # Development server (port 5173)
npm run build        # Production build
npm run preview      # Preview production build
```

### Backend
```bash
cd backend
go run main.go       # Start server (port 3000)
go build            # Build binary
```

## API Endpoints

Base URL: `http://localhost:3000/api`

### Projects
- GET    /api/projects
- POST   /api/projects
- GET    /api/projects/:id
- PUT    /api/projects/:id
- DELETE /api/projects/:id

### Tasks
- GET    /api/projects/:project_id/tasks
- POST   /api/projects/:project_id/tasks
- GET    /api/tasks/:id
- PUT    /api/tasks/:id
- DELETE /api/tasks/:id

### Time Logs
- GET    /api/tasks/:task_id/timelogs
- POST   /api/tasks/:task_id/timelogs
- GET    /api/timelogs/:id
- PUT    /api/timelogs/:id
- DELETE /api/timelogs/:id

## Testing Strategy

Currently: Manual testing workflow

**Frontend Testing**:
1. Build verification (`npm run build`)
2. Development server verification (`npm run dev`)
3. Manual user workflow testing
4. Browser console verification (no errors/warnings)

**Backend Testing**:
No test framework configured (future enhancement)

## Code Style Guidelines

### Svelte 5 Best Practices
- Use runes for all reactive state
- Prefer `$derived` over manual calculations
- Use `$effect` only for side effects (not derived state)
- Keep lifecycle hooks (`onMount`, etc.) when appropriate
- Maintain store usage with `$` prefix in templates
- Use event dispatchers for component communication

### Go Best Practices
- Layered architecture: handler → service → repository
- Use struct pointers for optional fields
- UUID primary keys
- JSON tags on all model fields
- Error handling with proper HTTP status codes

## Common Issues & Solutions

### Svelte 5 Migration
- **Issue**: Props not working
  - **Solution**: Use `let { prop } = $props()` not `export let prop`
  
- **Issue**: State not reactive
  - **Solution**: Use `$state()` for all mutable state
  
- **Issue**: Derived value not updating
  - **Solution**: Use `$derived()` not regular assignment

### Build Issues
- **Issue**: HMR not working
  - **Solution**: Clear Vite cache (`rm -rf node_modules/.vite`)

## Resources

- Svelte 5 Docs: https://svelte.dev/docs/svelte/overview
- Svelte 5 Migration: https://svelte.dev/docs/svelte/v5-migration-guide
- Feature Specs: `specs/[###-feature-name]/`

---

**Note**: This file is auto-updated during `/speckit.plan` workflow. Manual edits between markers will be preserved.
