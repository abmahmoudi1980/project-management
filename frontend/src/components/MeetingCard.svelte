<script>
  import Avatar from './Avatar.svelte';
  import { formatJalaliDate } from '../lib/utils';
  
  let { meeting } = $props();
</script>

{#if meeting}
<div class="bg-gradient-to-br from-indigo-600 to-purple-700 p-6 rounded-xl shadow-lg text-white">
  <div class="flex justify-between items-start mb-6">
    <div>
      <h3 class="text-lg font-bold mb-1">جلسه بعدی</h3>
      <p class="text-indigo-100 text-sm">{formatJalaliDate(meeting.meeting_date, 'relative')} ساعت {formatJalaliDate(meeting.meeting_date, 'time')}</p>
    </div>
    <div class="bg-white/20 p-2 rounded-lg backdrop-blur-sm">
      <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
      </svg>
    </div>
  </div>

  <h4 class="text-xl font-bold mb-2">{meeting.title}</h4>
  <p class="text-indigo-100 text-sm mb-6 line-clamp-2">{meeting.description || 'بدون توضیحات'}</p>

  <div class="flex items-center justify-between">
    <div class="flex -space-x-2">
      {#each meeting.attendees || [] as attendee}
        <div class="border-2 border-indigo-600 rounded-full">
          <Avatar user={attendee} size="sm" />
        </div>
      {/each}
      {#if meeting.total_attendees > 3}
        <div class="w-8 h-8 rounded-full bg-white/20 border-2 border-indigo-600 flex items-center justify-center text-[10px] font-medium backdrop-blur-sm">
          +{meeting.total_attendees - 3}
        </div>
      {/if}
    </div>
    <button class="bg-white text-indigo-600 px-4 py-2 rounded-lg text-sm font-bold hover:bg-indigo-50 transition-colors">
      ورود به جلسه
    </button>
  </div>
</div>
{/if}
