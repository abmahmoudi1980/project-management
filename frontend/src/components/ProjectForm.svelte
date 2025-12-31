<script>
  import { projects } from "../stores/projectStore";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  let title = $state("");
  let description = $state("");
  let status = $state("active");
  let identifier = $state("");
  let homepage = $state("");
  let is_public = $state(false);
  let error = $state("");
  let identifierError = $state("");
  let homepageError = $state("");

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
      });

      // Reset form
      title = "";
      description = "";
      status = "active";
      identifier = "";
      homepage = "";
      is_public = false;
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
