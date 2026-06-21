<script>
  import Avatar from '../Avatar.svelte';
  import { formatJalaliDate } from '../../lib/utils.js';

  let { meeting, onJoin = () => {} } = $props();

  let countdown = $derived.by(() => {
    if (!meeting?.meeting_date) return null;
    const target = new Date(meeting.meeting_date);
    const now = new Date();
    const diffMs = target - now;
    if (diffMs < 0) return 'در حال برگزاری';
    const mins = Math.floor(diffMs / 60000);
    if (mins < 60) return `در ${mins} دقیقه`;
    const hrs = Math.floor(mins / 60);
    if (hrs < 24) return `در ${hrs} ساعت`;
    const days = Math.floor(hrs / 24);
    return `${days} روز دیگر`;
  });

  let time = $derived(meeting ? formatJalaliDate(meeting.meeting_date, 'time') : '');
  let date = $derived(meeting ? formatJalaliDate(meeting.meeting_date, 'short') : '');
</script>

{#if meeting}
  <article
    class="relative overflow-hidden bg-surface border border-border-subtle rounded-bento p-6 md:p-7 shadow-soft hover:shadow-lift transition-shadow"
  >
    <!-- Soft terracotta corner wash, no gradient text -->
    <div
      aria-hidden="true"
      class="pointer-events-none absolute -top-12 -right-12 w-48 h-48 rounded-full bg-brand-50/70 blur-2xl"
    ></div>

    <div class="relative flex flex-col h-full">
      <header class="flex items-start justify-between gap-4 mb-5">
        <div>
          <div class="flex items-center gap-2 text-xs font-medium text-brand-600 mb-2">
            <span class="w-1.5 h-1.5 rounded-full bg-brand-500 animate-pulse-dot"></span>
            <span>جلسه بعدی</span>
            {#if countdown}
              <span class="text-ink-subtle">·</span>
              <span class="text-ink-muted">{countdown}</span>
            {/if}
          </div>
          <h3 class="text-xl md:text-2xl font-semibold text-ink leading-snug tracking-tight">
            {meeting.title}
          </h3>
        </div>

        <div class="hidden sm:flex w-11 h-11 rounded-xl bg-brand-50 text-brand-600 items-center justify-center flex-shrink-0">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
        </div>
      </header>

      {#if meeting.description}
        <p class="text-sm text-ink-muted leading-relaxed line-clamp-2 mb-5">
          {meeting.description}
        </p>
      {/if}

      <div class="mt-auto flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 pt-4 border-t border-border-subtle">
        <div class="flex items-center gap-3">
          <div class="flex -space-x-2 rtl:space-x-reverse">
            {#each (meeting.attendees || []).slice(0, 3) as attendee}
              <div class="ring-2 ring-surface rounded-full">
                <Avatar user={attendee} size="sm" />
              </div>
            {/each}
            {#if meeting.total_attendees > 3}
              <div class="w-8 h-8 rounded-full bg-surface-muted border-2 border-surface flex items-center justify-center text-[10px] font-medium text-ink-muted">
                +{meeting.total_attendees - 3}
              </div>
            {/if}
          </div>
          <div class="text-xs text-ink-muted leading-tight">
            <div class="font-medium text-ink tabular-nums">{time}</div>
            <div>{date}</div>
          </div>
        </div>

        <button
          type="button"
          onclick={onJoin}
          class="inline-flex items-center justify-center gap-2 px-4 py-2.5 text-sm font-medium text-white bg-brand-500 hover:bg-brand-600 active:bg-brand-700 rounded-full transition-colors shadow-soft hover:shadow-lift focus:outline-none focus:ring-2 focus:ring-brand-500 focus:ring-offset-2 focus:ring-offset-canvas min-h-[44px]"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
          <span>ورود به جلسه</span>
        </button>
      </div>
    </div>
  </article>
{/if}
