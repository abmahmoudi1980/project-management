<script>
  import { onMount, onDestroy } from 'svelte';
  import { api } from '../lib/api';
  import StatCard from './StatCard.svelte';
  import ProjectCard from './ProjectCard.svelte';
  import TaskListItem from './TaskListItem.svelte';
  import MeetingCard from './MeetingCard.svelte';

  let dashboardData = $state(null);
  let loading = $state(true);
  let error = $state(null);
  let refreshInterval;

  async function loadDashboard(silent = false) {
    if (!silent) loading = true;
    try {
      dashboardData = await api.dashboard.get();
      error = null;
    } catch (err) {
      console.error('Failed to load dashboard:', err);
      error = 'خطا در بارگذاری اطلاعات داشبورد';
    } finally {
      if (!silent) loading = false;
    }
  }

  async function handleTaskComplete(taskId) {
    try {
      await api.tasks.toggleComplete(taskId);
      // Optimistically update UI
      if (dashboardData && dashboardData.user_tasks) {
        dashboardData.user_tasks = dashboardData.user_tasks.filter(t => t.id !== taskId);
        if (dashboardData.statistics.pending_tasks) {
          dashboardData.statistics.pending_tasks.current--;
          dashboardData.statistics.pending_tasks.change--;
        }
      }
    } catch (err) {
      console.error('Failed to complete task:', err);
    }
  }

  function navigateToProject(projectId) {
    window.location.hash = `#/projects/${projectId}`;
  }

  onMount(() => {
    loadDashboard();
    refreshInterval = setInterval(() => loadDashboard(true), 30000);
  });

  onDestroy(() => {
    if (refreshInterval) clearInterval(refreshInterval);
  });
</script>

<div class="p-8 max-w-7xl mx-auto" dir="rtl">
  <div class="flex justify-between items-center mb-8">
    <div>
      <h1 class="text-3xl font-bold text-gray-900">داشبورد مدیریتی</h1>
      <p class="text-gray-500">خوش آمدید! خلاصه وضعیت پروژه‌ها و وظایف شما.</p>
    </div>
    <div class="flex items-center space-x-4 rtl:space-x-reverse">
      <button 
        class="p-2 text-gray-400 hover:text-indigo-600 transition-colors"
        onclick={() => loadDashboard()}
        title="بروزرسانی"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
      </button>
    </div>
  </div>

  {#if loading && !dashboardData}
    <div class="flex justify-center items-center h-64">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>
  {:else if error}
    <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-xl mb-8">
      {error}
    </div>
  {/if}

  {#if dashboardData}
    <!-- Statistics Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <StatCard 
        title="پروژه‌های فعال" 
        value={dashboardData.statistics.active_projects?.current || 0} 
        change={dashboardData.statistics.active_projects?.change || 0}
        iconColor="text-blue-600"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
        </svg>
      </StatCard>
      <StatCard 
        title="وظایف منتظر" 
        value={dashboardData.statistics.pending_tasks?.current || 0} 
        change={dashboardData.statistics.pending_tasks?.change || 0}
        iconColor="text-orange-600"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
        </svg>
      </StatCard>
      <StatCard 
        title="اعضای تیم" 
        value={dashboardData.statistics.team_members?.current || 0} 
        change={dashboardData.statistics.team_members?.change || 0}
        iconColor="text-indigo-600"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
      </StatCard>
      <StatCard 
        title="ضرب‌الاجل‌های نزدیک" 
        value={dashboardData.statistics.upcoming_deadlines?.current || 0} 
        change={dashboardData.statistics.upcoming_deadlines?.change || 0}
        iconColor="text-red-600"
      >
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
      </StatCard>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
      <!-- Recent Projects -->
      <div class="lg:col-span-2">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-bold text-gray-900">پروژه‌های اخیر</h2>
          <a href="#/projects" class="text-indigo-600 text-sm font-medium hover:underline">مشاهده همه</a>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          {#each dashboardData.recent_projects || [] as project}
            <ProjectCard {project} onclick={navigateToProject} />
          {/each}
        </div>
      </div>

      <!-- Sidebar: Tasks & Meeting -->
      <div class="space-y-8">
        <!-- Next Meeting -->
        {#if dashboardData.next_meeting}
          <MeetingCard meeting={dashboardData.next_meeting} />
        {/if}

        <!-- User Tasks -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-100 overflow-hidden">
          <div class="p-5 border-b border-gray-100 flex justify-between items-center">
            <h3 class="font-bold text-gray-900">وظایف من</h3>
            <span class="bg-indigo-100 text-indigo-700 text-xs font-bold px-2 py-1 rounded-full">
              {dashboardData.user_tasks?.length || 0}
            </span>
          </div>
          <div class="divide-y divide-gray-100">
            {#each dashboardData.user_tasks || [] as task}
              <TaskListItem {task} onComplete={handleTaskComplete} />
            {:else}
              <div class="p-8 text-center text-gray-500 text-sm">
                همه وظایف انجام شده‌اند! 🎉
              </div>
            {/each}
          </div>
          <div class="p-4 bg-gray-50 text-center">
            <a href="#/tasks" class="text-indigo-600 text-sm font-medium hover:underline">مشاهده لیست کامل</a>
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>

