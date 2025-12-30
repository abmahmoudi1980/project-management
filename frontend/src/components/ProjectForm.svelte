<script>
  import { projects } from "../stores/projectStore";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();

  let title = "";
  let description = "";
  let status = "active";
  let identifier = "";
  let homepage = "";
  let is_public = false;
  let error = "";
  let identifierError = "";
  let homepageError = "";

  // Validate identifier format (alphanumeric, underscore, hyphen only)
  function validateIdentifier() {
    if (!identifier.trim()) {
      identifierError = "Identifier is required";
      return false;
    }
    const regex = /^[a-zA-Z0-9_-]+$/;
    if (!regex.test(identifier)) {
      identifierError =
        "Only alphanumeric characters, underscores, and hyphens allowed";
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
        homepageError = "Invalid URL format";
        return false;
      }
    }
    homepageError = "";
    return true;
  }

  async function handleSubmit() {
    error = "";

    if (!title.trim()) {
      error = "Title is required";
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
      error = err.message || "Failed to create project";
    }
  }
</script>

<form
  on:submit|preventDefault={handleSubmit}
  class="space-y-5"
>
  {#if error}
    <div class="p-3 bg-red-100 text-red-700 rounded-lg text-sm">
      {error}
    </div>
  {/if}

  <div>
    <label for="title" class="block text-sm font-medium text-slate-700 mb-1.5"
      >Title <span class="text-red-500">*</span></label
    >
    <input
      type="text"
      id="title"
      bind:value={title}
      class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
      placeholder="My Awesome Project"
      required
    />
  </div>

  <div>
    <label
      for="identifier"
      class="block text-sm font-medium text-slate-700 mb-1.5"
    >
      Identifier <span class="text-red-500">*</span>
    </label>
    <input
      type="text"
      id="identifier"
      bind:value={identifier}
      on:blur={validateIdentifier}
      class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
      class:border-red-500={identifierError}
      class:ring-2={identifierError}
      class:ring-red-500={identifierError}
      placeholder="my-awesome-project"
      required
    />
    {#if identifierError}
      <p class="text-red-600 text-xs mt-1.5">{identifierError}</p>
    {:else}
      <p class="text-slate-500 text-xs mt-1.5">
        Used in URLs and API calls. Alphanumeric, underscore, and hyphen only.
      </p>
    {/if}
  </div>

  <div>
    <label
      for="description"
      class="block text-sm font-medium text-slate-700 mb-1.5">Description</label
    >
    <textarea
      id="description"
      bind:value={description}
      rows="3"
      class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent resize-none"
      placeholder="A brief description of your project..."
    />
  </div>

  <div class="grid grid-cols-2 gap-4">
    <div>
      <label for="status" class="block text-sm font-medium text-slate-700 mb-1.5"
        >Status</label
      >
      <select
        id="status"
        bind:value={status}
        class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
      >
        <option value="active">Active</option>
        <option value="completed">Completed</option>
        <option value="archived">Archived</option>
      </select>
    </div>

    <div>
      <label class="block text-sm font-medium text-slate-700 mb-1.5"
        >Visibility</label
      >
      <div class="flex items-center h-10">
        <input
          type="checkbox"
          id="is_public"
          bind:checked={is_public}
          class="w-4 h-4 text-indigo-600 border-slate-300 rounded focus:ring-indigo-500"
        />
        <label for="is_public" class="ml-2 text-sm text-slate-700">
          Public project
        </label>
      </div>
    </div>
  </div>

  <div>
    <label for="homepage" class="block text-sm font-medium text-slate-700 mb-1.5"
      >Homepage URL</label
    >
    <input
      type="url"
      id="homepage"
      bind:value={homepage}
      on:blur={validateHomepage}
      class="w-full px-3 py-2 border border-slate-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
      class:border-red-500={homepageError}
      class:ring-2={homepageError}
      class:ring-red-500={homepageError}
      placeholder="https://github.com/username/project"
    />
    {#if homepageError}
      <p class="text-red-600 text-xs mt-1.5">{homepageError}</p>
    {:else}
      <p class="text-slate-500 text-xs mt-1.5">
        Optional URL to project homepage or repository
      </p>
    {/if}
  </div>

  <button
    type="submit"
    disabled={!title.trim() || !identifier.trim()}
    class="w-full bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-white px-4 py-2.5 rounded-lg transition-colors font-medium"
  >
    Create Project
  </button>
</form>
