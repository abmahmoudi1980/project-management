# Implementation Plan: Upgrade Frontend to Svelte 5

**Branch**: `002-svelte5-upgrade` | **Date**: 2025-12-30 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/002-svelte5-upgrade/spec.md`

## Summary

This feature upgrades the frontend from Svelte 4.2.0 to Svelte 5, leveraging the new runes system for improved reactivity and developer experience. The upgrade must maintain 100% backward compatibility with existing functionality while modernizing the component syntax.

**Primary requirements:**
- Upgrade Svelte and related dependencies to version 5
- Migrate all 8 components to Svelte 5 runes syntax (`$state`, `$derived`, `$props`, `$effect`)
- Ensure zero breaking changes to user-facing functionality
- Update build configuration for Svelte 5 compatibility

**Technical approach:** Package upgrades followed by systematic component migration using Svelte 5's runes system, with comprehensive testing at each stage to ensure functionality preservation.

## Technical Context

**Language/Version**: Node.js 18+ with Svelte 5 (upgrading from Svelte 4.2.0)
**Primary Dependencies**: Svelte 5.x, @sveltejs/vite-plugin-svelte 5.x, Vite 5.0.0, Tailwind CSS 3.4.0, jalali-moment 3.3.11
**Storage**: N/A (frontend uses backend API)
**Testing**: Manual testing workflow (no automated tests currently)
**Target Platform**: Modern web browsers (Chrome, Firefox, Safari, Edge latest versions)
**Project Type**: Web application frontend (SPA with Vite)
**Performance Goals**: Bundle size ≤ current size, HMR <2s, initial load <1s
**Constraints**: Zero breaking changes to existing features, maintain all current functionality
**Scale/Scope**: 8 components, 3 stores, 1 API service, <100KB total bundle size

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

[No violations - Constitution file is a template placeholder. This is a straightforward upgrade that maintains existing patterns and architecture.]

## Project Structure

### Documentation (this feature)

```text
specs/002-svelte5-upgrade/
├── plan.md              # This file
├── spec.md              # Feature specification
├── research.md          # Svelte 5 migration guide
├── quickstart.md        # Quick reference for upgrade
└── checklists/
    └── requirements.md   # Requirements checklist
```

### Source Code (repository root)

```text
frontend/
├── package.json                    # UPDATE: Svelte 5 dependencies
├── vite.config.js                 # UPDATE: Plugin configuration
├── src/
│   ├── App.svelte                 # MIGRATE: to Svelte 5 runes
│   ├── app.css                    # NO CHANGE
│   ├── main.js                    # NO CHANGE
│   ├── components/
│   │   ├── Modal.svelte           # MIGRATE: to Svelte 5 runes
│   │   ├── ProjectForm.svelte     # MIGRATE: to Svelte 5 runes
│   │   ├── ProjectList.svelte     # MIGRATE: to Svelte 5 runes
│   │   ├── TaskForm.svelte        # MIGRATE: to Svelte 5 runes
│   │   ├── TaskList.svelte        # MIGRATE: to Svelte 5 runes
│   │   ├── TimeLogForm.svelte     # MIGRATE: to Svelte 5 runes
│   │   └── JalaliDatePicker.svelte # MIGRATE: to Svelte 5 runes
│   ├── stores/
│   │   ├── projectStore.js        # VERIFY: Svelte 5 store compatibility
│   │   ├── taskStore.js           # VERIFY: Svelte 5 store compatibility
│   │   └── timeLogStore.js        # VERIFY: Svelte 5 store compatibility
│   └── lib/
│       └── api.js                 # NO CHANGE
├── index.html                     # NO CHANGE
├── postcss.config.js              # NO CHANGE
└── tailwind.config.js             # NO CHANGE
```

**Structure Decision**: This is a web application frontend using Svelte + Vite. The upgrade focuses on the Svelte components and build configuration. All components need migration to the new runes syntax, while stores and utilities remain compatible.

## Phase-by-Phase Implementation

### Phase 0: Dependency Upgrades & Research

**Goal**: Update npm packages and understand Svelte 5 migration patterns.

**Steps**:
1. Read Svelte 5 migration guide and document key changes in `research.md`
2. Update `package.json` dependencies:
   - `svelte`: `^4.2.0` → `^5.0.0`
   - `@sveltejs/vite-plugin-svelte`: `^3.0.0` → `^5.0.0`
3. Run `npm install` to update lock file
4. Verify build still works: `npm run build`
5. Verify dev server starts: `npm run dev`

**Key Research Areas**:
- Svelte 5 runes system (`$state`, `$derived`, `$props`, `$effect`)
- Props migration from `export let` to `$props()`
- Reactive declarations migration from `$:` to `$derived`
- Side effects migration from `$:` statements to `$effect`
- Store compatibility with Svelte 5
- Event handling changes (if any)
- Lifecycle hook changes
- Two-way binding with `bind:` directive

**Validation**:
- Dependencies installed without errors
- Build completes successfully
- Dev server starts without errors
- No console errors on initial load

### Phase 1: Component Migration - Core App

**Goal**: Migrate the main App.svelte component to establish migration pattern.

**File**: `frontend/src/App.svelte`

**Migration Steps**:
1. Replace `export let` props with `let { propName } = $props()`
2. Convert reactive declarations (`$: value = ...`) to `$derived`
3. Convert side effects (`$: { ... }`) to `$effect`
4. Update `onMount` usage if needed
5. Test component renders correctly
6. Test all user interactions work

**Example Migrations**:

**Before (Svelte 4)**:
```svelte
<script>
  export let initialValue = 0;
  let count = 0;
  
  $: doubled = count * 2;
  
  $: {
    console.log('count changed:', count);
  }
