# Quickstart: Svelte 5 Upgrade

**Feature**: Upgrade Frontend to Svelte 5  
**Branch**: `002-svelte5-upgrade`

This quickstart provides a condensed reference for executing the Svelte 5 upgrade. For detailed explanations, see [plan.md](./plan.md) and [research.md](./research.md).

## Pre-Flight Checklist

- [ ] Current branch: `002-svelte5-upgrade`
- [ ] Backend is working (for API testing)
- [ ] Node.js 18+ installed
- [ ] npm working correctly
- [ ] Git working directory clean

## Step 1: Upgrade Dependencies (5 minutes)

```bash
cd frontend

# Backup current package-lock.json
cp package-lock.json package-lock.json.backup

# Update package.json
npm install svelte@^5.0.0 @sveltejs/vite-plugin-svelte@^5.0.0

# Verify installation
npm list svelte @sveltejs/vite-plugin-svelte
```

**Expected versions**:
- `svelte`: 5.x.x
- `@sveltejs/vite-plugin-svelte`: 5.x.x

**Test build**:
```bash
npm run build
npm run dev
```

Both should complete without errors. You can stop dev server after verification.

## Step 2: Migration Syntax Quick Reference

### Props Migration

**Before**:
```svelte
export let title = 'Default';
export let count;
```

**After**:
```svelte
let { title = 'Default', count } = $props();
```

### State Migration

**Before**:
```svelte
let value = 0;
```

**After**:
```svelte
let value = $state(0);
```

### Derived Values Migration

**Before**:
```svelte
$: doubled = value * 2;
```

**After**:
```svelte
let doubled = $derived(value * 2);
```

### Side Effects Migration

**Before**:
```svelte
$: {
  console.log('value changed:', value);
}
```

**After**:
```svelte
$effect(() => {
  console.log('value changed:', value);
});
```

## Step 3: Component Migration Order

Migrate in this order to maintain working state:

### 3.1 App.svelte (15 minutes)

**Location**: `src/App.svelte`

**Changes**:
1. Convert `selectedProject` to `$state(null)`
2. Keep `onMount` as-is (or convert to `$effect`)
3. Test: App loads, project list displays, selection works

**Test command**:
```bash
npm run dev
# Open http://localhost:5173
# Click on projects, verify task list displays
```

### 3.2 Modal.svelte (10 minutes)

**Location**: `src/components/Modal.svelte`

**Changes**:
1. Convert props: `let { show, onClose } = $props()`
2. Any internal state to `$state`
3. Test: Modal opens/closes, backdrop click works, ESC key works

### 3.3 ProjectForm.svelte (20 minutes)

**Location**: `src/components/ProjectForm.svelte`

**Changes**:
1. Convert all props to `$props()`
2. Convert form state variables to `$state`
3. Keep event dispatcher as-is
4. Test: Open form, fill fields, submit, verify creation

**Test checklist**:
- [ ] Form opens
- [ ] All fields accept input
- [ ] Submit creates project
- [ ] Form closes after submit
- [ ] New project appears in list

### 3.4 TaskForm.svelte (20 minutes)

**Location**: `src/components/TaskForm.svelte`

**Changes**: Same pattern as ProjectForm

**Test checklist**:
- [ ] Form opens
- [ ] All fields work (including date picker)
- [ ] Submit creates task
- [ ] New task appears in list

### 3.5 TimeLogForm.svelte (20 minutes)

**Location**: `src/components/TimeLogForm.svelte`

**Changes**: Same pattern as ProjectForm/TaskForm

**Test checklist**:
- [ ] Form opens
- [ ] Hours field accepts numbers
- [ ] Description field works
- [ ] Submit creates time log

### 3.6 ProjectList.svelte (15 minutes)

**Location**: `src/components/ProjectList.svelte`

**Changes**:
1. Convert props to `$props()`
2. Convert any local state to `$state`
3. Keep event dispatcher
4. Test: List displays, click selection works

### 3.7 TaskList.svelte (15 minutes)

**Location**: `src/components/TaskList.svelte`

**Changes**: Same pattern as ProjectList

**Test checklist**:
- [ ] Tasks display correctly
- [ ] Task details show
- [ ] Time logs display
- [ ] Edit/delete buttons work

### 3.8 JalaliDatePicker.svelte (25 minutes)

**Location**: `src/components/JalaliDatePicker.svelte`

**Changes**:
1. Convert props to `$props()`
2. Convert calendar state to `$state`
3. Migrate any derived calendar values to `$derived`
4. Test thoroughly with multiple date selections

**Test checklist**:
- [ ] Calendar opens
- [ ] Dates are selectable
- [ ] Selected date updates form
- [ ] Jalali conversion works correctly
- [ ] Calendar closes after selection

## Step 4: Verification Testing (30 minutes)

### Full Workflow Test

