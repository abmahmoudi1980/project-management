<script>
  import { projects } from "../stores/projectStore";
  import ProjectForm from "./ProjectForm.svelte";
  import Modal from "./Modal.svelte";
  import { createEventDispatcher } from "svelte";

  let { selectedProject = $bindable(null) } = $props();
  const dispatch = createEventDispatcher();

  let showModal = $state(false);
  let showDeleteModal = $state(false);
  let projectToDelete = $state(null);

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

  function confirmDelete(project) {
    showDeleteModal = true;
    projectToDelete = project;
  }

  async function handleDelete() {
    if (!projectToDelete) return;

    try {
      const projectId = projectToDelete.id;
      await projects.delete(projectId);
      showDeleteModal = false;
      projectToDelete = null;
      if (selectedProject?.id === projectId) {
        selectedProject = null;
        dispatch("select", null);
      }
    } catch (error) {
      alert(error.message);
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
  <nav class="flex-1 px-2 sm:px-3 space-y-1">
    {#each $projects || [] as project}
      <div class="group relative">
        <button
          onclick={() => handleProjectSelect(project)}
          class="w-full text-right sm:text-left px-3 py-3 sm:py-2.5 min-h-[56px] sm:min-h-0 rounded-lg transition-all relative
            {selectedProject?.id === project.id
            ? 'bg-indigo-50 text-indigo-700 border-r-4 sm:border-l-4 sm:border-r-0 border-indigo-600 pr-2 sm:pr-8 sm:pl-2.5'
            : 'text-slate-700 hover:bg-slate-50 border-r-4 sm:border-l-4 sm:border-r-0 border-transparent'}"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1 min-w-0 pr-8 sm:pr-0 pl-2 sm:pl-0">
              <div class="flex items-center gap-1.5 sm:gap-2 mb-1 flex-wrap">
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
              <div class="flex items-center gap-1.5 sm:gap-2 mt-0.5 flex-wrap">
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
                    class="text-xs text-blue-500 hover:text-blue-700 flex items-center p-1 hover:bg-blue-50 rounded"
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
          onclick={(e) => { e.stopPropagation(); confirmDelete(project); }}
          class="absolute left-2 sm:right-2 sm:left-auto top-2 opacity-100 sm:opacity-0 sm:group-hover:opacity-100 p-2 sm:p-1 hover:bg-rose-50 rounded transition-opacity min-w-[36px] min-h-[36px] sm:min-w-0 sm:min-h-0"
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
  <div class="p-3 sm:p-4 border-t border-slate-200">
    <button
      onclick={openModal}
      class="w-full px-4 py-3 min-h-[44px] text-sm font-medium rounded-lg transition-colors bg-indigo-600 text-white hover:bg-indigo-700"
    >
      + پروژه جدید
    </button>
  </div>
</div>

<!-- Modal for New Project -->
<Modal show={showModal} title="ایجاد پروژه جدید" maxWidth="lg" on:close={closeModal}>
  {#snippet children()}
    <ProjectForm
      on:created={closeModal}
    />
  {/snippet}
</Modal>

  <Modal show={showDeleteModal} fullScreen={false} on:close={() => { showDeleteModal = false; projectToDelete = null; }}>
    <div class="p-4 sm:p-6">
      <h3 class="text-lg font-semibold text-slate-900 mb-2">
        حذف پروژه
      </h3>
      <p class="text-slate-600 mb-4">
        آیا مطمئن هستید که می‌خواهید این پروژه را حذف کنید؟
      </p>
      <div class="flex flex-col sm:flex-row gap-3 justify-end sm:justify-end">
        <button
          onclick={() => { showDeleteModal = false; projectToDelete = null; }}
          class="w-full sm:w-auto px-4 py-3 min-h-[44px] sm:min-h-0 bg-slate-200 text-slate-700 rounded-lg hover:bg-slate-300 font-medium"
        >
          لغو
        </button>
        <button
          onclick={handleDelete}
          class="w-full sm:w-auto px-4 py-3 min-h-[44px] sm:min-h-0 bg-rose-600 text-white rounded-lg hover:bg-rose-700 font-medium"
        >
          حذف
        </button>
      </div>
    </div>
  </Modal>

