<script>
  let {
    title,
    value,
    change = 0,
    icon,
    tone = 'default',
    href = undefined,
  } = $props();

  let trend = $derived.by(() => {
    if (change > 0) return { dir: 'up', cls: 'text-success-600 bg-success-50', label: `+${change}` };
    if (change < 0) return { dir: 'down', cls: 'text-danger-600 bg-danger-50', label: `${change}` };
    return { dir: 'flat', cls: 'text-ink-muted bg-surface-muted', label: '۰' };
  });

  let iconWrap = $derived.by(() => {
    switch (tone) {
      case 'accent':  return 'bg-brand-50 text-brand-600';
      case 'warning': return 'bg-warning-50 text-warning-600';
      case 'success': return 'bg-success-50 text-success-600';
      case 'info':    return 'bg-info-50 text-info-600';
      default:        return 'bg-surface-muted text-ink-muted';
    }
  });

  let Wrapper = $derived(href ? 'a' : 'div');
</script>

<svelte:element
  this={Wrapper}
  href={href}
  class="group block bg-surface border border-border-subtle rounded-bento p-5 shadow-soft hover:shadow-lift hover:border-border transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-brand-500 focus:ring-offset-2 focus:ring-offset-canvas"
>
  <div class="flex items-start justify-between mb-4">
    <div class="w-10 h-10 rounded-xl {iconWrap} flex items-center justify-center transition-colors">
      {#if icon}{@render icon()}{/if}
    </div>
    <span class="inline-flex items-center gap-1 text-[11px] font-medium px-2 py-0.5 rounded-full {trend.cls}">
      {#if trend.dir === 'up'}
        <svg class="w-3 h-3" viewBox="0 0 12 12" fill="none" aria-hidden="true">
          <path d="M3 8l3-3 3 3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      {:else if trend.dir === 'down'}
        <svg class="w-3 h-3" viewBox="0 0 12 12" fill="none" aria-hidden="true">
          <path d="M3 4l3 3 3-3" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
      {/if}
      {trend.label}
    </span>
  </div>

  <div class="flex items-baseline gap-1.5">
    <span class="text-3xl font-semibold text-ink tracking-tight tabular-nums">{value}</span>
  </div>
  <p class="text-sm text-ink-muted mt-1">{title}</p>
</svelte:element>
