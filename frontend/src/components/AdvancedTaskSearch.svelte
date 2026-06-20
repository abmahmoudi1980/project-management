<script>
  import JalaliDatePicker from './JalaliDatePicker.svelte';
  import { dateToJalaliString } from '../lib/utils.js';
  import { hasActiveFilters, getActiveFilterCount } from '../lib/filterUtils.js';

  // Props - stores applied filters
  let { filters = $bindable({
    start_date_from: null,
    start_date_to: null,
    due_date_from: null,
    due_date_to: null
  }) } = $props();

  let isOpen = $state(false);

  // Local state for form inputs (not applied until button click)
  let localFilters = $state({
    start_date_from: null,
    start_date_to: null,
    due_date_from: null,
    due_date_to: null
  });

  // Initialize local state when panel opens
  function openPanel() {
    // Copy current filters to local state when opening
    localFilters = { ...filters };
    isOpen = true;
  }

  // Display values for JalaliDatePicker (derived from local state)
  let start_date_from_display = $derived(dateToJalaliString(localFilters.start_date_from));
  let start_date_to_display = $derived(dateToJalaliString(localFilters.start_date_to));
  let due_date_from_display = $derived(dateToJalaliString(localFilters.due_date_from));
  let due_date_to_display = $derived(dateToJalaliString(localFilters.due_date_to));

  // No need for $effect blocks anymore since display values are $derived

  function parseGregorianDate(yyyyMmDd) {
    if (!yyyyMmDd || typeof yyyyMmDd !== 'string') {
      return null;
    }

    const match = /^([0-9]{4})-([0-9]{2})-([0-9]{2})$/.exec(yyyyMmDd.trim());
    if (!match) {
      return null;
    }

    const year = Number(match[1]);
    const monthIndex = Number(match[2]) - 1;
    const day = Number(match[3]);

    const date = new Date(year, monthIndex, day);
    return Number.isNaN(date.getTime()) ? null : date;
  }

  /**
   * Handle start date 'from' changes
   */
  function handleStartDateFromChange(gregorianYyyyMmDd) {
    if (!gregorianYyyyMmDd) {
      localFilters.start_date_from = null;
    } else {
      localFilters.start_date_from = parseGregorianDate(gregorianYyyyMmDd);
    }
  }

  /**
   * Handle start date 'to' changes
   */
  function handleStartDateToChange(gregorianYyyyMmDd) {
    if (!gregorianYyyyMmDd) {
      localFilters.start_date_to = null;
    } else {
      localFilters.start_date_to = parseGregorianDate(gregorianYyyyMmDd);
    }
  }

  /**
   * Handle due date 'from' changes
   */
  function handleDueDateFromChange(gregorianYyyyMmDd) {
    if (!gregorianYyyyMmDd) {
      localFilters.due_date_from = null;
    } else {
      localFilters.due_date_from = parseGregorianDate(gregorianYyyyMmDd);
    }
  }

  /**
   * Handle due date 'to' changes
   */
  function handleDueDateToChange(gregorianYyyyMmDd) {
    if (!gregorianYyyyMmDd) {
      localFilters.due_date_to = null;
    } else {
      localFilters.due_date_to = parseGregorianDate(gregorianYyyyMmDd);
    }
  }

  /**
   * Apply filters - copy local state to parent filters and close panel
   */
  function applyFilters() {
    filters.start_date_from = localFilters.start_date_from;
    filters.start_date_to = localFilters.start_date_to;
    filters.due_date_from = localFilters.due_date_from;
    filters.due_date_to = localFilters.due_date_to;
    isOpen = false;
  }

  /**
   * Clear local filters only
   */
  function clearFilters() {
    localFilters = {
      start_date_from: null,
      start_date_to: null,
      due_date_from: null,
      due_date_to: null
    };
  }

  /**
   * Check if any local filters are active
   */
  function hasLocalFilters() {
    return localFilters.start_date_from || localFilters.start_date_to || localFilters.due_date_from || localFilters.due_date_to;
  }

  // Derived state for active filters count (from applied filters)
  let activeFilterCount = $derived(getActiveFilterCount({ ...filters, text: '' }));
  let isFiltersActive = $derived(hasActiveFilters(filters));

  function closePanel() {
    isOpen = false;
  }
