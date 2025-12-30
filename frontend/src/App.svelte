<script>
  import { onMount } from "svelte";
  import { projects } from "./stores/projectStore";
  import { tasks } from "./stores/taskStore";
  import ProjectList from "./components/ProjectList.svelte";
  import TaskList from "./components/TaskList.svelte";

  let selectedProject = $state(null);

  onMount(async () => {
    await projects.load();
  });

  async function handleProjectSelect(event) {
    selectedProject = event.detail;
    if (event.detail) {
      await tasks.load(event.detail.id);
    }
  }
</script>

<div class="flex h-screen bg-slate-50">
  <!-- Sidebar: Fixed width project list -->
  <aside class="w-64 bg-white border-r border-slate-200 flex flex-col">
    <!-- App Header -->
    <div class="px-6 py-5 border-b border-slate-200">
      <h1 class="text-xl font-semibold text-slate-900">جریان کار</h1>
      <p class="text-xs text-slate-500 mt-0.5">مدیریت پروژه و وظایف</p>
    </div>

    <!-- Project List -->
    <div class="flex-1 overflow-y-auto">
      <ProjectList bind:selectedProject on:select={handleProjectSelect} />
    </div>
  </aside>

  <!-- Main Content Area: Fluid width -->
  <main class="flex-1 overflow-y-auto">
    {#if selectedProject}
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
