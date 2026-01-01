# Quickstart Guide: Task Search Implementation

**Phase**: 1 (Design)  
**Date**: 2026-01-01  
**Target Audience**: Developers  
**Complexity**: Beginner to Intermediate

## Overview

This guide walks through implementing the task search feature from specification to working code. Expected time: **2-4 hours**.

---

## Pre-Implementation Checklist

Before starting, ensure:

- [ ] Svelte 5 dev environment running (`npm run dev` in `frontend/`)
- [ ] Backend running (`go run main.go` in `backend/`)
- [ ] Understanding of Svelte 5 runes (`$state`, `$derived`, `$effect`)
- [ ] Understanding of JalaliDatePicker component
- [ ] Read [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md)

---

## Implementation Steps

### Step 1: Create TaskSearch Component (30 minutes)

Create new file: `frontend/src/components/TaskSearch.svelte`

**Purpose**: Standalone component for search UI and filter state management

**Key Features**:
- Text input field for keyword search
- Two Jalali date picker fields (start date, due date)
- Clear button to reset filters
- Real-time state updates using Svelte 5 runes

**Skeleton**:
```svelte
<script>
  import JalaliDatePicker from './JalaliDatePicker.svelte';

  // Filter state (Svelte 5 runes)
  let filters = $state({
    text: '',
    start_date_from: null,
    start_date_to: null,
    due_date_from: null,
    due_date_to: null
  });

  // Helper to convert JalaliDatePicker string to Date
  function parseDateString(dateStr) {
    if (!dateStr) return null;
    const [year, month, day] = dateStr.split('/').map(Number);
    // Convert Jalali to Gregorian and return Date object
    return convertJalaliToGregorian(year, month, day);
  }

  // Clear all filters
  function clearFilters() {
    filters = {
      text: '',
      start_date_from: null,
      start_date_to: null,
      due_date_from: null,
      due_date_to: null
    };
  }
</script>

<div class="search-panel">
  <!-- Text Search Input -->
  <input 
    type="text" 
    placeholder="جستجو در وظایف..."
    bind:value={filters.text}
  />
  
  <!-- Date Range Inputs -->
  <JalaliDatePicker 
    bind:value={filters.start_date_from_display}
    placeholder="تاریخ شروع (از)"
  />
  
  <!-- More date pickers... -->
  
  <!-- Clear Button -->
  <button onclick={clearFilters}>پاک کردن</button>
</div>

<style>
  .search-panel {
    display: flex;
    gap: 1rem;
    padding: 1rem;
    background: white;
    border-radius: 0.75rem;
    border: 1px solid #e2e8f0;
  }
</style>
```

