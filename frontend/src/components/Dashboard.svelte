<script>
	import { onMount } from 'svelte';
	import StatCard from './StatCard.svelte';
	import ProjectCard from './ProjectCard.svelte';
	import TaskListItem from './TaskListItem.svelte';
	import MeetingCard from './MeetingCard.svelte';
	import { api } from '../lib/api.js';

	let dashboardData = $state(null);
	let loading = $state(true);
	let error = $state(null);
	let autoRefreshInterval = $state(null);

	async function loadDashboard() {
		try {
			loading = true;
			error = null;
			const data = await api.dashboard.get();
			dashboardData = data;
		} catch (err) {
			console.error('Failed to load dashboard:', err);
			error = err.message || 'Failed to load dashboard';
		} finally {
			loading = false;
		}
	}

	function handleTaskComplete(taskId) {
		// Mark task as complete and remove from list after animation
		if (dashboardData) {
			dashboardData.user_tasks = dashboardData.user_tasks.filter((t) => t.id !== taskId);
			// Update pending_tasks count
			if (dashboardData.statistics) {
				dashboardData.statistics.pending_tasks.current = Math.max(
					0,
					dashboardData.statistics.pending_tasks.current - 1
				);
			}
		}
	}

	function handleProjectClick(projectId) {
		// Navigate to project details
		window.location.href = `/projects/${projectId}`;
	}

	onMount(() => {
		// Load dashboard on mount
		loadDashboard();

		// Set up auto-refresh every 30 seconds
		autoRefreshInterval = setInterval(() => {
			loadDashboard();
		}, 30000);

		// Cleanup on unmount
		return () => {
			if (autoRefreshInterval) {
				clearInterval(autoRefreshInterval);
			}
		};
	});
</script>

<div class="min-h-screen bg-slate-50 p-4 md:p-6">
	<div class="max-w-7xl mx-auto">
		{#if error}
			<div
				class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg text-red-800 flex justify-between items-center"
			>
				<span>{error}</span>
				<button
					on:click={loadDashboard}
					class="text-red-700 hover:text-red-900 font-medium underline"
				>
					Retry
				</button>
			</div>
		{/if}

		<!-- Header -->
		<div class="mb-8">
			<h1 class="text-3xl font-bold text-slate-900">Dashboard</h1>
			<p class="text-slate-600 mt-1">Welcome back! Here's your project overview.</p>
		</div>

		{#if loading && !dashboardData}
			<div class="flex items-center justify-center h-64">
				<div class="text-slate-500">Loading dashboard...</div>
			</div>
		{:else if dashboardData}
			<!-- Statistics Grid -->
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
				<StatCard
					title="Active Projects"
					value={dashboardData.statistics.active_projects.current}
					change={dashboardData.statistics.active_projects.change}
					icon="folder"
					iconColor="blue"
				/>
				<StatCard
					title="Pending Tasks"
					value={dashboardData.statistics.pending_tasks.current}
					change={dashboardData.statistics.pending_tasks.change}
					icon="check-circle"
					iconColor="purple"
				/>
				<StatCard
					title="Team Members"
					value={dashboardData.statistics.team_members.current}
					change={dashboardData.statistics.team_members.change}
					icon="users"
					iconColor="green"
				/>
				<StatCard
					title="Upcoming Deadlines"
					value={dashboardData.statistics.upcoming_deadlines.current}
					change={dashboardData.statistics.upcoming_deadlines.change}
					icon="calendar"
					iconColor="orange"
				/>
			</div>

			<!-- Main Content Grid -->
			<div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
				<!-- Projects and Tasks (2 columns) -->
				<div class="lg:col-span-2 space-y-6">
					<!-- Recent Projects -->
					{#if dashboardData.recent_projects.length > 0}
						<div>
							<h2 class="text-lg font-semibold text-slate-900 mb-4">Recent Projects</h2>
							<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
								{#each dashboardData.recent_projects as project}
									<ProjectCard {project} onclick={handleProjectClick} />
								{/each}
							</div>
						</div>
					{/if}

					<!-- User Tasks -->
					{#if dashboardData.user_tasks.length > 0}
						<div>
							<h2 class="text-lg font-semibold text-slate-900 mb-4">Your Tasks</h2>
							<div class="space-y-2">
								{#each dashboardData.user_tasks as task}
									<TaskListItem {task} onComplete={handleTaskComplete} />
								{/each}
							</div>
						</div>
					{:else}
						<div class="bg-white p-8 rounded-2xl border border-slate-100 text-center">
							<p class="text-slate-600">You're all caught up! 🎉</p>
						</div>
					{/if}
				</div>

				<!-- Sidebar (1 column) -->
				<div class="space-y-6">
					<!-- Next Meeting -->
					{#if dashboardData.next_meeting}
						<div>
							<h2 class="text-lg font-semibold text-slate-900 mb-4">Upcoming Meeting</h2>
							<MeetingCard meeting={dashboardData.next_meeting} />
						</div>
					{/if}
				</div>
			</div>
		{/if}
	</div>
</div>

<style>
	:global(body) {
		background-color: #f8fafc;
	}
</style>
