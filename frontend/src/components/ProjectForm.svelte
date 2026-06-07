<script>
  import { projects } from "../stores/projectStore";
  import { createEventDispatcher, onMount } from "svelte";

  const dispatch = createEventDispatcher();

  let title = $state("");
  let description = $state("");
  let status = $state("active");
  let identifier = $state("");
  let homepage = $state("");
  let is_public = $state(false);
  let parent_id = $state("");
  let parentQuery = $state("");
  let parentDropdownOpen = $state(false);
  let error = $state("");
  let identifierError = $state("");
  let homepageError = $state("");

  onMount(() => {
    // Ensure the project list is loaded so the parent picker has options.
    if (!$projects || $projects.length === 0) {
      projects.load().catch(() => {});
    }
  });

  let filteredParents = $derived(
    ($projects || [])
      .filter((p) => {
        if (!parentQuery.trim()) return true;
        const q = parentQuery.toLowerCase();
        return (
          (p.title || "").toLowerCase().includes(q) ||
          (p.identifier || "").toLowerCase().includes(q)
        );
      })
      .slice(0, 50)
  );

  let selectedParent = $derived(
    ($projects || []).find((p) => p.id === parent_id) || null
  );

  function pickParent(p) {
    parent_id = p ? p.id : "";
    parentQuery = p ? p.title : "";
    parentDropdownOpen = false;
  }

  function clearParent() {
    parent_id = "";
    parentQuery = "";
    parentDropdownOpen = false;
  }

  // Validate identifier format (alphanumeric, underscore, hyphen only)
  function validateIdentifier() {
    if (!identifier.trim()) {
      identifierError = "شناسه الزامی است";
      return false;
    }
    const regex = /^[a-zA-Z0-9_-]+$/;
    if (!regex.test(identifier)) {
      identifierError =
        "فقط حروف، اعداد، خط تیره و زیرخط مجاز است";
      return false;
    }
    identifierError = "";
    return true;
  }

  // Validate URL format
  function validateHomepage() {
    if (homepage.trim() && homepage.trim() !== "") {
      try {
        new URL(homepage);
        homepageError = "";
        return true;
      } catch (e) {
        homepageError = "فرمت آدرس نامعتبر است";
        return false;
      }
    }
    homepageError = "";
    return true;
  }

  async function handleSubmit() {
    error = "";

    if (!title.trim()) {
      error = "عنوان الزامی است";
      return;
    }

    const isIdentifierValid = validateIdentifier();
    const isHomepageValid = validateHomepage();

    if (!isIdentifierValid || !isHomepageValid) {
      return;
    }

    try {
      await projects.create({
        title: title.trim(),
        description: description.trim(),
        status,
        identifier: identifier.trim(),
        homepage: homepage.trim() || null,
        is_public,
        parent_id: parent_id || null,
      });

      // Reset form
      title = "";
      description = "";
      status = "active";
      identifier = "";
      homepage = "";
      is_public = false;
      parent_id = "";
      parentQuery = "";
      parentDropdownOpen = false;
      error = "";
      identifierError = "";
      homepageError = "";

      dispatch("created");
    } catch (err) {
      error = err.message || "ایجاد پروژه با خطا مواجه شد";
    }
  }
</script>

<form
  onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}
  class="space-y-5"
