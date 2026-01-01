<script>
	let { title, value, change, icon, iconColor = 'blue' } = $props();

	const changeColor = $derived(
		change > 0
			? 'text-green-600 bg-green-50'
			: change < 0
				? 'text-red-600 bg-red-50'
				: 'text-slate-600 bg-slate-100'
	);
	const changePrefix = $derived(change > 0 ? '+' : '');
</script>

<div
	class="bg-white p-6 rounded-2xl border border-slate-100 shadow-sm hover:shadow-md transition-shadow duration-200"
>
	<div class="flex items-start justify-between mb-4">
		<div class="p-3 rounded-xl bg-{iconColor}-50">
			<svg
				class="w-6 h-6 text-{iconColor}-600"
				fill="none"
				stroke="currentColor"
				viewBox="0 0 24 24"
			>
				<!-- Icon placeholder - will be replaced by lucide-svelte in actual implementation -->
				<circle cx="12" cy="12" r="10" />
			</svg>
		</div>
		<span class="text-xs font-medium px-2 py-1 rounded-full {changeColor}">
			{changePrefix}{change}
		</span>
	</div>
	<div>
		<h3 class="text-2xl font-bold text-slate-800">{value}</h3>
		<p class="text-sm text-slate-500 font-medium">{title}</p>
	</div>
</div>
