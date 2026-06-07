<script>
  import { projects } from "../stores/projectStore";
  import ProjectForm from "./ProjectForm.svelte";
  import ProjectTree from "./ProjectTree.svelte";
  import Modal from "./Modal.svelte";
  import { createEventDispatcher, onMount } from "svelte";

  let { selectedProject = $bindable(null) } = $props();
  const dispatch = createEventDispatcher();

  let showModal = $state(false);
  let showDeleteModal = $state(false);
  let projectToDelete = $state(null);
  let treeProjects = $state([]);

  async function loadTree() {
    try {
      // Prefer the dedicated tree endpoint; fall back to the flat list if the
      // backend is older (no /tree route) so the UI never breaks.
      const tree = await projects.tree();
      treeProjects = (tree && Array.isArray(tree.nodes)) ? tree.nodes.map((n) => n) : [];
    } catch (err) {
      console.warn('[ProjectList] /tree unavailable, falling back to /projects', err);
      try {
        await projects.load();
        treeProjects = $projects || [];
      } catch (e) {
        console.error('[ProjectList] failed to load projects', e);
        treeProjects = [];
      }
    }
  }

  onMount(() => {
    loadTree();
  });

  function openModal() {
    showModal = true;
  }

  function closeModal() {
    showModal = false;
    // Reload the tree so the newly created project appears immediately.
    loadTree();
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
      // Refresh the tree so the deleted node disappears immediately.
      loadTree();
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

  <!-- Project List (Tree) -->
  <nav class="flex-1 px-2 sm:px-3 py-1 overflow-y-auto">
    <ProjectTree
      projects={treeProjects}
      selectedId={selectedProject?.id}
      onDelete={confirmDelete}
      on:select={(e) => handleProjectSelect(e.detail)}
    />
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

