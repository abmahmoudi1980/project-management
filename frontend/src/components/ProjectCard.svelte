<script>
  import Avatar from './Avatar.svelte';
  import { formatJalaliDate } from '../lib/utils';
  
  let { project, onclick } = $props();

  const statusColors = {
    'Planning': 'bg-gray-100 text-gray-700',
    'In Progress': 'bg-blue-100 text-blue-700',
    'On Track': 'bg-green-100 text-green-700',
    'Review': 'bg-purple-100 text-purple-700',
    'active': 'bg-blue-100 text-blue-700'
  };

  let statusClass = $derived(statusColors[project.status] || 'bg-gray-100 text-gray-700');
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div 
  class="bg-white p-5 rounded-xl shadow-sm border border-gray-100 hover:border-indigo-300 transition-colors cursor-pointer"
  onclick={() => onclick(project.id)}
>
  <div class="flex justify-between items-start mb-4">
    <span class="px-2.5 py-0.5 rounded-full text-xs font-medium {statusClass}">
      {project.status}
    </span>
    <span class="text-xs text-gray-400">
      {formatJalaliDate(project.due_date, 'short')}
    </span>
  </div>

  <h4 class="text-lg font-bold text-gray-900 mb-1 truncate">{project.name}</h4>
  <p class="text-sm text-gray-500 mb-4 truncate">{project.client}</p>

  <div class="mb-4">
    <div class="flex justify-between text-xs mb-1">
      <span class="text-gray-500">پیشرفت</span>
      <span class="font-medium text-gray-900">{project.progress}%</span>
    </div>
    <div class="w-full bg-gray-100 rounded-full h-1.5">
      <div class="bg-indigo-600 h-1.5 rounded-full" style="width: {project.progress}%"></div>
    </div>
  </div>

  <div class="flex items-center justify-between">
    <div class="flex -space-x-2">
      {#each project.team_members || [] as member}
        <Avatar user={member} size="sm" />
      {/each}
      {#if project.total_members > 3}
        <div class="w-8 h-8 rounded-full bg-gray-50 border-2 border-white flex items-center justify-center text-[10px] font-medium text-gray-500">
          +{project.total_members - 3}
        </div>
      {/if}
    </div>
    <button class="text-indigo-600 text-sm font-medium hover:text-indigo-800">
      مشاهده
    </button>
  </div>
</div>
