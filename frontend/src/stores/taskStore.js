import { writable } from 'svelte/store';
import { api } from '../lib/api.js';

function createTaskStore() {
  const { subscribe, set, update } = writable([]);

  return {
    subscribe,
    load: async (projectId) => {
      try {
        const tasks = await api.tasks.getByProject(projectId);
        set(Array.isArray(tasks) ? tasks : []);
      } catch (error) {
        console.error('Failed to load tasks:', error);
        set([]);
      }
    },
    create: async (projectId, taskData) => {
      const task = await api.tasks.create(projectId, taskData);
      update(currentTasks => [task, ...currentTasks]);
      return task;
    },
    update: async (id, taskData) => {
      const task = await api.tasks.update(id, taskData);
      update(currentTasks => currentTasks.map(t => t.id === id ? task : t));
      return task;
    },
    toggleComplete: async (id) => {
      const task = await api.tasks.toggleComplete(id);
      update(currentTasks => currentTasks.map(t => t.id === id ? task : t));
      return task;
    },
    delete: async (id) => {
      await api.tasks.delete(id);
      update(currentTasks => currentTasks.filter(t => t.id !== id));
    },
    getById: async (id) => {
      const task = await api.tasks.get(id);
      update(currentTasks => {
        const index = currentTasks.findIndex(t => t.id === id);
        if (index >= 0) {
          currentTasks[index] = task;
        }
        return currentTasks;
      });
      return task;
    }
  };
}

export const tasks = createTaskStore();
