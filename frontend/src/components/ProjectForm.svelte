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
  class="space-y-4 p-4 border rounded-lg bg-white"
>
  <h3 class="text-lg font-semibold text-gray-800">Create New Project</h3>

  {#if error}
    <div class="p-3 bg-red-100 text-red-700 rounded-lg text-sm">
      {error}
    </div>
  {/if}

  <div>
    <label for="title" class="block text-sm font-medium text-gray-700 mb-1"
      >Title</label
    >
    <input
      type="text"
      id="title"
      bind:value={title}
      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      placeholder="Project title"
      required
    />
  </div>

  <div>
    <label
      for="identifier"
      class="block text-sm font-medium text-gray-700 mb-1"
    >
      Identifier <span class="text-red-500">*</span>
    </label>
    <input
      type="text"
      id="identifier"
      bind:value={identifier}
      on:blur={validateIdentifier}
      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      class:border-red-500={identifierError}
      placeholder="project-identifier"
      required
    />
    {#if identifierError}
      <p class="text-red-500 text-xs mt-1">{identifierError}</p>
    {:else}
      <p class="text-gray-500 text-xs mt-1">
        Alphanumeric, underscore, and hyphen only
      </p>
    {/if}
  </div>

  <div>
    <label
      for="description"
      class="block text-sm font-medium text-gray-700 mb-1">Description</label
    >
    <textarea
      id="description"
      bind:value={description}
      rows="3"
      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      placeholder="Project description (optional)"
    />
  </div>

  <div>
    <label for="homepage" class="block text-sm font-medium text-gray-700 mb-1"
      >Homepage URL</label
    >
    <input
      type="url"
      id="homepage"
      bind:value={homepage}
      on:blur={validateHomepage}
      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      class:border-red-500={homepageError}
      placeholder="https://example.com"
    />
    {#if homepageError}
      <p class="text-red-500 text-xs mt-1">{homepageError}</p>
    {:else}
      <p class="text-gray-500 text-xs mt-1">
        Optional project homepage or repository URL
      </p>
    {/if}
  </div>

  <div>
    <label for="status" class="block text-sm font-medium text-gray-700 mb-1"
      >Status</label
    >
    <select
      id="status"
      bind:value={status}
      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
    >
      <option value="active">Active</option>
      <option value="completed">Completed</option>
      <option value="archived">Archived</option>
    </select>
  </div>

  <div class="flex items-center">
    <input
      type="checkbox"
      id="is_public"
      bind:checked={is_public}
      class="w-4 h-4 text-blue-500 border-gray-300 rounded focus:ring-blue-500"
    />
    <label for="is_public" class="ml-2 text-sm font-medium text-gray-700">
      Make this project public
    </label>
  </div>

  <button
    type="submit"
    disabled={!title.trim() || !identifier.trim()}
    class="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 text-white px-4 py-2 rounded-lg transition-colors"
  >
    Create Project
  </button>
</form>
