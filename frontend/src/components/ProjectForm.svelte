<script>
  import { projects } from '../stores/projectStore';
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  let title = '';
  let description = '';
  let status = 'active';

  async function handleSubmit() {
    if (!title.trim()) return;

    await projects.create({
      title: title.trim(),
      description: description.trim(),
      status
    });

    title = '';
    description = '';
    status = 'active';
    dispatch('created');
  }
</script>

<form on:submit|preventDefault={handleSubmit} class="space-y-4 p-4 border rounded-lg bg-white">
  <h3 class="text-lg font-semibold text-gray-800">Create New Project</h3>

  <div>
    <label for="title" class="block text-sm font-medium text-gray-700 mb-1">Title</label>
    <input
      type="text"
      id="title"
      bind:value={title}
      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      placeholder="Project title"
    />
  </div>

  <div>
    <label for="description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
    <textarea
      id="description"
      bind:value={description}
      rows="3"
      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      placeholder="Project description (optional)"
    />
  </div>

  <div>
    <label for="status" class="block text-sm font-medium text-gray-700 mb-1">Status</label>
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

  <button
    type="submit"
    disabled={!title.trim()}
    class="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 text-white px-4 py-2 rounded-lg transition-colors"
  >
    Create Project
  </button>
</form>
