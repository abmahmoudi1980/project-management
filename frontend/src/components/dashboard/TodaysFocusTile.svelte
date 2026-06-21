<script>
  import { formatJalaliDate } from '../../lib/utils.js';

  let { tasks = [] } = $props();

  let focus = $derived.by(() => {
    if (!tasks || tasks.length === 0) return null;
    const rank = { Critical: 0, High: 1, Medium: 2, Low: 3 };
    return [...tasks].sort((a, b) => {
      const r = (rank[a.priority_label] ?? 9) - (rank[b.priority_label] ?? 9);
      if (r !== 0) return r;
      return new Date(a.due_date || 0) - new Date(b.due_date || 0);
    })[0];
  });

  const priorityTone = {
    'Critical': { wrap: 'bg-danger-50',  text: 'text-danger-600', dot: 'bg-danger-500',  label: 'فوری' },
    'High':     { wrap: 'bg-brand-50',   text: 'text-brand-600',  dot: 'bg-brand-500',   label: 'مهم' },
    'Medium':   { wrap: 'bg-info-50',    text: 'text-info-600',   dot: 'bg-info-500',    label: 'متوسط' },
    'Low':      { wrap: 'bg-surface-muted', text: 'text-ink-muted', dot: 'bg-ink-subtle', label: 'کم' },
  };
</script>

<article class="relative overflow-hidden bg-surface border border-border-subtle rounded-bento p-5 md:p-6 shadow-soft hover:shadow-lift transition-shadow h-full flex flex-col">
  {#if focus}
    {@const tone = priorityTone[focus.priority_label] || priorityTone.Low}
    <div
      aria-hidden="true"
      class="pointer-events-none absolute -bottom-16 -left-16 w-56 h-56 rounded-full {tone.wrap} blur-2xl opacity-60"
    ></div>

    <div class="relative flex flex-col h-full">
      <header class="flex items-center justify-between mb-4">
        <h3 class="text-base font-semibold text-ink">تمرکز امروز</h3>
        <span class="inline-flex items-center gap-1.5 text-[11px] font-medium px-2.5 py-1 rounded-full {tone.wrap} {tone.text}">
          <span class="w-1.5 h-1.5 rounded-full {tone.dot}" aria-hidden="true"></span>
          {tone.label}
        </span>
      </header>

      <p class="text-base md:text-lg font-semibold text-ink leading-snug tracking-tight mb-3">
        {focus.title}
      </p>

      <div class="flex items-center gap-2 text-xs text-ink-muted mb-6">
        <span class="truncate">{focus.project_name || '—'}</span>
        <span class="text-ink-subtle">·</span>
        <span class="tabular-nums">{formatJalaliDate(focus.due_date, 'short')}</span>
      </div>

      <a
        href={`#/projects/${focus.project_id}`}
        class="mt-auto inline-flex items-center justify-center gap-2 w-full px-4 py-2.5 text-sm font-medium text-white bg-ink hover:bg-ink-muted active:bg-ink rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-ink focus:ring-offset-2 focus:ring-offset-canvas min-h-[44px]"
      >
        <span>شروع کن</span>
        <svg class="w-4 h-4 rtl:rotate-180" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
        </svg>
      </a>
    </div>
  {:else}
    <div class="flex-1 flex flex-col items-center justify-center text-center py-8">
      <div class="w-12 h-12 rounded-full bg-success-50 text-success-600 flex items-center justify-center mb-3">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
      </div>
      <p class="text-sm font-medium text-ink mb-1">لیست خالی</p>
      <p class="text-xs text-ink-muted">وظیفه فعالی برای امروز ندارید</p>
    </div>
  {/if}
</article>