>
  {#if error}
    <div class="p-3 bg-red-100 text-red-700 rounded-lg text-sm">
      {error}
    </div>
  {/if}

  <div>
    <label for="title" class="block text-sm font-medium text-slate-700 mb-1.5"
      >عنوان <span class="text-red-500">*</span></label
    >
    <input
      type="text"
      id="title"
      bind:value={title}
      class="w-full px-3 py-3 min-h-[44px] border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
      placeholder="پروژه من"
      required
    />
  </div>

  <div>
    <label
      for="identifier"
      class="block text-sm font-medium text-slate-700 mb-1.5"
    >
      شناسه <span class="text-red-500">*</span>
    </label>
    <input
      type="text"
      id="identifier"
      bind:value={identifier}
      onblur={validateIdentifier}
      class="w-full px-3 py-3 min-h-[44px] border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
      class:border-red-500={identifierError}
      class:ring-2={identifierError}
      class:ring-red-500={identifierError}
      placeholder="my-project"
      required
    />
    {#if identifierError}
      <p class="text-red-600 text-xs mt-1.5">{identifierError}</p>
    {:else}
      <p class="text-slate-500 text-xs mt-1.5">
        در آدرس‌ها و APIها استفاده می‌شود. فقط حروف، اعداد، خط تیره و زیرخط.
      </p>
    {/if}
  </div>

  <div>
    <label
      for="description"
      class="block text-sm font-medium text-slate-700 mb-1.5">توضیحات</label
    >
    <textarea
      id="description"
      bind:value={description}
      rows="3"
      class="w-full px-3 py-3 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent resize-none"
      placeholder="توضیحات مختصری درباره پروژه..."
    ></textarea>
  </div>

  <div class="relative">
    <label
      for="parent"
      class="block text-sm font-medium text-slate-700 mb-1.5"
    >
      پروژه والد
    </label>
    <div class="relative">
      <input
        type="text"
        id="parent"
        bind:value={parentQuery}
        onfocus={() => (parentDropdownOpen = true)}
        onblur={() => setTimeout(() => (parentDropdownOpen = false), 150)}
        class="w-full px-3 py-3 min-h-[44px] border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
        placeholder="جستجو در پروژه‌ها..."
        autocomplete="off"
      />
      {#if parent_id}
        <button
          type="button"
          onclick={clearParent}
          class="absolute inset-y-0 left-0 flex items-center px-3 text-slate-400 hover:text-slate-600"
          aria-label="حذف پروژه والد"
        >
          ✕
        </button>
      {/if}
    </div>
    {#if selectedParent && !parentDropdownOpen}
      <p class="text-slate-500 text-xs mt-1.5">
        والد انتخاب‌شده: <span class="text-slate-700">{selectedParent.title}</span>
      </p>
    {:else}
      <p class="text-slate-500 text-xs mt-1.5">
        خالی بگذارید تا پروژه به‌عنوان پروژه سطح بالا ایجاد شود.
      </p>
    {/if}
    {#if parentDropdownOpen}
      <div
        class="absolute z-10 mt-1 w-full max-h-60 overflow-auto bg-white border border-slate-200 rounded-lg shadow-lg"
      >
        <button
          type="button"
          onclick={() => pickParent(null)}
          class="w-full text-right px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
        >
          <span class="text-slate-500">(بدون والد — سطح بالا)</span>
        </button>
        {#if filteredParents.length === 0}
          <div class="px-3 py-2 text-sm text-slate-500">موردی یافت نشد</div>
        {:else}
          {#each filteredParents as p (p.id)}
            <button
              type="button"
              onclick={() => pickParent(p)}
              class="w-full text-right px-3 py-2 text-sm text-slate-700 hover:bg-indigo-50"
            >
              <span class="font-medium">{p.title}</span>
              <span class="text-slate-400 mr-2">({p.identifier})</span>
            </button>
          {/each}
        {/if}
      </div>
    {/if}
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    <div>
      <label for="status" class="block text-sm font-medium text-slate-700 mb-1.5"
        >وضعیت</label
      >
      <select
        id="status"
        bind:value={status}
        class="w-full px-3 py-3 min-h-[44px] border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
      >
        <option value="active">فعال</option>
        <option value="completed">تکمیل شده</option>
        <option value="archived">بایگانی شده</option>
      </select>
    </div>

    <div>
      <label for="is_public" class="block text-sm font-medium text-slate-700 mb-1.5"
        >دسترسی</label
      >
      <div class="flex items-center h-11">
        <input
          type="checkbox"
          id="is_public"
          bind:checked={is_public}
          class="w-5 h-5 text-indigo-600 border-slate-300 rounded focus:ring-indigo-500 cursor-pointer"
        />
        <label for="is_public" class="ml-2 text-sm text-slate-700 cursor-pointer">
          پروژه عمومی
        </label>
      </div>
    </div>
  </div>

  <div>
    <label for="homepage" class="block text-sm font-medium text-slate-700 mb-1.5"
      >آدرس صفحه اصلی</label
    >
    <input
      type="url"
      id="homepage"
      bind:value={homepage}
      onblur={validateHomepage}
      class="w-full px-3 py-3 min-h-[44px] border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
      class:border-red-500={homepageError}
      class:ring-2={homepageError}
      class:ring-red-500={homepageError}
      placeholder="https://github.com/username/project"
    />
    {#if homepageError}
      <p class="text-red-600 text-xs mt-1.5">{homepageError}</p>
    {:else}
      <p class="text-slate-500 text-xs mt-1.5">
        آدرس اختیاری به صفحه اصلی یا مخزن پروژه
      </p>
    {/if}
  </div>

  <button
    type="submit"
    disabled={!title.trim() || !identifier.trim()}
    class="w-full min-h-[44px] bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-white px-4 py-3 rounded-lg transition-colors font-medium"
  >
    ایجاد پروژه
  </button>
</form>
