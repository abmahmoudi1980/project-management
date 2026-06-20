<script>
  import Card from "./ui/Card.svelte";

  let { title, value, change, iconColor = 'text-brand-600', children } = $props();

  let changeColor = $derived(change > 0 ? 'text-success-600' : change < 0 ? 'text-danger-600' : 'text-slate-500');
  let changePrefix = $derived(change > 0 ? '+' : '');
</script>

<Card padding="lg">
  <div class="flex items-center justify-between mb-4">
    <div class="p-2 bg-slate-50 rounded-lg {iconColor}">
      {@render children?.()}
    </div>
    <span class="text-sm font-medium {changeColor}">
      {changePrefix}{change}
    </span>
  </div>
  <h3 class="text-2xl font-bold text-slate-900 mb-1">{value}</h3>
  <p class="text-sm text-slate-500">{title}</p>
</Card>
