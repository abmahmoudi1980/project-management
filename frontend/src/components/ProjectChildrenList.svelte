<script>
  import { api } from "../lib/api.js";
  import { createEventDispatcher } from "svelte";
  import Skeleton from "./ui/Skeleton.svelte";

  let { projectId, onSelect = null } = $props();

  const dispatch = createEventDispatcher();

  let children = $state([]);
  let count = $state(0);
  let state = $state("loading"); // 'loading' | 'loaded' | 'notfound' | 'error'
  let errorMessage = $state("");

  async function load() {
    if (!projectId) return;
    state = "loading";
    errorMessage = "";
    try {
      const result = await api.projects.getChildren(projectId);
      children = (result && Array.isArray(result.children)) ? result.children : [];
      count = (result && typeof result.count === "number") ? result.count : children.length;
      state = "loaded";
    } catch (err) {
      const msg = (err && err.message) || "خطا در بارگذاری زیرمجموعه‌ها";
      if (/not found/i.test(msg) || /not found/i.test(msg.toLowerCase()) || msg.includes("404")) {
        state = "notfound";
      } else {
        state = "error";
        errorMessage = msg;
      }
    }
  }

  $effect(() => {
    if (projectId) load();
  });

  function select(child) {
    if (onSelect) onSelect(child);
    dispatch("select", child);
  }
</script>

<section class="bg-white rounded-xl border border-slate-200 p-5">
  <header class="flex items-center justify-between mb-4">
    <h2 class="text-base font-semibold text-slate-800">
      زیرمجموعه‌ها
    </h2>
    {#if state === "loaded" && count > 0}
      <span class="text-xs text-slate-500">{count} مورد</span>
    {/if}
  </header>

  {#if state === "loading"}
    <div class="space-y-3" aria-busy="true" aria-label="در حال بارگذاری زیرمجموعه‌ها">
      {#each Array(3) as _}
        <div class="flex items-center gap-3 py-2">
          <div class="flex-1 space-y-1.5">
            <Skeleton width="w-1/2" height="h-4" />
            <Skeleton width="w-1/4" height="h-3" />
          </div>
          <Skeleton width="w-16" height="h-5" rounded="rounded-full" />
        </div>
      {/each}
    </div>
  {:else if state === "notfound"}
    <div class="text-sm text-rose-600 py-3" role="alert">پروژه یافت نشد.</div>
  {:else if state === "error"}
    <div class="text-sm text-rose-600 py-3" role="alert">{errorMessage}</div>
  {:else if state === "loaded" && count === 0}
    <div class="text-sm text-slate-500 py-3">
      هنوز زیرمجموعه‌ای ایجاد نشده است.
    </div>
  {:else if state === "loaded"}
    <ul class="divide-y divide-slate-100">
      {#each children as child (child.id)}
        <li>
          <button
            type="button"
            onclick={() => select(child)}
            class="w-full text-right py-2.5 px-2 rounded-md hover:bg-slate-50 transition-colors flex items-center justify-between gap-3"
          >
            <div class="flex-1 min-w-0">
              <div class="font-medium text-sm text-slate-800 truncate">{child.title}</div>
              {#if child.identifier}
                <div class="text-xs text-slate-500 font-mono truncate">{child.identifier}</div>
              {/if}
            </div>
            <div class="flex items-center gap-2 text-xs text-slate-500">
              <span class="px-2 py-0.5 rounded-full bg-slate-100">{child.status}</span>
              <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </div>
          </button>
        </li>
      {/each}
    </ul>
  {/if}
</section>
