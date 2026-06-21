<script>
  import { onMount, onDestroy } from 'svelte';
  import { fly } from 'svelte/transition';
  import { cubicOut } from 'svelte/easing';
  import { api } from '../lib/api';
  import { authStore } from '../stores/authStore.js';
  import DashboardHeader from './dashboard/DashboardHeader.svelte';
  import StatTile from './dashboard/StatTile.svelte';
  import NextMeetingTile from './dashboard/NextMeetingTile.svelte';
  import ActiveProjectsTile from './dashboard/ActiveProjectsTile.svelte';
  import MyTasksTile from './dashboard/MyTasksTile.svelte';
  import TodaysFocusTile from './dashboard/TodaysFocusTile.svelte';
  import DashboardSkeleton from './dashboard/DashboardSkeleton.svelte';

  let dashboardData = $state(null);
  let loading = $state(true);
  let error = $state(null);
  let refreshInterval;

  let reducedMotion = $state(false);

  onMount(() => {
    if (typeof window !== 'undefined') {
      reducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    }
    loadDashboard();
    refreshInterval = setInterval(() => loadDashboard(true), 30000);
  });

  onDestroy(() => {
    if (refreshInterval) clearInterval(refreshInterval);
  });

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
      if (dashboardData && dashboardData.user_tasks) {
        dashboardData.user_tasks = dashboardData.user_tasks.filter(t => t.id !== taskId);
        if (dashboardData.statistics.pending_tasks) {
          dashboardData.statistics.pending_tasks.current = Math.max(0, dashboardData.statistics.pending_tasks.current - 1);
          dashboardData.statistics.pending_tasks.change -= 1;
        }
      }
    } catch (err) {
      console.error('Failed to complete task:', err);
    }
  }

  function navigateToProject(projectId) {
    window.location.hash = `#/projects/${projectId}`;
  }

  function tileTransition(index = 0) {
    if (reducedMotion) return { duration: 0 };
    return { y: 12, duration: 420, delay: index * 50, easing: cubicOut };
  }
</script>

<div class="min-h-[100dvh] bg-canvas" dir="rtl">
  <div class="max-w-[1400px] mx-auto px-5 md:px-8 py-8 md:py-10">
    {#if loading && !dashboardData}
      <DashboardSkeleton />
    {:else if error}
      <div
        role="alert"
        class="bg-danger-50 border border-danger-100 text-danger-700 px-4 py-3 rounded-2xl mb-6 text-sm"
      >
        {error}
      </div>
    {/if}

    {#if dashboardData}
      <div in:fly={tileTransition(0)}>
        <DashboardHeader
          user={$authStore.user}
          onRefresh={() => loadDashboard()}
          isRefreshing={loading}
        />
      </div>

      <!-- Row 1: Stats -->
      <section
        aria-label="آمار کلی"
        class="grid grid-cols-2 lg:grid-cols-4 gap-4 md:gap-5 mb-5"
      >
        <div in:fly={tileTransition(1)}>
          <StatTile
            title="پروژه‌های فعال"
            value={dashboardData.statistics.active_projects?.current || 0}
            change={dashboardData.statistics.active_projects?.change || 0}
            tone="accent"
          >
            {#snippet icon()}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7h18M3 12h18M3 17h18" />
              </svg>
            {/snippet}
          </StatTile>
        </div>

        <div in:fly={tileTransition(2)}>
          <StatTile
            title="وظایف منتظر"
            value={dashboardData.statistics.pending_tasks?.current || 0}
            change={dashboardData.statistics.pending_tasks?.change || 0}
            tone="accent"
          >
            {#snippet icon()}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
              </svg>
            {/snippet}
          </StatTile>
        </div>

        <div in:fly={tileTransition(3)}>
          <StatTile
            title="اعضای تیم"
            value={dashboardData.statistics.team_members?.current || 0}
            change={dashboardData.statistics.team_members?.change || 0}
            tone="info"
          >
            {#snippet icon()}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
              </svg>
            {/snippet}
          </StatTile>
        </div>

        <div in:fly={tileTransition(4)}>
          <StatTile
            title="ضرب‌الاجل‌های نزدیک"
            value={dashboardData.statistics.upcoming_deadlines?.current || 0}
            change={dashboardData.statistics.upcoming_deadlines?.change || 0}
            tone="warning"
          >
            {#snippet icon()}
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
            {/snippet}
          </StatTile>
        </div>
      </section>

      <!-- Row 2: Next meeting (8) + Active projects (4) -->
      <section
        aria-label="جلسه و پروژه‌ها"
        class="grid grid-cols-12 gap-4 md:gap-5 mb-5"
      >
        <div class="col-span-12 lg:col-span-8" in:fly={tileTransition(5)}>
          <NextMeetingTile meeting={dashboardData.next_meeting} />
        </div>
        <div class="col-span-12 lg:col-span-4" in:fly={tileTransition(6)}>
          <ActiveProjectsTile
            projects={dashboardData.recent_projects || []}
            onSelect={navigateToProject}
          />
        </div>
      </section>

      <!-- Row 3: My tasks (8) + Today's focus (4) -->
      <section
        aria-label="وظایف و تمرکز"
        class="grid grid-cols-12 gap-4 md:gap-5"
      >
        <div class="col-span-12 lg:col-span-8" in:fly={tileTransition(7)}>
          <MyTasksTile
            tasks={dashboardData.user_tasks || []}
            onComplete={handleTaskComplete}
            onSelect={navigateToProject}
          />
        </div>
        <div class="col-span-12 lg:col-span-4" in:fly={tileTransition(8)}>
          <TodaysFocusTile tasks={dashboardData.user_tasks || []} />
        </div>
      </section>
    {/if}
  </div>
</div>
