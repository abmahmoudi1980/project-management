<script>
  let { title, value, change, iconColor = 'text-indigo-600', children } = $props();

  let changeColor = $derived(change > 0 ? 'text-green-600' : change < 0 ? 'text-red-600' : 'text-gray-500');
  let changePrefix = $derived(change > 0 ? '+' : '');
</script>

<div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
  <div class="flex items-center justify-between mb-4">
    <div class="p-2 bg-gray-50 rounded-lg {iconColor}">
      {@render children?.()}
    </div>
    <span class="text-sm font-medium {changeColor}">
      {changePrefix}{change}
    </span>
  </div>
  <h3 class="text-2xl font-bold text-gray-900 mb-1">{value}</h3>
  <p class="text-sm text-gray-500">{title}</p>
</div>