1. **Create Project**:
   ```
   - Click "افزودن پروژه" (Add Project)
   - Fill in title, identifier, description
   - Submit form
   - Verify project appears in list
   ```

2. **Add Tasks**:
   ```
   - Select project
   - Click "افزودن وظیفه" (Add Task)
   - Fill in task details
   - Select dates using date picker
   - Submit
   - Verify task appears
   ```

3. **Log Time**:
   ```
   - Select task
   - Add time log
   - Enter hours and description
   - Submit
   - Verify time log appears
   ```

4. **Edit Operations**:
   ```
   - Edit project
   - Edit task
   - Edit time log
   - Verify updates persist
   ```

5. **Delete Operations**:
   ```
   - Delete time log
   - Delete task
   - Delete project
   - Verify deletions work
   ```

### Console Check

Open browser DevTools and verify:
- [ ] No console errors
- [ ] No console warnings
- [ ] Network requests working (check Network tab)
- [ ] No 404s or failed requests

## Step 5: Build Verification (10 minutes)

```bash
# Production build
npm run build

# Check output
ls -lh dist/

# Test production build
npm run preview
# Open http://localhost:4173
# Run quick smoke test (create project, add task)
```

**Verify**:
- [ ] Build completes successfully
- [ ] No build warnings
- [ ] Bundle size reasonable (<150KB total)
- [ ] Preview works correctly

## Step 6: Store Verification (10 minutes)

**No code changes needed**, but verify stores work:

**Test**:
1. Open app, verify projects load from API
2. Create project, verify store updates
3. Switch between projects, verify tasks load
4. Check browser DevTools for any subscription issues

**Expected**: Everything works identically to Svelte 4 version.

## Step 7: Final Commit

```bash
cd /home/abolfazl/apps/project-management

# Check status
git status

# Add all changes
git add frontend/

# Commit
git commit -m "feat: upgrade frontend to Svelte 5 with runes migration

- Upgraded svelte from 4.2.0 to 5.0.0
- Upgraded @sveltejs/vite-plugin-svelte to 5.0.0
- Migrated all 8 components to Svelte 5 runes syntax
- Converted props from export let to $props()
- Converted state to $state()
- Converted reactive declarations to $derived()
- Converted side effects to $effect()
- Verified all functionality works identically to Svelte 4
- Stores remain compatible, no changes needed
- All tests passing"

# Push branch
git push origin 002-svelte5-upgrade
```

## Rollback Procedure

If critical issues found:

```bash
# Revert to previous commit
git log --oneline  # Find commit hash before upgrade
git revert <commit-hash>
git push

# Or reset and reinstall
cd frontend
npm install svelte@^4.2.0 @sveltejs/vite-plugin-svelte@^3.0.0
git checkout -- package-lock.json
npm install
```

## Troubleshooting

### Issue: Build fails after dependency upgrade

**Solution**:
```bash
rm -rf node_modules package-lock.json
npm cache clean --force
npm install
```

### Issue: HMR not working

**Solution**: 
- Stop dev server
- Clear Vite cache: `rm -rf node_modules/.vite`
- Restart: `npm run dev`

### Issue: Component not rendering

**Check**:
1. Console for errors
2. Props are properly destructured with `$props()`
3. State variables use `$state()`
4. No mixing of old and new syntax in same component

### Issue: Store subscription not working

**Solution**: Verify you're using `$storeName` syntax in template, not in `<script>`.

**Correct**:
```svelte
<script>
  import { projects } from './stores/projectStore';
</script>

{#each $projects as project}
  <!-- ... -->
{/each}
```

### Issue: Date picker broken

**Check**:
1. jalali-moment dependency still installed
2. Date prop passed correctly
3. Event handler receiving date updates
4. Console for any date conversion errors

## Success Criteria

Before considering upgrade complete:

- [ ] All 8 components migrated to runes
- [ ] Zero console errors
- [ ] Zero build warnings
- [ ] All user workflows tested and working
- [ ] Production build succeeds
- [ ] Bundle size acceptable
- [ ] Stores working correctly
- [ ] Date picker functional
- [ ] Forms submit successfully
- [ ] No regressions found

## Time Estimate

- **Dependency Upgrade**: 5 minutes
- **Component Migration**: 2-3 hours
- **Testing**: 30-45 minutes
- **Build Verification**: 10 minutes
- **Documentation**: Included in this guide

**Total**: 3-4 hours for complete upgrade

## Next Steps

After successful migration:
1. Monitor production for any edge cases
2. Update development documentation if needed
3. Consider future enhancements using Svelte 5 features
4. Keep dependencies up to date with minor/patch updates

## References

- [Full Plan](./plan.md)
- [Research Document](./research.md)
- [Feature Spec](./spec.md)
- [Svelte 5 Migration Guide](https://svelte.dev/docs/svelte/v5-migration-guide)
