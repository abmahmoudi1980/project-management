<script>
  import { projects } from "../stores/projectStore";
  import ProjectForm from "./ProjectForm.svelte";
  import ProjectTree from "./ProjectTree.svelte";
  import Modal from "./Modal.svelte";
  import ConfirmDialog from "./ConfirmDialog.svelte";
  import { toasts } from "../lib/toastStore.js";
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

    const projectId = projectToDelete.id;
    showDeleteModal = false;
    projectToDelete = null;

    try {
      await projects.delete(projectId);
      toasts.success('پروژه با موفقیت حذف شد');
      if (selectedProject?.id === projectId) {
        selectedProject = null;
        dispatch("select", null);
      }
      // Refresh the tree so the deleted node disappears immediately.
      loadTree();
    } catch (error) {
      toasts.error(error.message);
    }
  }

  function closeDeleteDialog() {
    showDeleteModal = false;
    projectToDelete = null;
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
      class="w-full px-4 py-3 min-h-[44px] text-sm font-medium rounded-lg transition-colors bg-brand-600 text-white hover:bg-brand-700"
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

  <ConfirmDialog
    show={showDeleteModal}
    title="حذف پروژه"
    message="آیا مطمئن هستید که می‌خواهید این پروژه را حذف کنید؟"
    confirmText="حذف"
    variant="danger"
    on:confirm={handleDelete}
    on:cancel={closeDeleteDialog}
  />

