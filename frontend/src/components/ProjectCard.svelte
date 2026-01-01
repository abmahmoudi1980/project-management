<script>
	import Avatar from './Avatar.svelte';
	import { dateToJalaliString } from '../lib/dateUtils';

	let { project, onclick } = $props();

	const statusColors = {
		Planning: 'bg-slate-100 text-slate-700',
		'In Progress': 'bg-blue-100 text-blue-700',
		'On Track': 'bg-green-100 text-green-700',
		Review: 'bg-purple-100 text-purple-700',
		Completed: 'bg-emerald-100 text-emerald-700'
	};

	const statusColor = $derived(statusColors[project.status] || 'bg-slate-100 text-slate-700');
	const visibleMembers = $derived(project.team_members?.slice(0, 3) || []);
	const remainingMembers = $derived(Math.max(0, (project.total_members || 0) - visibleMembers.length));
</script>

<button
	on:click={() => onclick?.(project.id)}
	class="bg-white p-6 rounded-2xl border border-slate-100 shadow-sm hover:shadow-lg hover:border-slate-200 transition-all duration-200 text-left w-full"
>
	<div class="flex items-start justify-between mb-3">
		<h3 class="font-semibold text-slate-900 text-lg flex-1 truncate">{project.name}</h3>
		<span class={`text-xs font-medium px-3 py-1 rounded-full ${statusColor} whitespace-nowrap ml-2`}>
			{project.status}
		</span>
	</div>

	{#if project.client}
		<p class="text-sm text-slate-600 mb-3">{project.client}</p>
	{/if}

	{#if project.progress !== undefined}
		<div class="mb-4">
			<div class="flex justify-between items-center mb-1">
				<span class="text-xs font-medium text-slate-600">Progress</span>
				<span class="text-xs font-bold text-slate-900">{project.progress}%</span>
			</div>
			<div class="w-full bg-slate-200 rounded-full h-2">
				<div class="bg-blue-600 h-2 rounded-full" style="width: {project.progress}%" />
			</div>
		</div>
	{/if}

	<div class="flex items-center justify-between">
		<div class="flex -space-x-2">
			{#each visibleMembers as member}
				<div class="relative z-0">
					<Avatar user={member} size="sm" />
				</div>
			{/each}
			{#if remainingMembers > 0}
				<div class="w-7 h-7 rounded-full bg-slate-300 flex items-center justify-center text-xs font-bold text-slate-700">
					+{remainingMembers}
				</div>
			{/if}
		</div>

		{#if project.due_date}
			<span class="text-xs text-slate-500">{dateToJalaliString(new Date(project.due_date))}</span>
		{/if}
	</div>
</button>
