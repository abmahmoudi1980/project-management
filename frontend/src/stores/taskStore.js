import { writable } from 'svelte/store';
import { api } from '../lib/api.js';

function createTaskStore() {
  const { subscribe, set } = writable([]);

  return {
    subscribe,
    load: async (projectId) => {
      const tasks = await api.tasks.getByProject(projectId);
      set(tasks);
    },
    create: async (projectId, taskData) => {
      const task = await api.tasks.create(projectId, taskData);
      set(currentTasks => [task, ...currentTasks]);
      return task;
    },
    update: async (id, taskData) => {
      const task = await api.tasks.update(id, taskData);
      set(currentTasks => currentTasks.map(t => t.id === id ? task : t));
      return task;
    },
    toggleComplete: async (id) => {
      const task = await api.tasks.toggleComplete(id);
      set(currentTasks => currentTasks.map(t => t.id === id ? task : t));
      return task;
    },
    delete: async (id) => {
      await api.tasks.delete(id);
      set(currentTasks => currentTasks.filter(t => t.id !== id));
    }
  };
}

export const tasks = createTaskStore();