</script>

<!-- Advanced Search Button -->
<div class="flex items-center gap-2">
  <button
    onclick={openPanel}
    class="px-4 py-2 text-sm font-medium rounded-lg border transition-colors
      {isFiltersActive
        ? 'bg-brand-50 border-brand-300 text-brand-700 hover:bg-brand-100'
        : 'bg-white border-slate-300 text-slate-700 hover:bg-slate-50'}"
  >
    <span>🔍 جستجوی پیشرفته</span>
    {#if isFiltersActive}
      <span class="ml-2 inline-block bg-brand-600 text-white text-xs rounded-full w-5 h-5 flex items-center justify-center">
        {activeFilterCount}
      </span>
    {/if}
  </button>
</div>

<!-- Advanced Search Panel (Slide-in from top) -->
{#if isOpen}
  <div class="advanced-search-panel bg-white rounded-xl shadow-lg border border-slate-200 p-6 mb-4 animate-in slide-in-from-top">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <h3 class="text-lg font-semibold text-slate-900">جستجوی پیشرفته</h3>
      <button
        onclick={closePanel}
        class="text-slate-500 hover:text-slate-700 transition-colors"
        aria-label="بستن"
      >
        ✕
      </button>
    </div>

    <!-- Filter Controls Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <!-- Start Date Range -->
      <div>
        <label for="start-date-from" class="block text-sm font-medium text-slate-700 mb-2">
          تاریخ شروع (از)
        </label>
        <JalaliDatePicker
          id="start-date-from"
          value={start_date_from_display}
          placeholder="1403/01/01"
          on:change={(e) => handleStartDateFromChange(e.detail)}
        />
      </div>

      <div>
        <label for="start-date-to" class="block text-sm font-medium text-slate-700 mb-2">
          تاریخ شروع (تا)
        </label>
        <JalaliDatePicker
          id="start-date-to"
          value={start_date_to_display}
          placeholder="1403/12/29"
          on:change={(e) => handleStartDateToChange(e.detail)}
        />
      </div>

      <!-- Due Date Range -->
      <div>
        <label for="due-date-from" class="block text-sm font-medium text-slate-700 mb-2">
          تاریخ پایان (از)
        </label>
        <JalaliDatePicker
          id="due-date-from"
          value={due_date_from_display}
          placeholder="1403/01/01"
          on:change={(e) => handleDueDateFromChange(e.detail)}
        />
      </div>

      <div>
        <label for="due-date-to" class="block text-sm font-medium text-slate-700 mb-2">
          تاریخ پایان (تا)
        </label>
        <JalaliDatePicker
          id="due-date-to"
          value={due_date_to_display}
          placeholder="1403/12/29"
          on:change={(e) => handleDueDateToChange(e.detail)}
        />
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="flex items-center justify-between pt-4 border-t border-slate-200">
      <button
        onclick={clearFilters}
        disabled={!hasLocalFilters()}
        class="px-4 py-2 text-sm font-medium rounded-lg transition-colors
          {hasLocalFilters()
            ? 'bg-slate-100 text-slate-700 hover:bg-slate-200'
            : 'bg-slate-50 text-slate-400 cursor-not-allowed'}"
      >
        پاک کردن فیلترها
      </button>

      <button
        onclick={applyFilters}
        class="px-6 py-2 text-sm font-medium bg-brand-600 text-white rounded-lg
          hover:bg-brand-700 transition-colors"
      >
        اعمال فیلترها
      </button>
    </div>
  </div>
{/if}

<style>
  .advanced-search-panel {
    animation: slideInFromTop 0.3s ease-out;
  }

  @keyframes slideInFromTop {
    from {
      opacity: 0;
      transform: translateY(-20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
