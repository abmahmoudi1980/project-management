<script>
  import { formatJalaliDate } from '../../lib/utils.js';

  let { tasks = [], onComplete = () => {}, onSelect = () => {} } = $props();

  let filter = $state('all');

  const filters = [
    { id: 'all',     label: 'همه' },
    { id: 'urgent',  label: 'فوری' },
    { id: 'today',   label: 'امروز' },
  ];

  const priorityTone = {
    'Critical': { dot: 'bg-danger-500',  text: 'text-danger-600' },
    'High':     { dot: 'bg-brand-500',   text: 'text-brand-600' },
    'Medium':   { dot: 'bg-info-500',    text: 'text-info-600' },
    'Low':      { dot: 'bg-ink-subtle',  text: 'text-ink-muted' },
  };

  function isToday(d) {
    if (!d) return false;
    const t = new Date(d);
    const n = new Date();
    return t.getFullYear() === n.getFullYear() && t.getMonth() === n.getMonth() && t.getDate() === n.getDate();
  }

  let visible = $derived.by(() => {
    if (filter === 'urgent') return tasks.filter(t => t.priority_label === 'Critical' || t.priority_label === 'High');
    if (filter === 'today')  return tasks.filter(t => isToday(t.due_date));
    return tasks;
  });

  let completed = $state(new Set());

  async function toggle(taskId) {
    if (completed.has(taskId)) return;
    const next = new Set(completed);
    next.add(taskId);
    completed = next;
    setTimeout(() => onComplete(taskId), 350);
  }
</script>

<article class="bg-surface border border-border-subtle rounded-bento p-5 md:p-6 shadow-soft hover:shadow-lift transition-shadow h-full flex flex-col">
  <header class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 mb-4">
    <div class="flex items-center gap-3">
      <h3 class="text-base font-semibold text-ink">وظایف من</h3>
      <span class="inline-flex items-center justify-center min-w-[24px] h-6 px-2 rounded-full text-xs font-semibold bg-brand-50 text-brand-700 tabular-nums">
        {tasks.length}
      </span>
    </div>
    <div role="tablist" class="inline-flex items-center gap-1 p-1 bg-surface-muted rounded-full self-start sm:self-auto">
      {#each filters as f}
        <button
          type="button"
          role="tab"
          aria-selected={filter === f.id}
          onclick={() => (filter = f.id)}
          class="px-3 py-1 text-xs font-medium rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-brand-500 focus:ring-offset-2 focus:ring-offset-canvas {filter === f.id ? 'bg-surface text-ink shadow-soft' : 'text-ink-muted hover:text-ink'}"
        >{f.label}</button>
      {/each}
    </div>
  </header>

  {#if visible.length === 0}
    <div class="flex-1 flex flex-col items-center justify-center text-center py-10">
      <div class="w-12 h-12 rounded-full bg-success-50 text-success-600 flex items-center justify-center mb-3">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
      </div>
      <p class="text-sm font-medium text-ink mb-1">همه چیز مرتبه</p>
      <p class="text-xs text-ink-muted">وظیفه‌ای در این فیلتر وجود ندارد</p>
    </div>
  {:else}
    <ul class="divide-y divide-border-subtle -mx-2 flex-1">
      {#each visible as task (task.id)}
        {@const tone = priorityTone[task.priority_label] || priorityTone.Low}
        {@const done = completed.has(task.id)}
        <li class="px-2 transition-opacity {done ? 'opacity-50' : ''}">
          <div class="flex items-center gap-3 py-3 group">
            <button
              type="button"
              role="checkbox"
              aria-checked={done}
              aria-label={done ? `لغو تکمیل ${task.title}` : `تکمیل ${task.title}`}
              onclick={() => toggle(task.id)}
              class="relative w-5 h-5 rounded-md border-2 flex-shrink-0 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-brand-500 focus:ring-offset-2 focus:ring-offset-canvas {done ? 'bg-brand-500 border-brand-500' : 'border-border-strong hover:border-brand-400 group-active:scale-90'}"
            >
              {#if done}
                <svg class="absolute inset-0 w-full h-full text-white p-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                </svg>
              {/if}
            </button>

            <button
              type="button"
              onclick={() => onSelect(task.project_id)}
              class="flex-1 min-w-0 text-right rtl:text-right"
            >
              <p class="text-sm font-medium text-ink truncate {done ? 'line-through' : ''}">{task.title}</p>
              <p class="text-xs text-ink-muted truncate mt-0.5">{task.project_name || '—'}</p>
            </button>

            <div class="flex items-center gap-3 flex-shrink-0">
              <div class="hidden sm:flex items-center gap-1.5">
                <span class="w-1.5 h-1.5 rounded-full {tone.dot}" aria-hidden="true"></span>
                <span class="text-[11px] font-medium {tone.text}">{task.priority_label}</span>
              </div>
              <span class="text-[11px] text-ink-subtle tabular-nums">
                {formatJalaliDate(task.due_date, 'relative')}
              </span>
            </div>
          </div>
        </li>
      {/each}
    </ul>

    <footer class="pt-4 mt-2 border-t border-border-subtle text-center">
      <a
        href="#/projects"
        class="text-xs font-medium text-brand-600 hover:text-brand-700 transition-colors"
      >مشاهده لیست کامل</a>
    </footer>
  {/if}
</article>