</script>
```

**After (Svelte 5)**:
```svelte
<script>
  let { initialValue = 0 } = $props();
  let count = $state(0);
  
  let doubled = $derived(count * 2);
  
  $effect(() => {
    console.log('count changed:', count);
  });
</script>
```

**Validation**:
- App loads without errors
- Project list displays correctly
- Project selection works
- Task list displays when project selected

### Phase 2: Component Migration - Forms

**Goal**: Migrate form components to Svelte 5 runes.

**Files**:
- `frontend/src/components/ProjectForm.svelte`
- `frontend/src/components/TaskForm.svelte`
- `frontend/src/components/TimeLogForm.svelte`

**Migration Focus**:
- Convert form state variables to `$state`
- Update two-way bindings (verify they work with `$state`)
- Migrate form validation logic
- Update event handlers
- Test form submission

**Validation for each component**:
- Form opens correctly
- All fields accept input
- Validation works as before
- Form submission creates/updates entities
- Form closes after successful submission
- Error messages display correctly

### Phase 3: Component Migration - Lists

**Goal**: Migrate list/display components to Svelte 5 runes.

**Files**:
- `frontend/src/components/ProjectList.svelte`
- `frontend/src/components/TaskList.svelte`

**Migration Focus**:
- Convert reactive state to `$state`
- Migrate derived values to `$derived`
- Update list filtering/sorting logic
- Verify event dispatching works

**Validation for each component**:
- Lists display data correctly
- Sorting works (if applicable)
- Filtering works (if applicable)
- Item selection works
- Click handlers work
- Events dispatch correctly to parent

### Phase 4: Component Migration - Shared Components

**Goal**: Migrate shared utility components to Svelte 5 runes.

**Files**:
- `frontend/src/components/Modal.svelte`
- `frontend/src/components/JalaliDatePicker.svelte`

**Migration Focus**:
- Update props handling
- Migrate modal open/close state
- Update date picker reactivity
- Verify keyboard/accessibility features still work

**Validation**:
- Modal opens/closes correctly
- Modal backdrop click closes modal
- Escape key closes modal
- Date picker displays calendar
- Date selection updates form
- Jalali to Gregorian conversion works

### Phase 5: Store Compatibility Verification

**Goal**: Ensure Svelte stores work correctly with Svelte 5 components.

**Files**:
- `frontend/src/stores/projectStore.js`
- `frontend/src/stores/taskStore.js`
- `frontend/src/stores/timeLogStore.js`

**Verification Steps**:
1. Confirm stores use standard Svelte store API (writable, readable, derived)
2. Test store subscriptions in migrated components
3. Verify store auto-unsubscription works
4. Test store update/set operations
5. Verify derived stores update correctly

**Note**: Svelte stores are backward compatible in Svelte 5. If any issues arise, consider migrating to runes-based state management.

**Validation**:
- Store data loads correctly
- Store updates trigger component re-renders
- No memory leaks from subscriptions
- Store methods (load, create, update, delete) work

### Phase 6: Build Configuration & Optimization

**Goal**: Finalize Svelte 5 plugin configuration and optimize build output.

**File**: `frontend/vite.config.js`

**Configuration Updates**:
1. Verify Svelte plugin options for Svelte 5
2. Check if any new plugin options improve performance
3. Update any deprecated options
4. Test production build
5. Compare bundle sizes (before/after)

**Optimization Checks**:
- Tree shaking working correctly
- CSS properly extracted
- No duplicate dependencies
- Source maps generated correctly

**Validation**:
- `npm run build` succeeds
- Production build runs without errors
- Bundle size acceptable (no significant increase)
- App works in production mode (`npm run preview`)

### Phase 7: Comprehensive Testing & Validation

**Goal**: Validate all functionality works identically to Svelte 4 version.

**Test Scenarios** (manual testing workflow):

1. **Project Management**:
   - Create new project
   - Edit project details
   - Delete project
   - Select project from list

2. **Task Management**:
   - Add task to project
   - Edit task details
   - Update task status
   - Delete task
   - View task list

3. **Time Logging**:
   - Add time log
   - Edit time log
   - Delete time log
   - View time logs by task

4. **Date Picker**:
   - Open date picker
   - Select date
   - Verify Jalali calendar displays correctly
   - Verify selected date updates form

5. **UI Interactions**:
   - Modal open/close
   - Form validation
   - Error messages
   - Loading states
   - Responsive layout

6. **Browser Compatibility**:
   - Test in Chrome
   - Test in Firefox
   - Test in Safari (if available)
   - Test in Edge (if available)

**Success Criteria**:
- All test scenarios pass
- No console errors
- No console warnings
- Identical behavior to pre-upgrade version

## Migration Checklist

- [ ] Phase 0: Dependencies upgraded
- [ ] Phase 0: Research completed
- [ ] Phase 0: Build verification passed
- [ ] Phase 1: App.svelte migrated
- [ ] Phase 1: App functionality validated
- [ ] Phase 2: ProjectForm migrated & tested
- [ ] Phase 2: TaskForm migrated & tested
- [ ] Phase 2: TimeLogForm migrated & tested
- [ ] Phase 3: ProjectList migrated & tested
- [ ] Phase 3: TaskList migrated & tested
- [ ] Phase 4: Modal migrated & tested
- [ ] Phase 4: JalaliDatePicker migrated & tested
- [ ] Phase 5: Stores verified compatible
- [ ] Phase 6: Build config updated
- [ ] Phase 6: Production build validated
- [ ] Phase 7: Full integration testing passed
- [ ] Phase 7: Browser compatibility verified

## Risk Mitigation

**Risk 1**: Breaking changes in Svelte 5 not covered by migration guide
- **Mitigation**: Test each component immediately after migration before proceeding

**Risk 2**: Third-party library (jalali-moment) incompatibility
- **Mitigation**: Test date picker early in migration; have fallback plan to vendor or replace

**Risk 3**: Subtle reactivity behavior changes
- **Mitigation**: Comprehensive manual testing of all user interactions; compare side-by-side with Svelte 4 version

**Risk 4**: Build performance regression
- **Mitigation**: Measure and document build times and bundle sizes at each phase

**Risk 5**: Store compatibility issues
- **Mitigation**: Test store functionality early; consider migrating stores to runes if problems arise

## Success Metrics

- **Code Quality**: All components use Svelte 5 runes (100% migration)
- **Functionality**: Zero regressions in user-facing features
- **Performance**: Bundle size ≤ Svelte 4 version
- **Build**: Zero errors, zero warnings in build output
- **Browser Support**: Works in all previously supported browsers
- **Developer Experience**: HMR as fast or faster than Svelte 4

## Rollback Plan

If critical issues discovered after deployment:
1. Revert to previous commit on branch 001-enhance-entities-with-redmine-fields
2. Document the issue
3. Research solution
4. Re-attempt upgrade after fix identified

Git workflow:
```bash
# If issues found
git revert HEAD
git push

# After fix identified
git checkout 002-svelte5-upgrade
# Continue with fix
```
