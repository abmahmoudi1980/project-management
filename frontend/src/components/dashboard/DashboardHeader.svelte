<script>
  import { onMount } from 'svelte';
  import { formatJalaliDate } from '../../lib/utils.js';

  let { user, onRefresh, isRefreshing = false } = $props();

  let today = $state('');

  onMount(() => {
    today = formatJalaliDate(new Date(), 'full');
  });

  let greeting = $derived.by(() => {
    if (!user?.username) return 'سلام';
    return `سلام، ${user.username}`;
  });
</script>

<header class="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4 mb-8 md:mb-10">
  <div class="space-y-1.5">
    <h1 class="text-3xl md:text-4xl font-semibold text-ink tracking-tight">
      {greeting}
    </h1>
    <div class="flex items-center gap-3 text-sm text-ink-muted">
      <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full bg-surface border border-border-subtle">
        <span class="w-1.5 h-1.5 rounded-full bg-brand-500 animate-pulse-dot"></span>
        <span>{today || '...'}</span>
      </span>
      <span class="hidden sm:inline">گزارش لحظه‌ای فعالیت‌ها</span>
    </div>
  </div>

  <div class="flex items-center gap-2">
    <button
      type="button"
      onclick={onRefresh}
      disabled={isRefreshing}
      aria-label="بروزرسانی داشبورد"
      class="inline-flex items-center gap-2 px-3.5 py-2 text-sm font-medium text-ink-muted bg-surface border border-border rounded-full hover:bg-surface-muted hover:text-ink active:scale-[0.98] transition-all disabled:opacity-50 shadow-soft focus:outline-none focus:ring-2 focus:ring-brand-500 focus:ring-offset-2 focus:ring-offset-canvas"
    >
      <svg
        class="w-4 h-4 {isRefreshing ? 'animate-spin' : ''}"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
        aria-hidden="true"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
        />
      </svg>
      <span>بروزرسانی</span>
    </button>
  </div>
</header>
