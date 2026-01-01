<script>
	import { dateToJalaliString } from '../lib/dateUtils';

	let { task, onComplete } = $props();

	const priorityColors = {
		4: 'bg-red-100 text-red-700',
		3: 'bg-orange-100 text-orange-700',
		2: 'bg-blue-100 text-blue-700',
		1: 'bg-slate-100 text-slate-700'
	};

	const priorityLabel = {
		4: 'Critical',
		3: 'High',
		2: 'Medium',
		1: 'Low'
	};

	const priorityColor = $derived(priorityColors[task.priority] || 'bg-slate-100 text-slate-700');
	const label = $derived(priorityLabel[task.priority] || 'Medium');
	const isCompleted = $derived(task.status === 'done');
</script>

<div class="flex items-start gap-3 p-3 bg-white rounded-lg border border-slate-100 hover:border-slate-200 transition-colors">
	<input
		type="checkbox"
		checked={isCompleted}
		on:change={() => onComplete?.(task.id)}
		class="mt-1 w-4 h-4 rounded border-slate-300 text-blue-600 focus:ring-2 focus:ring-blue-500"
		aria-label="Mark task complete"
	/>

	<div class="flex-1 min-w-0">
		<div class="flex items-center gap-2 mb-1">
			<h4
				class={`font-medium text-sm ${
					isCompleted ? 'line-through text-slate-400' : 'text-slate-900'
				}`}
			>
				{task.title}
			</h4>
			<span class={`text-xs font-medium px-2 py-0.5 rounded-full whitespace-nowrap ${priorityColor}`}>
				{label}
			</span>
		</div>
		<p class="text-xs text-slate-500 truncate">{task.project_name}</p>
		{#if task.due_date}
			<p class="text-xs text-slate-500 mt-1">{dateToJalaliString(new Date(task.due_date))}</p>
		{/if}
	</div>
</div>
