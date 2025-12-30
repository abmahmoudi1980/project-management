<script>
  import { projects } from "../stores/projectStore";
  import ProjectForm from "./ProjectForm.svelte";
  import { createEventDispatcher } from "svelte";

  export let selectedProject = null;
  const dispatch = createEventDispatcher();

  $: showForm = false;

  function toggleForm() {
    showForm = !showForm;
  }

  async function handleProjectSelect(project) {
    selectedProject = project;
    dispatch("select", project);
  }

  async function handleProjectDelete(projectId) {
    if (confirm("Are you sure you want to delete this project?")) {
      await projects.delete(projectId);
      if (selectedProject?.id === projectId) {
        selectedProject = null;
        dispatch("select", null);
      }
    }
  }
</script>

<div class="flex flex-col h-full">
  <!-- Projects Header -->
  <div class="px-6 py-4">
    <h2 class="text-xs font-semibold text-slate-500 uppercase tracking-wider">
      Projects
    </h2>
  </div>

  <!-- Project List -->
  <nav class="flex-1 px-3 space-y-1">
    {#each $projects as project}
      <button
        on:click={() => handleProjectSelect(project)}
        class="group w-full text-left px-3 py-2.5 rounded-lg transition-all relative
          {selectedProject?.id === project.id
          ? 'bg-indigo-50 text-indigo-700 border-l-4 border-indigo-600 pl-2.5'
          : 'text-slate-700 hover:bg-slate-50 border-l-4 border-transparent'}"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1 min-w-0">
            <h3 class="font-medium text-sm truncate">{project.title}</h3>
            <p
              class="text-xs mt-0.5 {selectedProject?.id === project.id
                ? 'text-indigo-600'
                : 'text-slate-500'}"
            >
              {project.status}
            </p>
          </div>
          <button
            on:click|stopPropagation={() => handleProjectDelete(project.id)}
            class="opacity-0 group-hover:opacity-100 ml-2 p-1 hover:bg-slate-200 rounded transition-opacity"
            title="Delete project"
          >
            <svg
              class="w-4 h-4 text-slate-400 hover:text-rose-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
              />
            </svg>
          </button>
        </div>
      </button>
    {/each}
  </nav>

  <!-- New Project Button (Fixed at bottom) -->
  <div class="p-4 border-t border-slate-200">
    <button
      on:click={toggleForm}
      class="w-full px-4 py-2.5 text-sm font-medium rounded-lg transition-colors
        {showForm
        ? 'bg-slate-100 text-slate-700 hover:bg-slate-200'
        : 'bg-indigo-600 text-white hover:bg-indigo-700'}"
    >
      {showForm ? "Cancel" : "+ New Project"}
    </button>
  </div>

  {#if showForm}
    <div class="absolute inset-0 bg-white z-10 overflow-y-auto">
      <div class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-slate-900">New Project</h3>
          <button
            on:click={toggleForm}
            class="text-slate-400 hover:text-slate-600"
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
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>
        <ProjectForm
          on:created={() => {
            showForm = false;
          }}
        />
      </div>
    </div>
  {/if}
</div>
