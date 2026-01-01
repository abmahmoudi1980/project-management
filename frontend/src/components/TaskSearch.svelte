<script>
  import JalaliDatePicker from './JalaliDatePicker.svelte';
  import { jalaliStringToDate, dateToJalaliString } from '../lib/utils.js';
  import { hasActiveFilters, getActiveFilterCount } from '../lib/filterUtils.js';

  // Props
  let { filters = $bindable({
    text: '',
    start_date_from: null,
    start_date_to: null,
    due_date_from: null,
    due_date_to: null
  }) } = $props();

  // Display values for JalaliDatePicker (string format)
  let start_date_from_display = $state(dateToJalaliString(filters.start_date_from));
  let start_date_to_display = $state(dateToJalaliString(filters.start_date_to));
  let due_date_from_display = $state(dateToJalaliString(filters.due_date_from));
  let due_date_to_display = $state(dateToJalaliString(filters.due_date_to));

  // Update display values when filters change
  $effect(() => {
    start_date_from_display = dateToJalaliString(filters.start_date_from);
  });

  $effect(() => {
    start_date_to_display = dateToJalaliString(filters.start_date_to);
  });

  $effect(() => {
    due_date_from_display = dateToJalaliString(filters.due_date_from);
  });

  $effect(() => {
    due_date_to_display = dateToJalaliString(filters.due_date_to);
  });

  /**
   * Handle start date 'from' changes
   */
  function handleStartDateFromChange(jalaliStr) {
    if (!jalaliStr) {
      filters.start_date_from = null;
    } else {
      const parsedDate = jalaliStringToDate(jalaliStr);
      filters.start_date_from = parsedDate;
    }
  }

  /**
   * Handle start date 'to' changes
   */
  function handleStartDateToChange(jalaliStr) {
    if (!jalaliStr) {
      filters.start_date_to = null;
    } else {
      const parsedDate = jalaliStringToDate(jalaliStr);
      filters.start_date_to = parsedDate;
    }
  }

  /**
   * Handle due date 'from' changes
   */
  function handleDueDateFromChange(jalaliStr) {
    if (!jalaliStr) {
      filters.due_date_from = null;
    } else {
      const parsedDate = jalaliStringToDate(jalaliStr);
      filters.due_date_from = parsedDate;
    }
  }

  /**
   * Handle due date 'to' changes
   */
  function handleDueDateToChange(jalaliStr) {
    if (!jalaliStr) {
      filters.due_date_to = null;
    } else {
      const parsedDate = jalaliStringToDate(jalaliStr);
      filters.due_date_to = parsedDate;
    }
  }

  /**
   * Clear all filters
   */
  function clearFilters() {
    filters = {
      text: '',
      start_date_from: null,
      start_date_to: null,
      due_date_from: null,
      due_date_to: null
    };
  }

  // Derived state for active filters count
  let activeFilterCount = $derived(getActiveFilterCount(filters));
  let isFiltersActive = $derived(hasActiveFilters(filters));
</script>

<div class="task-search-panel bg-white rounded-xl shadow-sm border border-slate-200 p-4 mb-4">
  <!-- Header -->
  <div class="flex items-center justify-between mb-4">
    <h3 class="text-sm font-medium text-slate-900">جستجو و فیلتر</h3>
    {#if isFiltersActive}
      <span class="text-xs text-slate-600">
        فیلترهای فعال: {activeFilterCount}
      </span>
    {/if}
  </div>

  <!-- Search and Filter Controls Grid -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-3">
    <!-- Text Search Input -->
    <div class="lg:col-span-2">
      <label for="text-search" class="block text-xs font-medium text-slate-700 mb-1">
        جستجو
      </label>
      <input
        id="text-search"
        type="text"
        bind:value={filters.text}
        placeholder="جستجو در عنوان و توضیح..."
        class="w-full px-3 py-2 border border-slate-300 rounded-lg text-sm
          placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500
          focus:border-transparent transition-colors"
      />
    </div>

    <!-- Start Date Range -->
    <div>
      <label for="start-date-from" class="block text-xs font-medium text-slate-700 mb-1">
        شروع (از)
      </label>
      <JalaliDatePicker
        id="start-date-from"
        value={start_date_from_display}
        placeholder="1403/01/01"
        onchange={(e) => handleStartDateFromChange(e.detail)}
      />
    </div>

    <div>
      <label for="start-date-to" class="block text-xs font-medium text-slate-700 mb-1">
        شروع (تا)
      </label>
      <JalaliDatePicker
        id="start-date-to"
        value={start_date_to_display}
        placeholder="1403/12/29"
        onchange={(e) => handleStartDateToChange(e.detail)}
      />
    </div>

    <!-- Due Date Range -->
    <div>
      <label for="due-date-from" class="block text-xs font-medium text-slate-700 mb-1">
        مهلت (از)
      </label>
      <JalaliDatePicker
        id="due-date-from"
        value={due_date_from_display}
        placeholder="1403/01/01"
        onchange={(e) => handleDueDateFromChange(e.detail)}
      />
    </div>

    <div>
      <label for="due-date-to" class="block text-xs font-medium text-slate-700 mb-1">
        مهلت (تا)
      </label>
      <JalaliDatePicker
        id="due-date-to"
        value={due_date_to_display}
        placeholder="1403/12/29"
        onchange={(e) => handleDueDateToChange(e.detail)}
      />
    </div>

    <!-- Clear Filters Button -->
    <div class="flex items-end">
      <button
        onclick={clearFilters}
        disabled={!isFiltersActive}
        class="w-full px-3 py-2 rounded-lg text-sm font-medium transition-colors
          {isFiltersActive
            ? 'bg-slate-100 text-slate-700 hover:bg-slate-200 active:bg-slate-300'
            : 'bg-slate-50 text-slate-400 cursor-not-allowed'}"
      >
        پاک کردن
      </button>
    </div>
  </div>
</div>

<style>
  .task-search-panel {
    direction: rtl;
  }
</style>
