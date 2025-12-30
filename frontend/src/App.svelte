<script>
  import { onMount } from "svelte";
  import { authStore } from "./stores/authStore.js";
  import { projects } from "./stores/projectStore";
  import { tasks } from "./stores/taskStore";
  import ProjectList from "./components/ProjectList.svelte";
  import TaskList from "./components/TaskList.svelte";
  import RegisterForm from "./components/RegisterForm.svelte";
  import LoginForm from "./components/LoginForm.svelte";
  import ForgotPasswordForm from "./components/ForgotPasswordForm.svelte";
  import ResetPasswordForm from "./components/ResetPasswordForm.svelte";
  import UserManagement from "./components/UserManagement.svelte";

  let selectedProject = $state(null);
  let currentRoute = $state('login'); // 'login', 'register', 'forgot-password', 'reset-password', 'app', 'users'
  let resetToken = $state('');
  let showUserManagement = $state(false);

  onMount(async () => {
    // Handle hash-based routing
    handleRoute();
    window.addEventListener('hashchange', handleRoute);
    
    await authStore.checkAuth();
    if ($authStore.isAuthenticated) {
      currentRoute = 'app';
      await projects.load();
    }
  });
  
  function handleRoute() {
    const hash = window.location.hash.slice(1); // Remove #
    
    if (hash.startsWith('/reset-password')) {
      const params = new URLSearchParams(hash.split('?')[1]);
      resetToken = params.get('token') || '';
      currentRoute = 'reset-password';
    } else if (hash === '/forgot-password') {
      currentRoute = 'forgot-password';
    } else if (hash === '/register') {
      currentRoute = 'register';
    } else if (hash === '/login') {
      currentRoute = 'login';
    } else if (hash === '/users') {
      showUserManagement = true;
      currentRoute = 'app';
    } else {
      showUserManagement = false;
    }
  }

  async function handleProjectSelect(event) {
    selectedProject = event.detail;
    if (event.detail) {
      await tasks.load(event.detail.id);
    }
  }
  
  async function handleLogout() {
    await authStore.logout();
    selectedProject = null;
    currentRoute = 'login';
    window.location.hash = '#/login';
  }
  
  function navigateTo(route) {
    currentRoute = route;
    window.location.hash = `#/${route}`;
  }
</script>

<!-- Authentication Screen -->
{#if $authStore.isLoading}
  <div class="flex items-center justify-center h-screen bg-slate-50">
    <div class="text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
      <p class="text-slate-600">در حال بارگذاری...</p>
    </div>
  </div>
{:else if !$authStore.isAuthenticated}
  <div class="min-h-screen bg-slate-50">
    {#if currentRoute === 'register'}
      <RegisterForm />
      <p class="text-center mt-4 text-sm text-gray-600" dir="rtl">
        قبلاً ثبت‌نام کرده‌اید؟
        <button
          onclick={() => navigateTo('login')}
          class="text-blue-600 hover:text-blue-700 font-medium"
        >
          ورود
        </button>
      </p>
    {:else if currentRoute === 'forgot-password'}
      <ForgotPasswordForm />
    {:else if currentRoute === 'reset-password'}
      <ResetPasswordForm token={resetToken} />
    {:else}
      <div class="flex items-center justify-center min-h-screen">
        <div class="w-full max-w-md">
          <LoginForm />
          <div class="text-center mt-4 text-sm text-gray-600" dir="rtl">
            <p class="mb-2">
              حساب کاربری ندارید؟
              <button
                onclick={() => navigateTo('register')}
                class="text-blue-600 hover:text-blue-700 font-medium"
              >
                ثبت‌نام
              </button>
            </p>
            <p>
              <button
                onclick={() => navigateTo('forgot-password')}
                class="text-blue-600 hover:text-blue-700 font-medium"
              >
                رمز عبور را فراموش کرده‌اید؟
              </button>
            </p>
          </div>
        </div>
      </div>
    {/if}
  </div>
{:else}
  <!-- Main Application -->
  <div class="flex h-screen bg-slate-50">
    <!-- Sidebar: Fixed width project list -->
    <aside class="w-64 bg-white border-r border-slate-200 flex flex-col">
      <!-- App Header -->
      <div class="px-6 py-5 border-b border-slate-200">
        <h1 class="text-xl font-semibold text-slate-900">جریان کار</h1>
        <p class="text-xs text-slate-500 mt-0.5">مدیریت پروژه و وظایف</p>
      </div>
      
      <!-- User Info & Navigation -->
      <div class="px-6 py-3 border-b border-slate-200 bg-slate-50" dir="rtl">
        <div class="flex items-center justify-between mb-3">
          <div>
            <p class="text-sm font-medium text-slate-700">{$authStore.user?.username || 'کاربر'}</p>
            <p class="text-xs text-slate-500">{$authStore.user?.role === 'admin' ? 'ادمین' : 'کاربر'}</p>
          </div>
          <button
            onclick={handleLogout}
            class="text-xs text-red-600 hover:text-red-700 font-medium"
          >
            خروج
          </button>
        </div>
        
        <!-- Admin Menu -->
        {#if $authStore.user?.role === 'admin'}
          <div class="mt-2 space-y-1">
            <button
              onclick={() => { showUserManagement = false; window.location.hash = ''; }}
              class="w-full text-right px-3 py-2 text-sm rounded {!showUserManagement ? 'bg-blue-50 text-blue-700' : 'text-slate-700 hover:bg-slate-100'}"
            >
              پروژه‌ها
            </button>
            <button
              onclick={() => { showUserManagement = true; window.location.hash = '#/users'; }}
              class="w-full text-right px-3 py-2 text-sm rounded {showUserManagement ? 'bg-blue-50 text-blue-700' : 'text-slate-700 hover:bg-slate-100'}"
            >
              مدیریت کاربران
            </button>
          </div>
        {/if}
      </div>

      <!-- Project List (only show when not in user management) -->
      {#if !showUserManagement}
        <div class="flex-1 overflow-y-auto">
          <ProjectList bind:selectedProject on:select={handleProjectSelect} />
        </div>
      {/if}
    </aside>

    <!-- Main Content Area: Fluid width -->
    <main class="flex-1 overflow-y-auto">
      {#if showUserManagement}
        <UserManagement />
      {:else if selectedProject}
        <div class="max-w-5xl mx-auto px-8 py-8">
          <!-- Page Header -->
          <div class="mb-8">
            <h2 class="text-3xl font-semibold text-slate-900">
              {selectedProject.title}
            </h2>
            {#if selectedProject.description}
              <p class="text-slate-500 mt-2">{selectedProject.description}</p>
            {/if}
          </div>

          <!-- Task List -->
          <TaskList project={selectedProject} />
        </div>
      {:else}
        <div class="flex items-center justify-center h-full">
          <div class="text-center">
            <svg
              class="w-16 h-16 mx-auto text-slate-300 mb-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="1.5"
                d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
              />
            </svg>
            <p class="text-slate-500 text-lg">یک پروژه را برای شروع انتخاب کنید</p>
            <p class="text-slate-400 text-sm mt-1">
              از نوار کناری انتخاب کنید یا پروژه جدیدی ایجاد کنید
            </p>
          </div>
        </div>
      {/if}
    </main>
  </div>
{/if}
