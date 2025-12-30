<script>
  import { projects } from "../stores/projectStore";
  import ProjectForm from "./ProjectForm.svelte";
  import Modal from "./Modal.svelte";
  import { createEventDispatcher } from "svelte";

  let { selectedProject = $bindable(null) } = $props();
  const dispatch = createEventDispatcher();

  let showModal = $state(false);

  function openModal() {
    showModal = true;
  }

  function closeModal() {
    showModal = false;
  }

  async function handleProjectSelect(project) {
    selectedProject = project;
    dispatch("select", project);
  }

  async function handleProjectDelete(projectId) {
    if (confirm("آیا مطمئن هستید که می‌خواهید این پروژه را حذف کنید؟")) {
      await projects.delete(projectId);
      if (selectedProject?.id === projectId) {
        selectedProject = null;
        dispatch("select", null);
      }
    }
  }
</script>

<div class="flex flex-col h-full relative">
  <!-- Projects Header -->
  <div class="px-6 py-4">
    <h2 class="text-xs font-semibold text-slate-500 uppercase tracking-wider">
      پروژه‌ها
    </h2>
  </div>

  <!-- Project List -->
  <nav class="flex-1 px-3 space-y-1">
    {#each $projects || [] as project}
      <div class="group relative">
        <button
          onclick={() => handleProjectSelect(project)}
          class="w-full text-left px-3 py-2.5 rounded-lg transition-all relative
            {selectedProject?.id === project.id
            ? 'bg-indigo-50 text-indigo-700 border-l-4 border-indigo-600 pl-2.5'
            : 'text-slate-700 hover:bg-slate-50 border-l-4 border-transparent'}"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1 min-w-0 pr-8">
              <div class="flex items-center gap-2 mb-1">
                <h3 class="font-medium text-sm truncate">{project.title}</h3>
                {#if project.is_public}
                  <span
                    class="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800"
                  >
                    عمومی
                  </span>
                {/if}
              </div>
              {#if project.identifier}
                <p class="text-xs text-slate-500 font-mono truncate">
                  {project.identifier}
                </p>
              {/if}
              <div class="flex items-center gap-2 mt-0.5">
                <p
                  class="text-xs {selectedProject?.id === project.id
                    ? 'text-indigo-600'
                    : 'text-slate-500'}"
                >
                  {project.status}
                </p>
                {#if project.homepage}
                  <a
                    href={project.homepage}
                    target="_blank"
                    rel="noopener noreferrer"
                    onclick={(e) => e.stopPropagation()}
                    class="text-xs text-blue-500 hover:text-blue-700 flex items-center"
                    title="Visit homepage"
                  >
                    <svg
                      class="w-3 h-3"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                      />
                    </svg>
                  </a>
                {/if}
              </div>
            </div>
          </div>
        </button>
        <button
          onclick={(e) => { e.stopPropagation(); handleProjectDelete(project.id); }}
          class="absolute left-2 top-2 opacity-0 group-hover:opacity-100 p-1 hover:bg-slate-200 rounded transition-opacity"
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
    {/each}
  </nav>

  <!-- New Project Button (Fixed at bottom) -->
  <div class="p-4 border-t border-slate-200">
    <button
      onclick={openModal}
      class="w-full px-4 py-2.5 text-sm font-medium rounded-lg transition-colors bg-indigo-600 text-white hover:bg-indigo-700"
    >
      + پروژه جدید
    </button>
  </div>
</div>

<!-- Modal for New Project -->
<Modal show={showModal} title="ایجاد پروژه جدید" maxWidth="lg" on:close={closeModal}>
  <ProjectForm
    on:created={closeModal}
  />
</Modal>