**Next**: See [Step 2](#step-2-integrate-with-tasklist-component) before styling.

---

### Step 2: Integrate with TaskList Component (40 minutes)

**File**: `frontend/src/components/TaskList.svelte`

**Changes**:

1. **Import TaskSearch component** (at top of script):
```javascript
import TaskSearch from './TaskSearch.svelte';
```

2. **Add filter state** to TaskList (same structure as TaskSearch):
```javascript
let filters = $state({
  text: '',
  start_date_from: null,
  start_date_to: null,
  due_date_from: null,
  due_date_to: null
});
```

3. **Create filtering functions**:
```javascript
function evaluateTextFilter(task) {
  if (!filters.text.trim()) return true;
  const searchTerm = filters.text.toLowerCase();
  return (
    task.title?.toLowerCase().includes(searchTerm) ||
    task.description?.toLowerCase().includes(searchTerm)
  );
}

function evaluateDateFilter(task, dateField, fromFilter, toFilter) {
  if (fromFilter === null && toFilter === null) return true;
  const taskDate = task[dateField];
  if (!taskDate) return false;
  
  const date = new Date(taskDate);
  if (fromFilter && date < fromFilter) return false;
  if (toFilter && date > toFilter) return false;
  return true;
}
```

4. **Create derived filtered task list**:
```javascript
let filteredTasks = $derived.by(() => {
  return ($tasks.tasks || []).filter(task => {
    const textMatch = evaluateTextFilter(task);
    const startMatch = evaluateDateFilter(task, 'start_date', 
      filters.start_date_from, filters.start_date_to);
    const dueMatch = evaluateDateFilter(task, 'due_date', 
      filters.due_date_from, filters.due_date_to);
    return textMatch && startMatch && dueMatch;
  });
});

let resultCount = $derived(filteredTasks.length);
```

5. **Add TaskSearch component above task list** (in template):
```svelte
<TaskSearch bind:filters />
```

6. **Update task count display**:
```svelte
{#if $tasks.total > 0}
  <div class="text-sm text-slate-500">
    نتایج: {resultCount} / {$tasks.total}
  </div>
{/if}
```

7. **Update task iteration loop**:
```svelte
<!-- Change from: -->
{#each $tasks.tasks || [] as task}

<!-- To: -->
{#each filteredTasks as task}
```

8. **Add empty state message**:
```svelte
{#if filteredTasks.length === 0 && $tasks.tasks.length > 0}
  <div class="empty-state">
    <p>هیچ وظیفه‌ای یافت نشد</p>
    {#if hasActiveFilters()}
      <button onclick={clearFilters}>پاک کردن فیلترها</button>
    {/if}
  </div>
{/if}
```

---

### Step 3: Complete TaskSearch Component Styling (30 minutes)

**File**: `frontend/src/components/TaskSearch.svelte`

Focus areas:
- Responsive layout (mobile first)
- RTL support (Persian)
- Tailwind CSS classes (match project style)
- Accessibility (labels, ARIA attributes)

**Template Structure**:
```svelte
<div class="search-panel bg-white rounded-xl shadow-sm border border-slate-200 p-4">
  <!-- Title -->
  <h3 class="text-sm font-medium text-slate-900 mb-4">جستجو و فیلتر</h3>
  
  <!-- Search Grid -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-3">
    <!-- Text Search -->
    <div class="lg:col-span-2">
      <label class="block text-xs font-medium text-slate-700 mb-1">
        جستجو
      </label>
      <input 
        type="text" 
        bind:value={filters.text}
        placeholder="جستجو در عنوان و توضیح..."
        class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm"
      />
    </div>
    
    <!-- Start Date From -->
    <div>
      <label class="block text-xs font-medium text-slate-700 mb-1">
        شروع (از)
      </label>
      <JalaliDatePicker 
        value={filters.start_date_from_display}
        placeholder="1403/01/01"
      />
    </div>
    
    <!-- Additional date pickers with similar structure -->
    
    <!-- Clear Button -->
    <div class="flex items-end">
      <button 
        onclick={clearFilters}
        class="w-full px-3 py-2 bg-slate-100 text-slate-700 rounded-lg 
               hover:bg-slate-200 text-sm font-medium transition-colors"
      >
        پاک کردن
      </button>
    </div>
  </div>
</div>
```

**Key Styling Considerations**:
- Mobile: Stack vertically
- Tablet: 2 columns
- Desktop: 5 columns with text search spanning 2
- RTL text alignment (text-right)
- Consistent color palette with TaskList

---

### Step 4: Add Jalali Date Conversion Logic (20 minutes)

**File**: `frontend/src/components/TaskSearch.svelte` (utilities section)

```javascript
import moment from 'jalali-moment';

// Convert Jalali date string (1403/10/15) to JavaScript Date (Gregorian)
function jalaliStringToDate(jalaliStr) {
  if (!jalaliStr) return null;
  try {
    // Parse Jalali date and convert to Gregorian
    const momentObj = moment.from(jalaliStr, 'fa', 'YYYY/MM/DD');
    return momentObj.toDate();
  } catch (e) {
    console.error('Invalid Jalali date:', jalaliStr);
    return null;
  }
}

// Convert JavaScript Date (Gregorian) to Jalali string (1403/10/15)
function dateToJalaliString(date) {
  if (!date) return '';
  try {
    return moment(date).locale('fa').format('YYYY/MM/DD');
  } catch (e) {
    console.error('Invalid date:', date);
    return '';
  }
}

// Handle date picker changes
function handleStartDateFromChange(event) {
  const jalaliStr = event.target.value; // From JalaliDatePicker
  filters.start_date_from = jalaliStringToDate(jalaliStr);
}
```

---

### Step 5: Test the Implementation (30 minutes)

**Manual Test Cases**:

1. **Text Search**:
   - [ ] Type "API" → only tasks with "API" shown
   - [ ] Type "api" (lowercase) → case-insensitive match
   - [ ] Type "xyz" → "No tasks found" message

2. **Date Range**:
   - [ ] Select start date range → only tasks in range shown
   - [ ] Select due date range → only tasks in range shown
   - [ ] Set invalid range (from > to) → show all tasks (no match)

3. **Combined Filters**:
   - [ ] Set text + both dates → tasks matching all criteria
   - [ ] Modify one filter → results update immediately

4. **Clear Filters**:
   - [ ] Click "Clear" button → all filters reset
   - [ ] Clear individual field → only that filter removed

5. **UI/UX**:
   - [ ] Mobile layout (< 768px) - stack vertically
   - [ ] Desktop layout - side-by-side
   - [ ] Button hovers - visual feedback
   - [ ] Empty state message - clear and helpful

---

## File Structure After Implementation

```
frontend/src/
├── components/
│   ├── TaskSearch.svelte         # NEW
│   ├── TaskList.svelte           # MODIFIED
│   ├── JalaliDatePicker.svelte   # (unchanged, reused)
│   └── [other components...]
└── [other directories...]
```

## No Backend Changes Required

✅ **MVP does not modify**:
- `backend/` files
- Database schema
- API endpoints
- Any Go code

---

## Expected Results

After completing all steps, you should have:

- ✅ Search panel appearing above task list
- ✅ Text search filtering title and description
- ✅ Date range filtering by start and due dates
- ✅ Combined filters with AND logic
- ✅ Real-time updates as filters change
- ✅ "No tasks found" message when appropriate
- ✅ Clear button to reset all filters
- ✅ Responsive design (mobile, tablet, desktop)
- ✅ Persian language support with RTL layout

---

## Common Issues & Solutions

### Issue: Dates not filtering correctly

**Cause**: Date comparison timezone mismatch

**Solution**: 
```javascript
function normalizeDate(date) {
  // Strip time portion to compare only dates
  return new Date(date.getFullYear(), date.getMonth(), date.getDate());
}
```

### Issue: Clear button not updating UI

**Cause**: Derived values not recalculating

**Solution**: Ensure `$derived.by()` is used for `filteredTasks`

### Issue: Jalali date picker not working

**Cause**: Missing `jalali-moment` import or binding

**Solution**: Verify import and use `bind:value={}` on component

---

## Next Steps After MVP

1. **Testing**: Add unit tests for filter functions
2. **Performance**: Profile with 1000+ tasks
3. **Phase 2**: Implement backend search endpoint if needed
4. **Analytics**: Track most common search terms
5. **UX**: Add search suggestions or autocomplete

---

## Resources

- [Specification](spec.md)
- [Research Findings](research.md)
- [Data Model](data-model.md)
- [Svelte 5 Runes](https://svelte.dev/docs/svelte/runes)
- [JalaliDatePicker Component](../../../frontend/src/components/JalaliDatePicker.svelte)
- [TaskList Component](../../../frontend/src/components/TaskList.svelte)

---

## Questions?

Refer to:
1. `research.md` for technical context
2. `data-model.md` for state structure
3. `contracts/API.md` for integration points
4. Existing components for code patterns
