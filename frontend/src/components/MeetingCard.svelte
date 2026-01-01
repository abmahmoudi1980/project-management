<script>
	import Avatar from './Avatar.svelte';
	import { dateToJalaliString } from '../lib/dateUtils';

	let { meeting } = $props();

	if (!meeting) {
		// Component is hidden when meeting is null
	}

	const visibleAttendees = $derived(meeting?.attendees?.slice(0, 3) || []);
	const remainingAttendees = $derived(
		meeting ? Math.max(0, (meeting.total_attendees || 0) - visibleAttendees.length) : 0
	);

	function formatMeetingTime(dateString) {
		const date = new Date(dateString);
		return dateToJalaliString(date) + ' ' + date.toLocaleTimeString('en-US', {
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

{#if meeting}
	<div class="bg-gradient-to-br from-indigo-600 to-purple-700 p-6 rounded-2xl text-white shadow-lg">
		<h3 class="font-semibold text-lg mb-1">{meeting.title}</h3>

		{#if meeting.description}
			<p class="text-indigo-100 text-sm mb-4">{meeting.description}</p>
		{/if}

		<div class="flex items-center gap-2 mb-4 text-indigo-100">
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
				/>
			</svg>
			<span class="text-sm">{formatMeetingTime(meeting.meeting_date)}</span>
			<span class="text-sm">({meeting.duration_minutes} min)</span>
		</div>

		<div class="flex items-center justify-between">
			<div class="flex -space-x-2">
				{#each visibleAttendees as attendee}
					<div class="relative z-0">
						<Avatar user={attendee} size="sm" />
					</div>
				{/each}
				{#if remainingAttendees > 0}
					<div class="w-7 h-7 rounded-full bg-indigo-500 bg-opacity-75 flex items-center justify-center text-xs font-bold text-white">
						+{remainingAttendees}
					</div>
				{/if}
			</div>

			<button
				class="text-xs bg-white bg-opacity-20 hover:bg-opacity-30 text-white px-3 py-1 rounded-full transition-all"
			>
				Join
			</button>
		</div>
	</div>
{/if}
