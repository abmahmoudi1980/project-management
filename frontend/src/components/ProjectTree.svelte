<script>
  import { createEventDispatcher, onMount } from "svelte";

  let { projects = [], selectedId = null, onDelete = null } = $props();

  const dispatch = createEventDispatcher();

  let expanded = $state(new Set());

  onMount(() => {
    // Expand the first level (root projects) by default per the spec.
    const roots = projects.filter((p) => !p.parent_id).map((p) => p.id);
    expanded = new Set(roots);
  });

  // Build a lookup: parent_id -> children, plus a quick id->project map.
  let tree = $derived(buildTree(projects));

  function buildTree(list) {
    const byId = new Map();
    const childrenOf = new Map();
    for (const p of list) {
      byId.set(p.id, p);
      const key = p.parent_id || "__root__";
      if (!childrenOf.has(key)) childrenOf.set(key, []);
      childrenOf.get(key).push(p);
    }
    return { byId, childrenOf };
  }

  function toggle(id) {
    const next = new Set(expanded);
    if (next.has(id)) next.delete(id);
    else next.add(id);
    expanded = next;
  }

  function hasChildren(id) {
    return (tree.childrenOf.get(id) || []).length > 0;
  }

  function select(project) {
    dispatch("select", project);
  }

  function handleDelete(project, e) {
    e.stopPropagation();
    if (onDelete) onDelete(project);
  }
</script>

{#snippet renderNode(project, depth)}
  {@const children = tree.childrenOf.get(project.id) || []}
  {@const isExpanded = expanded.has(project.id)}
  {@const isSelected = selectedId === project.id}
  <div class="group relative">
    <div
      class="w-full text-right sm:text-left flex items-start gap-1 py-2.5 sm:py-2 min-h-[44px] rounded-lg transition-all cursor-pointer border-r-4 sm:border-l-4 sm:border-r-0
        {isSelected
        ? 'bg-indigo-50 text-indigo-700 border-indigo-600'
        : 'text-slate-700 hover:bg-slate-50 border-transparent'}"
      style="padding-inline-start: {0.5 + depth * 0.75}rem; padding-inline-end: 2.5rem;"
      onclick={() => select(project)}
      onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); select(project); } }}
      role="button"
      tabindex="0"
    >
      {#if children.length > 0}
        <button
          type="button"
          onclick={(e) => { e.stopPropagation(); toggle(project.id); }}
          aria-expanded={isExpanded}
          aria-label={isExpanded ? 'بستن زیرمجموعه' : 'باز کردن زیرمجموعه'}
          class="flex-shrink-0 w-5 h-5 flex items-center justify-center text-slate-400 hover:text-slate-700"
        >
          <svg
            class="w-3 h-3 transition-transform {isExpanded ? 'rotate-90' : ''}"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path d="M6 4l8 6-8 6V4z" />
          </svg>
        </button>
      {:else}
        <span class="flex-shrink-0 w-5 h-5" aria-hidden="true"></span>
      {/if}
      <div class="flex-1 min-w-0">
        <div class="flex items-center gap-1.5 mb-0.5 flex-wrap">
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
        <p class="text-xs {isSelected ? 'text-indigo-600' : 'text-slate-500'}">
          {project.status}
        </p>
      </div>
    </div>
    {#if onDelete}
      <button
        type="button"
        onclick={(e) => handleDelete(project, e)}
        class="absolute end-2 top-2 opacity-100 sm:opacity-0 sm:group-hover:opacity-100 p-2 sm:p-1 hover:bg-rose-50 rounded transition-opacity min-w-[36px] min-h-[36px] sm:min-w-0 sm:min-h-0"
        title="حذف پروژه"
        aria-label="حذف پروژه"
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
    {/if}
  </div>
  {#if isExpanded && children.length > 0}
    <div role="group">
      {#each children as child (child.id)}
        {@render renderNode(child, depth + 1)}
      {/each}
    </div>
  {/if}
{/snippet}

{#if projects.length === 0}
  <div class="px-4 py-8 text-center text-sm text-slate-500">
    هیچ پروژه‌ای یافت نشد
  </div>
{:else}
  <div class="space-y-0.5">
    {#each (tree.childrenOf.get('__root__') || []) as root (root.id)}
      {@render renderNode(root, 0)}
    {/each}
  </div>
{/if}
