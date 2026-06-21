<script>
  import Avatar from '../Avatar.svelte';
  import { formatJalaliDate } from '../../lib/utils.js';

  let { projects = [], onSelect = () => {} } = $props();

  const statusMap = {
    'Planning':   { dot: 'bg-ink-subtle',     label: 'برنامه‌ریزی' },
    'In Progress':{ dot: 'bg-brand-500',      label: 'در حال اجرا' },
    'On Track':   { dot: 'bg-success-500',    label: 'طبق برنامه' },
    'Review':     { dot: 'bg-info-500',       label: 'بازبینی' },
    'active':     { dot: 'bg-brand-500',      label: 'فعال' },
  };

  function statusFor(s) {
    return statusMap[s] || { dot: 'bg-ink-subtle', label: s || '—' };
  }
</script>

<article class="bg-surface border border-border-subtle rounded-bento p-5 shadow-soft hover:shadow-lift transition-shadow h-full flex flex-col">
  <header class="flex items-center justify-between mb-4">
    <h3 class="text-base font-semibold text-ink">پروژه‌های اخیر</h3>
    <a
      href="#/projects"
      class="text-xs font-medium text-brand-600 hover:text-brand-700 transition-colors"
    >مشاهده همه</a>
  </header>

  {#if projects.length === 0}
    <div class="flex-1 flex flex-col items-center justify-center text-center py-8">
      <div class="w-10 h-10 rounded-full bg-surface-muted flex items-center justify-center mb-3 text-ink-subtle">
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7h18M3 12h18M3 17h18" />
        </svg>
      </div>
      <p class="text-sm text-ink-muted">پروژه‌ای برای نمایش وجود ندارد</p>
    </div>
  {:else}
    <ul class="space-y-1 -mx-2 flex-1">
      {#each projects.slice(0, 3) as project (project.id)}
        {@const status = statusFor(project.status)}
        <li>
          <button
            type="button"
            onclick={() => onSelect(project.id)}
            class="w-full text-right rtl:text-right px-3 py-3 rounded-xl hover:bg-surface-muted transition-colors focus:outline-none focus:ring-2 focus:ring-brand-500 focus:ring-inset"
          >
            <div class="flex items-center gap-3">
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2 mb-1">
                  <span class="w-1.5 h-1.5 rounded-full {status.dot}" aria-hidden="true"></span>
                  <span class="text-sm font-medium text-ink truncate">{project.name}</span>
                </div>
                <div class="flex items-center justify-between gap-3">
                  <p class="text-xs text-ink-muted truncate">{project.client || '—'}</p>
                  <span class="text-[10px] text-ink-subtle tabular-nums flex-shrink-0">
                    {formatJalaliDate(project.due_date, 'short')}
                  </span>
                </div>
                <div class="mt-2 h-1 rounded-full bg-surface-muted overflow-hidden">
                  <div
                    class="h-full bg-brand-500 rounded-full transition-all duration-500"
                    style="width: {project.progress || 0}%"
                  ></div>
                </div>
              </div>

              <div class="flex flex-col items-end gap-1 flex-shrink-0">
                <span class="text-xs font-semibold text-ink tabular-nums">{project.progress || 0}%</span>
                {#if (project.team_members || []).length > 0}
                  <div class="flex -space-x-1.5 rtl:space-x-reverse">
                    {#each (project.team_members || []).slice(0, 2) as member}
                      <div class="ring-1 ring-surface rounded-full">
                        <Avatar user={member} size="xs" />
                      </div>
                    {/each}
                  </div>
                {/if}
              </div>
            </div>
          </button>
        </li>
      {/each}
    </ul>
  {/if}
</article>
