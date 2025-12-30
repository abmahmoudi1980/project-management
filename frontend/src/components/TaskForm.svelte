<script>
  import { tasks } from '../stores/taskStore';
  import { createEventDispatcher } from 'svelte';

  export let project;
  const dispatch = createEventDispatcher();

  let title = '';
  let priority = 'Medium';

  async function handleSubmit() {
    if (!title.trim()) return;

    await tasks.create(project.id, {
      title: title.trim(),
      priority
    });

    title = '';
    priority = 'Medium';
    dispatch('created');
  }
</script>

<form on:submit|preventDefault={handleSubmit} class="space-y-4 p-4 border rounded-lg bg-white">
  <h3 class="text-lg font-semibold text-gray-800">Create New Task</h3>

  <div>
    <label for="task-title" class="block text-sm font-medium text-gray-700 mb-1">Title</label>
    <input
      type="text"
      id="task-title"
      bind:value={title}
      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
      placeholder="Task title"
    />
  </div>

  <div>
    <label for="priority" class="block text-sm font-medium text-gray-700 mb-1">Priority</label>
    <select
      id="priority"
      bind:value={priority}
      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
    >
      <option value="Low">Low</option>
      <option value="Medium">Medium</option>
      <option value="High">High</option>
    </select>
  </div>

  <button
    type="submit"
    disabled={!title.trim()}
    class="w-full bg-blue-500 hover:bg-blue-600 disabled:bg-gray-300 text-white px-4 py-2 rounded-lg transition-colors"
  >
    Create Task
  </button>
</form>
