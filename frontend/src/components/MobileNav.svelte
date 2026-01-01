<script>
  import { createEventDispatcher } from "svelte";
  import { fade, fly } from "svelte/transition";
  import ProjectList from "./ProjectList.svelte";

  let { show = false, user = null, isAdmin = false, selectedProject = null } = $props();
  const dispatch = createEventDispatcher();

  function closeDrawer() {
    dispatch("close");
  }

  function handleBackdropClick(event) {
    if (event.target === event.currentTarget) {
      closeDrawer();
    }
  }

  function handleKeydown(event) {
    if (event.key === "Escape") {
      closeDrawer();
    }
  }
</script>

<svelte:window onkeydown={handleKeydown} />

{#if show}
  <!-- Backdrop -->
  <div
    class="fixed inset-0 bg-black bg-opacity-50 z-50 md:hidden"
    onclick={handleBackdropClick}
    transition:fade={{ duration: 200 }}
  />

  <!-- Drawer -->
  <div
    class="fixed top-0 right-0 h-full w-72 bg-white shadow-2xl z-50 md:hidden"
    transition:fly={{ x: 400, duration: 300, easing: (t) => t * (2 - t) }}
    dir="rtl"
  >
    <div class="flex flex-col h-full">
      <!-- Header -->
      <div class="px-6 py-5 border-b border-slate-200 bg-gradient-to-l from-indigo-600 to-purple-600">
        <h1 class="text-xl font-semibold text-white">جریان کار</h1>
        <p class="text-xs text-indigo-100 mt-0.5">مدیریت پروژه و وظایف</p>
        
        <button
          onclick={closeDrawer}
          class="absolute top-5 left-4 text-white hover:text-indigo-200 transition-colors"
          aria-label="Close menu"
        >
          <svg
            class="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>

      <!-- User Info -->
      <div class="px-6 py-4 border-b border-slate-200 bg-slate-50">
        <div class="flex items-center gap-3 mb-3">
          <div class="w-10 h-10 rounded-full bg-indigo-100 flex items-center justify-center">
            <svg
              class="w-5 h-5 text-indigo-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
              />
            </svg>
          </div>
          <div>
            <p class="text-sm font-medium text-slate-700">{user?.username || 'کاربر'}</p>
            <p class="text-xs text-slate-500">{user?.role === 'admin' ? 'ادمین' : 'کاربر'}</p>
          </div>
        </div>
      </div>

      <!-- Navigation -->
      <div class="flex-1 overflow-y-auto">
        <div class="px-4 py-3 border-b border-slate-200">
          <button
            onclick={() => {
              window.location.hash = "#/dashboard";
              closeDrawer();
            }}
            class="w-full flex items-center gap-3 px-3 py-2.5 rounded-lg text-slate-600 hover:bg-indigo-50 hover:text-indigo-600 transition-all"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
            </svg>
            <span class="text-sm font-medium">داشبورد</span>
          </button>
        </div>

        <!-- Admin Menu -->
        {#if isAdmin}
          <div class="px-4 py-3 border-b border-slate-200">
            <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2 px-2">
              مدیریت
            </p>
            <button
              onclick={() => { dispatch("navigate", "projects"); closeDrawer(); }}
              class="w-full text-right px-3 py-2.5 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-100 transition-colors flex items-center gap-3"
            >
              <svg class="w-5 h-5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"/>
              </svg>
              پروژه‌ها
            </button>
            <button
              onclick={() => { dispatch("navigate", "users"); closeDrawer(); }}
              class="w-full text-right px-3 py-2.5 rounded-lg text-sm font-medium text-slate-700 hover:bg-slate-100 transition-colors flex items-center gap-3"
            >
              <svg class="w-5 h-5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"/>
              </svg>
              مدیریت کاربران
            </button>
          </div>
        {/if}

        <!-- Projects Section -->
        <div class="px-4 py-3">
          <p class="text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2 px-2">
            پروژه‌ها
          </p>
          <ProjectList
            bind:selectedProject
            on:select={(event) => {
              dispatch("select", event.detail);
              closeDrawer();
            }}
          />
        </div>
      </div>

      <!-- Logout Button -->
      <div class="p-4 border-t border-slate-200">
        <button
          onclick={() => { dispatch("logout"); closeDrawer(); }}
          class="w-full px-4 py-3 text-sm font-medium text-red-600 hover:bg-red-50 rounded-lg transition-colors flex items-center justify-center gap-2"
        >
          <svg
            class="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
            />
          </svg>
          خروج
        </button>
      </div>
    </div>
  </div>
{/if}
